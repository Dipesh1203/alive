'use client'

import { useEffect, useMemo, useState } from 'react'
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
import { fetchAllWebsiteDetails } from '../../lib/api'
import type { Incident, WebsiteDetails } from '../../lib/mock-data'

type IncidentListRow = Incident & { websiteName: string; websiteUrl: string }

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
  const [websites, setWebsites] = useState<WebsiteDetails[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const loadIncidents = async () => {
      setIsLoading(true)
      setError(null)

      try {
        const websiteDetails = await fetchAllWebsiteDetails()
        setWebsites(websiteDetails)
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load incidents')
      } finally {
        setIsLoading(false)
      }
    }

    void loadIncidents()
  }, [])

  const incidents = useMemo<IncidentListRow[]>(() => {
    const rows = websites.flatMap((website) =>
      website.incidents.map((incident) => ({
        ...incident,
        websiteName: website.name,
        websiteUrl: website.url,
      })),
    )

    return rows.sort((a, b) =>
      new Date(b.startTime).getTime() - new Date(a.startTime).getTime(),
    )
  }, [websites])

  const totalIncidents = incidents.length
  const downtimeIncidents = incidents.filter(i => i.type === 'downtime').length
  const avgDuration = incidents.length
    ? Math.round(incidents.reduce((sum, i) => sum + i.duration, 0) / incidents.length)
    : 0

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

        {isLoading && <p className="text-sm text-muted-foreground">Loading incidents...</p>}
        {error && <p className="text-sm text-down">{error}</p>}

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
