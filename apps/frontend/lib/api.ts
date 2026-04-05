import type {
  Incident,
  LatencyDataPoint,
  Region,
  RegionLatency,
  UptimeDay,
  Website,
  WebsiteDetails,
} from './mock-data'

const API_BASE_URL = (process.env.NEXT_PUBLIC_API_BASE_URL ?? 'http://localhost:8000').replace(/\/$/, '')

interface ApiWebsite {
  id: string
  websiteName: string
  url: string
  status: 'Up' | 'Down' | 'Unknown'
  isMonitoringEnabled: boolean
  createdAt: string
  updatedAt: string
}

interface ApiRegion {
  regionId: string
  regionName: string
  createdAt: string
  updatedAt: string
}

interface ApiTick {
  id: number
  websiteId: string
  upStatus: 'Up' | 'Down' | 'Unknown'
  latency: number | null
  websiteRegionId: string
  createdAt: string
  updatedAt: string
}

interface CreateWebsitePayload {
  websiteName: string
  url: string
  organizationId: string
}

interface ApiAuthUser {
  id: string
  email: string
  createdAt?: string
  updatedAt?: string
}

interface ApiAuthResponse {
  message?: string
  token?: string
  error?: string
  user?: ApiAuthUser
}

export interface AuthPayload {
  token: string
  user: ApiAuthUser
}

function getAuthToken(): string | null {
  if (typeof window === 'undefined') {
    return null
  }

  return localStorage.getItem('auth_token')
}

async function apiFetch<T>(path: string, init?: RequestInit): Promise<T> {
  const authToken = getAuthToken()

  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...init,
    headers: {
      'Content-Type': 'application/json',
      ...(authToken ? { Authorization: `Bearer ${authToken}` } : {}),
      ...(init?.headers ?? {}),
    },
    cache: 'no-store',
  })

  if (!response.ok) {
    const errorBody = await response.text()

    try {
      const parsed = JSON.parse(errorBody) as { error?: string; message?: string }
      throw new Error(parsed.error || parsed.message || `Request failed (${response.status})`)
    } catch {
      throw new Error(errorBody || `Request failed (${response.status})`)
    }
  }

  return response.json() as Promise<T>
}

function mapAuthPayload(response: ApiAuthResponse): AuthPayload {
  if (!response.token || !response.user) {
    throw new Error(response.error || response.message || 'Authentication failed')
  }

  return {
    token: response.token,
    user: response.user,
  }
}

function normalizeStatus(status: ApiWebsite['status'] | ApiTick['upStatus']): Website['status'] {
  if (status === 'Up') return 'up'
  if (status === 'Down') return 'down'
  return 'degraded'
}

function buildRegionCode(name: string): string {
  return name
    .split(/\s+/)
    .filter(Boolean)
    .slice(0, 2)
    .map((chunk) => chunk[0]?.toUpperCase() ?? '')
    .join('') || 'RG'
}

function getWebsiteTicksPath(websiteId: string): string {
  const now = new Date()
  const start = new Date(now)
  start.setDate(start.getDate() - 90)

  const startDate = start.toISOString().slice(0, 10)
  const endDate = now.toISOString().slice(0, 10)

  return `/api/websites/${websiteId}/details?take=1000&startDate=${startDate}&endDate=${endDate}`
}

function mapRegions(apiRegions: ApiRegion[]): Region[] {
  return apiRegions.map((region) => ({
    id: region.regionId,
    name: region.regionName,
    code: buildRegionCode(region.regionName),
    available: true,
  }))
}

function calculateUptimePercent(ticks: ApiTick[], fallbackStatus: ApiWebsite['status']): number {
  if (ticks.length === 0) {
    if (fallbackStatus === 'Up') return 100
    if (fallbackStatus === 'Down') return 0
    return 99
  }

  const upTicks = ticks.filter((tick) => tick.upStatus === 'Up').length
  return Math.round((upTicks / ticks.length) * 10000) / 100
}

