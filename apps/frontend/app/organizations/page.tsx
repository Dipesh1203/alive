'use client'

import { useState, useEffect } from 'react'
import { Plus, QrCode, Users, Settings } from 'lucide-react'
import Link from 'next/link'
import { DashboardLayout } from '../../components/dashboard-layout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../../components/ui/card'
import { Button } from '../../components/ui/button'
import { Input } from '../../components/ui/input'
import { Label } from '../../components/ui/label'
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from '../../components/ui/dialog'
import { Badge } from '../../components/ui/badge'
import { useToast } from '../../hooks/use-toast'
import { createOrganization, listUserOrganizations, type Organization } from '../../lib/api'

export default function OrganizationsPage() {
    const [organizations, setOrganizations] = useState<Organization[]>([])
    const [loading, setLoading] = useState(false)
    const [createDialogOpen, setCreateDialogOpen] = useState(false)
    const [newOrgName, setNewOrgName] = useState('')
    const { toast } = useToast()

    useEffect(() => {
        loadOrganizations()
    }, [])

    const loadOrganizations = async () => {
        try {
            setLoading(true)
            const orgs = await listUserOrganizations()
            setOrganizations(orgs)
        } catch (error) {
            toast({
                title: 'Error',
                description: 'Failed to load organizations',
                variant: 'destructive',
            })
        } finally {
            setLoading(false)
        }
    }

    const handleCreateOrganization = async () => {
        if (!newOrgName.trim()) {
            toast({
                title: 'Error',
                description: 'Organization name is required',
                variant: 'destructive',
            })
            return
        }

        try {
            await createOrganization(newOrgName)
            toast({
                title: 'Success',
                description: 'Organization created successfully',
            })
            setNewOrgName('')
            setCreateDialogOpen(false)
            loadOrganizations()
        } catch (error) {
            toast({
                title: 'Error',
                description: 'Failed to create organization',
                variant: 'destructive',
            })
        }
    }

    return (
        <DashboardLayout title="Organizations">
            <div className="space-y-6">
                <div className="flex items-center justify-between">
                    <div>
                        <h1 className="text-2xl font-bold tracking-tight">Organizations</h1>
                        <p className="text-muted-foreground">
                            Manage your organizations and team access
                        </p>
                    </div>
                    <Dialog open={createDialogOpen} onOpenChange={setCreateDialogOpen}>
                        <DialogTrigger asChild>
                            <Button>
                                <Plus className="mr-2 size-4" />
                                New Organization
                            </Button>
                        </DialogTrigger>
                        <DialogContent>
                            <DialogHeader>
                                <DialogTitle>Create Organization</DialogTitle>
                                <DialogDescription>
                                    Create a new organization to manage your monitoring resources.
                                </DialogDescription>
                            </DialogHeader>
                            <div className="grid gap-4 py-4">
                                <div className="grid gap-2">
                                    <Label htmlFor="org-name">Organization Name</Label>
                                    <Input
                                        id="org-name"
                                        placeholder="e.g., My Company"
                                        value={newOrgName}
                                        onChange={(e) => setNewOrgName(e.target.value)}
                                    />
                                </div>
                            </div>
                            <DialogFooter>
                                <Button variant="outline" onClick={() => setCreateDialogOpen(false)}>
                                    Cancel
                                </Button>
                                <Button onClick={handleCreateOrganization}>Create</Button>
                            </DialogFooter>
                        </DialogContent>
                    </Dialog>
                </div>

                {loading ? (
                    <div className="flex items-center justify-center py-12">
                        <p className="text-muted-foreground">Loading organizations...</p>
                    </div>
                ) : organizations.length === 0 ? (
                    <Card>
                        <CardContent className="flex flex-col items-center justify-center py-12">
                            <Users className="mb-4 size-12 text-muted-foreground" />
                            <h3 className="text-lg font-medium">No organizations yet</h3>
                            <p className="mb-4 text-sm text-muted-foreground">
                                Create your first organization to get started with monitoring.
                            </p>
                            <Dialog open={createDialogOpen} onOpenChange={setCreateDialogOpen}>
                                <DialogTrigger asChild>
                                    <Button>
                                        <Plus className="mr-2 size-4" />
                                        Create Organization
                                    </Button>
                                </DialogTrigger>
                                <DialogContent>
                                    <DialogHeader>
                                        <DialogTitle>Create Organization</DialogTitle>
                                        <DialogDescription>
                                            Create a new organization to manage your monitoring resources.
                                        </DialogDescription>
                                    </DialogHeader>
                                    <div className="grid gap-4 py-4">
                                        <div className="grid gap-2">
                                            <Label htmlFor="org-name">Organization Name</Label>
                                            <Input
                                                id="org-name"
                                                placeholder="e.g., My Company"
                                                value={newOrgName}
                                                onChange={(e) => setNewOrgName(e.target.value)}
                                            />
                                        </div>
                                    </div>
                                    <DialogFooter>
                                        <Button variant="outline" onClick={() => setCreateDialogOpen(false)}>
                                            Cancel
                                        </Button>
                                        <Button onClick={handleCreateOrganization}>Create</Button>
                                    </DialogFooter>
                                </DialogContent>
                            </Dialog>
                        </CardContent>
                    </Card>
                ) : (
                    <div className="grid gap-4 md:grid-cols-2">
                        {organizations.map((org) => (
                            <Card key={org.id} className="hover:border-primary/50 transition-colors">
                                <CardHeader>
                                    <div className="flex items-start justify-between">
                                        <div>
                                            <CardTitle>{org.name}</CardTitle>
                                            <CardDescription className="mt-1">
                                                Created {new Date(org.createdAt).toLocaleDateString()}
                                            </CardDescription>
                                        </div>
                                        <Badge variant="outline">Organization</Badge>
                                    </div>
                                </CardHeader>
                                <CardContent>
                                    <div className="flex gap-2">
                                        <Link href={`/organizations/${org.id}`} className="flex-1">
                                            <Button variant="outline" size="sm" className="w-full">
                                                <Users className="mr-2 size-4" />
                                                Members
                                            </Button>
                                        </Link>
                                        {/* <Button variant="outline" size="sm" className="flex-1">
                                            <Settings className="mr-2 size-4" />
                                            Settings
                                        </Button> */}
                                    </div>
                                </CardContent>
                            </Card>
                        ))}
                    </div>
                )}
            </div>
        </DashboardLayout>
    )
}
