import { useEffect, useRef, useState } from 'react'
import { Filter } from 'lucide-react'
import { Button } from '@/components/ui/button'
import TokenInput from '@/components/TokenInput'
import { type Artist } from '@/features/artists/api/useArtists'
import ArtistCombobox from '@/features/artists/components/ArtistCombobox'
import { useCharacters } from '@/features/characters/api/useCharacters'

interface ArtpiecesFilterPopoverProps {
  selectedArtist: Artist | null
  onArtistChange: (artist: Artist | null) => void
  selectedCharacters: { id: string; name: string }[]
  onCharactersChange: (chars: { id: string; name: string }[]) => void
  characterMatch: 'all' | 'any'
  onCharacterMatchChange: (mode: 'all' | 'any') => void
}

function buildSummaryLabel(
  selectedArtist: Artist | null,
  selectedCharacters: { id: string; name: string }[],
  characterMatch: 'all' | 'any',
): string {
  const parts: string[] = []
  if (selectedArtist) parts.push('1 artist')
  if (selectedCharacters.length > 0) {
    parts.push(`${selectedCharacters.length} character${selectedCharacters.length > 1 ? 's' : ''}, ${characterMatch}`)
  }
  if (parts.length === 0) return 'Filter'
  return parts.join(' · ')
}

export default function ArtpiecesFilterPopover({
  selectedArtist,
  onArtistChange,
  selectedCharacters,
  onCharactersChange,
  characterMatch,
  onCharacterMatchChange,
}: ArtpiecesFilterPopoverProps) {
  const [open, setOpen] = useState(false)
  const containerRef = useRef<HTMLDivElement>(null)

  const charactersQuery = useCharacters()

  const allCharacters = charactersQuery.data ?? []

  const availableCharacters = allCharacters
    .map((c) => ({ id: c.id, name: c.name }))
    .filter((c) => !selectedCharacters.some((sc) => sc.id === c.id))

  useEffect(() => {
    if (!open) return
    function handleMouseDown(e: MouseEvent) {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setOpen(false)
      }
    }
    function handleKeyDown(e: KeyboardEvent) {
      if (e.key === 'Escape') setOpen(false)
    }
    document.addEventListener('mousedown', handleMouseDown)
    document.addEventListener('keydown', handleKeyDown)
    return () => {
      document.removeEventListener('mousedown', handleMouseDown)
      document.removeEventListener('keydown', handleKeyDown)
    }
  }, [open])

  const summaryLabel = buildSummaryLabel(selectedArtist, selectedCharacters, characterMatch)

  return (
    <div ref={containerRef} className="relative">
      <Button variant="outline" onClick={() => setOpen((v) => !v)}>
        <Filter className="size-4" />
        {summaryLabel}
      </Button>

      {open && (
        <div className="absolute left-0 top-full z-50 mt-1 w-72 rounded-lg border border-border bg-popover p-4 shadow-md">
          <div className="flex flex-col gap-4">
            <div className="flex flex-col gap-1.5">
              <span className="text-sm font-medium">Artist</span>
              <ArtistCombobox value={selectedArtist} onChange={onArtistChange} />
            </div>

            <div className="flex flex-col gap-1.5">
              <span className="text-sm font-medium">Characters</span>
              {charactersQuery.isError ? (
                <p className="text-sm text-destructive">Failed to load characters.</p>
              ) : (
                <TokenInput
                  items={selectedCharacters}
                  onChange={onCharactersChange}
                  suggestions={availableCharacters}
                  placeholder="Search characters…"
                />
              )}
            </div>

            {selectedCharacters.length > 0 && (
              <div className="flex flex-col gap-1.5">
                <span className="text-sm font-medium">Match</span>
                <div className="flex gap-1">
                  <Button
                    type="button"
                    size="sm"
                    variant={characterMatch === 'all' ? 'default' : 'outline'}
                    onClick={() => onCharacterMatchChange('all')}
                  >
                    All
                  </Button>
                  <Button
                    type="button"
                    size="sm"
                    variant={characterMatch === 'any' ? 'default' : 'outline'}
                    onClick={() => onCharacterMatchChange('any')}
                  >
                    Any
                  </Button>
                </div>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  )
}
