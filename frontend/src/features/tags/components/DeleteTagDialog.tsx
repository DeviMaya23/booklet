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
import type { Tag } from '../api/useTags'

interface DeleteTagDialogProps {
  item: Tag | null
  onCancel: () => void
  onConfirm: () => void
  isPending?: boolean
}

function usageSentence(count: number) {
  if (count === 0) return 'No artpieces use this tag.'
  const noun = count === 1 ? 'artpiece uses' : 'artpieces use'
  const pronoun = count === 1 ? 'It' : 'They'
  return `${count} ${noun} this tag. ${pronoun} will lose it.`
}

export default function DeleteTagDialog({ item, onCancel, onConfirm, isPending = false }: DeleteTagDialogProps) {
  // Keep the last item so the content doesn't blank out during the close animation.
  const [shown, setShown] = useState(item)
  if (item && item !== shown) setShown(item)

  return (
    <AlertDialog open={item !== null} onOpenChange={(open) => { if (!open && !isPending) onCancel() }}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Delete “{shown?.name}”?</AlertDialogTitle>
          <AlertDialogDescription>
            {usageSentence(shown?.artpiece_count ?? 0)} Nothing else about them changes.
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
            {isPending ? 'Deleting…' : 'Delete'}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}
