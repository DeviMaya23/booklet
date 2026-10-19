import { useState } from 'react'
import { ImageIcon, Bookmark } from 'lucide-react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog'
import { type ArtpieceSummary } from '../api/useArtpieces'
import { useCreateSavedFilter, type SavedFilter, type DuplicateNameError } from '../api/useCreateSavedFilter'
import { type Artist } from '@/features/artists/api/useArtists'

interface SaveFilterModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  filteredArtpieces: ArtpieceSummary[]
  filterState: {
    artist: Artist | null
    characters: { id: string; name: string }[]
    matchMode: 'all' | 'any'
  }
  onSuccess: (savedFilter: SavedFilter) => void
}

function buildFilterSummary(filterState: SaveFilterModalProps['filterState']): string {
  const parts: string[] = []
  if (filterState.artist) parts.push(`Artist: ${filterState.artist.name}`)
  if (filterState.characters.length > 0) {
    parts.push(`Characters: ${filterState.characters.map((c) => c.name).join(', ')} (${filterState.matchMode})`)
  }
  return parts.length > 0 ? parts.join(' · ') : 'No filters'
}

export default function SaveFilterModal({
  open,
  onOpenChange,
  filteredArtpieces,
  filterState,
  onSuccess,
}: SaveFilterModalProps) {
  const thumbnailCandidates = filteredArtpieces.slice(0, 10)
  const [name, setName] = useState('')
  const [nameError, setNameError] = useState<string | null>(null)
  const [selectedIndex, setSelectedIndex] = useState(0)

  const createMutation = useCreateSavedFilter()

  const selectedArtpiece = thumbnailCandidates[selectedIndex] ?? null
  const hasResults = thumbnailCandidates.length > 0
  const canSave = hasResults && name.trim().length > 0 && !createMutation.isPending

  function handleClose() {
    setName('')
    setNameError(null)
    setSelectedIndex(0)
    onOpenChange(false)
  }

  async function handleSave() {
    if (!canSave) return
    setNameError(null)
    try {
      const result = await createMutation.mutateAsync({
        name: name.trim(),
        thumbnailArtpieceId: selectedArtpiece?.id ?? null,
        filterPayload: {
          artist_id: filterState.artist?.id ?? null,
          character_ids: filterState.characters.map((c) => c.id),
          match_mode: filterState.matchMode,
        },
      })
      toast.success(`Filter "${result.name}" saved`)
      onSuccess(result)
      handleClose()
    } catch (err) {
      const isDuplicate = err instanceof Error && (err as Error & DuplicateNameError).type === 'duplicate_name'
      if (isDuplicate) {
        setNameError(`A saved filter named "${name.trim()}" already exists.`)
      } else {
        toast.error('Failed to save filter')
      }
    }
  }

  const filterSummary = buildFilterSummary(filterState)

  return (
    <Dialog open={open} onOpenChange={(isOpen) => { if (!isOpen) handleClose() }}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Save filter</DialogTitle>
        </DialogHeader>

        <div className="flex flex-col gap-4">
          {/* Filter summary chip */}
          <div className="inline-flex items-center gap-1.5 rounded-full border bg-muted px-3 py-1 text-sm text-muted-foreground self-start">
            <Bookmark className="size-3.5 shrink-0" />
            {filterSummary}
          </div>

          {/* Name input */}
          <div className="flex flex-col gap-1.5">
            <label htmlFor="save-filter-name" className="text-sm font-medium">
              Name
            </label>
            <Input
              id="save-filter-name"
              value={name}
              onChange={(e) => {
                setName(e.target.value)
                if (nameError) setNameError(null)
              }}
              placeholder="e.g. Chev portraits"
              aria-invalid={nameError != null}
              aria-describedby={nameError ? 'save-filter-name-error' : undefined}
            />
            {nameError && (
              <p id="save-filter-name-error" className="text-sm text-destructive" role="alert">
                {nameError}
              </p>
            )}
          </div>

          {/* Thumbnail picker or empty state */}
          {!hasResults ? (
            <p className="text-sm text-muted-foreground text-center py-4">
              No artpieces match this filter
            </p>
          ) : (
            <div className="flex flex-col gap-1.5">
              <span className="text-sm font-medium">Thumbnail</span>
              <div className="grid grid-cols-5 gap-1.5">
                {thumbnailCandidates.map((ap, idx) => (
                  <button
                    key={ap.id}
                    type="button"
                    aria-label={ap.title ?? 'Untitled'}
                    aria-pressed={idx === selectedIndex}
                    onClick={() => setSelectedIndex(idx)}
                    className={`aspect-square overflow-hidden rounded border-2 transition-colors ${
                      idx === selectedIndex
                        ? 'border-primary'
                        : 'border-transparent hover:border-muted-foreground/40'
                    }`}
                  >
                    {ap.thumbnail_url ? (
                      <img
                        src={ap.thumbnail_url}
                        alt={ap.title ?? 'Untitled'}
                        className="h-full w-full object-cover"
                      />
                    ) : (
                      <div className="flex h-full w-full items-center justify-center bg-muted">
                        <ImageIcon className="size-4 text-muted-foreground" aria-hidden />
                      </div>
                    )}
                  </button>
                ))}
              </div>
            </div>
          )}

          {/* Dashboard preview */}
          <div className="flex flex-col gap-1.5">
            <span className="text-sm font-medium text-muted-foreground">Preview</span>
            <div className="flex items-center gap-3 rounded-lg border p-3">
              <div className="size-12 shrink-0 overflow-hidden rounded">
                {selectedArtpiece?.thumbnail_url ? (
                  <img
                    src={selectedArtpiece.thumbnail_url}
                    alt=""
                    className="h-full w-full object-cover"
                  />
                ) : (
                  <div className="flex h-full w-full items-center justify-center bg-muted">
                    <ImageIcon className="size-5 text-muted-foreground" aria-hidden />
                  </div>
                )}
              </div>
              <span className="text-sm font-medium truncate">
                {name.trim() || <span className="text-muted-foreground italic">Filter name</span>}
              </span>
            </div>
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={handleClose} disabled={createMutation.isPending}>
            Cancel
          </Button>
          <Button onClick={handleSave} disabled={!canSave}>
            {createMutation.isPending ? 'Saving…' : 'Save'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
