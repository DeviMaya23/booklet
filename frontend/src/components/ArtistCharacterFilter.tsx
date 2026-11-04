import { type ReactNode } from 'react'
import { Filter } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import { Popover, PopoverTrigger, PopoverContent } from '@/components/ui/popover'
import TokenInput from '@/components/TokenInput'
import { type Artist } from '@/features/artists/api/useArtists'
import ArtistCombobox from '@/features/artists/components/ArtistCombobox'
import { useCharacters } from '@/features/characters/api/useCharacters'
import { useTags } from '@/features/tags/api/useTags'
import { type Verse } from '@/features/verses/api/useVerses'
import VerseCombobox from '@/features/verses/components/VerseCombobox'

type CharacterToken = { id: string; name: string }
type TagToken = { id: string; name: string }

interface ArtistCharacterFilterProps {
  artist: Artist | null
  onArtistChange: (artist: Artist | null) => void
  characters: CharacterToken[]
  onCharactersChange: (chars: CharacterToken[]) => void
  characterMatch?: 'all' | 'any'
  onCharacterMatchChange?: (mode: 'all' | 'any') => void
  tags?: TagToken[]
  onTagsChange?: (tags: TagToken[]) => void
  tagMatch?: 'all' | 'any'
  onTagMatchChange?: (mode: 'all' | 'any') => void
  isReference?: boolean
  onIsReferenceChange?: (value: boolean) => void
  selectedVerse?: Verse | null
  onVerseChange?: (verse: Verse | null) => void
  footerAction?: ReactNode
}

function buildSummaryLabel(
  artist: Artist | null,
  characters: CharacterToken[],
  characterMatch?: 'all' | 'any',
  isReference?: boolean,
  selectedVerse?: Verse | null,
  tags?: TagToken[],
  tagMatch?: 'all' | 'any',
): string {
  const parts: string[] = []
  if (artist) parts.push('1 artist')
  if (characters.length > 0) {
    const charLabel = `${characters.length} character${characters.length > 1 ? 's' : ''}`
    parts.push(characterMatch ? `${charLabel}, ${characterMatch}` : charLabel)
  }
  if (tags && tags.length > 0) {
    const tagLabel = `${tags.length} tag${tags.length > 1 ? 's' : ''}`
    parts.push(tagMatch ? `${tagLabel}, ${tagMatch}` : tagLabel)
  }
  if (selectedVerse) parts.push('1 verse')
  if (isReference) parts.push('references')
  return parts.length === 0 ? 'Filter' : parts.join(' · ')
}

export default function ArtistCharacterFilter({
  artist,
  onArtistChange,
  characters,
  onCharactersChange,
  characterMatch,
  onCharacterMatchChange,
  tags,
  onTagsChange,
  tagMatch,
  onTagMatchChange,
  isReference,
  onIsReferenceChange,
  selectedVerse,
  onVerseChange,
  footerAction,
}: ArtistCharacterFilterProps) {
  const charactersQuery = useCharacters({ sort: 'recent' })
  const allCharacters = (charactersQuery.data ?? []).map((c) => ({ id: c.id, name: c.name }))
  const availableCharacters = allCharacters.filter((c) => !characters.some((sc) => sc.id === c.id))

  const showTags = tags !== undefined && onTagsChange !== undefined
  // The tag list is only requested when the caller opts in to the Tags field.
  const tagsQuery = useTags({ sort: 'recent', enabled: showTags })
  const availableTags = (tagsQuery.data ?? []).filter((t) => !(tags ?? []).some((st) => st.id === t.id))

  const summaryLabel = buildSummaryLabel(artist, characters, characterMatch, isReference, selectedVerse, tags, tagMatch)
  const showToggle = characterMatch !== undefined && onCharacterMatchChange !== undefined && characters.length >= 2
  const showVerse = selectedVerse !== undefined && onVerseChange !== undefined
  const hasActiveFilter =
    artist !== null || characters.length > 0 || (tags?.length ?? 0) > 0 || isReference || !!selectedVerse

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

        {tags !== undefined && onTagsChange !== undefined && (
          <div className="flex flex-col gap-1.5">
            <div className="flex items-center justify-between">
              <span className="text-sm font-medium">Tags</span>
              {tagMatch !== undefined && onTagMatchChange !== undefined && tags.length >= 2 && (
                <div className="flex gap-1">
                  <Button
                    type="button"
                    size="sm"
                    variant={tagMatch === 'any' ? 'default' : 'outline'}
                    onClick={() => onTagMatchChange('any')}
                  >
                    Any
                  </Button>
                  <Button
                    type="button"
                    size="sm"
                    variant={tagMatch === 'all' ? 'default' : 'outline'}
                    onClick={() => onTagMatchChange('all')}
                  >
                    All
                  </Button>
                </div>
              )}
            </div>
            {tagsQuery.isError ? (
              <p className="text-sm text-destructive">Failed to load tags.</p>
            ) : (
              <TokenInput
                items={tags}
                onChange={onTagsChange}
                suggestions={availableTags}
                placeholder={tagsQuery.isPending ? 'Loading…' : 'Search tags…'}
                disabled={tagsQuery.isPending}
              />
            )}
          </div>
        )}

        {showVerse && (
          <div className="flex flex-col gap-1.5">
            <span className="text-sm font-medium">Verse</span>
            <VerseCombobox value={selectedVerse ?? null} onChange={onVerseChange!} />
          </div>
        )}

        {isReference !== undefined && onIsReferenceChange !== undefined && (
          <div className="flex items-center justify-between">
            <label htmlFor="filter-is-reference" className="text-sm font-medium cursor-pointer">References only</label>
            <Checkbox
              id="filter-is-reference"
              checked={isReference}
              onCheckedChange={(v) => onIsReferenceChange(v === true)}
            />
          </div>
        )}

        {(hasActiveFilter || footerAction) && (
          <div className="flex items-center justify-between">
            {hasActiveFilter ? (
              <button
                type="button"
                className="text-sm underline text-muted-foreground hover:text-foreground"
                onClick={() => {
                  onArtistChange(null)
                  onCharactersChange([])
                  onIsReferenceChange?.(false)
                  onVerseChange?.(null)
                  onTagsChange?.([])
                }}
              >
                Clear all
              </button>
            ) : <span />}
            {footerAction}
          </div>
        )}
      </PopoverContent>
    </Popover>
  )
}
