import { useState } from 'react'
import { Link } from 'react-router-dom'
import { ArrowLeft, Image, Layers, Pencil } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { useVerses } from '@/features/verses/api/useVerses'
import { useCharacter } from '../api/useCharacter'
import { usePublicFolders } from '../api/usePublicFolders'
import CharacterEditForm from './CharacterEditForm'
import CharacterArtpieceTabs from './CharacterArtpieceTabs'
import { FolderChipsView } from './CharacterFolderChips'

interface CharacterDetailViewProps {
  characterId: string
}

function BackToCharacters() {
  return (
    <Link
      to="/app/characters"
      className="flex w-fit items-center gap-1 text-sm text-muted-foreground hover:text-foreground"
    >
      <ArrowLeft className="size-4" /> Characters
    </Link>
  )
}

export default function CharacterDetailView({ characterId }: CharacterDetailViewProps) {
  const { data: character, isLoading } = useCharacter(characterId)
  const versesQuery = useVerses()
  const publicFoldersQuery = usePublicFolders()
  const [editing, setEditing] = useState(false)

  if (!character) {
    return (
      <div className="flex flex-col gap-4">
        <BackToCharacters />
        {isLoading ? (
          <p className="text-sm text-muted-foreground">Loading…</p>
        ) : (
          <p className="text-sm text-destructive">Failed to load character.</p>
        )}
      </div>
    )
  }

  if (editing) {
    return <CharacterEditForm character={character} onClose={() => setEditing(false)} />
  }

  const verse = (versesQuery.data ?? []).find((v) => v.id === character.verse_id)

  return (
    <div className="flex flex-col gap-6">
      <BackToCharacters />

      <div className="flex items-start gap-6">
        <div className="size-28 shrink-0 overflow-hidden rounded-2xl bg-muted">
          {character.avatar_url ? (
            <img src={character.avatar_url} alt={character.name} className="size-full object-cover" />
          ) : (
            <div
              className="flex size-full items-center justify-center text-muted-foreground"
              aria-label="No avatar"
            >
              <Image className="size-8" />
            </div>
          )}
        </div>

        <div className="flex min-w-0 flex-1 flex-col gap-3">
          <div className="flex flex-wrap items-center gap-3">
            <h1 className="text-3xl font-bold">{character.name}</h1>
            {verse && (
              <span className="inline-flex items-center gap-1.5 rounded-lg bg-muted px-2.5 py-1 text-sm font-medium">
                <Layers className="size-3.5 text-muted-foreground" />
                {verse.name}
              </span>
            )}
            <Button variant="secondary" size="sm" className="gap-1.5" onClick={() => setEditing(true)}>
              <Pencil className="size-3.5" />
              Edit
            </Button>
          </div>
          {character.notes && <p className="whitespace-pre-wrap">{character.notes}</p>}
          <FolderChipsView folders={character.folders} liveFolders={publicFoldersQuery.data} />
        </div>
      </div>

      <CharacterArtpieceTabs characterId={character.id} characterName={character.name} />
    </div>
  )
}
