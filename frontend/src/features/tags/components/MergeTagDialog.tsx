import { useState } from 'react'
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
import type { MergeTarget } from '../api/useMergeTag'

interface MergeTagDialogProps {
  item: MergeTarget | null
  onCancel: () => void
  onConfirm: () => void
  isPending?: boolean
}

export default function MergeTagDialog({ item, onCancel, onConfirm, isPending = false }: MergeTagDialogProps) {
  // Keep the last item so the content doesn't blank out during the close animation.
  const [shown, setShown] = useState(item)
  if (item && item !== shown) setShown(item)

  return (
    <AlertDialog open={item !== null} onOpenChange={(open) => { if (!open && !isPending) onCancel() }}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Merge into “{shown?.target.name}”?</AlertDialogTitle>
          <AlertDialogDescription>
            Artpieces tagged “{shown?.source.name}” will get “{shown?.target.name}” instead. Artpieces that
            have both keep just “{shown?.target.name}”. “{shown?.source.name}” will be deleted. This can’t be
            undone.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel disabled={isPending}>Cancel</AlertDialogCancel>
          <AlertDialogAction
            onClick={(e) => {
              // Radix closes on click; keep it open until the mutation settles.
              e.preventDefault()
              onConfirm()
            }}
            disabled={isPending}
          >
            {isPending ? 'Merging…' : 'Merge tags'}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}
