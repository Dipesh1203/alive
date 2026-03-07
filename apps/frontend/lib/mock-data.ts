// Types for the Alive uptime monitoring application
export interface Website {
  id: string
  name: string
  url: string
  status: 'up' | 'down' | 'degraded'
  uptime30d: number
  currentLatency: number
  avgLatency: number
  monitoringEnabled: boolean
  regions: string[]
  lastChecked: string
  createdAt: string
}

export interface WebsiteDetails extends Website {
  latencyHistory: LatencyDataPoint[]
  uptimeCalendar: UptimeDay[]
  regionLatencies: RegionLatency[]
  incidents: Incident[]
}

export interface LatencyDataPoint {
  timestamp: string
  latency: number
}

export interface UptimeDay {
  date: string
  status: 'up' | 'down' | 'degraded' | 'partial'
  uptime: number
}

export interface RegionLatency {
  region: string
  code: string
  latency: number
  status: 'up' | 'down' | 'degraded'
}

export interface Incident {
  id: string
  websiteId: string
  startTime: string
  endTime: string | null
  duration: number
  type: 'downtime' | 'degraded' | 'timeout'
  description: string
  resolved: boolean
}

export interface NotificationChannel {
  id: string
  type: 'slack' | 'email' | 'webhook' | 'pagerduty'
  name: string
  enabled: boolean
  config: Record<string, string>
}

export interface TeamMember {
  id: string
  name: string
  email: string
  role: 'owner' | 'admin' | 'member'
  avatar: string
}

export interface Region {
  id: string
  name: string
  code: string
  available: boolean
}

// Mock data generators
function generateLatencyHistory(hours: number = 24): LatencyDataPoint[] {
  const data: LatencyDataPoint[] = []
  const now = new Date()
  
  for (let i = hours; i >= 0; i--) {
    const timestamp = new Date(now.getTime() - i * 60 * 60 * 1000)
    const baseLatency = 120 + Math.random() * 80
    const spike = Math.random() > 0.95 ? Math.random() * 200 : 0
    
    data.push({
      timestamp: timestamp.toISOString(),
      latency: Math.round(baseLatency + spike),
    })
  }
  
  return data
}

function generateUptimeCalendar(days: number = 90): UptimeDay[] {
  const data: UptimeDay[] = []
  const now = new Date()
  
  for (let i = days - 1; i >= 0; i--) {
    const date = new Date(now.getTime() - i * 24 * 60 * 60 * 1000)
    const rand = Math.random()
    
    let status: UptimeDay['status']
    let uptime: number
    
    if (rand > 0.98) {
      status = 'down'
      uptime = 85 + Math.random() * 10
    } else if (rand > 0.92) {
      status = 'degraded'
      uptime = 95 + Math.random() * 4
    } else if (rand > 0.88) {
      status = 'partial'
      uptime = 99 + Math.random() * 0.9
    } else {
      status = 'up'
      uptime = 99.9 + Math.random() * 0.1
    }
    
    data.push({
      date: date.toISOString().split('T')[0] ?? date.toISOString(),
      status,
      uptime: Math.round(uptime * 100) / 100,
    })
  }
  
  return data
}

// Mock websites data
export const mockWebsites: Website[] = [
  {
    id: '1',
    name: 'Production API',
    url: 'https://api.example.com',
    status: 'up',
    uptime30d: 99.98,
    currentLatency: 145,
    avgLatency: 132,
    monitoringEnabled: true,
    regions: ['us-east-1', 'eu-west-1', 'ap-southeast-1'],
    lastChecked: new Date(Date.now() - 30000).toISOString(),
    createdAt: '2024-01-15T10:00:00Z',
  },
  {
    id: '2',
    name: 'Marketing Website',
    url: 'https://www.example.com',
    status: 'up',
    uptime30d: 99.95,
    currentLatency: 89,
    avgLatency: 95,
    monitoringEnabled: true,
    regions: ['us-east-1', 'eu-west-1'],
    lastChecked: new Date(Date.now() - 45000).toISOString(),
    createdAt: '2024-02-01T08:30:00Z',
  },
  {
    id: '3',
    name: 'Dashboard App',
    url: 'https://dashboard.example.com',
    status: 'degraded',
    uptime30d: 98.75,
    currentLatency: 450,
    avgLatency: 180,
    monitoringEnabled: true,
    regions: ['us-east-1', 'us-west-2', 'eu-west-1', 'ap-northeast-1'],
    lastChecked: new Date(Date.now() - 20000).toISOString(),
    createdAt: '2024-01-20T14:00:00Z',
  },
  {
    id: '4',
    name: 'Auth Service',
    url: 'https://auth.example.com',
    status: 'up',
    uptime30d: 99.99,
    currentLatency: 52,
    avgLatency: 48,
    monitoringEnabled: true,
    regions: ['us-east-1', 'us-west-2', 'eu-west-1', 'eu-central-1', 'ap-southeast-1'],
    lastChecked: new Date(Date.now() - 15000).toISOString(),
    createdAt: '2024-01-10T09:00:00Z',
  },
  {
    id: '5',
    name: 'CDN Origin',
    url: 'https://cdn.example.com',
    status: 'down',
    uptime30d: 97.50,
    currentLatency: 0,
    avgLatency: 35,
    monitoringEnabled: true,
    regions: ['us-east-1'],
    lastChecked: new Date(Date.now() - 10000).toISOString(),
    createdAt: '2024-03-01T11:00:00Z',
  },
  {
    id: '6',
    name: 'Staging Environment',
    url: 'https://staging.example.com',
    status: 'up',
    uptime30d: 99.85,
    currentLatency: 210,
    avgLatency: 195,
    monitoringEnabled: false,
    regions: ['us-east-1'],
    lastChecked: new Date(Date.now() - 60000).toISOString(),
    createdAt: '2024-02-15T16:00:00Z',
  },
]

