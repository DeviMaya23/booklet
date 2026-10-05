import { useState } from 'react'
import { toast } from 'sonner'
import ResourceCard from '@/components/ResourceCard'
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
import { type ArtpieceSummary } from '../api/useArtpieces'
import { useDeleteArtpiece } from '../api/useDeleteArtpiece'

interface ArtpiecesGridProps {
  artpieces: ArtpieceSummary[]
}

export default function ArtpiecesGrid({ artpieces }: ArtpiecesGridProps) {
  const [pendingDeleteId, setPendingDeleteId] = useState<string | null>(null)
  const deleteMutation = useDeleteArtpiece()

  function handleDeleteConfirm() {
    if (!pendingDeleteId) return
    deleteMutation.mutate(pendingDeleteId, {
      onSuccess: () => {
        toast.success('Artpiece deleted')
        setPendingDeleteId(null)
      },
      onError: () => {
        toast.error('Failed to delete artpiece')
        setPendingDeleteId(null)
      },
    })
  }

  if (artpieces.length === 0) {
    return <p className="text-sm text-muted-foreground">No artpieces yet.</p>
  }

  return (
    <>
      <div className="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6">
        {artpieces.map((artpiece) => (
          <ResourceCard
            key={artpiece.id}
            imageUrl={artpiece.thumbnail_url}
            label={artpiece.title ?? 'Untitled'}
            sublabel={artpiece.artist_name ?? undefined}
            onDeleteClick={() => setPendingDeleteId(artpiece.id)}
          />
        ))}
      </div>

      <AlertDialog
        open={pendingDeleteId !== null}
        onOpenChange={(open) => { if (!open) setPendingDeleteId(null) }}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete artpiece?</AlertDialogTitle>
            <AlertDialogDescription>
              This cannot be undone. Files attached to this artpiece will remain in your inbox.
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
