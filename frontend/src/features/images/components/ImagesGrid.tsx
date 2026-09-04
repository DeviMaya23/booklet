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
import { type Image } from '../api/useImages'
import { useDeleteImage } from '../api/useDeleteImage'

interface ImagesGridProps {
  images: Image[]
}

export default function ImagesGrid({ images }: ImagesGridProps) {
  const [pendingDeleteId, setPendingDeleteId] = useState<string | null>(null)
  const deleteMutation = useDeleteImage()

  function handleDeleteConfirm() {
    if (!pendingDeleteId) return
    deleteMutation.mutate(pendingDeleteId, {
      onSuccess: () => {
        toast.success('Image deleted')
        setPendingDeleteId(null)
      },
      onError: () => {
        toast.error('Failed to delete image')
        setPendingDeleteId(null)
      },
    })
  }

  return (
    <>
      <div className="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6">
        {images.map((image) => (
          <ResourceCard
            key={image.id}
            imageUrl={image.thumbnail_url}
            label={image.title ?? 'Untitled'}
            onDeleteClick={() => setPendingDeleteId(image.id)}
          />
        ))}
      </div>

      <AlertDialog
        open={pendingDeleteId !== null}
        onOpenChange={(open) => { if (!open) setPendingDeleteId(null) }}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete image?</AlertDialogTitle>
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
