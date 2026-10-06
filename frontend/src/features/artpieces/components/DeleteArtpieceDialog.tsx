import { useState } from 'react'
import { toast } from 'sonner'
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
import { useDeleteArtpiece } from '../api/useDeleteArtpiece'

interface DeleteArtpieceDialogProps {
  artpieceId: string
  open: boolean
  onOpenChange: (open: boolean) => void
  onSuccess?: () => void
  fileCount?: number
}

export default function DeleteArtpieceDialog({
  artpieceId,
  open,
  onOpenChange,
  onSuccess,
  fileCount,
}: DeleteArtpieceDialogProps) {
  const [deleteFiles, setDeleteFiles] = useState(false)
  const deleteMutation = useDeleteArtpiece()

  function handleConfirm() {
    deleteMutation.mutate({ id: artpieceId, deleteFiles }, {
      onSuccess: () => {
        toast.success('Artpiece deleted')
        setDeleteFiles(false)
        onSuccess?.()
      },
      onError: () => {
        toast.error('Failed to delete artpiece')
      },
    })
  }

  function handleOpenChange(next: boolean) {
    if (!next) setDeleteFiles(false)
    onOpenChange(next)
  }

  const showFileCheckbox = fileCount === undefined || fileCount > 0
  const checkboxLabel =
    fileCount !== undefined
      ? `Also delete ${fileCount} attached file${fileCount === 1 ? '' : 's'}`
      : 'Also delete attached files'

  return (
    <AlertDialog open={open} onOpenChange={handleOpenChange}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Delete artpiece?</AlertDialogTitle>
          <AlertDialogDescription>
            This cannot be undone.
          </AlertDialogDescription>
        </AlertDialogHeader>

        {showFileCheckbox && (
          <label className="flex items-center gap-2 text-sm">
            <input
              type="checkbox"
              checked={deleteFiles}
              onChange={(e) => setDeleteFiles(e.target.checked)}
              className="accent-destructive"
            />
            {checkboxLabel}
          </label>
        )}

        <AlertDialogFooter>
          <AlertDialogCancel disabled={deleteMutation.isPending}>Cancel</AlertDialogCancel>
          <AlertDialogAction
            variant="destructive"
            onClick={handleConfirm}
            disabled={deleteMutation.isPending}
          >
            Delete
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}
