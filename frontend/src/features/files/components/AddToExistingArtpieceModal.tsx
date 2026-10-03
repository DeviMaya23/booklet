import { useState } from 'react'
import { ChevronDown, ChevronUp, FileIcon, ImageIcon, X } from 'lucide-react'
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
import {
  Combobox,
  ComboboxInput,
  ComboboxInputGroup,
  ComboboxTrigger,
  ComboboxPopup,
  ComboboxItem,
  ComboboxEmpty,
  useComboboxFilter,
} from '@/components/ui/combobox'
import TokenInput from '@/components/TokenInput'
import { useArtists, type Artist } from '@/features/artists/api/useArtists'
import { useCharacters } from '@/features/characters/api/useCharacters'
import { useFiles, FILES_QUERY_KEY } from '../api/useFiles'
import { useArtpieces } from '@/features/artpieces/api/useArtpieces'
import { useArtpiece } from '@/features/artpieces/api/useArtpiece'
import { useAttachFilesToArtpiece } from '@/features/artpieces/api/useAttachFilesToArtpiece'

interface CharacterToken {
  id: string
  name: string
}

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
  const [filtersOpen, setFiltersOpen] = useState(false)
  const [selectedArtist, setSelectedArtist] = useState<Artist | null>(null)
  const [artistSearch, setArtistSearch] = useState('')
  const [selectedCharacters, setSelectedCharacters] = useState<CharacterToken[]>([])
  const [pickedArtpieceId, setPickedArtpieceId] = useState<string | null>(null)
  const [isSubmitting, setIsSubmitting] = useState(false)

  const { data: inboxFiles = [] } = useFiles()
  const { data: allArtpieces = [] } = useArtpieces()
  const { data: pickedArtpiece } = useArtpiece(pickedArtpieceId)
  const artistsQuery = useArtists()
  const charactersQuery = useCharacters()
  const attachFiles = useAttachFilesToArtpiece()
  const queryClient = useQueryClient()

  const artists = artistsQuery.data ?? []
  const allCharacters = charactersQuery.data ?? []

  const artistFilter = useComboboxFilter()
  const filteredArtists = artists.filter((a) => artistFilter.contains(a, artistSearch, (a) => a.name))

  const availableCharacters = allCharacters
    .map((c) => ({ id: c.id, name: c.name }))
    .filter((c) => !selectedCharacters.some((sc) => sc.id === c.id))

  const filteredArtpieces = allArtpieces.filter((a) => {
    if (titleSearch.trim()) {
      const title = a.title ?? ''
      if (!title.toLowerCase().includes(titleSearch.trim().toLowerCase())) return false
    }
    if (selectedArtist && a.artist_id !== selectedArtist.id) return false
    if (selectedCharacters.length > 0) {
      const artpieceCharIds = a.characters.map((c) => c.id)
      if (!selectedCharacters.every((sc) => artpieceCharIds.includes(sc.id))) return false
    }
    return true
  })

  const selectedFiles = fileIds.map((id) => inboxFiles.find((f) => f.id === id)).filter(Boolean)

  function resetState() {
    setTitleSearch('')
    setFiltersOpen(false)
    setSelectedArtist(null)
    setArtistSearch('')
    setSelectedCharacters([])
    setPickedArtpieceId(null)
    setIsSubmitting(false)
  }

  function handleOpenChange(next: boolean) {
    if (!next) resetState()
    onOpenChange(next)
  }

  async function handleSave() {
    if (!pickedArtpiece || isSubmitting) return
    setIsSubmitting(true)
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
    } finally {
      setIsSubmitting(false)
    }
  }

  const canSave = pickedArtpiece !== undefined && !isSubmitting

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className="flex max-h-[calc(100dvh-2rem)] max-w-md flex-col overflow-hidden">
        <DialogHeader>
          <DialogTitle>Add {fileIds.length} file{fileIds.length !== 1 ? 's' : ''} to artpiece</DialogTitle>
        </DialogHeader>

        <div className="flex min-h-0 flex-1 flex-col gap-4 overflow-y-auto pr-1">
          {/* Selected file thumbnail strip */}
          <div className="flex max-h-40 flex-wrap gap-2 overflow-y-auto">
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

          <hr className="border-border" />

          {/* Artpiece search section */}
          <div className="flex flex-col gap-2">
            <label className="text-sm font-medium">Artpiece</label>
            <Input
              placeholder="Search by title"
              value={titleSearch}
              onChange={(e) => setTitleSearch(e.target.value)}
            />

            {/* Search results — only when query is typed */}
            {titleSearch.trim() && filteredArtpieces.length > 0 && (
              <div className="flex max-h-48 flex-col overflow-y-auto rounded-md border">
                {filteredArtpieces.map((a) => (
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
                ))}
              </div>
            )}

            {titleSearch.trim() && filteredArtpieces.length === 0 && (
              <p className="text-center text-sm text-muted-foreground">No artpieces found</p>
            )}

            {/* Picked artpiece preview */}
            {pickedArtpiece && !titleSearch.trim() && (
              <div className="relative flex gap-4 pt-1">
                <button
                  type="button"
                  className="absolute right-0 top-0 rounded p-0.5 text-muted-foreground hover:text-foreground"
                  onClick={() => setPickedArtpieceId(null)}
                  aria-label="Clear picked artpiece"
                >
                  <X className="h-4 w-4" />
                </button>
                <div className="flex h-36 w-36 shrink-0 items-center justify-center overflow-hidden rounded-md bg-muted">
                  {pickedArtpiece.thumbnail_url ? (
                    <img src={pickedArtpiece.thumbnail_url} alt="" className="h-full w-full object-cover" />
                  ) : (
                    <ImageIcon className="h-6 w-6 text-muted-foreground" />
                  )}
                </div>
                <div className="flex flex-col gap-2 text-sm">
                  <div>
                    <p className="text-xs text-muted-foreground">Title</p>
                    <p>{pickedArtpiece.title ?? '—'}</p>
                  </div>
                  <div>
                    <p className="text-xs text-muted-foreground">Artist</p>
                    <p>{pickedArtpiece.artist_name ?? '—'}</p>
                  </div>
                  <div>
                    <p className="text-xs text-muted-foreground">Character(s)</p>
                    <p>{pickedArtpiece.characters.length > 0 ? pickedArtpiece.characters.map((c) => c.name).join(', ') : '—'}</p>
                  </div>
                </div>
              </div>
            )}
          </div>

          {/* Filters toggle */}
          <div className="flex flex-col gap-2">
            <button
              type="button"
              className="flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground"
              onClick={() => setFiltersOpen((v) => !v)}
            >
              Filters
              {filtersOpen ? <ChevronUp className="h-3.5 w-3.5" /> : <ChevronDown className="h-3.5 w-3.5" />}
            </button>

            {filtersOpen && (
              <div className="flex flex-col gap-3">
                <div className="flex flex-col gap-1.5">
                  <label className="text-sm font-medium">Artist</label>
                  <Combobox
                    value={selectedArtist}
                    onValueChange={(v) => setSelectedArtist(v as Artist | null)}
                    itemToStringLabel={(a) => (a as Artist).name}
                    onInputValueChange={(v) => setArtistSearch(v)}
                  >
                    <ComboboxInputGroup>
                      <ComboboxInput placeholder="Filter by artist…" />
                      {selectedArtist ? (
                        <button
                          type="button"
                          className="absolute right-2 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
                          onClick={() => { setSelectedArtist(null); setArtistSearch('') }}
                          aria-label="Clear artist filter"
                        >
                          <X className="size-4" />
                        </button>
                      ) : (
                        <ComboboxTrigger>
                          <ChevronDown className="size-4" />
                        </ComboboxTrigger>
                      )}
                    </ComboboxInputGroup>
                    <ComboboxPopup>
                      {filteredArtists.length === 0 && <ComboboxEmpty>No artists found</ComboboxEmpty>}
                      {filteredArtists.map((artist) => (
                        <ComboboxItem key={artist.id} value={artist}>
                          {artist.name}
                        </ComboboxItem>
                      ))}
                    </ComboboxPopup>
                  </Combobox>
                </div>

                <div className="flex flex-col gap-1.5">
                  <label className="text-sm font-medium">Characters</label>
                  <TokenInput
                    items={selectedCharacters}
                    onChange={setSelectedCharacters}
                    suggestions={availableCharacters}
                    placeholder="Filter by character…"
                  />
                </div>
              </div>
            )}
          </div>
        </div>

        <DialogFooter>
          <Button
            variant="outline"
            type="button"
            onClick={() => handleOpenChange(false)}
            disabled={isSubmitting}
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
