import { useEffect, useRef, useState } from 'react'
import { toast } from 'sonner'
import { ChevronDown, Download, Image, Plus, X } from 'lucide-react'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
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
  AlertDialog,
  AlertDialogContent,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogCancel,
  AlertDialogAction,
} from '@/components/ui/alert-dialog'
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
import { apiFetch } from '@/lib/api'
import { type Image as ImageType, IMAGES_QUERY_KEY } from '../api/useImages'
import { useInitImageUpload } from '../api/useInitImageUpload'
import { useCompleteImageUpload } from '../api/useCompleteImageUpload'
import { useUpdateImage } from '../api/useUpdateImage'
import { useDeleteImage } from '../api/useDeleteImage'
import { useArtists, type Artist } from '@/features/artists/api/useArtists'
import ArtistFormModal from '@/features/artists/components/ArtistFormModal'
import { useCharacters } from '@/features/characters/api/useCharacters'

interface CharacterToken {
  id: string
  name: string
}

interface ImageFormModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  image?: ImageType
}

export default function ImageFormModal({ open, onOpenChange, image }: ImageFormModalProps) {
  const isEditMode = image !== undefined

  const [title, setTitle] = useState('')
  const [notes, setNotes] = useState('')
  const [selectedArtist, setSelectedArtist] = useState<Artist | null>(null)
  const [artistSearch, setArtistSearch] = useState('')
  const [selectedCharacters, setSelectedCharacters] = useState<CharacterToken[]>([])
  const [localFile, setLocalFile] = useState<File | null>(null)
  const [localPreviewUrl, setLocalPreviewUrl] = useState<string | null>(null)
  const [dragActive, setDragActive] = useState(false)
  const [imageUrl, setImageUrl] = useState<string | null>(null)
  const [isLoadingImage, setIsLoadingImage] = useState(false)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
  const [artistModalOpen, setArtistModalOpen] = useState(false)

  const fileInputRef = useRef<HTMLInputElement>(null)
  const { getToken } = useKindeAuth()
  const queryClient = useQueryClient()

  const initUploadMutation = useInitImageUpload()
  const completeUploadMutation = useCompleteImageUpload()
  const updateMutation = useUpdateImage()
  const deleteMutation = useDeleteImage()

  const artistsQuery = useArtists()
  const charactersQuery = useCharacters()

  const artists = artistsQuery.data ?? []
  const allCharacters = charactersQuery.data ?? []

  const artistFilter = useComboboxFilter()
  const filteredArtists = artists.filter((a) =>
    artistFilter.contains(a, artistSearch, (a) => a.name),
  )

  const isPending = isSubmitting || deleteMutation.isPending || isLoadingImage

  const availableCharacters = allCharacters
    .map((c) => ({ id: c.id, name: c.name }))
    .filter((c) => !selectedCharacters.some((sc) => sc.id === c.id))

  useEffect(() => {
    if (!open) return

    if (isEditMode && image) {
      setTitle(image.title ?? '')
      setNotes(image.notes ?? '')
      setSelectedCharacters(image.characters.map((c) => ({ id: c.id, name: c.name })))

      if (image.artist_id && image.artist_name) {
        const found = artists.find((a) => a.id === image.artist_id)
        setSelectedArtist(
          found ?? {
            id: image.artist_id,
            name: image.artist_name,
            notes: null,
            artist_link: null,
            created_at: '',
            updated_at: '',
          },
        )
      } else {
        setSelectedArtist(null)
      }

      setIsLoadingImage(true)
      apiFetch(`/images/${image.id}`, getToken)
        .then((res) => res.json())
        .then((data: { image_url?: string | null }) => setImageUrl(data.image_url ?? null))
        .catch(() => setImageUrl(null))
        .finally(() => setIsLoadingImage(false))
    } else if (!isEditMode) {
      resetForm()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open, image?.id])

  useEffect(() => {
    return () => {
      if (localPreviewUrl) URL.revokeObjectURL(localPreviewUrl)
    }
  }, [localPreviewUrl])

  function resetForm() {
    setTitle('')
    setNotes('')
    setSelectedArtist(null)
    setArtistSearch('')
    setSelectedCharacters([])
    setLocalFile(null)
    if (localPreviewUrl) URL.revokeObjectURL(localPreviewUrl)
    setLocalPreviewUrl(null)
    setImageUrl(null)
    setIsLoadingImage(false)
    setDeleteDialogOpen(false)
  }

  function handleOpenChange(nextOpen: boolean) {
    if (!nextOpen) resetForm()
    onOpenChange(nextOpen)
  }

  function handleFileSet(file: File) {
    if (localPreviewUrl) URL.revokeObjectURL(localPreviewUrl)
    setLocalFile(file)
    setLocalPreviewUrl(URL.createObjectURL(file))
  }

  function handleFileChange(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0]
    if (!file) return
    handleFileSet(file)
    e.target.value = ''
  }

  function handleDrop(e: React.DragEvent) {
    e.preventDefault()
    setDragActive(false)
    const file = e.dataTransfer.files?.[0]
    if (file) handleFileSet(file)
  }

  function handleClearFile() {
    if (localPreviewUrl) URL.revokeObjectURL(localPreviewUrl)
    setLocalFile(null)
    setLocalPreviewUrl(null)
  }

  function handleDownload() {
    if (!imageUrl) return
    const a = document.createElement('a')
    a.href = imageUrl
    a.download = image?.title ?? 'image'
    a.click()
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (isEditMode) {
      await handleEditSubmit()
    } else {
      await handleCreateSubmit()
    }
  }

  async function handleCreateSubmit() {
    if (!localFile) return
    setIsSubmitting(true)
    try {
      const initResult = await initUploadMutation.mutateAsync({
        mimeType: localFile.type,
        title: title.trim() || undefined,
        notes: notes.trim() || undefined,
        artistId: selectedArtist?.id || undefined,
        characterIds: selectedCharacters.map((c) => c.id),
      })

      const r2Res = await fetch(initResult.upload_url, {
        method: 'PUT',
        headers: { 'Content-Type': localFile.type },
        body: localFile,
      })

      if (!r2Res.ok) {
        await apiFetch(`/images/${initResult.id}`, getToken, { method: 'DELETE' })
        toast.error('Failed to upload image')
        return
      }

      try {
        await completeUploadMutation.mutateAsync(initResult.id)
      } catch {
        await apiFetch(`/images/${initResult.id}`, getToken, { method: 'DELETE' })
        toast.error('Failed to complete image upload')
        return
      }

      queryClient.invalidateQueries({ queryKey: IMAGES_QUERY_KEY })
      toast.success('Image created')
      handleOpenChange(false)
    } catch {
      toast.error('Failed to create image')
    } finally {
      setIsSubmitting(false)
    }
  }

  async function handleEditSubmit() {
    if (!image) return
    setIsSubmitting(true)
    try {
      await updateMutation.mutateAsync({
        id: image.id,
        title: title.trim() || null,
        notes: notes.trim() || null,
        artistId: selectedArtist?.id ?? null,
        characterIds: selectedCharacters.map((c) => c.id),
      })
      toast.success('Image updated')
      handleOpenChange(false)
    } catch {
      toast.error('Failed to update image')
    } finally {
      setIsSubmitting(false)
    }
  }

  function handleDeleteConfirm() {
    if (!image) return
    deleteMutation.mutate(image.id, {
      onSuccess: () => {
        toast.success('Image deleted')
        setDeleteDialogOpen(false)
        handleOpenChange(false)
      },
      onError: () => {
        toast.error('Failed to delete image')
        setDeleteDialogOpen(false)
      },
    })
  }

  const previewUrl = isEditMode ? (image?.thumbnail_url ?? null) : localPreviewUrl
  const canSave = isEditMode ? true : localFile !== null

  return (
    <>
      <Dialog open={open} onOpenChange={handleOpenChange}>
        <DialogContent className="flex max-h-[calc(100dvh-2rem)] max-w-md flex-col overflow-hidden">
          <DialogHeader>
            <DialogTitle>{isEditMode ? 'Edit image' : 'New image'}</DialogTitle>
          </DialogHeader>

          <form onSubmit={handleSubmit} className="flex min-h-0 flex-1 flex-col gap-4 overflow-y-auto pr-1">
            {/* Image preview / file picker */}
            <div className="relative">
              {isEditMode ? (
                <div className={`flex w-full items-center justify-center overflow-hidden rounded-lg bg-muted${!previewUrl ? ' h-40' : ''}`}>
                  {previewUrl ? (
                    <img src={previewUrl} alt="Image preview" className="max-h-[35vh] max-w-full" />
                  ) : (
                    <Image className="size-10 text-muted-foreground" />
                  )}
                </div>
              ) : (
                <>
                  <button
                    type="button"
                    className={`flex w-full cursor-pointer items-center justify-center overflow-hidden rounded-lg bg-muted transition-colors${!localPreviewUrl ? ' h-40' : ''}${dragActive ? ' ring-2 ring-ring' : ''}`}
                    onClick={() => fileInputRef.current?.click()}
                    onDragOver={(e) => { e.preventDefault(); setDragActive(true) }}
                    onDragLeave={() => setDragActive(false)}
                    onDrop={handleDrop}
                    disabled={isPending}
                  >
                    {localPreviewUrl ? (
                      <img src={localPreviewUrl} alt="Preview" className="max-h-[35vh] max-w-full" />
                    ) : (
                      <div className="flex flex-col items-center gap-2 text-muted-foreground">
                        <Image className="size-10" />
                        <span className="text-sm">{dragActive ? 'Drop to upload' : 'Click or drag to upload'}</span>
                      </div>
                    )}
                  </button>
                  {localPreviewUrl && (
                    <button
                      type="button"
                      className="absolute right-2 top-2 inline-flex size-6 items-center justify-center rounded-full bg-black/50 text-white hover:bg-black/70"
                      onClick={handleClearFile}
                      disabled={isPending}
                    >
                      <X className="size-3" />
                    </button>
                  )}
                  <input
                    ref={fileInputRef}
                    type="file"
                    accept="image/jpeg,image/png"
                    className="hidden"
                    onChange={handleFileChange}
                  />
                </>
              )}
            </div>

            {/* Title */}
            <div className="flex flex-col gap-1.5">
              <label className="text-sm font-medium" htmlFor="image-title">
                Title
              </label>
              <Input
                id="image-title"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                placeholder="Image title"
                disabled={isPending}
              />
            </div>

            {/* Notes */}
            <div className="flex flex-col gap-1.5">
              <label className="text-sm font-medium" htmlFor="image-notes">
                Notes
              </label>
              <textarea
                id="image-notes"
                value={notes}
                onChange={(e) => setNotes(e.target.value)}
                placeholder="Notes about this image..."
                rows={3}
                disabled={isPending}
                className="flex w-full rounded-md border border-input bg-transparent px-3 py-2 text-sm shadow-xs outline-none placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px] disabled:cursor-not-allowed disabled:opacity-50"
              />
            </div>

            {/* Artist */}
            <div className="flex flex-col gap-1.5">
              <label className="text-sm font-medium">Artist</label>
              <div className="flex gap-2">
                <div className="flex-1">
                  <Combobox
                    value={selectedArtist}
                    onValueChange={(v) => setSelectedArtist(v as Artist | null)}
                    itemToStringLabel={(a) => (a as Artist).name}
                    onInputValueChange={(v) => setArtistSearch(v)}
                    disabled={isPending}
                  >
                    <ComboboxInputGroup>
                      <ComboboxInput placeholder="Search artists…" />
                      <ComboboxTrigger>
                        <ChevronDown className="size-4" />
                      </ComboboxTrigger>
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
                <Button
                  type="button"
                  variant="outline"
                  size="icon"
                  onClick={() => setArtistModalOpen(true)}
                  disabled={isPending}
                >
                  <Plus className="size-4" />
                </Button>
              </div>
            </div>

            {/* Characters */}
            <div className="flex flex-col gap-1.5">
              <label className="text-sm font-medium">Characters</label>
              <TokenInput
                items={selectedCharacters}
                onChange={setSelectedCharacters}
                suggestions={availableCharacters}
                placeholder="Search characters…"
                disabled={isPending}
              />
            </div>
          </form>

          <DialogFooter>
            <div className="flex w-full items-center gap-2">
              {isEditMode && (
                <Button
                  variant="destructive"
                  type="button"
                  onClick={() => setDeleteDialogOpen(true)}
                  disabled={isPending}
                >
                  Delete
                </Button>
              )}
              {isEditMode && imageUrl && (
                <Button
                  variant="outline"
                  type="button"
                  onClick={handleDownload}
                  disabled={isPending}
                >
                  <Download className="size-4" />
                  Download
                </Button>
              )}
              <div className="ml-auto flex gap-2">
                <Button
                  variant="outline"
                  type="button"
                  onClick={() => handleOpenChange(false)}
                  disabled={isPending}
                >
                  Cancel
                </Button>
                <Button
                  type="submit"
                  onClick={handleSubmit}
                  disabled={isPending || !canSave}
                >
                  Save
                </Button>
              </div>
            </div>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <AlertDialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete image?</AlertDialogTitle>
            <AlertDialogDescription>This action cannot be undone.</AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              variant="destructive"
              onClick={handleDeleteConfirm}
              disabled={deleteMutation.isPending}
            >
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      <ArtistFormModal
        open={artistModalOpen}
        onOpenChange={setArtistModalOpen}
        onCreated={(artist) => setSelectedArtist(artist)}
      />
    </>
  )
}
