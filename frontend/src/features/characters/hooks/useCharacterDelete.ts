import { useState } from 'react'
import { toast } from 'sonner'
import { useDeleteCharacter } from '../api/useDeleteCharacter'

interface UseCharacterDeleteOptions {
  characterId: string
  onDeleted: () => void
}

/** Confirm-then-delete lifecycle for one character. */
export function useCharacterDelete({ characterId, onDeleted }: UseCharacterDeleteOptions) {
  const [confirmOpen, setConfirmOpen] = useState(false)
  const mutation = useDeleteCharacter()

  function confirm() {
    mutation.mutate(characterId, {
      onSuccess: () => {
        toast.success('Character deleted')
        onDeleted()
      },
      onError: () => {
        toast.error('Failed to delete character')
        setConfirmOpen(false)
      },
    })
  }

  return {
    confirmOpen,
    askToDelete: () => setConfirmOpen(true),
    cancel: () => setConfirmOpen(false),
    confirm,
    isPending: mutation.isPending,
  }
}
