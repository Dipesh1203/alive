'use client'

import { useState, useEffect } from 'react'
import { Settings, User, Key, Globe, Clock, Trash2, Bell } from 'lucide-react'
import { DashboardLayout } from '../../components/dashboard-layout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../../components/ui/card'
import { Button } from '../../components/ui/button'
import { Input } from '../../components/ui/input'
import { Label } from '../../components/ui/label'
import { Switch } from '../../components/ui/switch'
import { Separator } from '../../components/ui/separator'
import { useToast } from '../../hooks/use-toast'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '../../components/ui/select'
import { getMyProfile, updateMyProfile, getMyPreferences, updateMyPreferences } from '../../lib/api'

type Preferences = {
  emailNotifications: boolean
  pushNotifications: boolean
  downtimeAlerts: boolean
  degradationAlerts: boolean
  recoveryNotifications: boolean
  dailyDigest: boolean
}

const defaultPreferences: Preferences = {
  emailNotifications: true,
  pushNotifications: true,
  downtimeAlerts: true,
  degradationAlerts: true,
  recoveryNotifications: true,
  dailyDigest: false,
}

export default function SettingsPage() {
  const { toast } = useToast()
  const [formData, setFormData] = useState({ firstName: '', lastName: '', phone: '', bio: '' })
  const [preferences, setPreferences] = useState<Preferences>(defaultPreferences)
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    const loadProfile = async () => {
      try {
        const profile = await getMyProfile()
        setFormData({ firstName: profile.firstName || '', lastName: profile.lastName || '', phone: profile.phone || '', bio: profile.bio || '' })
        const prefs = await getMyPreferences() as Record<string, unknown>
        setPreferences({
          emailNotifications: typeof prefs.emailNotifications === 'boolean' ? prefs.emailNotifications : defaultPreferences.emailNotifications,
          pushNotifications: typeof prefs.pushNotifications === 'boolean' ? prefs.pushNotifications : defaultPreferences.pushNotifications,
          downtimeAlerts: typeof prefs.downtimeAlerts === 'boolean' ? prefs.downtimeAlerts : defaultPreferences.downtimeAlerts,
          degradationAlerts: typeof prefs.degradationAlerts === 'boolean' ? prefs.degradationAlerts : defaultPreferences.degradationAlerts,
          recoveryNotifications: typeof prefs.recoveryNotifications === 'boolean' ? prefs.recoveryNotifications : defaultPreferences.recoveryNotifications,
          dailyDigest: typeof prefs.dailyDigest === 'boolean' ? prefs.dailyDigest : defaultPreferences.dailyDigest,
        })
      } catch (error) {
        toast({ title: 'Error', description: 'Failed to load profile settings', variant: 'destructive' })
      } finally {
        setLoading(false)
      }
    }
    loadProfile()
  }, [])

  const handleSaveProfile = async () => {
    setSaving(true)
    try {
      await updateMyProfile(formData.firstName, formData.lastName, formData.phone, formData.bio)
      toast({ title: 'Success', description: 'Profile updated successfully' })
    } catch (error) {
      toast({ title: 'Error', description: 'Failed to save profile', variant: 'destructive' })
    } finally {
      setSaving(false)
    }
  }

  const handleSavePreferences = async () => {
    setSaving(true)
    try {
      await updateMyPreferences(preferences)
      toast({ title: 'Success', description: 'Notification preferences updated' })
    } catch (error) {
      toast({ title: 'Error', description: 'Failed to save preferences', variant: 'destructive' })
    } finally {
      setSaving(false)
    }
  }

  if (loading) return (<DashboardLayout title="Settings"><div className="space-y-6 max-w-3xl"><div>Loading settings...</div></div></DashboardLayout>)

  return (
    <DashboardLayout title="Settings">
      <div className="space-y-6 max-w-3xl">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Settings</h1>
          <p className="text-muted-foreground">
            Manage your account and workspace preferences
          </p>
        </div>

        {/* Profile settings */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <User className="size-5" />
              Profile
            </CardTitle>
            <CardDescription>
              Your personal account information
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid gap-2">
              <Label htmlFor="firstName">First Name</Label>
              <Input id="firstName" value={formData.firstName} onChange={(e) => setFormData({ ...formData, firstName: e.target.value })} placeholder="Enter your first name" />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="lastName">Last Name</Label>
              <Input id="lastName" value={formData.lastName} onChange={(e) => setFormData({ ...formData, lastName: e.target.value })} placeholder="Enter your last name" />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="phone">Phone Number</Label>
              <Input id="phone" type="tel" value={formData.phone} onChange={(e) => setFormData({ ...formData, phone: e.target.value })} placeholder="Enter your phone number" />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="bio">Bio</Label>
              <Input id="bio" value={formData.bio} onChange={(e) => setFormData({ ...formData, bio: e.target.value })} placeholder="Tell us about yourself" />
            </div>
            <Button onClick={handleSaveProfile} disabled={saving}>{saving ? 'Saving...' : 'Save Profile'}</Button>
          </CardContent>
        </Card>

        {/* Notification Preferences */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Bell className="size-5" />
              Notification Preferences
            </CardTitle>
            <CardDescription>
              Manage how you receive alerts and updates
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center justify-between p-4 rounded-lg bg-secondary/50">
              <div>
                <div className="font-medium">Email Notifications</div>
                <div className="text-sm text-muted-foreground">Receive alerts via email</div>
              </div>
              <Switch checked={preferences.emailNotifications} onCheckedChange={(checked) => setPreferences({ ...preferences, emailNotifications: checked })} />
            </div>
            <div className="flex items-center justify-between p-4 rounded-lg bg-secondary/50">
              <div>
                <div className="font-medium">Push Notifications</div>
                <div className="text-sm text-muted-foreground">Receive browser push notifications</div>
              </div>
              <Switch checked={preferences.pushNotifications} onCheckedChange={(checked) => setPreferences({ ...preferences, pushNotifications: checked })} />
            </div>
            <div className="flex items-center justify-between p-4 rounded-lg bg-secondary/50">
              <div>
                <div className="font-medium">Downtime Alerts</div>
                <div className="text-sm text-muted-foreground">Alert when a website goes down</div>
              </div>
              <Switch checked={preferences.downtimeAlerts} onCheckedChange={(checked) => setPreferences({ ...preferences, downtimeAlerts: checked })} />
            </div>
            <div className="flex items-center justify-between p-4 rounded-lg bg-secondary/50">
              <div>
                <div className="font-medium">Degradation Alerts</div>
                <div className="text-sm text-muted-foreground">Alert when response times are slow</div>
              </div>
              <Switch checked={preferences.degradationAlerts} onCheckedChange={(checked) => setPreferences({ ...preferences, degradationAlerts: checked })} />
            </div>
            <div className="flex items-center justify-between p-4 rounded-lg bg-secondary/50">
              <div>
                <div className="font-medium">Recovery Notifications</div>
                <div className="text-sm text-muted-foreground">Alert when a website comes back online</div>
              </div>
              <Switch checked={preferences.recoveryNotifications} onCheckedChange={(checked) => setPreferences({ ...preferences, recoveryNotifications: checked })} />
            </div>
            <div className="flex items-center justify-between p-4 rounded-lg bg-secondary/50">
              <div>
                <div className="font-medium">Daily Digest</div>
                <div className="text-sm text-muted-foreground">Receive daily summary of monitoring activity</div>
              </div>
              <Switch checked={preferences.dailyDigest} onCheckedChange={(checked) => setPreferences({ ...preferences, dailyDigest: checked })} />
            </div>
            <Button onClick={handleSavePreferences} disabled={saving}>{saving ? 'Saving...' : 'Save Preferences'}</Button>
          </CardContent>
        </Card>

        {/* Workspace settings */}
        {/* <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Globe className="size-5" />
              Workspace
            </CardTitle>
            <CardDescription>
              Configure your monitoring workspace
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid gap-2">
              <Label htmlFor="workspace-name">Workspace Name</Label>
              <Input id="workspace-name" defaultValue="My Company" />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="timezone">Timezone</Label>
              <Select defaultValue="utc">
                <SelectTrigger>
                  <SelectValue placeholder="Select timezone" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="utc">UTC</SelectItem>
                  <SelectItem value="est">Eastern Time (EST)</SelectItem>
                  <SelectItem value="pst">Pacific Time (PST)</SelectItem>
                  <SelectItem value="cet">Central European Time (CET)</SelectItem>
                  <SelectItem value="jst">Japan Standard Time (JST)</SelectItem>
                </SelectContent>
              </Select>
              <p className="text-xs text-muted-foreground">
                Used for reports and incident timestamps
              </p>
            </div>
            <Button>Save Changes</Button>
          </CardContent>
        </Card> */}

        {/* Monitoring defaults */}
        {/* <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Clock className="size-5" />
              Monitoring Defaults
            </CardTitle>
            <CardDescription>
              Default settings for new monitors
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid gap-2">
              <Label htmlFor="check-interval">Check Interval</Label>
              <Select defaultValue="60">
                <SelectTrigger>
                  <SelectValue placeholder="Select interval" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="30">Every 30 seconds</SelectItem>
                  <SelectItem value="60">Every 1 minute</SelectItem>
                  <SelectItem value="300">Every 5 minutes</SelectItem>
                  <SelectItem value="600">Every 10 minutes</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="grid gap-2">
              <Label htmlFor="timeout">Request Timeout</Label>
              <Select defaultValue="30">
                <SelectTrigger>
                  <SelectValue placeholder="Select timeout" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="10">10 seconds</SelectItem>
                  <SelectItem value="30">30 seconds</SelectItem>
                  <SelectItem value="60">60 seconds</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="flex items-center justify-between p-4 rounded-lg bg-secondary/50">
              <div>
                <div className="font-medium">SSL Verification</div>
                <div className="text-sm text-muted-foreground">
                  Verify SSL certificates on HTTPS endpoints
                </div>
              </div>
              <Switch defaultChecked />
            </div>
            <div className="flex items-center justify-between p-4 rounded-lg bg-secondary/50">
              <div>
                <div className="font-medium">Follow Redirects</div>
                <div className="text-sm text-muted-foreground">
                  Automatically follow HTTP redirects
                </div>
              </div>
              <Switch defaultChecked />
            </div>
            <Button>Save Changes</Button>
          </CardContent>
        </Card> */}

        {/* API settings */}
        {/* <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Key className="size-5" />
              API Access
            </CardTitle>
            <CardDescription>
              Manage your API keys for programmatic access
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center justify-between p-4 rounded-lg border">
              <div>
                <div className="font-medium font-mono text-sm">live_sk_...a3f8</div>
                <div className="text-xs text-muted-foreground">
                  Created Jan 15, 2024 - Last used 2 hours ago
                </div>
              </div>
              <div className="flex items-center gap-2">
                <Button variant="outline" size="sm">Reveal</Button>
                <Button variant="ghost" size="sm" className="text-destructive">
                  Revoke
                </Button>
              </div>
            </div>
            <Button variant="outline">
              <Key className="mr-2 size-4" />
              Generate New API Key
            </Button>
          </CardContent>
        </Card> */}

        {/* Danger zone */}
        {/* <Card className="border-destructive/50">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-destructive">
              <Trash2 className="size-5" />
              Danger Zone
            </CardTitle>
            <CardDescription>
              Irreversible actions that affect your workspace
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center justify-between p-4 rounded-lg border border-destructive/30 bg-destructive/5">
              <div>
                <div className="font-medium">Delete Workspace</div>
                <div className="text-sm text-muted-foreground">
                  Permanently delete this workspace and all its data
                </div>
              </div>
              <Button variant="destructive" size="sm">
                Delete Workspace
              </Button>
            </div>
          </CardContent>
        </Card> */}
      </div>
    </DashboardLayout>
  )
}
