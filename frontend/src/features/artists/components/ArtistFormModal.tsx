import { useState } from 'react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog'
import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogCancel,
  AlertDialogAction,
} from '@/components/ui/alert-dialog'
import { type Artist } from '../api/useArtists'
import { useCreateArtist } from '../api/useCreateArtist'
import { useUpdateArtist } from '../api/useUpdateArtist'
import { useDeleteArtist } from '../api/useDeleteArtist'

interface ArtistFormModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  artist?: Artist
}

export default function ArtistFormModal({
  open,
  onOpenChange,
  artist,
}: ArtistFormModalProps) {
  const isEditMode = artist !== undefined

  const [name, setName] = useState(artist?.name ?? '')
  const [artistLink, setArtistLink] = useState(artist?.artist_link ?? '')
  const [notes, setNotes] = useState(artist?.notes ?? '')
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)

  const createMutation = useCreateArtist()
  const updateMutation = useUpdateArtist()
  const deleteMutation = useDeleteArtist()

  const isPending =
    createMutation.isPending || updateMutation.isPending || deleteMutation.isPending

  function resetForm() {
    setName(artist?.name ?? '')
    setArtistLink(artist?.artist_link ?? '')
    setNotes(artist?.notes ?? '')
  }

  function handleOpenChange(nextOpen: boolean) {
    if (!nextOpen) resetForm()
    onOpenChange(nextOpen)
  }

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!name.trim()) return

    if (isEditMode) {
      updateMutation.mutate(
        {
          id: artist.id,
          name: name.trim(),
          artist_link: artistLink.trim() || null,
          notes: notes.trim() || null,
        },
        {
          onSuccess: () => {
            toast.success('Artist updated')
            onOpenChange(false)
          },
          onError: (err) => {
            const status = (err as Error & { status?: number }).status
            if (status === 409) {
              toast.error('An artist with this name already exists')
            } else if (status === 422) {
              toast.error('Link must be a valid URL')
            } else {
              toast.error('Failed to update artist')
            }
          },
        },
      )
    } else {
      createMutation.mutate(
        {
          name: name.trim(),
          ...(artistLink.trim() ? { artist_link: artistLink.trim() } : {}),
          ...(notes.trim() ? { notes: notes.trim() } : {}),
        },
        {
          onSuccess: () => {
            toast.success('Artist created')
            onOpenChange(false)
          },
          onError: (err) => {
            const status = (err as Error & { status?: number }).status
            if (status === 409) {
              toast.error('An artist with this name already exists')
            } else if (status === 422) {
              toast.error('Link must be a valid URL')
            } else {
              toast.error('Failed to create artist')
            }
          },
        },
      )
    }
  }

  function handleDeleteConfirm() {
    if (!artist) return
    deleteMutation.mutate(artist.id, {
      onSuccess: () => {
        toast.success('Artist deleted')
        setDeleteDialogOpen(false)
        onOpenChange(false)
      },
      onError: () => {
        toast.error('Failed to delete artist')
        setDeleteDialogOpen(false)
      },
    })
  }

  return (
    <>
      <Dialog open={open} onOpenChange={handleOpenChange}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle>{isEditMode ? 'Edit artist' : 'New artist'}</DialogTitle>
          </DialogHeader>

          <form onSubmit={handleSubmit} className="flex flex-col gap-4">
            <div className="flex flex-col gap-1.5">
              <label className="text-sm font-medium" htmlFor="artist-name">
                Name <span className="text-destructive">*</span>
              </label>
              <Input
                id="artist-name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="Artist name"
                required
              />
            </div>

            <div className="flex flex-col gap-1.5">
              <label className="text-sm font-medium" htmlFor="artist-link">
                Link
              </label>
              <Input
                id="artist-link"
                value={artistLink}
                onChange={(e) => setArtistLink(e.target.value)}
                placeholder="https://..."
              />
            </div>

            <div className="flex flex-col gap-1.5">
              <label className="text-sm font-medium" htmlFor="artist-notes">
                Blurb
              </label>
              <textarea
                id="artist-notes"
                value={notes}
                onChange={(e) => setNotes(e.target.value)}
                placeholder="Notes about this artist..."
                rows={4}
                className="flex w-full rounded-md border border-input bg-transparent px-3 py-2 text-sm shadow-xs outline-none placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px] disabled:cursor-not-allowed disabled:opacity-50"
              />
            </div>
          </form>

          <DialogFooter>
            {isEditMode && (
              <Button
                variant="destructive"
                type="button"
                onClick={() => setDeleteDialogOpen(true)}
                disabled={isPending}
                className="mr-auto"
              >
                Delete
              </Button>
            )}
            <Button
              variant="outline"
              type="button"
              onClick={() => handleOpenChange(false)}
              disabled={isPending}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              onClick={handleSubmit}
              disabled={isPending || !name.trim()}
            >
              Save
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <AlertDialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete artist?</AlertDialogTitle>
            <AlertDialogDescription>
              This action cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              variant="destructive"
              onClick={handleDeleteConfirm}
              disabled={deleteMutation.isPending}
            >
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  )
}
