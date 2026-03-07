'use client'

import { AlertCircle, Clock, CheckCircle2, AlertTriangle } from 'lucide-react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './ui/card'
import { Badge } from './ui/badge'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from './ui/table'
import { cn } from '../lib/utils'
import type { Incident } from '../lib/mock-data'

interface IncidentLogProps {
  incidents: Incident[]
  title?: string
  description?: string
}

export function IncidentLog({ incidents, title = "Recent Incidents", description = "Latest downtime events" }: IncidentLogProps) {
  const getTypeConfig = (type: Incident['type']) => {
    switch (type) {
      case 'downtime':
        return {
          icon: AlertCircle,
          label: 'Downtime',
          className: 'bg-down/10 text-down border-down/20',
        }
      case 'degraded':
        return {
          icon: AlertTriangle,
          label: 'Degraded',
          className: 'bg-warning/10 text-warning border-warning/20',
        }
      case 'timeout':
        return {
          icon: Clock,
          label: 'Timeout',
          className: 'bg-warning/10 text-warning border-warning/20',
        }
      default:
        return {
          icon: AlertCircle,
          label: type,
          className: 'bg-muted text-muted-foreground',
        }
    }
  }

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
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

  if (incidents.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>{title}</CardTitle>
          <CardDescription>{description}</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col items-center justify-center py-8 text-center">
            <CheckCircle2 className="size-12 text-uptime mb-3" />
            <h3 className="font-medium">No incidents</h3>
            <p className="text-sm text-muted-foreground">
              This service has been running smoothly
            </p>
          </div>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>{title}</CardTitle>
        <CardDescription>{description}</CardDescription>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Type</TableHead>
              <TableHead>Description</TableHead>
              <TableHead>Started</TableHead>
              <TableHead>Duration</TableHead>
              <TableHead className="text-right">Status</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {incidents.map((incident) => {
              const typeConfig = getTypeConfig(incident.type)
              const TypeIcon = typeConfig.icon

              return (
                <TableRow key={incident.id}>
                  <TableCell>
                    <Badge variant="outline" className={cn("gap-1", typeConfig.className)}>
                      <TypeIcon className="size-3" />
                      {typeConfig.label}
                    </Badge>
                  </TableCell>
                  <TableCell className="max-w-[300px] truncate">
                    {incident.description}
                  </TableCell>
                  <TableCell className="text-muted-foreground">
                    {formatDate(incident.startTime)}
                  </TableCell>
                  <TableCell className="font-mono">
                    {formatDuration(incident.duration)}
                  </TableCell>
                  <TableCell className="text-right">
                    {incident.resolved ? (
                      <Badge variant="outline" className="bg-uptime/10 text-uptime border-uptime/20">
                        Resolved
                      </Badge>
                    ) : (
                      <Badge variant="outline" className="bg-down/10 text-down border-down/20">
                        Ongoing
                      </Badge>
                    )}
                  </TableCell>
                </TableRow>
              )
            })}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  )
}