function calculateLatencyMetrics(ticks: ApiTick[]): { currentLatency: number; avgLatency: number } {
  const latencyTicks = ticks.map((tick) => tick.latency).filter((value): value is number => value !== null)

  if (latencyTicks.length === 0) {
    return { currentLatency: 0, avgLatency: 0 }
  }

  const currentLatency = latencyTicks[0] ?? 0
  const avgLatency = Math.round(latencyTicks.reduce((sum, value) => sum + value, 0) / latencyTicks.length)
  return { currentLatency, avgLatency }
}

function mapWebsite(
  website: ApiWebsite,
  ticks: ApiTick[],
  regionsById: Map<string, Region>,
): Website {
  const latestTick = ticks[0]
  const effectiveStatus = latestTick ? normalizeStatus(latestTick.upStatus) : normalizeStatus(website.status)
  const latency = calculateLatencyMetrics(ticks)

  const uniqueRegionIds = Array.from(new Set(ticks.map((tick) => tick.websiteRegionId)))
  const regionNames = uniqueRegionIds.map((regionId) => regionsById.get(regionId)?.name ?? regionId)

  return {
    id: website.id,
    name: website.websiteName,
    url: website.url,
    status: effectiveStatus,
    uptime30d: calculateUptimePercent(ticks, website.status),
    currentLatency: latency.currentLatency,
    avgLatency: latency.avgLatency,
    monitoringEnabled: website.isMonitoringEnabled,
    regions: regionNames,
    lastChecked: latestTick?.createdAt ?? website.updatedAt,
    createdAt: website.createdAt,
  }
}

function buildLatencyHistory(ticks: ApiTick[]): LatencyDataPoint[] {
  const points = ticks
    .filter((tick) => tick.latency !== null)
    .map((tick) => ({
      timestamp: tick.createdAt,
      latency: tick.latency ?? 0,
    }))

  return points.reverse()
}

function buildUptimeCalendar(ticks: ApiTick[]): UptimeDay[] {
  const groupedByDay = new Map<string, ApiTick[]>()

  for (const tick of ticks) {
    const key = tick.createdAt.slice(0, 10)
    const existing = groupedByDay.get(key) ?? []
    existing.push(tick)
    groupedByDay.set(key, existing)
  }

  const calendar: UptimeDay[] = []
  const today = new Date()

  for (let i = 89; i >= 0; i -= 1) {
    const day = new Date(today)
    day.setDate(today.getDate() - i)
    const dateKey = day.toISOString().slice(0, 10)
    const dayTicks = groupedByDay.get(dateKey) ?? []

    if (dayTicks.length === 0) {
      calendar.push({ date: dateKey, status: 'partial', uptime: 100 })
      continue
    }

    const upCount = dayTicks.filter((tick) => tick.upStatus === 'Up').length
    const uptime = (upCount / dayTicks.length) * 100

    let status: UptimeDay['status'] = 'up'
    if (uptime < 90) status = 'down'
    else if (uptime < 99) status = 'degraded'
    else if (uptime < 100) status = 'partial'

    calendar.push({
      date: dateKey,
      status,
      uptime: Math.round(uptime * 100) / 100,
    })
  }

  return calendar
}

function buildRegionLatencies(ticks: ApiTick[], regionsById: Map<string, Region>): RegionLatency[] {
  const latestByRegion = new Map<string, ApiTick>()

  for (const tick of ticks) {
    if (!latestByRegion.has(tick.websiteRegionId)) {
      latestByRegion.set(tick.websiteRegionId, tick)
    }
  }

  return Array.from(latestByRegion.entries()).map(([regionId, tick]) => {
    const region = regionsById.get(regionId)

    return {
      region: region?.name ?? regionId,
      code: region?.code ?? regionId,
      latency: tick.latency ?? 0,
      status: normalizeStatus(tick.upStatus),
    }
  })
}

function buildIncidents(ticks: ApiTick[]): Incident[] {
  return ticks
    .filter((tick) => tick.upStatus !== 'Up')
    .slice(0, 20)
    .map((tick) => ({
      id: String(tick.id),
      websiteId: tick.websiteId,
      startTime: tick.createdAt,
      endTime: tick.updatedAt,
      duration: 5,
      type: tick.upStatus === 'Down' ? 'downtime' : 'degraded',
      description:
        tick.upStatus === 'Down'
          ? 'Monitor detected a failed availability check'
          : 'Monitor detected degraded or unknown status',
      resolved: true,
    }))
}

