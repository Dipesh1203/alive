'use client'

import { useState } from 'react'
import { Bell, Slack, Mail, Webhook, AlertTriangle, Plus, MoreVertical, Trash2, Pencil } from 'lucide-react'
import { DashboardLayout } from '../../components/dashboard-layout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../../components/ui/card'
import { Button } from '../../components/ui/button'
import { Switch } from '../../components/ui/switch'
import { Badge } from '../../components/ui/badge'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '../../components/ui/dropdown-menu'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '../../components/ui/dialog'
import { Input } from '../../components/ui/input'
import { Label } from '../../components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '../../components/ui/select'
import { cn } from '../../lib/utils'
import { getNotificationChannels, type NotificationChannel } from '../../lib/mock-data'

const channelIcons = {
  slack: Slack,
  email: Mail,
  webhook: Webhook,
  pagerduty: AlertTriangle,
}

const channelLabels = {
  slack: 'Slack',
  email: 'Email',
  webhook: 'Webhook',
  pagerduty: 'PagerDuty',
}

export default function NotificationsPage() {
  const [channels, setChannels] = useState<NotificationChannel[]>(getNotificationChannels())
  const [addDialogOpen, setAddDialogOpen] = useState(false)

  const toggleChannel = (id: string) => {
    setChannels(prev =>
      prev.map(ch =>
        ch.id === id ? { ...ch, enabled: !ch.enabled } : ch
      )
    )
  }

  return (
    <DashboardLayout title="Notifications">
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold tracking-tight">Notification Settings</h1>
            <p className="text-muted-foreground">
              Configure how and where you receive alerts when issues occur
            </p>
          </div>
          <Dialog open={addDialogOpen} onOpenChange={setAddDialogOpen}>
            <DialogTrigger asChild>
              <Button>
                <Plus className="mr-2 size-4" />
                Add Channel
              </Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Add Notification Channel</DialogTitle>
                <DialogDescription>
                  Set up a new channel to receive alerts about your monitored services.
                </DialogDescription>
              </DialogHeader>
              <div className="grid gap-4 py-4">
                <div className="grid gap-2">
                  <Label htmlFor="channel-type">Channel Type</Label>
                  <Select defaultValue="slack">
                    <SelectTrigger>
                      <SelectValue placeholder="Select channel type" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="slack">Slack</SelectItem>
                      <SelectItem value="email">Email</SelectItem>
                      <SelectItem value="webhook">Webhook</SelectItem>
                      <SelectItem value="pagerduty">PagerDuty</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="channel-name">Display Name</Label>
                  <Input id="channel-name" placeholder="e.g., Engineering Team" />
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="channel-config">Webhook URL / Configuration</Label>
                  <Input id="channel-config" placeholder="https://..." />
                </div>
              </div>
              <DialogFooter>
                <Button variant="outline" onClick={() => setAddDialogOpen(false)}>
                  Cancel
                </Button>
                <Button onClick={() => setAddDialogOpen(false)}>
                  Add Channel
                </Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>
        </div>

        {/* Alert types configuration */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Bell className="size-5" />
              Alert Types
            </CardTitle>
            <CardDescription>
              Choose which events trigger notifications
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center justify-between p-4 rounded-lg bg-secondary/50">
              <div>
                <div className="font-medium">Downtime Alerts</div>
                <div className="text-sm text-muted-foreground">
                  Notify immediately when a service goes down
                </div>
              </div>
              <Switch defaultChecked />
            </div>
            <div className="flex items-center justify-between p-4 rounded-lg bg-secondary/50">
              <div>
                <div className="font-medium">Performance Degradation</div>
                <div className="text-sm text-muted-foreground">
                  Alert when response times exceed thresholds
                </div>
              </div>
              <Switch defaultChecked />
            </div>
            <div className="flex items-center justify-between p-4 rounded-lg bg-secondary/50">
              <div>
                <div className="font-medium">Recovery Notifications</div>
                <div className="text-sm text-muted-foreground">
                  Notify when services recover from issues
                </div>
              </div>
              <Switch defaultChecked />
            </div>
            <div className="flex items-center justify-between p-4 rounded-lg bg-secondary/50">
              <div>
                <div className="font-medium">SSL Certificate Expiry</div>
                <div className="text-sm text-muted-foreground">
                  Warn before SSL certificates expire
                </div>
              </div>
              <Switch />
            </div>
            <div className="flex items-center justify-between p-4 rounded-lg bg-secondary/50">
              <div>
                <div className="font-medium">Daily Summary</div>
                <div className="text-sm text-muted-foreground">
                  Receive a daily digest of all monitoring activity
                </div>
              </div>
              <Switch />
            </div>
          </CardContent>
        </Card>

        {/* Notification channels */}
        <Card>
          <CardHeader>
            <CardTitle>Notification Channels</CardTitle>
            <CardDescription>
              Manage your configured notification destinations
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {channels.map((channel) => {
                const Icon = channelIcons[channel.type]

                return (
                  <div
                    key={channel.id}
                    className={cn(
                      "flex items-center justify-between p-4 rounded-lg border transition-colors",
                      channel.enabled ? "bg-card" : "bg-muted/50 opacity-60"
                    )}
                  >
                    <div className="flex items-center gap-4">
                      <div className={cn(
                        "flex size-10 items-center justify-center rounded-lg",
                        channel.type === 'slack' && "bg-[#4A154B] text-white",
                        channel.type === 'email' && "bg-blue-600 text-white",
                        channel.type === 'webhook' && "bg-orange-600 text-white",
                        channel.type === 'pagerduty' && "bg-green-600 text-white"
                      )}>
                        <Icon className="size-5" />
                      </div>
                      <div>
                        <div className="flex items-center gap-2">
                          <span className="font-medium">{channel.name}</span>
                          <Badge variant="outline" className="text-xs">
                            {channelLabels[channel.type]}
                          </Badge>
                        </div>
                        <div className="text-sm text-muted-foreground">
                          {channel.type === 'email'
                            ? channel.config.emails
                            : channel.config.webhook || channel.config.url || channel.config.serviceKey}
                        </div>
                      </div>
                    </div>
                    <div className="flex items-center gap-3">
                      <Switch
                        checked={channel.enabled}
                        onCheckedChange={() => toggleChannel(channel.id)}
                      />
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button variant="ghost" size="icon" className="size-8">
                            <MoreVertical className="size-4" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          <DropdownMenuItem>
                            <Pencil className="mr-2 size-4" />
                            Edit
                          </DropdownMenuItem>
                          <DropdownMenuItem className="text-destructive">
                            <Trash2 className="mr-2 size-4" />
                            Delete
                          </DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </div>
                  </div>
                )
              })}
            </div>
          </CardContent>
        </Card>
      </div>
    </DashboardLayout>
  )
}
