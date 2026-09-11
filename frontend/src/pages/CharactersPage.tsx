import { useState } from 'react'
import { Plus } from 'lucide-react'
import { toast } from 'sonner'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
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
import { useCharacters, type Character } from '@/features/characters/api/useCharacters'
import { useDeleteCharacter } from '@/features/characters/api/useDeleteCharacter'
import CharactersGrid from '@/features/characters/components/CharactersGrid'
import CharacterFormModal from '@/features/characters/components/CharacterFormModal'

export default function CharactersPage() {
  const [search, setSearch] = useState('')
  const [modalOpen, setModalOpen] = useState(false)
  const [selectedCharacter, setSelectedCharacter] = useState<Character | undefined>(undefined)
  const [pendingDeleteId, setPendingDeleteId] = useState<string | null>(null)

  const { data: characters = [], isLoading } = useCharacters()
  const deleteMutation = useDeleteCharacter()

  const filtered = search.trim()
    ? characters.filter((c) =>
        c.name.toLowerCase().includes(search.toLowerCase())
      )
    : characters

  function handleNewClick() {
    setSelectedCharacter(undefined)
    setModalOpen(true)
  }

  function handleEditClick(character: Character) {
    setSelectedCharacter(character)
    setModalOpen(true)
  }

  function handleModalOpenChange(open: boolean) {
    setModalOpen(open)
    if (!open) setSelectedCharacter(undefined)
  }

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
    <div className="flex flex-col gap-4">
      <div className="flex items-center gap-2">
        <Input
          placeholder="Search characters..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="max-w-sm"
        />
        <Button variant="outline" className="ml-auto" onClick={handleNewClick}>
          <Plus />
          New
        </Button>
      </div>

      {isLoading ? null : (
        <CharactersGrid
          characters={filtered}
          onDeleteClick={setPendingDeleteId}
          onEditClick={handleEditClick}
        />
      )}

      <CharacterFormModal
        key={selectedCharacter?.id ?? 'new'}
        open={modalOpen}
        onOpenChange={handleModalOpenChange}
        character={selectedCharacter}
      />

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
    </div>
  )
}
