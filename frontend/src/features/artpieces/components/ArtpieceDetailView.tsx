import { useRef, useState } from 'react'
import { useQueryClient } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { ChevronLeft, MoreHorizontal, X } from 'lucide-react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import TokenInput from '@/components/TokenInput'
import { useArtpiece, artpieceQueryKey } from '../api/useArtpiece'
import { useUpdateArtpiece } from '../api/useUpdateArtpiece'
import { useSetCover } from '../api/useSetCover'
import { useAttachFilesToArtpiece } from '../api/useAttachFilesToArtpiece'
import { ARTISTS_QUERY_KEY, type Artist } from '@/features/artists/api/useArtists'
import ArtistCombobox from '@/features/artists/components/ArtistCombobox'
import { useCharacters } from '@/features/characters/api/useCharacters'
import { useInitFileUpload } from '@/features/files/api/useInitFileUpload'
import { useCompleteFileUpload } from '@/features/files/api/useCompleteFileUpload'
import { apiFetch } from '@/lib/api'
import ArtpieceDetailFileGrid from './ArtpieceDetailFileGrid'
import DeleteArtpieceDialog from './DeleteArtpieceDialog'

interface CharacterToken {
  id: string
  name: string
}

interface ArtpieceDetailViewProps {
  artpieceId: string
  onClose: () => void
  onDeleted: () => void
}

