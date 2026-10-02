import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'

interface FileDeleteDialogProps {
  open: boolean
  count: number
  isPending: boolean
  onConfirm: () => void
  onCancel: () => void
}

export default function FileDeleteDialog({ open, count, isPending, onConfirm, onCancel }: FileDeleteDialogProps) {
  return (
    <AlertDialog open={open} onOpenChange={(o) => { if (!o) onCancel() }}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>
            Delete {count === 1 ? '1 file' : `${count} files`}?
          </AlertDialogTitle>
          <AlertDialogDescription>This action cannot be undone.</AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel disabled={isPending}>Cancel</AlertDialogCancel>
          <AlertDialogAction variant="destructive" onClick={onConfirm} disabled={isPending}>
            Delete
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}
