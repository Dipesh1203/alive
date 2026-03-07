'use client'

import { Globe2 } from 'lucide-react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './ui/card'
import { cn } from '../lib/utils'
import type { RegionLatency } from '../lib/mock-data'

interface RegionStatusProps {
  regions: RegionLatency[]
  title?: string
  description?: string
}

export function RegionStatus({ regions, title = "Region Status", description = "Latency by monitoring location" }: RegionStatusProps) {
  const getStatusColor = (status: RegionLatency['status']) => {
    switch (status) {
      case 'up':
        return 'bg-uptime'
      case 'down':
        return 'bg-down'
      case 'degraded':
        return 'bg-warning'
      default:
        return 'bg-muted'
    }
  }

  const getLatencyColor = (latency: number) => {
    if (latency === 0) return 'text-down'
    if (latency < 100) return 'text-uptime'
    if (latency < 300) return 'text-foreground'
    return 'text-warning'
  }

  const avgLatency = Math.round(
    regions.reduce((sum, r) => sum + r.latency, 0) / regions.length
  )

  return (
    <Card>
      <CardHeader className="pb-2">
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="flex items-center gap-2">
              <Globe2 className="size-5" />
              {title}
            </CardTitle>
            <CardDescription>{description}</CardDescription>
          </div>
          <div className="text-right">
            <div className="text-lg font-bold font-mono">{avgLatency}ms</div>
            <div className="text-xs text-muted-foreground">avg latency</div>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <div className="space-y-3">
          {regions.map((region) => (
            <div
              key={region.code}
              className="flex items-center justify-between p-3 rounded-lg bg-secondary/50"
            >
              <div className="flex items-center gap-3">
                <span className={cn("size-2.5 rounded-full", getStatusColor(region.status))} />
                <div>
                  <div className="font-medium text-sm">{region.code}</div>
                  <div className="text-xs text-muted-foreground">{region.region}</div>
                </div>
              </div>
              <div className={cn("font-mono font-medium", getLatencyColor(region.latency))}>
                {region.latency === 0 ? 'N/A' : `${region.latency}ms`}
              </div>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  )
}