export const mockRegions: Region[] = [
  { id: 'us-east-1', name: 'US East (N. Virginia)', code: 'US-E', available: true },
  { id: 'us-west-2', name: 'US West (Oregon)', code: 'US-W', available: true },
  { id: 'eu-west-1', name: 'Europe (Ireland)', code: 'EU-W', available: true },
  { id: 'eu-central-1', name: 'Europe (Frankfurt)', code: 'EU-C', available: true },
  { id: 'ap-southeast-1', name: 'Asia Pacific (Singapore)', code: 'AP-S', available: true },
  { id: 'ap-northeast-1', name: 'Asia Pacific (Tokyo)', code: 'AP-N', available: true },
  { id: 'sa-east-1', name: 'South America (São Paulo)', code: 'SA-E', available: false },
]

export const mockNotificationChannels: NotificationChannel[] = [
  {
    id: '1',
    type: 'slack',
    name: 'Engineering Alerts',
    enabled: true,
    config: { webhook: 'https://hooks.slack.com/services/xxx' },
  },
  {
    id: '2',
    type: 'email',
    name: 'Team Email',
    enabled: true,
    config: { emails: 'team@example.com,alerts@example.com' },
  },
  {
    id: '3',
    type: 'pagerduty',
    name: 'On-Call',
    enabled: false,
    config: { serviceKey: 'xxxx-xxxx-xxxx' },
  },
  {
    id: '4',
    type: 'webhook',
    name: 'Custom Integration',
    enabled: true,
    config: { url: 'https://webhooks.example.com/alerts' },
  },
]

export const mockTeamMembers: TeamMember[] = [
  {
    id: '1',
    name: 'Alex Johnson',
    email: 'alex@example.com',
    role: 'owner',
    avatar: '',
  },
  {
    id: '2',
    name: 'Sarah Chen',
    email: 'sarah@example.com',
    role: 'admin',
    avatar: '',
  },
  {
    id: '3',
    name: 'Mike Williams',
    email: 'mike@example.com',
    role: 'member',
    avatar: '',
  },
  {
    id: '4',
    name: 'Emily Davis',
    email: 'emily@example.com',
    role: 'member',
    avatar: '',
  },
]

// Mock API functions
export function getWebsites(): Website[] {
  return mockWebsites
}

export function getWebsiteDetails(id: string): WebsiteDetails | null {
  const website = mockWebsites.find(w => w.id === id)
  if (!website) return null
  
  const latencyHistory = generateLatencyHistory(24)
  const uptimeCalendar = generateUptimeCalendar(90)
  
  const regionLatencies: RegionLatency[] = website.regions.map(regionId => {
    const region = mockRegions.find(r => r.id === regionId)
    const baseLatency = website.currentLatency + (Math.random() - 0.5) * 50
    
    return {
      region: region?.name || regionId,
      code: region?.code || regionId,
      latency: Math.max(0, Math.round(baseLatency)),
      status: website.status,
    }
  })
  
  const incidents: Incident[] = [
    {
      id: '1',
      websiteId: id,
      startTime: new Date(Date.now() - 3 * 24 * 60 * 60 * 1000).toISOString(),
      endTime: new Date(Date.now() - 3 * 24 * 60 * 60 * 1000 + 15 * 60 * 1000).toISOString(),
      duration: 15,
      type: 'timeout',
      description: 'Connection timeout from EU-West region',
      resolved: true,
    },
    {
      id: '2',
      websiteId: id,
      startTime: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString(),
      endTime: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000 + 45 * 60 * 1000).toISOString(),
      duration: 45,
      type: 'downtime',
      description: 'Server returned 503 Service Unavailable',
      resolved: true,
    },
    {
      id: '3',
      websiteId: id,
      startTime: new Date(Date.now() - 14 * 24 * 60 * 60 * 1000).toISOString(),
      endTime: new Date(Date.now() - 14 * 24 * 60 * 60 * 1000 + 5 * 60 * 1000).toISOString(),
      duration: 5,
      type: 'degraded',
      description: 'High latency detected (>500ms)',
      resolved: true,
    },
  ]
  
  return {
    ...website,
    latencyHistory,
    uptimeCalendar,
    regionLatencies,
    incidents,
  }
}

export function getRegions(): Region[] {
  return mockRegions
}

export function getNotificationChannels(): NotificationChannel[] {
  return mockNotificationChannels
}

export function getTeamMembers(): TeamMember[] {
  return mockTeamMembers
}

export function calculateGlobalUptime(): number {
  const enabledWebsites = mockWebsites.filter(w => w.monitoringEnabled)
  if (enabledWebsites.length === 0) return 100
  
  const totalUptime = enabledWebsites.reduce((sum, w) => sum + w.uptime30d, 0)
  return Math.round((totalUptime / enabledWebsites.length) * 1000) / 1000
}