export async function fetchRegions(): Promise<Region[]> {
  const apiRegions = await apiFetch<ApiRegion[]>('/api/regions')
  return mapRegions(apiRegions)
}

export async function createWebsite(payload: CreateWebsitePayload): Promise<void> {
  await apiFetch<ApiWebsite>('/api/websites', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export async function toggleWebsiteMonitoring(websiteId: string): Promise<void> {
  await apiFetch<ApiWebsite>(`/api/monitoring/${websiteId}`, {
    method: 'POST',
  })
}

export async function fetchWebsites(): Promise<Website[]> {
  const [apiWebsites, regionOptions] = await Promise.all([
    apiFetch<ApiWebsite[]>('/api/websites'),
    fetchRegions().catch(() => [] as Region[]),
  ])

  const regionsById = new Map(regionOptions.map((region) => [region.id, region]))

  const mapped = await Promise.all(
    apiWebsites.map(async (website) => {
      try {
        const ticks = await apiFetch<ApiTick[]>(getWebsiteTicksPath(website.id))
        return mapWebsite(website, ticks, regionsById)
      } catch {
        return mapWebsite(website, [], regionsById)
      }
    }),
  )

  return mapped
}

export async function fetchWebsiteDetails(websiteId: string): Promise<WebsiteDetails> {
  const [apiWebsite, ticks, regionOptions] = await Promise.all([
    apiFetch<ApiWebsite>(`/api/websites/${websiteId}`),
    apiFetch<ApiTick[]>(getWebsiteTicksPath(websiteId)),
    fetchRegions().catch(() => [] as Region[]),
  ])

  const regionsById = new Map(regionOptions.map((region) => [region.id, region]))
  const website = mapWebsite(apiWebsite, ticks, regionsById)

  return {
    ...website,
    latencyHistory: buildLatencyHistory(ticks),
    uptimeCalendar: buildUptimeCalendar(ticks),
    regionLatencies: buildRegionLatencies(ticks, regionsById),
    incidents: buildIncidents(ticks),
  }
}

export async function fetchAllWebsiteDetails(): Promise<WebsiteDetails[]> {
  const websites = await fetchWebsites()
  const details = await Promise.all(
    websites.map(async (website) => {
      try {
        return await fetchWebsiteDetails(website.id)
      } catch {
        return {
          ...website,
          latencyHistory: [],
          uptimeCalendar: [],
          regionLatencies: [],
          incidents: [],
        }
      }
    }),
  )

  return details
}

// User Login and Signup
export async function login(email: string, password: string): Promise<AuthPayload> {
  const response = await apiFetch<ApiAuthResponse>('/api/auth/login', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  })

  return mapAuthPayload(response)
}

export async function signup(email: string, password: string): Promise<AuthPayload> {
  const response = await apiFetch<ApiAuthResponse>('/api/auth/signup', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  })

  return mapAuthPayload(response)
}

// Organization Management
export interface Organization {
  id: string
  name: string
  adminId: string
  createdAt: string
  updatedAt: string
}

export interface OrganizationMember {
  id: string
  userId: string
  organizationId: string
  role: 'admin' | 'viewer'
  email?: string
  createdAt: string
  updatedAt: string
}

export interface UserProfile {
  id: string
  userId: string
  firstName: string
  lastName: string
  phone: string
  bio: string
  avatar: string
  preferences?: Record<string, unknown>
  createdAt: string
  updatedAt: string
}

export interface LandingStat {
  label: string
  value: string
  hint: string
}

export interface LandingFeature {
  title: string
  description: string
}

export interface LandingTestimonial {
  name: string
  role: string
  company: string
  quote: string
  rating: number
  location: string
}

export interface LandingPricingPlan {
  name: string
  description: string
  price: number
  currency: string
  interval: 'monthly' | 'yearly'
  popular: boolean
  cta: string
  features: string[]
}