export default function ArtpieceDetailView({ artpieceId, onClose, onDeleted }: ArtpieceDetailViewProps) {
  const { getToken } = useKindeAuth()
  const queryClient = useQueryClient()

  const { data: artpiece, isLoading, isError } = useArtpiece(artpieceId)

  const [mode, setMode] = useState<'view' | 'edit'>('view')
  const [editTitle, setEditTitle] = useState('')
  const [editNotes, setEditNotes] = useState('')
  const [editArtist, setEditArtist] = useState<Artist | null>(null)
  const [editCharacters, setEditCharacters] = useState<CharacterToken[]>([])
  const [locallyRemovedIds, setLocallyRemovedIds] = useState<Set<string>>(new Set())
  const [pendingCoverFileId, setPendingCoverFileId] = useState<string | null>(null)
  const [uploadingCount, setUploadingCount] = useState(0)
  const [isSaving, setIsSaving] = useState(false)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
  const [dragOver, setDragOver] = useState(false)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const charactersQuery = useCharacters()
  const updateArtpiece = useUpdateArtpiece()
  const setCoverMutation = useSetCover()
  const replaceFiles = useAttachFilesToArtpiece()
  const initUpload = useInitFileUpload()
  const completeUpload = useCompleteFileUpload()

  const allCharacters = charactersQuery.data ?? []

  const availableCharacters = allCharacters
    .map((c) => ({ id: c.id, name: c.name }))
    .filter((c) => !editCharacters.some((ec) => ec.id === c.id))

  function handleEnterEdit() {
    if (!artpiece) return
    setEditTitle(artpiece.title ?? '')
    setEditNotes(artpiece.notes ?? '')
    const cachedArtists = queryClient.getQueryData<Artist[]>(ARTISTS_QUERY_KEY) ?? []
    setEditArtist(cachedArtists.find((a) => a.id === artpiece.artist_id) ?? null)
    setEditCharacters(artpiece.characters.map((c) => ({ id: c.id, name: c.name })))
    setLocallyRemovedIds(new Set())
    setPendingCoverFileId(null)
    setMode('edit')
  }

  function handleCancel() {
    setLocallyRemovedIds(new Set())
    setPendingCoverFileId(null)
    setMode('view')
  }

  async function handleSave() {
    if (!artpiece || isSaving) return
    setIsSaving(true)
    try {
      await updateArtpiece.mutateAsync({
        id: artpieceId,
        title: editTitle.trim() || null,
        notes: editNotes.trim() || null,
        artistId: editArtist?.id ?? null,
        characterIds: editCharacters.map((c) => c.id),
      })

      if (locallyRemovedIds.size > 0) {
        const workingFileIds = artpiece.files
          .map((f) => f.id)
          .filter((id) => !locallyRemovedIds.has(id))
        await replaceFiles.mutateAsync({ artpieceId, fileIds: workingFileIds })
      }

      if (pendingCoverFileId) {
        await setCoverMutation.mutateAsync({ artpieceId, fileId: pendingCoverFileId })
      }

      toast.success('Artpiece updated')
      setMode('view')
    } catch {
      toast.error('Failed to update artpiece')
    } finally {
      setIsSaving(false)
    }
  }

  async function uploadAndAttach(rawFiles: FileList | File[]) {
    const list = Array.from(rawFiles)
    await Promise.all(
      list.map(async (file) => {
        setUploadingCount((c) => c + 1)
        try {
          const result = await initUpload.mutateAsync({
            mimeType: file.type || 'application/octet-stream',
            name: file.name,
          })
          await fetch(result.upload_url, { method: 'PUT', body: file })
          await completeUpload.mutateAsync(result.id)
          const res = await apiFetch(`/artpieces/${artpieceId}/files/${result.id}`, getToken, {
            method: 'POST',
          })
          if (!res.ok) throw new Error('Failed to attach file')
          queryClient.invalidateQueries({ queryKey: artpieceQueryKey(artpieceId) })
        } catch {
          toast.error(`Failed to upload ${file.name}`)
        } finally {
          setUploadingCount((c) => c - 1)
        }
      }),
    )
  }

  function handleDrop(e: React.DragEvent) {
    e.preventDefault()
    setDragOver(false)
    if (e.dataTransfer.files.length > 0) void uploadAndAttach(e.dataTransfer.files)
  }

  function handleFileInputChange(e: React.ChangeEvent<HTMLInputElement>) {
    if (e.target.files && e.target.files.length > 0) {
      void uploadAndAttach(e.target.files)
      e.target.value = ''
    }
  }

  if (isLoading) {
    return (
      <div className="flex flex-col gap-4">
        <button onClick={onClose} className="flex w-fit items-center gap-1 text-sm text-muted-foreground hover:text-foreground">
          <ChevronLeft className="size-4" /> Back
        </button>
        <p className="text-sm text-muted-foreground">Loading…</p>
      </div>
    )
  }

  if (isError || !artpiece) {
    return (
      <div className="flex flex-col gap-4">
        <button onClick={onClose} className="flex w-fit items-center gap-1 text-sm text-muted-foreground hover:text-foreground">
          <ChevronLeft className="size-4" /> Back
        </button>
        <p className="text-sm text-destructive">Failed to load artpiece.</p>
      </div>
    )
  }

  const characterNames = artpiece.characters.map((c) => c.name).join(', ')

  return (
    <div className="flex flex-col gap-6">
      {/* Header */}
      {mode === 'view' ? (
        <div className="flex items-center gap-2">
          <button
            onClick={onClose}
            className="flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground"
          >
            <ChevronLeft className="size-4" />
          </button>
          <span className="text-lg font-semibold">{artpiece.title ?? 'Untitled'}</span>
          <DropdownMenu>
            <DropdownMenuTrigger className="inline-flex size-7 items-center justify-center rounded-md hover:bg-accent focus-visible:outline-none">
              <MoreHorizontal className="size-4" />
            </DropdownMenuTrigger>
            <DropdownMenuContent>
              <DropdownMenuItem onClick={handleEnterEdit}>Edit</DropdownMenuItem>
              <DropdownMenuItem variant="destructive" onClick={() => setDeleteDialogOpen(true)}>
                Delete
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      ) : (
        <div className="flex items-center gap-2">
          <Input
            value={editTitle}
            onChange={(e) => setEditTitle(e.target.value)}
            placeholder="Artpiece title"
            className="w-1/2"
            disabled={isSaving}
          />
          <Button onClick={handleSave} disabled={isSaving || uploadingCount > 0}>
            Save
          </Button>
          <button
            onClick={handleCancel}
            disabled={isSaving}
            className="inline-flex size-8 items-center justify-center rounded-md hover:bg-accent focus-visible:outline-none disabled:opacity-50"
            aria-label="Cancel"
          >
            <X className="size-4" />
          </button>
        </div>
      )}

      {/* Fields */}
      {mode === 'view' ? (
        <div className="grid gap-x-6 gap-y-4 sm:grid-cols-2">
          <div className="flex flex-col gap-4">
            <div>
              <p className="text-sm font-medium mb-1">Notes</p>
              <p className="text-sm text-muted-foreground whitespace-pre-wrap">{artpiece.notes ?? '—'}</p>
            </div>
            <div>
              <p className="text-sm font-medium mb-1">Characters</p>
              <p className="text-sm text-muted-foreground">{characterNames || '—'}</p>
            </div>
          </div>
          <div>
            <p className="text-sm font-medium mb-1">Artist</p>
            <p className="text-sm text-muted-foreground">{artpiece.artist_name ?? '—'}</p>
          </div>
        </div>
      ) : (
        <div className="grid gap-x-6 gap-y-4 sm:grid-cols-2">
          <div className="flex flex-col gap-4">
            {/* Notes */}
            <div className="flex flex-col gap-1.5">
              <label className="text-sm font-medium">Notes</label>
              <Textarea
                value={editNotes}
                onChange={(e) => setEditNotes(e.target.value)}
                placeholder="Notes about this artpiece…"
                rows={3}
                disabled={isSaving}
              />
            </div>
            {/* Characters */}
            <div className="flex flex-col gap-1.5">
              <label className="text-sm font-medium">Characters</label>
              <TokenInput
                items={editCharacters}
                onChange={setEditCharacters}
                suggestions={availableCharacters}
                placeholder="Search characters…"
                disabled={isSaving}
              />
            </div>
          </div>
          {/* Artist */}
          <div className="flex flex-col gap-1.5">
            <label className="text-sm font-medium">Artist</label>
            <ArtistCombobox value={editArtist} onChange={setEditArtist} disabled={isSaving} />
          </div>
        </div>
      )}

      {/* Files section */}
      <div className="flex flex-col gap-3">
        <p className="text-sm font-medium">Files</p>

        <div className="flex flex-col gap-3 rounded-lg border border-border p-4">
        {mode === 'edit' && (
          <>
            <button
              type="button"
              className={`flex w-full cursor-pointer items-center justify-center rounded-lg border border-dashed border-border py-3 text-sm text-muted-foreground transition-colors hover:bg-muted/50 ${dragOver ? 'ring-2 ring-ring bg-muted/50' : ''}`}
              onClick={() => fileInputRef.current?.click()}
              onDragOver={(e) => { e.preventDefault(); setDragOver(true) }}
              onDragLeave={() => setDragOver(false)}
              onDrop={handleDrop}
              disabled={isSaving}
            >
              {dragOver ? 'Drop to add' : uploadingCount > 0 ? `Uploading ${uploadingCount} file${uploadingCount > 1 ? 's' : ''}…` : 'Drag files or pick'}
            </button>
            <input
              ref={fileInputRef}
              type="file"
              multiple
              className="hidden"
              onChange={handleFileInputChange}
            />
          </>
        )}

        {mode === 'view' ? (
          <ArtpieceDetailFileGrid
            files={artpiece.files}
            coverFileId={artpiece.cover_file_id}
            mode="view"
          />
        ) : (
          <ArtpieceDetailFileGrid
            files={artpiece.files}
            coverFileId={artpiece.cover_file_id}
            mode="edit"
            locallyRemovedIds={locallyRemovedIds}
            pendingCoverFileId={pendingCoverFileId}
            onSetCover={(fileId) => setPendingCoverFileId(fileId)}
            onRemove={(fileId) =>
              setLocallyRemovedIds((prev) => new Set([...prev, fileId]))
            }
          />
        )}
        </div>
      </div>

      <DeleteArtpieceDialog
        artpieceId={artpieceId}
        open={deleteDialogOpen}
        onOpenChange={setDeleteDialogOpen}
        fileCount={artpiece.files.length}
        onSuccess={onDeleted}
      />
    </div>
  )
}
