'use client'

import { AlertCircle, CheckCircle2, Clock, AlertTriangle, Filter } from 'lucide-react'
import { DashboardLayout } from '../../components/dashboard-layout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../../components/ui/card'
import { Badge } from '../../components/ui/badge'
import { Button } from '../../components/ui/button'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '../../components/ui/select'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '../../components/ui/table'
import { cn } from '../../lib/utils'
import { getWebsites } from '../../lib/mock-data'

// Generate mock incidents for all websites
function generateIncidents() {
  const websites = getWebsites()
  const incidents = []

  for (const website of websites) {
    // Generate 2-4 incidents per website
    const numIncidents = Math.floor(Math.random() * 3) + 2

    for (let i = 0; i < numIncidents; i++) {
      const daysAgo = Math.floor(Math.random() * 30) + 1
      const duration = Math.floor(Math.random() * 120) + 5
      const types = ['downtime', 'degraded', 'timeout'] as const
      const type = types[Math.floor(Math.random() * types.length)]

      incidents.push({
        id: `${website.id}-${i}`,
        websiteId: website.id,
        websiteName: website.name,
        websiteUrl: website.url,
        startTime: new Date(Date.now() - daysAgo * 24 * 60 * 60 * 1000).toISOString(),
        endTime: new Date(Date.now() - daysAgo * 24 * 60 * 60 * 1000 + duration * 60 * 1000).toISOString(),
        duration,
        type,
        description: getDescription(type),
        resolved: true,
      })
    }
  }

  // Sort by start time, most recent first
  return incidents.sort((a, b) =>
    new Date(b.startTime).getTime() - new Date(a.startTime).getTime()
  )
}

function getDescription(type: 'downtime' | 'degraded' | 'timeout') {
  const descriptions = {
    downtime: [
      'Server returned 503 Service Unavailable',
      'Connection refused by remote server',
      'DNS resolution failed',
      'Server returned 500 Internal Server Error',
    ],
    degraded: [
      'High latency detected (>500ms)',
      'Slow response times across multiple regions',
      'Partial service degradation',
      'Elevated error rate detected',
    ],
    timeout: [
      'Connection timeout from EU-West region',
      'Request timeout after 30 seconds',
      'TCP connection timeout',
      'SSL handshake timeout',
    ],
  }

  const options = descriptions[type]
  return options[Math.floor(Math.random() * options.length)]
}

const typeConfig = {
  downtime: {
    icon: AlertCircle,
    label: 'Downtime',
    className: 'bg-down/10 text-down border-down/20',
  },
  degraded: {
    icon: AlertTriangle,
    label: 'Degraded',
    className: 'bg-warning/10 text-warning border-warning/20',
  },
  timeout: {
    icon: Clock,
    label: 'Timeout',
    className: 'bg-warning/10 text-warning border-warning/20',
  },
}

export default function IncidentsPage() {
  const incidents = generateIncidents()

  const totalIncidents = incidents.length
  const downtimeIncidents = incidents.filter(i => i.type === 'downtime').length
  const avgDuration = Math.round(incidents.reduce((sum, i) => sum + i.duration, 0) / incidents.length)

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  const formatDuration = (minutes: number) => {
    if (minutes < 60) return `${minutes}m`
    const hours = Math.floor(minutes / 60)
    const mins = minutes % 60
    return mins > 0 ? `${hours}h ${mins}m` : `${hours}h`
  }

  return (
    <DashboardLayout title="Incidents">
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold tracking-tight">Incidents</h1>
            <p className="text-muted-foreground">
              View and analyze all downtime events across your services
            </p>
          </div>
        </div>

        {/* Stats cards */}
        <div className="grid gap-4 md:grid-cols-4">
          <Card>
            <CardHeader className="pb-2">
              <CardDescription>Total Incidents (30d)</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold">{totalIncidents}</div>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="pb-2">
              <CardDescription>Downtime Events</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold text-down">{downtimeIncidents}</div>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="pb-2">
              <CardDescription>Avg Duration</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold font-mono">{avgDuration}m</div>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="pb-2">
              <CardDescription>MTTR</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold font-mono">12m</div>
              <p className="text-xs text-muted-foreground mt-1">Mean Time To Recovery</p>
            </CardContent>
          </Card>
        </div>

        {/* Incidents table */}
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle>All Incidents</CardTitle>
                <CardDescription>
                  Complete history of service disruptions
                </CardDescription>
              </div>
              <div className="flex items-center gap-2">
                <Select defaultValue="all">
                  <SelectTrigger className="w-[140px]">
                    <Filter className="mr-2 size-4" />
                    <SelectValue placeholder="Filter type" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">All Types</SelectItem>
                    <SelectItem value="downtime">Downtime</SelectItem>
                    <SelectItem value="degraded">Degraded</SelectItem>
                    <SelectItem value="timeout">Timeout</SelectItem>
                  </SelectContent>
                </Select>
                <Select defaultValue="30">
                  <SelectTrigger className="w-[140px]">
                    <SelectValue placeholder="Time range" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="7">Last 7 days</SelectItem>
                    <SelectItem value="30">Last 30 days</SelectItem>
                    <SelectItem value="90">Last 90 days</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Service</TableHead>
                  <TableHead>Type</TableHead>
                  <TableHead>Description</TableHead>
                  <TableHead>Started</TableHead>
                  <TableHead>Duration</TableHead>
                  <TableHead className="text-right">Status</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {incidents.slice(0, 20).map((incident) => {
                  const config = typeConfig[incident.type]
                  const TypeIcon = config.icon

                  return (
                    <TableRow key={incident.id}>
                      <TableCell>
                        <div>
                          <div className="font-medium">{incident.websiteName}</div>
                          <div className="text-xs text-muted-foreground truncate max-w-[200px]">
                            {incident.websiteUrl}
                          </div>
                        </div>
                      </TableCell>
                      <TableCell>
                        <Badge variant="outline" className={cn("gap-1", config.className)}>
                          <TypeIcon className="size-3" />
                          {config.label}
                        </Badge>
                      </TableCell>
                      <TableCell className="max-w-[250px] truncate">
                        {incident.description}
                      </TableCell>
                      <TableCell className="text-muted-foreground whitespace-nowrap">
                        {formatDate(incident.startTime)}
                      </TableCell>
                      <TableCell className="font-mono">
                        {formatDuration(incident.duration)}
                      </TableCell>
                      <TableCell className="text-right">
                        {incident.resolved ? (
                          <Badge variant="outline" className="bg-uptime/10 text-uptime border-uptime/20 gap-1">
                            <CheckCircle2 className="size-3" />
                            Resolved
                          </Badge>
                        ) : (
                          <Badge variant="outline" className="bg-down/10 text-down border-down/20 gap-1">
                            <AlertCircle className="size-3" />
                            Ongoing
                          </Badge>
                        )}
                      </TableCell>
                    </TableRow>
                  )
                })}
              </TableBody>
            </Table>

            <div className="flex items-center justify-center pt-4">
              <Button variant="outline">Load More</Button>
            </div>
          </CardContent>
        </Card>
      </div>
    </DashboardLayout>
  )
}
