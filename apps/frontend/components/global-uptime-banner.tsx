'use client'

import { Activity, CheckCircle2, AlertTriangle, XCircle } from 'lucide-react'
import { Card } from './ui/card'
import { cn } from '../lib/utils'

interface GlobalUptimeBannerProps {
  uptime: number
  totalWebsites: number
  websitesUp: number
  websitesDown: number
  websitesDegraded: number
}

export function GlobalUptimeBanner({
  uptime,
  totalWebsites,
  websitesUp,
  websitesDown,
  websitesDegraded,
}: GlobalUptimeBannerProps) {
  const getStatusColor = () => {
    if (websitesDown > 0) return 'from-down/20 to-down/5 border-down/30'
    if (websitesDegraded > 0) return 'from-warning/20 to-warning/5 border-warning/30'
    return 'from-uptime/20 to-uptime/5 border-uptime/30'
  }

  const getStatusIcon = () => {
    if (websitesDown > 0) return <XCircle className="size-6 text-down" />
    if (websitesDegraded > 0) return <AlertTriangle className="size-6 text-warning" />
    return <CheckCircle2 className="size-6 text-uptime" />
  }

  const getStatusMessage = () => {
    if (websitesDown > 0) return `${websitesDown} service${websitesDown > 1 ? 's' : ''} down`
    if (websitesDegraded > 0) return `${websitesDegraded} service${websitesDegraded > 1 ? 's' : ''} degraded`
    return 'All systems operational'
  }

  return (
    <Card className={cn(
      "relative overflow-hidden bg-gradient-to-r border",
      getStatusColor()
    )}>
      <div className="absolute inset-0 bg-grid-white/5" />
      <div className="relative flex flex-col md:flex-row items-center justify-between gap-4 p-6">
        <div className="flex items-center gap-4">
          <div className="flex size-14 items-center justify-center rounded-xl bg-background/80 backdrop-blur">
            <Activity className="size-7 text-primary" />
          </div>
          <div>
            <div className="flex items-baseline gap-2">
              <span className="text-4xl font-bold tracking-tight">{uptime.toFixed(3)}%</span>
              <span className="text-muted-foreground">uptime</span>
            </div>
            <p className="text-sm text-muted-foreground">
              Global availability across {totalWebsites} monitored services
            </p>
          </div>
        </div>

        <div className="flex items-center gap-6">
          <div className="flex items-center gap-6 px-6 py-3 rounded-lg bg-background/50 backdrop-blur">
            <div className="flex items-center gap-2">
              <div className="size-3 rounded-full bg-uptime" />
              <span className="text-sm font-medium">{websitesUp} Up</span>
            </div>
            {websitesDegraded > 0 && (
              <div className="flex items-center gap-2">
                <div className="size-3 rounded-full bg-warning" />
                <span className="text-sm font-medium">{websitesDegraded} Degraded</span>
              </div>
            )}
            {websitesDown > 0 && (
              <div className="flex items-center gap-2">
                <div className="size-3 rounded-full bg-down" />
                <span className="text-sm font-medium">{websitesDown} Down</span>
              </div>
            )}
          </div>

          <div className="flex items-center gap-2 text-sm">
            {getStatusIcon()}
            <span className="font-medium">{getStatusMessage()}</span>
          </div>
        </div>
      </div>
    </Card>
  )
}
