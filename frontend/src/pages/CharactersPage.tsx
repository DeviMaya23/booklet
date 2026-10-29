import { useState, useMemo } from 'react'
import { Plus, BookOpen } from 'lucide-react'
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
import { useVerses } from '@/features/verses/api/useVerses'
import VersesModal from '@/features/verses/components/VersesModal'

export default function CharactersPage() {
  const [search, setSearch] = useState('')
  const [modalOpen, setModalOpen] = useState(false)
  const [versesModalOpen, setVersesModalOpen] = useState(false)
  const [selectedCharacter, setSelectedCharacter] = useState<Character | undefined>(undefined)
  const [pendingDeleteId, setPendingDeleteId] = useState<string | null>(null)

  const { data: characters = [], isLoading } = useCharacters()
  const { data: verses = [], isLoading: versesLoading } = useVerses()
  const deleteMutation = useDeleteCharacter()

  const filtered = search.trim()
    ? characters.filter((c) =>
        c.name.toLowerCase().includes(search.toLowerCase())
      )
    : characters

  const verseSet = useMemo(() => new Set(verses.map((v) => v.id)), [verses])

  const groups = useMemo(() => {
    const result: { id: string; name: string; characters: Character[] }[] = []
    for (const verse of verses) {
      const members = filtered.filter((c) => c.verse_id === verse.id)
      if (members.length > 0) {
        result.push({ id: verse.id, name: verse.name, characters: members })
      }
    }
    const unsorted = filtered.filter(
      (c) => c.verse_id === null || !verseSet.has(c.verse_id)
    )
    if (unsorted.length > 0) {
      result.push({ id: '__unsorted__', name: 'Unsorted', characters: unsorted })
    }
    return result
  }, [filtered, verses, verseSet])

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
        <div className="flex-1" />
        <Button variant="outline" onClick={() => setVersesModalOpen(true)}>
          <BookOpen className="size-4" />
          Verses
        </Button>
        <Button onClick={handleNewClick}>
          <Plus className="size-4" />
          New Character
        </Button>
      </div>

      {isLoading || versesLoading ? null : verses.length === 0 ? (
        <CharactersGrid
          characters={filtered}
          onDeleteClick={setPendingDeleteId}
          onEditClick={handleEditClick}
        />
      ) : (
        <div className="flex flex-col gap-6">
          {groups.map((group, i) => (
            <div key={group.id}>
              {i > 0 && <div className="border-t mb-6" />}
              <h2 className="text-sm font-semibold text-muted-foreground mb-3">
                {group.name} <span className="font-normal">({group.characters.length})</span>
              </h2>
              <CharactersGrid
                characters={group.characters}
                onDeleteClick={setPendingDeleteId}
                onEditClick={handleEditClick}
              />
            </div>
          ))}
        </div>
      )}

      <CharacterFormModal
        key={selectedCharacter?.id ?? 'new'}
        open={modalOpen}
        onOpenChange={handleModalOpenChange}
        character={selectedCharacter}
      />

      <VersesModal
        open={versesModalOpen}
        onOpenChange={setVersesModalOpen}
        verses={verses}
        characters={characters}
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
