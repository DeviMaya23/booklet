import { useEffect, useRef, useState } from 'react'
import { toast } from 'sonner'
import { Image, X } from 'lucide-react'
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
import { apiFetch } from '@/lib/api'
import { type Character, CHARACTERS_QUERY_KEY } from '../api/useCharacters'
import { useCreateCharacter } from '../api/useCreateCharacter'
import { useUpdateCharacter } from '../api/useUpdateCharacter'
import { useDeleteCharacter } from '../api/useDeleteCharacter'
import { useInitAvatarUpload } from '../api/useInitAvatarUpload'
import { useCompleteAvatarUpload } from '../api/useCompleteAvatarUpload'
import { usePublicFolders, type PublicFolder } from '../api/usePublicFolders'
import { FolderPicker } from './FolderPicker'

interface CharacterFormModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  character?: Character
}

export default function CharacterFormModal({
  open,
  onOpenChange,
  character,
}: CharacterFormModalProps) {
  const isEditMode = character !== undefined

  const [name, setName] = useState(character?.name ?? '')
  const [notes, setNotes] = useState(character?.notes ?? '')
  const [localFile, setLocalFile] = useState<File | null>(null)
  const [localPreviewUrl, setLocalPreviewUrl] = useState<string | null>(null)
  const [avatarCleared, setAvatarCleared] = useState(false)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
  const [isSubmitting, setIsSubmitting] = useState(false)

  // Folder diff state: track adds/removes against the original character folders
  const [addedFolders, setAddedFolders] = useState<PublicFolder[]>([])
  const [removedIds, setRemovedIds] = useState<Set<string>>(new Set())
  const [renamedFolders, setRenamedFolders] = useState<Map<string, string>>(new Map())

  const fileInputRef = useRef<HTMLInputElement>(null)
  const { getToken } = useKindeAuth()
  const queryClient = useQueryClient()

  const createMutation = useCreateCharacter()
  const updateMutation = useUpdateCharacter()
  const deleteMutation = useDeleteCharacter()
  const initUploadMutation = useInitAvatarUpload()
  const completeUploadMutation = useCompleteAvatarUpload()
  const publicFoldersQuery = usePublicFolders()

  const isPending = isSubmitting || deleteMutation.isPending

  // Compute current selected folders: (original + added) - removed
  const originalFolders: PublicFolder[] = (character?.folders ?? []).map((f) => ({
    id: f.id,
    name: f.name,
  }))
  const selectedFolders: PublicFolder[] = [
    ...originalFolders
      .filter((f) => !removedIds.has(f.id))
      .map((f) => ({ ...f, name: renamedFolders.get(f.id) ?? f.name })),
    ...addedFolders,
  ]
  const selectedIds = new Set(selectedFolders.map((f) => f.id))
  const availableFolders: PublicFolder[] = (publicFoldersQuery.data ?? []).filter(
    (f) => !selectedIds.has(f.id),
  )
  const folderPickerDisabled = publicFoldersQuery.isLoading || publicFoldersQuery.isError

  useEffect(() => {
    return () => {
      if (localPreviewUrl) URL.revokeObjectURL(localPreviewUrl)
    }
  }, [localPreviewUrl])

  useEffect(() => {
    if (!publicFoldersQuery.data) return
    const liveMap = new Map<string, string>(publicFoldersQuery.data.map((f) => [f.id, f.name]))
    const nameOverrides = new Map<string, string>()
    const newRemovedIds = new Set<string>()
    for (const f of originalFolders) {
      const liveName = liveMap.get(f.id)
      if (liveName === undefined) {
        newRemovedIds.add(f.id)
      } else if (liveName !== f.name) {
        nameOverrides.set(f.id, liveName)
      }
    }
    if (nameOverrides.size > 0) {
      setRenamedFolders(nameOverrides)
    }
    if (newRemovedIds.size > 0) {
      setRemovedIds((prev) => new Set([...prev, ...newRemovedIds]))
    }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [publicFoldersQuery.data])

  const displayUrl = localPreviewUrl ?? (avatarCleared ? null : character?.avatar_url ?? null)
  const showClearButton = isEditMode && !avatarCleared && (localPreviewUrl !== null || character?.avatar_url != null)

  function resetForm() {
    setName(character?.name ?? '')
    setNotes(character?.notes ?? '')
    setLocalFile(null)
    setLocalPreviewUrl(null)
    setAvatarCleared(false)
    setDeleteDialogOpen(false)
    setAddedFolders([])
    setRemovedIds(new Set())
    setRenamedFolders(new Map())
  }

  function handleOpenChange(nextOpen: boolean) {
    if (!nextOpen) resetForm()
    onOpenChange(nextOpen)
  }

  function handleFileChange(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0]
    if (!file) return
    const url = URL.createObjectURL(file)
    setLocalFile(file)
    setLocalPreviewUrl(url)
    setAvatarCleared(false)
    e.target.value = ''
  }

  function handleClearAvatar() {
    setLocalFile(null)
    setLocalPreviewUrl(null)
    setAvatarCleared(true)
  }

  async function runAvatarUpload(characterId: string, file: File) {
    const initResult = await initUploadMutation.mutateAsync({ characterId, mimeType: file.type })
    const r2Res = await fetch(initResult.upload_url, {
      method: 'PUT',
      headers: { 'Content-Type': file.type },
      body: file,
    })
    if (!r2Res.ok) throw new Error('R2 upload failed')
    await completeUploadMutation.mutateAsync({ characterId, uploadId: initResult.id })
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!name.trim()) return

    setIsSubmitting(true)
    try {
      if (isEditMode) {
        await handleEditSubmit()
      } else {
        await handleCreateSubmit()
      }
    } catch {
      toast.error(isEditMode ? 'Failed to update character' : 'Failed to create character')
    } finally {
      setIsSubmitting(false)
    }
  }

  async function handleCreateSubmit() {
    const finalFolderIds = selectedFolders.map((f) => f.id)
    const created = await createMutation.mutateAsync({
      name: name.trim(),
      ...(notes.trim() ? { notes: notes.trim() } : {}),
      folder_ids: finalFolderIds,
    })

    if (localFile) {
      try {
        await runAvatarUpload(created.id, localFile)
      } catch {
        apiFetch(`/characters/${created.id}`, getToken, { method: 'DELETE' })
        toast.error('Failed to upload avatar')
        return
      }
      queryClient.invalidateQueries({ queryKey: CHARACTERS_QUERY_KEY })
    }

    toast.success('Character created')
    onOpenChange(false)
  }

  async function handleEditSubmit() {
    if (!character) return

    if (localFile) {
      try {
        await runAvatarUpload(character.id, localFile)
      } catch {
        toast.error('Failed to upload avatar')
        return
      }
    } else if (avatarCleared) {
      const res = await apiFetch(`/characters/${character.id}/avatar`, getToken, { method: 'DELETE' })
      if (!res.ok) {
        toast.error('Failed to remove avatar')
        return
      }
    }

    const finalFolderIds = selectedFolders.map((f) => f.id)
    await updateMutation.mutateAsync({
      id: character.id,
      name: name.trim(),
      notes: notes.trim() || null,
      folder_ids: finalFolderIds,
    })

    toast.success('Character updated')
    onOpenChange(false)
  }

  function handleDeleteConfirm() {
    if (!character) return
    deleteMutation.mutate(character.id, {
      onSuccess: () => {
        toast.success('Character deleted')
        setDeleteDialogOpen(false)
        onOpenChange(false)
      },
      onError: () => {
        toast.error('Failed to delete character')
        setDeleteDialogOpen(false)
      },
    })
  }

  return (
    <>
      <Dialog open={open} onOpenChange={handleOpenChange}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle>{isEditMode ? 'Edit character' : 'New character'}</DialogTitle>
          </DialogHeader>

          <form onSubmit={handleSubmit} className="flex flex-col gap-4">
            <div className="relative">
              <button
                type="button"
                className="flex h-40 w-full cursor-pointer items-center justify-center overflow-hidden rounded-lg bg-muted"
                onClick={() => fileInputRef.current?.click()}
                disabled={isPending}
              >
                {displayUrl ? (
                  <img src={displayUrl} alt="Avatar preview" className="size-full object-cover" />
                ) : (
                  <Image className="size-10 text-muted-foreground" />
                )}
              </button>
              {showClearButton && (
                <button
                  type="button"
                  className="absolute right-2 top-2 inline-flex size-6 items-center justify-center rounded-full bg-black/50 text-white hover:bg-black/70"
                  onClick={handleClearAvatar}
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
            </div>

            <div className="flex flex-col gap-1.5">
              <label className="text-sm font-medium" htmlFor="character-name">
                Name <span className="text-destructive">*</span>
              </label>
              <Input
                id="character-name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="Character name"
                required
                disabled={isPending}
              />
            </div>

            <div className="flex flex-col gap-1.5">
              <label className="text-sm font-medium" htmlFor="character-notes">
                Notes
              </label>
              <textarea
                id="character-notes"
                value={notes}
                onChange={(e) => setNotes(e.target.value)}
                placeholder="Notes about this character..."
                rows={4}
                disabled={isPending}
                className="flex w-full rounded-md border border-input bg-transparent px-3 py-2 text-sm shadow-xs outline-none placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px] disabled:cursor-not-allowed disabled:opacity-50"
              />
            </div>

            <FolderPicker
              selected={selectedFolders}
              available={availableFolders}
              onAdd={(folder) => setAddedFolders((prev) => [...prev, folder])}
              onRemove={(id) => {
                const isOriginal = originalFolders.some((f) => f.id === id)
                if (isOriginal) {
                  setRemovedIds((prev) => new Set([...prev, id]))
                } else {
                  setAddedFolders((prev) => prev.filter((f) => f.id !== id))
                }
              }}
              disabled={folderPickerDisabled}
              isError={publicFoldersQuery.isError}
              onRetry={() => publicFoldersQuery.refetch()}
            />
          </form>

          <DialogFooter>
            {isEditMode && (
              <Button
                variant="destructive"
                type="button"
                onClick={() => setDeleteDialogOpen(true)}
                disabled={isPending}
                className="mr-auto"
              >
                Delete
              </Button>
            )}
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
              disabled={isPending || !name.trim()}
            >
              Save
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <AlertDialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete character?</AlertDialogTitle>
            <AlertDialogDescription>
              This action cannot be undone.
            </AlertDialogDescription>
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
    </>
  )
}
