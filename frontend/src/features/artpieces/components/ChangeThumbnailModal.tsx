import { useState, useMemo } from 'react'
import { ImageOff, Loader2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog'
import { useArtpieces } from '../api/useArtpieces'
import { useUpdateSavedFilter } from '../api/useUpdateSavedFilter'
import { filterAndSortArtpieces } from '../hooks/useArtpiecesFilter'
import { parseFilterPayload } from '../lib/parseFilterPayload'
import { type SavedFilter } from '../api/useCreateSavedFilter'
import { useArtists } from '@/features/artists/api/useArtists'

interface ChangeThumbnailModalProps {
  filter: SavedFilter
  open: boolean
  onOpenChange: (open: boolean) => void
  onSaved: () => void
}

const MAX_PICKER_ITEMS = 10

export default function ChangeThumbnailModal({ filter, open, onOpenChange, onSaved }: ChangeThumbnailModalProps) {
  const artpiecesQuery = useArtpieces()
  const artistsQuery = useArtists()
  const updateFilter = useUpdateSavedFilter()

  const payload = useMemo(() => parseFilterPayload(filter.filter_payload), [filter.filter_payload])

  const filtered = useMemo(() => {
    if (!artpiecesQuery.data) return []
    if (!payload) return artpiecesQuery.data.slice(0, MAX_PICKER_ITEMS)
    const artist = payload.artist_id
      ? (artistsQuery.data?.find((a) => a.id === payload.artist_id) ?? null)
      : null
    return filterAndSortArtpieces(artpiecesQuery.data, {
      search: '',
      sort: 'newest',
      selectedArtist: artist,
      selectedCharacters: payload.character_ids.map((id) => ({ id, name: '' })),
      characterMatch: payload.match_mode,
      isReference: payload.is_reference,
    }).slice(0, MAX_PICKER_ITEMS)
  }, [artpiecesQuery.data, artistsQuery.data, payload])

  const [selectedId, setSelectedId] = useState<string | null>(null)
  const effectiveSelected = selectedId ?? filtered[0]?.id ?? null

  const filterChipLabel = useMemo(() => {
    if (!payload) return filter.name
    const parts: string[] = []
    if (payload.artist_id) parts.push('by artist')
    if (payload.character_ids.length > 0) parts.push(`${payload.character_ids.length} character${payload.character_ids.length > 1 ? 's' : ''}`)
    if (payload.is_reference) parts.push('references only')
    return parts.length > 0 ? parts.join(', ') : 'All artpieces'
  }, [payload, filter.name])

  function handleSave() {
    if (!effectiveSelected) return
    updateFilter.mutate(
      { id: filter.id, thumbnailArtpieceId: effectiveSelected },
      { onSuccess: onSaved },
    )
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-lg">
        <DialogHeader>
          <DialogTitle>Change thumbnail</DialogTitle>
        </DialogHeader>

        <div className="flex flex-col gap-3">
          <div className="flex items-center gap-1.5">
            <span className="text-xs text-muted-foreground">Filter:</span>
            <span className="rounded-full bg-muted px-2 py-0.5 text-xs font-medium">{filterChipLabel}</span>
          </div>

          {artpiecesQuery.isPending ? (
            <div className="flex justify-center py-6">
              <Loader2 className="size-5 animate-spin text-muted-foreground" />
            </div>
          ) : artpiecesQuery.isError ? (
            <p className="py-6 text-center text-sm text-muted-foreground">Failed to load artpieces.</p>
          ) : filtered.length === 0 ? (
            <p className="py-6 text-center text-sm text-muted-foreground">No artpieces match this filter.</p>
          ) : (
            <div className="grid grid-cols-5 gap-2">
              {filtered.map((a) => {
                const isSelected = a.id === effectiveSelected
                return (
                  <button
                    key={a.id}
                    type="button"
                    onClick={() => setSelectedId(a.id)}
                    className={`relative aspect-square overflow-hidden rounded-lg bg-muted transition-all ${isSelected ? 'ring-2 ring-primary ring-offset-2' : 'hover:opacity-80'}`}
                    aria-label={a.title ?? 'Untitled'}
                    aria-pressed={isSelected}
                  >
                    {a.thumbnail_url ? (
                      <img src={a.thumbnail_url} alt={a.title ?? ''} className="size-full object-cover" />
                    ) : (
                      <div className="flex size-full items-center justify-center">
                        <ImageOff className="size-5 text-muted-foreground" />
                      </div>
                    )}
                  </button>
                )
              })}
            </div>
          )}
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>Cancel</Button>
          <Button
            onClick={handleSave}
            disabled={effectiveSelected == null || updateFilter.isPending}
          >
            {updateFilter.isPending ? 'Saving…' : 'Save'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
