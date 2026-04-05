'use client'

import { useState, useEffect } from 'react'
import { useParams, useRouter } from 'next/navigation'
import { ArrowLeft, Plus, Trash2, Edit2 } from 'lucide-react'
import { DashboardLayout } from '../../../components/dashboard-layout'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../../../components/ui/card'
import { Button } from '../../../components/ui/button'
import { Input } from '../../../components/ui/input'
import { Label } from '../../../components/ui/label'
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '../../../components/ui/select'
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from '../../../components/ui/dialog'
import { useToast } from '../../../hooks/use-toast'
import { getOrganization, listOrganizationMembers, addOrganizationMember, updateMemberRole, removeOrganizationMember } from '../../../lib/api'

export default function OrganizationDetailPage() {
    type MemberRole = 'viewer' | 'admin'

    const params = useParams()
    const router = useRouter()
    const { toast } = useToast()
    const orgId = params.id as string

    const [organization, setOrganization] = useState<any>(null)
    const [members, setMembers] = useState<any[]>([])
    const [loading, setLoading] = useState(true)
    const [showAddMemberDialog, setShowAddMemberDialog] = useState(false)
    const [addingMember, setAddingMember] = useState(false)
    const [newMemberEmail, setNewMemberEmail] = useState('')
    const [newMemberRole, setNewMemberRole] = useState<MemberRole>('viewer')
    const [updatingMember, setUpdatingMember] = useState<string | null>(null)
    const [selectedMemberRole, setSelectedMemberRole] = useState('')
    const [removingMember, setRemovingMember] = useState(false)

    useEffect(() => {
        loadOrganizationDetails()
    }, [orgId])

    const loadOrganizationDetails = async () => {
        try {
            const [org, membersList] = await Promise.all([
                getOrganization(orgId),
                listOrganizationMembers(orgId)
            ])
            setOrganization(org)
            setMembers(membersList)
        } catch (error) {
            toast({
                title: 'Error',
                description: 'Failed to load organization details',
                variant: 'destructive'
            })
        } finally {
            setLoading(false)
        }
    }

    const handleAddMember = async () => {
        if (!newMemberEmail.trim()) {
            toast({
                title: 'Error',
                description: 'Please enter an email address'
            })
            return
        }

        setAddingMember(true)
        try {
            await addOrganizationMember(orgId, newMemberEmail, newMemberRole)
            toast({
                title: 'Success',
                description: `Added ${newMemberEmail} to organization`
            })
            setNewMemberEmail('')
            setNewMemberRole('viewer')
            setShowAddMemberDialog(false)
            await loadOrganizationDetails()
        } catch (error) {
            toast({
                title: 'Error',
                description: 'Failed to add member',
                variant: 'destructive'
            })
        } finally {
            setAddingMember(false)
        }
    }

    const handleUpdateMemberRole = async (memberId: string, newRole: MemberRole) => {
        setUpdatingMember(memberId)
        try {
            await updateMemberRole(orgId, memberId, newRole)
            toast({
                title: 'Success',
                description: 'Member role updated'
            })
            setUpdatingMember(null)
            await loadOrganizationDetails()
        } catch (error) {
            toast({
                title: 'Error',
                description: 'Failed to update member role',
                variant: 'destructive'
            })
        }
    }

    const handleRemoveMember = async (memberId: string) => {
        // if (!confirm('Are you sure you want to remove this member?')) return

        setRemovingMember(true)
        try {
            await removeOrganizationMember(orgId, memberId)
            toast({
                title: 'Success',
                description: 'Member removed from organization'
            })
            await loadOrganizationDetails()
        } catch (error) {
            toast({
                title: 'Error',
                description: 'Failed to remove member',
                variant: 'destructive'
            })
        } finally {
            setRemovingMember(false)
        }
    }

    if (loading) {
        return (
            <DashboardLayout title="Organization">
                <div className="space-y-6">Loading...</div>
            </DashboardLayout>
        )
    }

    if (!organization) {
        return (
            <DashboardLayout title="Organization">
                <div className="space-y-6">
                    <Button variant="outline" onClick={() => router.back()}>
                        <ArrowLeft className="mr-2 size-4" />
                        Back to Organizations
                    </Button>
                    <div className="text-center text-muted-foreground">Organization not found</div>
                </div>
            </DashboardLayout>
        )
    }

    return (
        <DashboardLayout title={organization.name}>
            <div className="space-y-6">
                <div>
                    <Button variant="outline" onClick={() => router.back()}>
                        <ArrowLeft className="mr-2 size-4" />
                        Back to Organizations
                    </Button>
                </div>

                {/* Organization Header */}
                <Card>
                    <CardHeader>
                        <CardTitle className="text-3xl">{organization.name}</CardTitle>
                        <CardDescription>
                            Created {new Date(organization.createdAt).toLocaleDateString()}
                        </CardDescription>
                    </CardHeader>
                </Card>

                {/* Members Section */}
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between">
                        <div>
                            <CardTitle>Team Members</CardTitle>
                            <CardDescription>
                                {members.length} member{members.length !== 1 ? 's' : ''} in this organization
                            </CardDescription>
                        </div>
                        <Dialog open={showAddMemberDialog} onOpenChange={setShowAddMemberDialog}>
                            <DialogTrigger asChild>
                                <Button>
                                    <Plus className="mr-2 size-4" />
                                    Add Member
                                </Button>
                            </DialogTrigger>
                            <DialogContent>
                                <DialogHeader>
                                    <DialogTitle>Add Team Member</DialogTitle>
                                    <DialogDescription>
                                        Invite a new member to your organization by email
                                    </DialogDescription>
                                </DialogHeader>
                                <div className="space-y-4">
                                    <div className="grid gap-2">
                                        <Label htmlFor="email">Email Address</Label>
                                        <Input
                                            id="email"
                                            type="email"
                                            value={newMemberEmail}
                                            onChange={(e) => setNewMemberEmail(e.target.value)}
                                            placeholder="member@example.com"
                                        />
                                    </div>
                                    <div className="grid gap-2">
                                        <Label htmlFor="role">Role</Label>
                                        <Select
                                            value={newMemberRole}
                                            onValueChange={(value) => setNewMemberRole(value as MemberRole)}
                                        >
                                            <SelectTrigger>
                                                <SelectValue />
                                            </SelectTrigger>
                                            <SelectContent>
                                                <SelectItem value="viewer">Viewer (Read-only)</SelectItem>
                                                <SelectItem value="admin">Admin</SelectItem>
                                            </SelectContent>
                                        </Select>
                                    </div>
                                    <div className="flex justify-end gap-2 pt-4">
                                        <Button variant="outline" onClick={() => setShowAddMemberDialog(false)}>
                                            Cancel
                                        </Button>
                                        <Button onClick={handleAddMember} disabled={addingMember}>
                                            {addingMember ? 'Adding...' : 'Add Member'}
                                        </Button>
                                    </div>
                                </div>
                            </DialogContent>
                        </Dialog>
                    </CardHeader>
                    <CardContent>
                        {members.length === 0 ? (
                            <div className="text-center text-muted-foreground py-8">
                                No members yet. Add the first member to get started.
                            </div>
                        ) : (
                            <div className="space-y-3">
                                {members.map((member) => (
                                    <div
                                        key={member.id}
                                        className="flex items-center justify-between p-4 rounded-lg border hover:bg-secondary/50 transition-colors"
                                    >
                                        <div className="flex-1">
                                            <div className="font-medium">{member.user?.email}</div>
                                            <div className="text-sm text-muted-foreground">
                                                {member?.name ? member.name : 'No name available'}
                                            </div>
                                        </div>
                                        <div className="flex items-center gap-2">
                                            <Select
                                                value={member.role}
                                                onValueChange={(value) => handleUpdateMemberRole(member.id, value as MemberRole)}
                                                disabled={updatingMember === member.id}
                                            >
                                                <SelectTrigger className="w-40">
                                                    <SelectValue />
                                                </SelectTrigger>
                                                <SelectContent>
                                                    <SelectItem value="viewer">Viewer</SelectItem>
                                                    <SelectItem value="admin">Admin</SelectItem>
                                                </SelectContent>
                                            </Select>
                                            <Button
                                                variant="ghost"
                                                size="sm"
                                                onClick={() => handleRemoveMember(member.id)}
                                                disabled={removingMember}
                                                className="text-destructive hover:text-destructive"
                                            >
                                                <Trash2 className="size-4" />
                                            </Button>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        )}
                    </CardContent>
                </Card>
            </div>
        </DashboardLayout>
    )
}
