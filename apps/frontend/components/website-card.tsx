'use client'

import Link from 'next/link'
import { Globe, MapPin, Clock, ExternalLink } from 'lucide-react'
import { Card, CardContent, CardHeader } from './ui/card'
import { Badge } from './ui/badge'
import { Button } from './ui/button'
import { cn } from '../lib/utils'
import type { Website } from '../lib/mock-data'

interface WebsiteCardProps {
  website: Website
}

export function WebsiteCard({ website }: WebsiteCardProps) {
  const statusConfig = {
    up: {
      label: 'UP',
      className: 'bg-uptime text-uptime-foreground',
      beaconClass: 'bg-uptime',
    },
    down: {
      label: 'DOWN',
      className: 'bg-down text-down-foreground',
      beaconClass: 'bg-down',
    },
    degraded: {
      label: 'DEGRADED',
      className: 'bg-warning text-warning-foreground',
      beaconClass: 'bg-warning',
    },
  }

  const status = statusConfig[website.status]

  // Generate a simple performance bar based on current latency
  const performanceSegments = Array.from({ length: 12 }, (_, i) => {
    const baseLatency = website.avgLatency
    const variation = (Math.sin(i * 0.5) * 0.3 + Math.random() * 0.2) * baseLatency
    const latency = baseLatency + variation
    const normalizedLatency = Math.min(latency / 500, 1)
    return normalizedLatency
  })

  return (
    <Card className="group relative overflow-hidden transition-all hover:shadow-lg hover:border-primary/50">
      <Link href={`/websites/${website.id}`} className="absolute inset-0 z-10">
        <span className="sr-only">View {website.name} details</span>
      </Link>

      <CardHeader className="pb-3">
        <div className="flex items-start justify-between">
          <div className="flex items-center gap-3">
            <div className="flex size-10 items-center justify-center rounded-lg bg-secondary">
              <Globe className="size-5 text-muted-foreground" />
            </div>
            <div>
              <h3 className="font-semibold leading-none">{website.name}</h3>
              <p className="mt-1 text-sm text-muted-foreground truncate max-w-[200px]">
                {website.url}
              </p>
            </div>
          </div>

          <div className="flex items-center gap-2">
            {/* Pulsing beacon */}
            <span className="relative flex size-3">
              <span className={cn(
                "absolute inline-flex h-full w-full animate-ping rounded-full opacity-75",
                status.beaconClass
              )} />
              <span className={cn(
                "relative inline-flex size-3 rounded-full",
                status.beaconClass
              )} />
            </span>
            <span className={cn(
              "text-sm font-bold",
              website.status === 'up' && 'text-uptime',
              website.status === 'down' && 'text-down',
              website.status === 'degraded' && 'text-warning'
            )}>
              {status.label}
            </span>
          </div>
        </div>
      </CardHeader>

      <CardContent className="space-y-4">
        {/* Metrics row */}
        <div className="flex items-center gap-4">
          <Badge variant="secondary" className="font-mono text-xs">
            {website.uptime30d.toFixed(2)}% uptime
          </Badge>
          <div className="flex items-center gap-1 text-xs text-muted-foreground">
            <Clock className="size-3" />
            <span>{website.currentLatency}ms</span>
          </div>
          <div className="flex items-center gap-1 text-xs text-muted-foreground">
            <MapPin className="size-3" />
            <span>{website.regions.length} regions</span>
          </div>
        </div>

        {/* Performance bar */}
        <div className="space-y-1.5">
          <div className="flex items-center justify-between text-xs">
            <span className="text-muted-foreground">Response Time (24h)</span>
            <span className="text-muted-foreground">avg {website.avgLatency}ms</span>
          </div>
          <div className="flex h-6 gap-0.5 rounded overflow-hidden">
            {performanceSegments.map((segment, i) => (
              <div
                key={i}
                className={cn(
                  "flex-1 rounded-sm transition-all",
                  segment < 0.3 && "bg-uptime/80",
                  segment >= 0.3 && segment < 0.6 && "bg-uptime/60",
                  segment >= 0.6 && segment < 0.8 && "bg-warning/70",
                  segment >= 0.8 && "bg-down/70"
                )}
                style={{ opacity: 0.5 + segment * 0.5 }}
              />
            ))}
          </div>
        </div>

        {/* Footer */}
        <div className="flex items-center justify-between pt-2 border-t border-border">
          <span className="text-xs text-muted-foreground">
            Last checked {getTimeAgo(website.lastChecked)}
          </span>
          <Button
            variant="ghost"
            size="sm"
            className="relative z-20 h-7 px-2 opacity-0 group-hover:opacity-100 transition-opacity"
            asChild
          >
            <a href={website.url} target="_blank" rel="noopener noreferrer" onClick={(e) => e.stopPropagation()}>
              <ExternalLink className="size-3 mr-1" />
              Visit
            </a>
          </Button>
        </div>
      </CardContent>
    </Card>
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
