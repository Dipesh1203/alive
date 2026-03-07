import { AddWebsiteModal } from "../components/add-website-modal"
import { DashboardLayout } from "../components/dashboard-layout"
import { GlobalUptimeBanner } from "../components/global-uptime-banner"
import { WebsiteCard } from "../components/website-card"
import { calculateGlobalUptime, getWebsites } from "../lib/mock-data"


export default function DashboardPage() {
  const websites = getWebsites()
  const globalUptime = calculateGlobalUptime()

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
          <AddWebsiteModal />
        </div>

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
