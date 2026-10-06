import { useRef, useState } from 'react'
import { useQueryClient } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { ArrowLeft, Download, ExternalLink } from 'lucide-react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import TokenInput from '@/components/TokenInput'
import { useArtpiece, artpieceQueryKey } from '../api/useArtpiece'
import { useDownloadArtpiece } from '../api/useDownloadArtpiece'
import { useUpdateArtpiece } from '../api/useUpdateArtpiece'
import { useSetCover } from '../api/useSetCover'
import { useAttachFilesToArtpiece } from '../api/useAttachFilesToArtpiece'
import { ARTISTS_QUERY_KEY, type Artist } from '@/features/artists/api/useArtists'
import ArtistCombobox from '@/features/artists/components/ArtistCombobox'
import { useCharacters } from '@/features/characters/api/useCharacters'
import { useInitFileUpload } from '@/features/files/api/useInitFileUpload'
import { useCompleteFileUpload } from '@/features/files/api/useCompleteFileUpload'
import { apiFetch } from '@/lib/api'
import { stripFileExtension } from '@/features/files/lib/files'
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
  const downloadArtpiece = useDownloadArtpiece()
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
            name: stripFileExtension(file.name),
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

  function unsavedChangesSummary(): string | null {
    if (!artpiece) return null
    const parts: string[] = []
    if (locallyRemovedIds.size > 0)
      parts.push(`${locallyRemovedIds.size} file${locallyRemovedIds.size > 1 ? 's' : ''} will be removed`)
    if (pendingCoverFileId && pendingCoverFileId !== artpiece.cover_file_id)
      parts.push('cover changed')
    const cachedArtists = queryClient.getQueryData<Artist[]>(ARTISTS_QUERY_KEY) ?? []
    const originalArtist = cachedArtists.find((a) => a.id === artpiece.artist_id) ?? null
    if ((editArtist?.id ?? null) !== (originalArtist?.id ?? null)) parts.push('artist changed')
    const originalCharIds = new Set(artpiece.characters.map((c) => c.id))
    const editCharIds = new Set(editCharacters.map((c) => c.id))
    const charsChanged =
      originalCharIds.size !== editCharIds.size ||
      [...originalCharIds].some((id) => !editCharIds.has(id))
    if (charsChanged) parts.push('characters changed')
    if ((editTitle.trim() || null) !== artpiece.title) parts.push('title changed')
    if ((editNotes.trim() || null) !== artpiece.notes) parts.push('notes changed')
    return parts.length > 0 ? `Unsaved changes: ${parts.join(', ')}` : null
  }

  if (isLoading) {
    return (
      <div className="flex flex-col gap-4">
        <button onClick={onClose} className="flex w-fit items-center gap-1 text-sm text-muted-foreground hover:text-foreground">
          <ArrowLeft className="size-4" /> Artpieces
        </button>
        <p className="text-sm text-muted-foreground">Loading…</p>
      </div>
    )
  }

  if (isError || !artpiece) {
    return (
      <div className="flex flex-col gap-4">
        <button onClick={onClose} className="flex w-fit items-center gap-1 text-sm text-muted-foreground hover:text-foreground">
          <ArrowLeft className="size-4" /> Artpieces
        </button>
        <p className="text-sm text-destructive">Failed to load artpiece.</p>
      </div>
    )
  }

  const cachedArtists = queryClient.getQueryData<Artist[]>(ARTISTS_QUERY_KEY) ?? []
  const viewArtist = cachedArtists.find((a) => a.id === artpiece.artist_id) ?? null

  const previewThumbnail =
    mode === 'edit' && pendingCoverFileId
      ? (artpiece.files.find((f) => f.id === pendingCoverFileId)?.thumbnail_url ?? artpiece.thumbnail_url)
      : artpiece.thumbnail_url

  const changesSummary = mode === 'edit' ? unsavedChangesSummary() : null

  return (
    <div className="flex flex-col gap-6">
      {/* Back link */}
      <button
        onClick={onClose}
        className="flex w-fit items-center gap-1 text-sm text-muted-foreground hover:text-foreground"
      >
        <ArrowLeft className="size-4" /> Artpieces
      </button>

      {/* Two-column layout: details | cover */}
      <div className="grid grid-cols-1 gap-6 sm:grid-cols-[1fr_360px]">
        {/* Left column */}
        <div className="flex flex-col gap-4">
          {mode === 'view' ? (
            /* View mode header */
            <div className="flex items-center gap-2">
              <span className="text-2xl font-bold">{artpiece.title ?? 'Untitled'}</span>
              <Button variant="outline" size="sm" onClick={handleEnterEdit}>
                Edit
              </Button>
            </div>
          ) : (
            /* Edit mode header */
            <div className="flex flex-col gap-1">
              <div className="flex items-center gap-2">
                <Input
                  value={editTitle}
                  onChange={(e) => setEditTitle(e.target.value)}
                  placeholder="Artpiece title"
                  className="text-xl font-bold"
                  disabled={isSaving}
                />
                <Button variant="outline" onClick={handleCancel} disabled={isSaving}>
                  Cancel
                </Button>
                <Button onClick={handleSave} disabled={isSaving || uploadingCount > 0}>
                  Save
                </Button>
              </div>
              {changesSummary && (
                <p className="flex items-center gap-1.5 text-sm text-muted-foreground">
                  <span className="size-2 shrink-0 rounded-full bg-amber-500" />
                  {changesSummary}
                </p>
              )}
            </div>
          )}

          {mode === 'view' ? (
            /* View mode fields */
            <div className="flex flex-col gap-4">
              <div className="flex flex-wrap gap-x-8 gap-y-3">
                <div>
                  <p className="text-xs font-medium uppercase tracking-wide text-muted-foreground">Artist</p>
                  {artpiece.artist_name ? (
                    <div className="mt-1 flex items-center gap-1">
                      <span className="text-sm">{artpiece.artist_name}</span>
                      {viewArtist?.artist_link && (
                        <a
                          href={viewArtist.artist_link}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="text-muted-foreground hover:text-foreground"
                        >
                          <ExternalLink className="size-3.5" />
                        </a>
                      )}
                    </div>
                  ) : (
                    <p className="mt-1 text-sm text-muted-foreground">—</p>
                  )}
                </div>

                <div>
                  <p className="text-xs font-medium uppercase tracking-wide text-muted-foreground">Characters</p>
                  {artpiece.characters.length > 0 ? (
                    <div className="mt-1 flex flex-wrap gap-1">
                      {artpiece.characters.map((c) => (
                        <span
                          key={c.id}
                          className="rounded-md bg-muted px-2 py-0.5 text-sm"
                        >
                          {c.name}
                        </span>
                      ))}
                    </div>
                  ) : (
                    <p className="mt-1 text-sm text-muted-foreground">—</p>
                  )}
                </div>
              </div>

              <div>
                <p className="text-xs font-medium uppercase tracking-wide text-muted-foreground">Notes</p>
                <p className="mt-1 text-sm whitespace-pre-wrap">{artpiece.notes ?? '—'}</p>
              </div>
            </div>
          ) : (
            /* Edit mode fields */
            <div className="flex flex-col gap-4">
              <div className="grid grid-cols-2 gap-3">
                <div className="flex flex-col gap-1.5">
                  <label className="text-sm font-medium">Artist</label>
                  <ArtistCombobox value={editArtist} onChange={setEditArtist} disabled={isSaving} />
                </div>
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
              <div className="flex flex-col gap-1.5">
                <label className="text-sm font-medium">Notes</label>
                <Textarea
                  value={editNotes}
                  onChange={(e) => setEditNotes(e.target.value)}
                  placeholder="Notes about this artpiece…"
                  rows={4}
                  disabled={isSaving}
                />
              </div>
            </div>
          )}
        </div>

        {/* Right column — cover tile */}
        <div className="flex flex-col gap-2">
          <div className="aspect-square w-full overflow-hidden rounded-xl bg-muted">
            {previewThumbnail ? (
              <img src={previewThumbnail} alt="" className="h-full w-full object-contain" />
            ) : (
              <div className="flex h-full w-full items-center justify-center text-muted-foreground text-sm">
                No cover
              </div>
            )}
          </div>
          {mode === 'edit' && (
            <p className="text-center text-xs text-muted-foreground">
              Cover preview. Change it with the ☆ on a file below.
            </p>
          )}
        </div>
      </div>

      {/* Files section — full width */}
      <div className="flex flex-col gap-3">
        <div className="flex items-center justify-between">
          <p className="text-sm font-medium">
            Files{' '}
            <span className="ml-1 rounded-md bg-muted px-1.5 py-0.5 text-xs font-normal text-muted-foreground">
              {artpiece.files.length}
            </span>
          </p>
          {mode === 'view' ? (
            <Button
              variant="outline"
              size="sm"
              className="gap-1.5"
              disabled={downloadArtpiece.isPending}
              onClick={() =>
                downloadArtpiece.mutate(
                  { id: artpiece.id, title: artpiece.title },
                  { onError: () => toast.error('Failed to download files') },
                )
              }
            >
              <Download className="size-3.5" />
              Download all
            </Button>
          ) : (
            <Button
              variant="outline"
              size="sm"
              className="gap-1.5"
              onClick={() => fileInputRef.current?.click()}
              disabled={isSaving}
            >
              + Add files
            </Button>
          )}
        </div>

        <input
          ref={fileInputRef}
          type="file"
          multiple
          className="hidden"
          onChange={handleFileInputChange}
        />

        <div className="rounded-lg border border-border p-4">
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
              onUndoRemove={(fileId) =>
                setLocallyRemovedIds((prev) => {
                  const next = new Set(prev)
                  next.delete(fileId)
                  return next
                })
              }
              onAddFiles={() => fileInputRef.current?.click()}
              dragHandlers={{
                dragOver,
                onDragOver: (e) => { e.preventDefault(); setDragOver(true) },
                onDragLeave: () => setDragOver(false),
                onDrop: handleDrop,
              }}
              uploadingCount={uploadingCount}
            />
          )}
        </div>
      </div>

      {/* Delete link — edit mode only */}
      {mode === 'edit' && (
        <button
          type="button"
          className="flex w-fit items-center gap-1.5 text-sm text-destructive hover:text-destructive/80 focus-visible:outline-none"
          onClick={() => setDeleteDialogOpen(true)}
          disabled={isSaving}
        >
          Delete artpiece
        </button>
      )}

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