export interface LandingFAQ {
  question: string
  answer: string
}

export interface LandingOverview {
  stats: LandingStat[]
  features: LandingFeature[]
  testimonials: LandingTestimonial[]
  pricing: LandingPricingPlan[]
  faqs: LandingFAQ[]
}

export async function createOrganization(name: string): Promise<Organization> {
  return apiFetch<Organization>('/api/organizations', {
    method: 'POST',
    body: JSON.stringify({ name }),
  })
}

export async function listUserOrganizations(): Promise<Organization[]> {
  return apiFetch<Organization[]>('/api/organizations', {
    method: 'GET',
  })
}

export async function getOrganization(orgId: string): Promise<Organization> {
  return apiFetch<Organization>(`/api/organizations/${orgId}`, {
    method: 'GET',
  })
}

export async function updateOrganization(orgId: string, name: string): Promise<Organization> {
  return apiFetch<Organization>(`/api/organizations/${orgId}`, {
    method: 'PUT',
    body: JSON.stringify({ name }),
  })
}

export async function listOrganizationMembers(orgId: string): Promise<OrganizationMember[]> {
  return apiFetch<OrganizationMember[]>(`/api/organizations/${orgId}/members`, {
    method: 'GET',
  })
}

export async function addOrganizationMember(orgId: string, email: string, role: 'admin' | 'viewer' = 'viewer'): Promise<OrganizationMember> {
  return apiFetch<OrganizationMember>(`/api/organizations/${orgId}/members`, {
    method: 'POST',
    body: JSON.stringify({ email, role }),
  })
}

export async function updateMemberRole(orgId: string, memberId: string, role: 'admin' | 'viewer'): Promise<OrganizationMember> {
  return apiFetch<OrganizationMember>(`/api/organizations/${orgId}/members/${memberId}/role`, {
    method: 'PUT',
    body: JSON.stringify({ role }),
  })
}

export async function removeOrganizationMember(orgId: string, memberId: string): Promise<{ message: string }> {
  return apiFetch<{ message: string }>(`/api/organizations/${orgId}/members/${memberId}`, {
    method: 'DELETE',
  })
}

// User Profile Management
export async function getMyProfile(): Promise<UserProfile> {
  return apiFetch<UserProfile>('/api/profile', {
    method: 'GET',
  })
}

export async function updateMyProfile(
  firstName: string,
  lastName: string,
  phone: string,
  bio: string,
  avatar?: string,
): Promise<UserProfile> {
  return apiFetch<UserProfile>('/api/profile', {
    method: 'PUT',
    body: JSON.stringify({ firstName, lastName, phone, bio, avatar }),
  })
}

export async function getMyPreferences(): Promise<Record<string, unknown>> {
  return apiFetch<Record<string, unknown>>('/api/profile/preferences', {
    method: 'GET',
  })
}

export async function updateMyPreferences(preferences: Record<string, unknown>): Promise<UserProfile> {
  return apiFetch<UserProfile>('/api/profile/preferences', {
    method: 'PUT',
    body: JSON.stringify({ preferences }),
  })
}

export async function fetchLandingOverview(billing: 'monthly' | 'yearly' = 'monthly'): Promise<LandingOverview> {
  return apiFetch<LandingOverview>(`/api/public/landing?billing=${billing}`, {
    method: 'GET',
  })
}

export async function fetchLandingTestimonials(): Promise<LandingTestimonial[]> {
  return apiFetch<LandingTestimonial[]>('/api/public/landing/testimonials', {
    method: 'GET',
  })
}

export async function fetchLandingPricing(billing: 'monthly' | 'yearly' = 'monthly'): Promise<LandingPricingPlan[]> {
  return apiFetch<LandingPricingPlan[]>(`/api/public/landing/pricing?billing=${billing}`, {
    method: 'GET',
  })
}

export async function fetchLandingFaqs(): Promise<LandingFAQ[]> {
  return apiFetch<LandingFAQ[]>('/api/public/landing/faqs', {
    method: 'GET',
  })
}
