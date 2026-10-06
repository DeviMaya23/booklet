import { useState } from 'react'
import { FileIcon, ImageIcon, X } from 'lucide-react'
import { toast } from 'sonner'
import { useQueryClient } from '@tanstack/react-query'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog'
import { type Artist } from '@/features/artists/api/useArtists'
import ArtistCharacterFilter from '@/components/ArtistCharacterFilter'
import { useFiles, FILES_QUERY_KEY } from '../api/useFiles'
import { useArtpieces } from '@/features/artpieces/api/useArtpieces'
import { useArtpiece, type ArtpieceDetail } from '@/features/artpieces/api/useArtpiece'
import { useAttachFilesToArtpiece } from '@/features/artpieces/api/useAttachFilesToArtpiece'

interface AddToExistingArtpieceModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  fileIds: string[]
  onSuccess?: () => void
}

export default function AddToExistingArtpieceModal({
  open,
  onOpenChange,
  fileIds,
  onSuccess,
}: AddToExistingArtpieceModalProps) {
  const [titleSearch, setTitleSearch] = useState('')
  const [selectedArtist, setSelectedArtist] = useState<Artist | null>(null)
  const [selectedCharacters, setSelectedCharacters] = useState<{ id: string; name: string }[]>([])
  const [pickedArtpieceId, setPickedArtpieceId] = useState<string | null>(null)

  const { data: inboxFiles = [] } = useFiles()
  const { data: allArtpieces = [], isLoading: artpiecesLoading } = useArtpieces()
  const { data: pickedArtpiece } = useArtpiece(pickedArtpieceId)
  const attachFiles = useAttachFilesToArtpiece()
  const queryClient = useQueryClient()

  const filteredArtpieces = allArtpieces.filter((a) => {
    if (titleSearch.trim()) {
      const title = a.title ?? ''
      if (!title.toLowerCase().includes(titleSearch.trim().toLowerCase())) return false
    }
    if (selectedArtist && a.artist_id !== selectedArtist.id) return false
    if (selectedCharacters.length > 0) {
      const artpieceCharIds = a.characters.map((c) => c.id)
      if (!selectedCharacters.some((sc) => artpieceCharIds.includes(sc.id))) return false
    }
    return true
  })

  const selectedFiles = fileIds.map((id) => inboxFiles.find((f) => f.id === id)).filter(Boolean)

  function resetState() {
    setTitleSearch('')
    setSelectedArtist(null)
    setSelectedCharacters([])
    setPickedArtpieceId(null)
  }

  function handleOpenChange(next: boolean) {
    if (!next) resetState()
    onOpenChange(next)
  }

  async function handleSave() {
    if (!pickedArtpiece || attachFiles.isPending) return
    try {
      const existingFileIds = pickedArtpiece.files.map((f) => f.id)
      const mergedFileIds = [...new Set([...existingFileIds, ...fileIds])]
      await attachFiles.mutateAsync({ artpieceId: pickedArtpiece.id, fileIds: mergedFileIds })
      await queryClient.invalidateQueries({ queryKey: FILES_QUERY_KEY })
      toast.success('Files added to artpiece')
      onSuccess?.()
      handleOpenChange(false)
    } catch {
      toast.error('Failed to add files to artpiece')
    }
  }

  const canSave = pickedArtpiece !== undefined && !attachFiles.isPending

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className="flex max-h-[calc(100dvh-2rem)] max-w-md flex-col overflow-hidden sm:h-[640px] sm:max-w-3xl">
        <DialogHeader>
          <DialogTitle>Add {fileIds.length} file{fileIds.length !== 1 ? 's' : ''} to artpiece</DialogTitle>
        </DialogHeader>

        {/* Responsive body: single column on mobile, two-column grid on desktop */}
        <div className="flex min-h-0 flex-1 flex-col gap-4 overflow-y-auto pr-1 sm:grid sm:grid-cols-[2fr_3fr] sm:gap-0 sm:overflow-hidden sm:pr-0">

          {/* LEFT COLUMN — file strip + preview (desktop), file strip only (mobile) */}
          <div className="flex flex-col gap-4 sm:overflow-y-auto sm:border-r sm:pr-6">
            {/* Selected file thumbnail strip */}
            <div className="flex max-h-40 flex-wrap gap-2 overflow-y-auto rounded-md border p-2">
              {selectedFiles.map((file) => {
                if (!file) return null
                return (
                  <div key={file.id} className="flex w-14 flex-col gap-1" title={file.name ?? undefined}>
                    <div className="flex h-14 w-14 items-center justify-center overflow-hidden rounded-md bg-muted">
                      {file.thumbnail_url ? (
                        <img src={file.thumbnail_url} alt="" className="h-full w-full object-cover" />
                      ) : (
                        <FileIcon className="h-4 w-4 text-muted-foreground" />
                      )}
                    </div>
                    <p className="truncate text-center text-xs text-muted-foreground">
                      {file.name ?? '—'}
                    </p>
                  </div>
                )
              })}
            </div>

            {/* Mobile-only divider between file strip and search */}
            <hr className="border-border sm:hidden" />

            {/* Preview — desktop only placement (left col), always show when picked */}
            <div className="hidden sm:block">
              {pickedArtpiece ? (
                <ArtpiecePreview artpiece={pickedArtpiece} onClear={() => setPickedArtpieceId(null)} />
              ) : (
                <p className="text-sm text-muted-foreground">No artpiece selected</p>
              )}
            </div>
          </div>

          {/* RIGHT COLUMN — search + results + filters + desktop footer */}
          <div className="flex min-h-0 flex-col gap-3 sm:pl-6">
            <div className="flex flex-col gap-2 sm:min-h-0 sm:flex-1">
              <label className="text-sm font-medium">Artpiece</label>
              <div className="flex gap-2">
                <Input
                  placeholder="Search by title"
                  value={titleSearch}
                  onChange={(e) => setTitleSearch(e.target.value)}
                />
                <ArtistCharacterFilter
                  artist={selectedArtist}
                  onArtistChange={setSelectedArtist}
                  characters={selectedCharacters}
                  onCharactersChange={setSelectedCharacters}
                />
              </div>

              {/* Results — mobile: only when typing; desktop: always visible, fills height */}
              <div className={`${titleSearch.trim() ? 'flex' : 'hidden'} max-h-48 flex-col overflow-y-auto rounded-md border sm:flex sm:max-h-none sm:flex-1`}>
                {artpiecesLoading ? (
                  <p className="py-4 text-center text-sm text-muted-foreground">Loading…</p>
                ) : filteredArtpieces.length === 0 ? (
                  <p className="py-4 text-center text-sm text-muted-foreground">No artpieces found</p>
                ) : (
                  filteredArtpieces.map((a) => (
                    <button
                      key={a.id}
                      type="button"
                      className={`flex items-center gap-2 px-3 py-2 text-left text-sm hover:bg-accent ${pickedArtpieceId === a.id ? 'bg-accent' : ''}`}
                      onClick={() => { setPickedArtpieceId(a.id); setTitleSearch('') }}
                    >
                      <div className="flex h-8 w-8 shrink-0 items-center justify-center overflow-hidden rounded bg-muted">
                        {a.thumbnail_url ? (
                          <img src={a.thumbnail_url} alt="" className="h-full w-full object-cover" />
                        ) : (
                          <ImageIcon className="h-3.5 w-3.5 text-muted-foreground" />
                        )}
                      </div>
                      <span className="truncate">{a.title ?? '—'}</span>
                    </button>
                  ))
                )}
              </div>

              {/* Preview — mobile only, hide while searching */}
              {!titleSearch.trim() && pickedArtpiece && (
                <ArtpiecePreview artpiece={pickedArtpiece} onClear={() => setPickedArtpieceId(null)} className="sm:hidden" />
              )}
            </div>


          </div>
        </div>

        <DialogFooter>
          <Button
            variant="outline"
            type="button"
            onClick={() => handleOpenChange(false)}
            disabled={attachFiles.isPending}
          >
            Cancel
          </Button>
          <Button type="button" onClick={handleSave} disabled={!canSave}>
            Save
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

function ArtpiecePreview({ artpiece, onClear, className }: { artpiece: ArtpieceDetail; onClear: () => void; className?: string }) {
  return (
    <div className={`relative flex gap-4 pt-1 ${className ?? ''}`}>
      <button
        type="button"
        className="absolute right-0 top-0 rounded p-0.5 text-muted-foreground hover:text-foreground"
        onClick={onClear}
        aria-label="Clear picked artpiece"
      >
        <X className="h-4 w-4" />
      </button>
      <div className="flex h-36 w-36 shrink-0 items-center justify-center overflow-hidden rounded-md bg-muted">
        {artpiece.thumbnail_url ? (
          <img src={artpiece.thumbnail_url} alt="" className="h-full w-full object-cover" />
        ) : (
          <ImageIcon className="h-6 w-6 text-muted-foreground" />
        )}
      </div>
      <div className="flex flex-col gap-2 text-sm">
        <div>
          <p className="text-xs text-muted-foreground">Title</p>
          <p>{artpiece.title ?? '—'}</p>
        </div>
        <div>
          <p className="text-xs text-muted-foreground">Artist</p>
          <p>{artpiece.artist_name ?? '—'}</p>
        </div>
        <div>
          <p className="text-xs text-muted-foreground">Character(s)</p>
          <p>{artpiece.characters.length > 0 ? artpiece.characters.map((c) => c.name).join(', ') : '—'}</p>
        </div>
      </div>
    </div>
  )
}
