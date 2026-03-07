'use client'

import { useState } from 'react'
import { Plus, Globe, Loader2 } from 'lucide-react'
import { Button } from './ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from './ui/dialog'
import { Input } from './ui/input'
import { Label } from './ui/label'
import { Checkbox } from './ui/checkbox'
import { mockRegions } from '../lib/mock-data'

interface AddWebsiteModalProps {
  onAdd?: (data: { name: string; url: string; regions: string[] }) => void
}

export function AddWebsiteModal({ onAdd }: AddWebsiteModalProps) {
  const [open, setOpen] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [name, setName] = useState('')
  const [url, setUrl] = useState('')
  const [selectedRegions, setSelectedRegions] = useState<string[]>(['us-east-1'])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsLoading(true)

    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1000))

    onAdd?.({ name, url, regions: selectedRegions })

    setIsLoading(false)
    setOpen(false)
    setName('')
    setUrl('')
    setSelectedRegions(['us-east-1'])
  }

  const toggleRegion = (regionId: string) => {
    setSelectedRegions(prev =>
      prev.includes(regionId)
        ? prev.filter(r => r !== regionId)
        : [...prev, regionId]
    )
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button>
          <Plus className="mr-2 size-4" />
          Add Website
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[500px]">
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <Globe className="size-5" />
              Add New Website
            </DialogTitle>
            <DialogDescription>
              Enter the details of the website you want to monitor. We'll start checking it immediately.
            </DialogDescription>
          </DialogHeader>

          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label htmlFor="name">Display Name</Label>
              <Input
                id="name"
                placeholder="e.g., Production API"
                value={name}
                onChange={(e) => setName(e.target.value)}
                required
              />
            </div>

            <div className="grid gap-2">
              <Label htmlFor="url">URL</Label>
              <Input
                id="url"
                type="url"
                placeholder="https://example.com"
                value={url}
                onChange={(e) => setUrl(e.target.value)}
                required
              />
              <p className="text-xs text-muted-foreground">
                Enter the full URL including the protocol (https://)
              </p>
            </div>

            <div className="grid gap-2">
              <Label>Monitoring Regions</Label>
              <div className="grid grid-cols-2 gap-2 rounded-lg border border-border p-3">
                {mockRegions.filter(r => r.available).map((region) => (
                  <div key={region.id} className="flex items-center space-x-2">
                    <Checkbox
                      id={region.id}
                      checked={selectedRegions.includes(region.id)}
                      onCheckedChange={() => toggleRegion(region.id)}
                    />
                    <label
                      htmlFor={region.id}
                      className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 cursor-pointer"
                    >
                      {region.code}
                      <span className="ml-1 text-muted-foreground font-normal">
                        ({region.name.split('(')[1]?.replace(')', '') || region.name})
                      </span>
                    </label>
                  </div>
                ))}
              </div>
              <p className="text-xs text-muted-foreground">
                Select at least one region. More regions provide better coverage.
              </p>
            </div>
          </div>

          <DialogFooter>
            <Button type="button" variant="outline" onClick={() => setOpen(false)}>
              Cancel
            </Button>
            <Button type="submit" disabled={isLoading || !name || !url || selectedRegions.length === 0}>
              {isLoading && <Loader2 className="mr-2 size-4 animate-spin" />}
              {isLoading ? 'Adding...' : 'Add Website'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
