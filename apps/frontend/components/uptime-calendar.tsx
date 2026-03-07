'use client'

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './ui/card'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from './ui/tooltip'
import { cn } from '../lib/utils'
import type { UptimeDay } from '../lib/mock-data'

interface UptimeCalendarProps {
  data: UptimeDay[]
  title?: string
  description?: string
}

export function UptimeCalendar({ data, title = "Uptime History", description = "Last 90 days" }: UptimeCalendarProps) {
  // Group data into weeks (7 days per row)
  const weeks: UptimeDay[][] = []
  for (let i = 0; i < data.length; i += 7) {
    weeks.push(data.slice(i, i + 7))
  }

  const getStatusColor = (status: UptimeDay['status']) => {
    switch (status) {
      case 'up':
        return 'bg-uptime hover:bg-uptime/80'
      case 'down':
        return 'bg-down hover:bg-down/80'
      case 'degraded':
        return 'bg-warning hover:bg-warning/80'
      case 'partial':
        return 'bg-uptime/60 hover:bg-uptime/50'
      default:
        return 'bg-muted hover:bg-muted/80'
    }
  }

  const getStatusLabel = (status: UptimeDay['status']) => {
    switch (status) {
      case 'up':
        return 'Operational'
      case 'down':
        return 'Outage'
      case 'degraded':
        return 'Degraded'
      case 'partial':
        return 'Partial outage'
      default:
        return 'Unknown'
    }
  }

  // Calculate overall uptime for the period
  const overallUptime = data.length
    ? data.reduce((sum, day) => sum + day.uptime, 0) / data.length
    : 100

  return (
    <Card>
      <CardHeader className="pb-2">
        <div className="flex items-center justify-between">
          <div>
            <CardTitle>{title}</CardTitle>
            <CardDescription>{description}</CardDescription>
          </div>
          <div className="text-right">
            <div className="text-2xl font-bold font-mono">{overallUptime.toFixed(2)}%</div>
            <div className="text-xs text-muted-foreground">90-day average</div>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <TooltipProvider>
          <div className="flex flex-col gap-1">
            {/* Month labels */}
            <div className="flex gap-1 mb-1 ml-8">
              {getMonthLabels(data).map((label, i) => (
                <div
                  key={i}
                  className="text-xs text-muted-foreground"
                  style={{ width: `${label.width}px`, minWidth: `${label.width}px` }}
                >
                  {label.month}
                </div>
              ))}
            </div>

            {/* Day labels and grid */}
            <div className="flex gap-1">
              <div className="flex flex-col gap-1 text-xs text-muted-foreground w-6">
                <span className="h-3"></span>
                <span className="h-3 leading-3">Mon</span>
                <span className="h-3"></span>
                <span className="h-3 leading-3">Wed</span>
                <span className="h-3"></span>
                <span className="h-3 leading-3">Fri</span>
                <span className="h-3"></span>
              </div>

              <div className="flex gap-1 overflow-x-auto">
                {weeks.map((week, weekIndex) => (
                  <div key={weekIndex} className="flex flex-col gap-1">
                    {week.map((day, dayIndex) => (
                      <Tooltip key={dayIndex}>
                        <TooltipTrigger asChild>
                          <div
                            className={cn(
                              "size-3 rounded-sm cursor-pointer transition-colors",
                              getStatusColor(day.status)
                            )}
                          />
                        </TooltipTrigger>
                        <TooltipContent side="top" className="text-xs">
                          <div className="font-medium">
                            {new Date(day.date).toLocaleDateString('en-US', {
                              weekday: 'short',
                              month: 'short',
                              day: 'numeric',
                            })}
                          </div>
                          <div className="flex items-center gap-2 mt-1">
                            <span className={cn(
                              "size-2 rounded-full",
                              getStatusColor(day.status)
                            )} />
                            <span>{getStatusLabel(day.status)}</span>
                            <span className="font-mono">{day.uptime.toFixed(2)}%</span>
                          </div>
                        </TooltipContent>
                      </Tooltip>
                    ))}
                  </div>
                ))}
              </div>
            </div>
          </div>
        </TooltipProvider>

        {/* Legend */}
        <div className="flex items-center justify-end gap-4 mt-4 text-xs text-muted-foreground">
          <span>Less</span>
          <div className="flex gap-1">
            <div className="size-3 rounded-sm bg-down" />
            <div className="size-3 rounded-sm bg-warning" />
            <div className="size-3 rounded-sm bg-uptime/60" />
            <div className="size-3 rounded-sm bg-uptime" />
          </div>
          <span>More</span>
        </div>
      </CardContent>
    </Card>
  )
}

function getMonthLabels(data: UptimeDay[]): { month: string; width: number }[] {
  const labels: { month: string; width: number }[] = []
  let currentMonth = ''
  let currentWidth = 0

  data.forEach((day) => {
    const month = new Date(day.date).toLocaleDateString('en-US', { month: 'short' })
    if (month !== currentMonth) {
      if (currentMonth) {
        labels.push({ month: currentMonth, width: currentWidth })
      }
      currentMonth = month
      currentWidth = 16 // size of one cell (12px) + gap (4px)
    } else {
      currentWidth += 16
    }
  })

  if (currentMonth) {
    labels.push({ month: currentMonth, width: currentWidth })
  }

  return labels
}
