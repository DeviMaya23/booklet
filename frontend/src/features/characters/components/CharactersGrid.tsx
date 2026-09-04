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
import { type Character } from '../api/useCharacters'
import { useDeleteCharacter } from '../api/useDeleteCharacter'

interface CharactersGridProps {
  characters: Character[]
}

export default function CharactersGrid({ characters }: CharactersGridProps) {
  const [pendingDeleteId, setPendingDeleteId] = useState<string | null>(null)
  const deleteMutation = useDeleteCharacter()

  function handleDeleteConfirm() {
    if (!pendingDeleteId) return
    deleteMutation.mutate(pendingDeleteId, {
      onSuccess: () => {
        toast.success('Character deleted')
        setPendingDeleteId(null)
      },
      onError: () => {
        toast.error('Failed to delete character')
        setPendingDeleteId(null)
      },
    })
  }

  return (
    <>
      <div className="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6">
        {characters.map((character) => (
          <ResourceCard
            key={character.id}
            imageUrl={character.avatar_url}
            label={character.name}
            onDeleteClick={() => setPendingDeleteId(character.id)}
          />
        ))}
      </div>

      <AlertDialog
        open={pendingDeleteId !== null}
        onOpenChange={(open) => { if (!open) setPendingDeleteId(null) }}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete character?</AlertDialogTitle>
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
