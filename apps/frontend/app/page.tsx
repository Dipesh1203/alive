"use client"

import { useEffect, useMemo, useState } from 'react'
import { AddWebsiteModal } from '../components/add-website-modal'
import { DashboardLayout } from '../components/dashboard-layout'
import { GlobalUptimeBanner } from '../components/global-uptime-banner'
import { WebsiteCard } from '../components/website-card'
import { fetchWebsites } from '../lib/api'
import type { Website } from '../lib/mock-data'


export default function DashboardPage() {
  const [websites, setWebsites] = useState<Website[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const loadWebsites = async () => {
    setIsLoading(true)
    setError(null)

    try {
      const response = await fetchWebsites()
      setWebsites(response)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load websites')
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    void loadWebsites()
  }, [])

  const globalUptime = useMemo(() => {
    const enabledWebsites = websites.filter((website) => website.monitoringEnabled)
    if (enabledWebsites.length === 0) return 100

    const sum = enabledWebsites.reduce((acc, website) => acc + website.uptime30d, 0)
    return Math.round((sum / enabledWebsites.length) * 1000) / 1000
  }, [websites])

  const websitesUp = websites.filter(w => w.status === 'up' && w.monitoringEnabled).length
  const websitesDown = websites.filter(w => w.status === 'down' && w.monitoringEnabled).length
  const websitesDegraded = websites.filter(w => w.status === 'degraded' && w.monitoringEnabled).length
  const totalMonitored = websites.filter(w => w.monitoringEnabled).length

  return (
    <DashboardLayout title="Dashboard">
      <div className="space-y-6">
        {/* Global uptime banner */}
        <GlobalUptimeBanner
          uptime={globalUptime}
          totalWebsites={totalMonitored}
          websitesUp={websitesUp}
          websitesDown={websitesDown}
          websitesDegraded={websitesDegraded}
        />

        {/* Header with actions */}
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold tracking-tight">Monitored Websites</h1>
            <p className="text-muted-foreground">
              Track the health and performance of your services
            </p>
          </div>
          <AddWebsiteModal onCreated={loadWebsites} />
        </div>

        {isLoading && (
          <p className="text-sm text-muted-foreground">Loading websites...</p>
        )}

        {error && (
          <p className="text-sm text-down">{error}</p>
        )}

        {/* Website grid */}
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {websites.map((website) => (
            <WebsiteCard key={website.id} website={website} />
          ))}
        </div>
      </div>
    </DashboardLayout>
  )
}
