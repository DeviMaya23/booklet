import { Filter } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Popover, PopoverTrigger, PopoverContent } from '@/components/ui/popover'
import TokenInput from '@/components/TokenInput'
import { type Artist } from '@/features/artists/api/useArtists'
import ArtistCombobox from '@/features/artists/components/ArtistCombobox'
import { useCharacters } from '@/features/characters/api/useCharacters'

type CharacterToken = { id: string; name: string }

interface ArtistCharacterFilterProps {
  artist: Artist | null
  onArtistChange: (artist: Artist | null) => void
  characters: CharacterToken[]
  onCharactersChange: (chars: CharacterToken[]) => void
  characterMatch?: 'all' | 'any'
  onCharacterMatchChange?: (mode: 'all' | 'any') => void
}

function buildSummaryLabel(
  artist: Artist | null,
  characters: CharacterToken[],
  characterMatch?: 'all' | 'any',
): string {
  const parts: string[] = []
  if (artist) parts.push('1 artist')
  if (characters.length > 0) {
    const charLabel = `${characters.length} character${characters.length > 1 ? 's' : ''}`
    parts.push(characterMatch ? `${charLabel}, ${characterMatch}` : charLabel)
  }
  return parts.length === 0 ? 'Filter' : parts.join(' · ')
}

export default function ArtistCharacterFilter({
  artist,
  onArtistChange,
  characters,
  onCharactersChange,
  characterMatch,
  onCharacterMatchChange,
}: ArtistCharacterFilterProps) {
  const charactersQuery = useCharacters()
  const allCharacters = (charactersQuery.data ?? []).map((c) => ({ id: c.id, name: c.name }))
  const availableCharacters = allCharacters.filter((c) => !characters.some((sc) => sc.id === c.id))

  const summaryLabel = buildSummaryLabel(artist, characters, characterMatch)
  const showToggle = characterMatch !== undefined && onCharacterMatchChange !== undefined && characters.length >= 2

  return (
    <Popover>
      <PopoverTrigger className="inline-flex shrink-0 items-center gap-2 rounded-md border border-input bg-background px-3 py-2 text-sm shadow-xs hover:bg-accent hover:text-accent-foreground whitespace-nowrap">
        <Filter className="size-4 shrink-0" />
        {summaryLabel}
      </PopoverTrigger>
      <PopoverContent align="end" className="flex w-72 flex-col gap-4 p-4">
        <div className="flex flex-col gap-1.5">
          <span className="text-sm font-medium">Artist</span>
          <ArtistCombobox value={artist} onChange={onArtistChange} />
        </div>

        <div className="flex flex-col gap-1.5">
          <div className="flex items-center justify-between">
            <span className="text-sm font-medium">Characters</span>
            {showToggle && (
              <div className="flex gap-1">
                <Button
                  type="button"
                  size="sm"
                  variant={characterMatch === 'any' ? 'default' : 'outline'}
                  onClick={() => onCharacterMatchChange!('any')}
                >
                  Any
                </Button>
                <Button
                  type="button"
                  size="sm"
                  variant={characterMatch === 'all' ? 'default' : 'outline'}
                  onClick={() => onCharacterMatchChange!('all')}
                >
                  All
                </Button>
              </div>
            )}
          </div>
          {charactersQuery.isError ? (
            <p className="text-sm text-destructive">Failed to load characters.</p>
          ) : (
            <TokenInput
              items={characters}
              onChange={onCharactersChange}
              suggestions={availableCharacters}
              placeholder={charactersQuery.isPending ? 'Loading…' : 'Search characters…'}
              disabled={charactersQuery.isPending}
            />
          )}
        </div>

        {(artist !== null || characters.length > 0) && (
          <div className="flex justify-end">
            <button
              type="button"
              className="text-sm underline text-muted-foreground hover:text-foreground"
              onClick={() => { onArtistChange(null); onCharactersChange([]) }}
            >
              Clear all
            </button>
          </div>
        )}
      </PopoverContent>
    </Popover>
  )
}
