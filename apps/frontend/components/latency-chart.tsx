'use client'

import { Area, AreaChart, CartesianGrid, XAxis, YAxis, ResponsiveContainer } from 'recharts'
import {
  ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from './ui/chart'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './ui/card'
import type { LatencyDataPoint } from '../lib/mock-data'

interface LatencyChartProps {
  data: LatencyDataPoint[]
  title?: string
  description?: string
}

const chartConfig = {
  latency: {
    label: "Latency",
    color: "var(--color-chart-1)",
  },
} satisfies ChartConfig

export function LatencyChart({ data, title = "Response Time", description = "Last 24 hours" }: LatencyChartProps) {
  const chartData = data.map((point) => ({
    time: new Date(point.timestamp).toLocaleTimeString('en-US', {
      hour: '2-digit',
      minute: '2-digit',
      hour12: false
    }),
    latency: point.latency,
  }))

  const avgLatency = data.length
    ? Math.round(data.reduce((sum, p) => sum + p.latency, 0) / data.length)
    : 0
  const maxLatency = data.length ? Math.max(...data.map(p => p.latency)) : 0
  const minLatency = data.length ? Math.min(...data.map(p => p.latency)) : 0

  return (
    <Card>
      <CardHeader className="pb-2">
        <div className="flex items-center justify-between">
          <div>
            <CardTitle>{title}</CardTitle>
            <CardDescription>{description}</CardDescription>
          </div>
          <div className="flex items-center gap-4 text-sm">
            <div className="text-center">
              <div className="text-muted-foreground">Avg</div>
              <div className="font-mono font-medium">{avgLatency}ms</div>
            </div>
            <div className="text-center">
              <div className="text-muted-foreground">Min</div>
              <div className="font-mono font-medium text-uptime">{minLatency}ms</div>
            </div>
            <div className="text-center">
              <div className="text-muted-foreground">Max</div>
              <div className="font-mono font-medium text-warning">{maxLatency}ms</div>
            </div>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <ChartContainer config={chartConfig} className="h-[300px] w-full">
          <AreaChart
            data={chartData}
            margin={{ top: 10, right: 10, left: 0, bottom: 0 }}
          >
            <defs>
              <linearGradient id="fillLatency" x1="0" y1="0" x2="0" y2="1">
                <stop
                  offset="5%"
                  stopColor="var(--color-chart-1)"
                  stopOpacity={0.8}
                />
                <stop
                  offset="95%"
                  stopColor="var(--color-chart-1)"
                  stopOpacity={0.1}
                />
              </linearGradient>
            </defs>
            <CartesianGrid strokeDasharray="3 3" vertical={false} className="stroke-muted" />
            <XAxis
              dataKey="time"
              tickLine={false}
              axisLine={false}
              tickMargin={8}
              minTickGap={32}
              className="text-xs"
            />
            <YAxis
              tickLine={false}
              axisLine={false}
              tickMargin={8}
              tickFormatter={(value) => `${value}ms`}
              className="text-xs"
            />
            <ChartTooltip
              cursor={false}
              content={
                <ChartTooltipContent
                  indicator="line"
                  formatter={(value) => (
                    <span className="font-mono">{value}ms</span>
                  )}
                />
              }
            />
            <Area
              dataKey="latency"
              type="monotone"
              fill="url(#fillLatency)"
              stroke="var(--color-chart-1)"
              strokeWidth={2}
            />
          </AreaChart>
        </ChartContainer>
      </CardContent>
    </Card>
  )
}
