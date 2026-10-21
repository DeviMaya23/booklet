import { useEffect, useRef, useState } from 'react'
import { toast } from 'sonner'
import { Image, X } from 'lucide-react'
import AvatarCropDialog from './AvatarCropDialog'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { useQueryClient } from '@tanstack/react-query'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
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
import { cn } from '@/lib/utils'
import { apiFetch } from '@/lib/api'
import { type Character, CHARACTERS_QUERY_KEY } from '../api/useCharacters'
import { useCreateCharacter } from '../api/useCreateCharacter'
import { useUpdateCharacter } from '../api/useUpdateCharacter'
import { useDeleteCharacter } from '../api/useDeleteCharacter'
import { useInitAvatarUpload } from '../api/useInitAvatarUpload'
import { useCompleteAvatarUpload } from '../api/useCompleteAvatarUpload'
import { usePublicFolders } from '../api/usePublicFolders'
import { useFolderDiff } from '../api/useFolderDiff'
import { FolderList } from './FolderList'
import { ChooseFoldersPanel } from './ChooseFoldersPanel'

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
  const [discardOpen, setDiscardOpen] = useState(false)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [cropSourceUrl, setCropSourceUrl] = useState<string | null>(null)
  const [dragActive, setDragActive] = useState(false)
  const [panelOpen, setPanelOpen] = useState(false)

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

  const originalFolders: { id: string; name: string }[] = character?.folders ?? []

  const {
    selectedFolders,
    selectedIds,
    invalidFolders,
    addFolder,
    removeFolder,
    dismissInvalid,
    reset: resetFolderDiff,
  } = useFolderDiff({ originalFolders, publicFolders: publicFoldersQuery.data })

  useEffect(() => {
    return () => {
      if (localPreviewUrl) URL.revokeObjectURL(localPreviewUrl)
    }
  }, [localPreviewUrl])

  const displayUrl = localPreviewUrl ?? (avatarCleared ? null : character?.avatar_url ?? null)
  const showClearButton = isEditMode && !avatarCleared && (localPreviewUrl !== null || character?.avatar_url != null)

  function resetForm() {
    setName(character?.name ?? '')
    setNotes(character?.notes ?? '')
    setLocalFile(null)
    setLocalPreviewUrl(null)
    setAvatarCleared(false)
    setDeleteDialogOpen(false)
    resetFolderDiff()
    setPanelOpen(false)
    if (cropSourceUrl) {
      URL.revokeObjectURL(cropSourceUrl)
      setCropSourceUrl(null)
    }
    setDragActive(false)
  }

  function isDirty(): boolean {
    if (isEditMode) {
      const originalFolderIds = originalFolders.map((f) => f.id)
      const currentFolderIds = selectedFolders.map((f) => f.id)
      return (
        name !== (character?.name ?? '') ||
        notes !== (character?.notes ?? '') ||
        localFile !== null ||
        avatarCleared ||
        currentFolderIds.length !== originalFolderIds.length ||
        currentFolderIds.some((id) => !originalFolderIds.includes(id))
      )
    }
    return (
      name.trim() !== '' ||
      notes.trim() !== '' ||
      localFile !== null ||
      selectedFolders.length > 0
    )
  }

  function handleOpenChange(nextOpen: boolean) {
    if (!nextOpen) resetForm()
    onOpenChange(nextOpen)
  }

  function stageFileForCrop(file: File) {
    if (cropSourceUrl) URL.revokeObjectURL(cropSourceUrl)
    setCropSourceUrl(URL.createObjectURL(file))
  }

  function handleFileChange(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0]
    if (!file) return
    stageFileForCrop(file)
    e.target.value = ''
  }

  function handleDrop(e: React.DragEvent) {
    e.preventDefault()
    setDragActive(false)
    const file = e.dataTransfer.files?.[0]
    if (file) stageFileForCrop(file)
  }

  function handlePaste(e: React.ClipboardEvent) {
    if (e.target instanceof HTMLInputElement || e.target instanceof HTMLTextAreaElement) return
    const item = Array.from(e.clipboardData?.items ?? []).find((i) => i.type.startsWith('image/'))
    const file = item?.getAsFile()
    if (file) stageFileForCrop(file)
  }

  function handleCropConfirm(croppedFile: File) {
    if (localPreviewUrl) URL.revokeObjectURL(localPreviewUrl)
    setLocalFile(croppedFile)
    setLocalPreviewUrl(URL.createObjectURL(croppedFile))
    setAvatarCleared(false)
    if (cropSourceUrl) URL.revokeObjectURL(cropSourceUrl)
    setCropSourceUrl(null)
  }

  function handleCropCancel() {
    if (cropSourceUrl) URL.revokeObjectURL(cropSourceUrl)
    setCropSourceUrl(null)
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
    handleOpenChange(false)
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
    handleOpenChange(false)
  }

  function handleDeleteConfirm() {
    if (!character) return
    deleteMutation.mutate(character.id, {
      onSuccess: () => {
        toast.success('Character deleted')
        setDeleteDialogOpen(false)
        handleOpenChange(false)
      },
      onError: () => {
        toast.error('Failed to delete character')
        setDeleteDialogOpen(false)
      },
    })
  }


  return (
    <>
      <Dialog
        open={open}
        onOpenChange={(nextOpen, eventDetails) => {
          if (!nextOpen && (eventDetails.reason === 'escape-key' || eventDetails.reason === 'outside-press') && isDirty()) {
            setDiscardOpen(true)
            return
          }
          handleOpenChange(nextOpen)
        }}
      >
        <DialogContent
          className={cn(
            'overflow-hidden transition-[max-width] duration-200',
            panelOpen ? 'max-w-3xl' : 'max-w-md',
          )}
          onPaste={handlePaste}
        >
          <div className="flex min-h-0">
            {/* Left column: the form */}
            <div className={cn('flex flex-col gap-4', panelOpen ? 'w-[400px] shrink-0 pr-6' : 'w-full')}>
              <DialogHeader>
                <DialogTitle>{isEditMode ? 'Edit character' : 'New character'}</DialogTitle>
              </DialogHeader>

              <form onSubmit={handleSubmit} className="flex flex-col gap-4">
                {/* Avatar + Name row */}
                <div className="flex items-start gap-4">
                  <div className="relative shrink-0">
                    <button
                      type="button"
                      className={cn(
                        'relative flex size-20 cursor-pointer items-center justify-center overflow-hidden rounded-lg bg-muted transition-colors',
                        dragActive && 'ring-2 ring-ring',
                      )}
                      onClick={() => fileInputRef.current?.click()}
                      onDragOver={(e) => { e.preventDefault(); setDragActive(true) }}
                      onDragLeave={() => setDragActive(false)}
                      onDrop={handleDrop}
                      disabled={isPending}
                    >
                      {displayUrl ? (
                        <>
                          <img
                            src={displayUrl}
                            alt="Avatar preview"
                            className="size-full object-cover"
                          />
                          <span className="absolute inset-x-0 bottom-0 bg-black/50 py-0.5 text-center text-xs font-medium text-white">
                            Change
                          </span>
                        </>
                      ) : (
                        <Image className="size-8 text-muted-foreground" />
                      )}
                    </button>
                    {showClearButton && (
                      <button
                        type="button"
                        className="absolute -right-2 -top-2 z-10 inline-flex size-5 items-center justify-center rounded-full bg-foreground text-background hover:opacity-80"
                        onClick={handleClearAvatar}
                        disabled={isPending}
                        aria-label="Remove avatar"
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

                  <div className="flex flex-1 flex-col gap-1.5">
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
                </div>

                <div className="flex flex-col gap-1.5">
                  <label className="text-sm font-medium" htmlFor="character-notes">
                    Notes
                  </label>
                  <Textarea
                    id="character-notes"
                    value={notes}
                    onChange={(e) => setNotes(e.target.value)}
                    placeholder="Notes about this character..."
                    rows={4}
                    disabled={isPending}
                  />
                </div>

                <FolderList
                  selectedFolders={selectedFolders}
                  invalidFolders={invalidFolders}
                  onRemove={removeFolder}
                  onDismissInvalid={dismissInvalid}
                  onOpenPanel={() => setPanelOpen(true)}
                  disabled={isPending}
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
            </div>

            {/* Right column: choose folders panel */}
            {panelOpen && (
              <div className="min-w-0 flex-1 border-l pl-6">
                <ChooseFoldersPanel
                  selectedIds={selectedIds}
                  onAdd={addFolder}
                  onRemove={removeFolder}
                  onClose={() => setPanelOpen(false)}
                  folders={publicFoldersQuery.data ?? null}
                  isLoading={publicFoldersQuery.isLoading}
                  isError={publicFoldersQuery.isError && !publicFoldersQuery.data}
                  onRefresh={() => publicFoldersQuery.refetch()}
                />
              </div>
            )}
          </div>
        </DialogContent>
      </Dialog>

      {cropSourceUrl && (
        <AvatarCropDialog
          open={!!cropSourceUrl}
          sourceUrl={cropSourceUrl}
          onCrop={handleCropConfirm}
          onCancel={handleCropCancel}
        />
      )}

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

      <AlertDialog open={discardOpen} onOpenChange={setDiscardOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Discard changes?</AlertDialogTitle>
            <AlertDialogDescription>Any unsaved changes will be lost.</AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Keep editing</AlertDialogCancel>
            <AlertDialogAction onClick={() => { setDiscardOpen(false); handleOpenChange(false) }}>
              Discard
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  )
}
