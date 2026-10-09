import { useState } from 'react'
import { toast } from 'sonner'
import { Link2, Star, X } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
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
import { useArtistLinkEditor } from '../hooks/useArtistLinkEditor'

interface ArtistFormModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  artist?: Artist
  onCreated?: (artist: Artist) => void
}

export default function ArtistFormModal({
  open,
  onOpenChange,
  artist,
  onCreated,
}: ArtistFormModalProps) {
  const isEditMode = artist !== undefined

  const [name, setName] = useState(artist?.name ?? '')
  const [notes, setNotes] = useState(artist?.notes ?? '')
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)

  const linkEditor = useArtistLinkEditor(artist?.links ?? [])

  const createMutation = useCreateArtist()
  const updateMutation = useUpdateArtist()
  const deleteMutation = useDeleteArtist()

  const isPending =
    createMutation.isPending || updateMutation.isPending || deleteMutation.isPending

  function resetForm() {
    setName(artist?.name ?? '')
    setNotes(artist?.notes ?? '')
    linkEditor.reset()
  }

  function handleOpenChange(nextOpen: boolean) {
    if (!nextOpen) resetForm()
    onOpenChange(nextOpen)
  }

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!name.trim()) return
    if (linkEditor.hasErrors()) return

    const filledLinks = linkEditor.buildLinks()

    if (isEditMode) {
      updateMutation.mutate(
        { id: artist.id, name: name.trim(), notes: notes.trim() || null, links: filledLinks },
        {
          onSuccess: () => {
            toast.success('Artist updated')
            handleOpenChange(false)
          },
          onError: (err) => {
            const status = (err as Error & { status?: number }).status
            toast.error(status === 409 ? 'An artist with this name already exists' : 'Failed to update artist')
          },
        },
      )
    } else {
      createMutation.mutate(
        { name: name.trim(), ...(notes.trim() ? { notes: notes.trim() } : {}), links: filledLinks },
        {
          onSuccess: (created) => {
            toast.success('Artist created')
            onCreated?.(created)
            handleOpenChange(false)
          },
          onError: (err) => {
            const status = (err as Error & { status?: number }).status
            toast.error(status === 409 ? 'An artist with this name already exists' : 'Failed to create artist')
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
        handleOpenChange(false)
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

            <fieldset className="flex flex-col gap-1.5 border-0 p-0 m-0">
              <legend className="text-sm font-medium mb-2">Links</legend>
              <div className="flex flex-col gap-2">
                {linkEditor.links.map((link, i) => (
                  <div key={link._key} className="flex flex-col gap-1">
                    <div className="flex items-center gap-1.5">
                      <Link2 className="size-4 text-muted-foreground shrink-0" />
                      <Input
                        value={link.url}
                        onChange={(e) => linkEditor.updateUrl(i, e.target.value)}
                        onBlur={() => linkEditor.handleBlur(i)}
                        placeholder="https://..."
                        className="flex-1"
                      />
                      {linkEditor.links.length > 1 && (
                        <button
                          type="button"
                          onClick={() => linkEditor.setPrimary(i)}
                          title={link.is_primary ? 'Main link' : 'Make main link'}
                          className={[
                            'flex size-8 shrink-0 items-center justify-center rounded border transition-colors',
                            link.is_primary
                              ? 'border-foreground bg-foreground text-background'
                              : 'border-input bg-background text-muted-foreground hover:text-foreground',
                          ].join(' ')}
                        >
                          <Star className="size-3.5" />
                        </button>
                      )}
                      <button
                        type="button"
                        onClick={() => linkEditor.removeLink(i)}
                        className="flex size-8 shrink-0 items-center justify-center rounded text-muted-foreground hover:text-foreground"
                        aria-label="Remove link"
                      >
                        <X className="size-4" />
                      </button>
                    </div>
                    {linkEditor.linkErrors[i] && (
                      <p className="text-xs text-destructive pl-6">{linkEditor.linkErrors[i]}</p>
                    )}
                  </div>
                ))}
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  onClick={linkEditor.addLink}
                  className="w-fit"
                >
                  + Add link
                </Button>
                {linkEditor.links.length > 1 && (
                  <p className="text-xs text-muted-foreground flex items-center gap-1">
                    <Star className="size-3 inline" />
                    Main link: the one opened from lists and the dashboard.
                  </p>
                )}
              </div>
            </fieldset>

            <div className="flex flex-col gap-1.5">
              <label className="text-sm font-medium" htmlFor="artist-notes">
                Notes
              </label>
              <Textarea
                id="artist-notes"
                value={notes}
                onChange={(e) => setNotes(e.target.value)}
                placeholder="Notes about this artist..."
                rows={4}
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
