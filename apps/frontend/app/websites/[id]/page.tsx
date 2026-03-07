'use client'

import { use } from 'react'
import { useEffect, useState } from 'react'
import Link from 'next/link'
import { ArrowLeft, ExternalLink, Clock, Activity } from 'lucide-react'
import { LatencyChart } from '../../../components/latency-chart'
import { UptimeCalendar } from '../../../components/uptime-calendar'
import { RegionStatus } from '../../../components/region-status'
import { IncidentLog } from '../../../components/incident-log'
import { Button } from '../../../components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../../../components/ui/card'
import { Switch } from '../../../components/ui/switch'
import { Badge } from '../../../components/ui/badge'
import { cn } from '../../../lib/utils'
import { fetchWebsiteDetails, toggleWebsiteMonitoring } from '../../../lib/api'
import type { WebsiteDetails } from '../../../lib/mock-data'
import { DashboardLayout } from '../../../components/dashboard-layout'

interface WebsiteDetailPageProps {
  params: Promise<{ id: string }>
}

export default function WebsiteDetailPage({ params }: WebsiteDetailPageProps) {
  const { id } = use(params)
  const [website, setWebsite] = useState<WebsiteDetails | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [isToggling, setIsToggling] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const loadWebsite = async () => {
      setIsLoading(true)
      setError(null)

      try {
        const details = await fetchWebsiteDetails(id)
        setWebsite(details)
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load website details')
        setWebsite(null)
      } finally {
        setIsLoading(false)
      }
    }

    void loadWebsite()
  }, [id])

  const handleToggleMonitoring = async () => {
    if (!website) return

    setIsToggling(true)
    setError(null)

    try {
      await toggleWebsiteMonitoring(website.id)
      const details = await fetchWebsiteDetails(website.id)
      setWebsite(details)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to toggle monitoring')
    } finally {
      setIsToggling(false)
    }
  }

  if (isLoading) {
    return (
      <DashboardLayout title="Website Details">
        <p className="text-sm text-muted-foreground">Loading website details...</p>
      </DashboardLayout>
    )
  }

  if (!website) {
    return (
      <DashboardLayout title="Website Details">
        <p className="text-sm text-down">{error ?? 'Website not found'}</p>
      </DashboardLayout>
    )
  }

  const statusConfig = {
    up: {
      label: 'Operational',
      className: 'bg-uptime text-uptime-foreground',
      textClass: 'text-uptime',
    },
    down: {
      label: 'Down',
      className: 'bg-down text-down-foreground',
      textClass: 'text-down',
    },
    degraded: {
      label: 'Degraded',
      className: 'bg-warning text-warning-foreground',
      textClass: 'text-warning',
    },
  }

  const status = statusConfig[website.status]

  return (
    <DashboardLayout
      breadcrumbs={[
        { label: 'Websites', href: '/' },
        { label: website.name },
      ]}
    >
      <div className="space-y-6">
        {/* Back button and header */}
        <div className="flex items-start justify-between">
          <div className="flex items-start gap-4">
            <Button variant="ghost" size="icon" asChild>
              <Link href="/">
                <ArrowLeft className="size-4" />
                <span className="sr-only">Back to dashboard</span>
              </Link>
            </Button>
            <div>
              <div className="flex items-center gap-3">
                <h1 className="text-2xl font-bold tracking-tight">{website.name}</h1>
                <Badge className={status.className}>{status.label}</Badge>
              </div>
              <div className="flex items-center gap-4 mt-1 text-sm text-muted-foreground">
                <a
                  href={website.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="flex items-center gap-1 hover:text-foreground transition-colors"
                >
                  {website.url}
                  <ExternalLink className="size-3" />
                </a>
                <span className="flex items-center gap-1">
                  <Clock className="size-3" />
                  Last checked {getTimeAgo(website.lastChecked)}
                </span>
              </div>
            </div>
          </div>

          <div className="flex items-center gap-4">
            <Button variant="outline" asChild>
              <a href={website.url} target="_blank" rel="noopener noreferrer">
                <ExternalLink className="mr-2 size-4" />
                Visit Site
              </a>
            </Button>
          </div>
        </div>

        {error && <p className="text-sm text-down">{error}</p>}

        {/* Overview cards */}
        <div className="grid gap-4 md:grid-cols-4">
          <Card>
            <CardHeader className="pb-2">
              <CardDescription>Current Latency</CardDescription>
            </CardHeader>
            <CardContent>
              <div className={cn("text-3xl font-bold font-mono", status.textClass)}>
                {website.currentLatency === 0 ? 'N/A' : `${website.currentLatency}ms`}
              </div>
              <p className="text-xs text-muted-foreground mt-1">
                avg {website.avgLatency}ms
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-2">
              <CardDescription>30-Day Uptime</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold font-mono">
                {website.uptime30d.toFixed(2)}%
              </div>
              <p className="text-xs text-muted-foreground mt-1">
                {website.uptime30d >= 99.9 ? 'Excellent' : website.uptime30d >= 99 ? 'Good' : 'Needs attention'}
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-2">
              <CardDescription>Monitoring Regions</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold font-mono">
                {website.regions.length}
              </div>
              <p className="text-xs text-muted-foreground mt-1">
                {website.regions.length === 1 ? 'location' : 'locations'}
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-2">
              <CardDescription>Monitoring Status</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <Activity className={cn(
                    "size-5",
                    website.monitoringEnabled ? 'text-uptime' : 'text-muted-foreground'
                  )} />
                  <span className="font-medium">
                    {website.monitoringEnabled ? 'Active' : 'Paused'}
                  </span>
                </div>
                <Switch
                  checked={website.monitoringEnabled}
                  onCheckedChange={handleToggleMonitoring}
                  disabled={isToggling}
                />
              </div>
              <p className="text-xs text-muted-foreground mt-2">
                Toggle to pause/resume monitoring
              </p>
            </CardContent>
          </Card>
        </div>

        {/* Latency chart */}
        <LatencyChart
          data={website.latencyHistory}
          title="Response Time History"
          description="Performance over the last 24 hours"
        />

        {/* Uptime calendar */}
        <UptimeCalendar
          data={website.uptimeCalendar}
          title="Uptime Calendar"
          description="Daily availability over the last 90 days"
        />

        {/* Bottom section: Region status and incidents */}
        <div className="grid gap-4 lg:grid-cols-2">
          <RegionStatus
            regions={website.regionLatencies}
            title="Region Performance"
            description="Latency from each monitoring location"
          />
          <IncidentLog
            incidents={website.incidents}
            title="Recent Incidents"
            description="Downtime and degradation events"
          />
        </div>
      </div>
    </DashboardLayout>
  )
}

function getTimeAgo(dateString: string): string {
  const date = new Date(dateString)
  const now = new Date()
  const seconds = Math.floor((now.getTime() - date.getTime()) / 1000)

  if (seconds < 60) return `${seconds}s ago`
  if (seconds < 3600) return `${Math.floor(seconds / 60)}m ago`
  if (seconds < 86400) return `${Math.floor(seconds / 3600)}h ago`
  return `${Math.floor(seconds / 86400)}d ago`
}
