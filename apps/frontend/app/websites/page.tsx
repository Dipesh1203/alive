import { DashboardLayout } from '../../components/dashboard-layout'
import { WebsiteCard } from '../../components/website-card'
import { AddWebsiteModal } from '../../components/add-website-modal'
import { getWebsites } from '../../lib/mock-data'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '../../components/ui/tabs'

export default function WebsitesPage() {
  const websites = getWebsites()

  const allWebsites = websites
  const activeWebsites = websites.filter(w => w.monitoringEnabled)
  const pausedWebsites = websites.filter(w => !w.monitoringEnabled)

  return (
    <DashboardLayout title="Websites">
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold tracking-tight">All Websites</h1>
            <p className="text-muted-foreground">
              Manage and monitor all your configured endpoints
            </p>
          </div>
          <AddWebsiteModal />
        </div>

        <Tabs defaultValue="all">
          <TabsList>
            <TabsTrigger value="all">
              All ({allWebsites.length})
            </TabsTrigger>
            <TabsTrigger value="active">
              Active ({activeWebsites.length})
            </TabsTrigger>
            <TabsTrigger value="paused">
              Paused ({pausedWebsites.length})
            </TabsTrigger>
          </TabsList>

          <TabsContent value="all" className="mt-4">
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
              {allWebsites.map((website) => (
                <WebsiteCard key={website.id} website={website} />
              ))}
            </div>
          </TabsContent>

          <TabsContent value="active" className="mt-4">
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
              {activeWebsites.map((website) => (
                <WebsiteCard key={website.id} website={website} />
              ))}
            </div>
          </TabsContent>

          <TabsContent value="paused" className="mt-4">
            {pausedWebsites.length > 0 ? (
              <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                {pausedWebsites.map((website) => (
                  <WebsiteCard key={website.id} website={website} />
                ))}
              </div>
            ) : (
              <div className="flex flex-col items-center justify-center py-12 text-center">
                <p className="text-muted-foreground">No paused monitors</p>
              </div>
            )}
          </TabsContent>
        </Tabs>
      </div>
    </DashboardLayout>
  )
}
