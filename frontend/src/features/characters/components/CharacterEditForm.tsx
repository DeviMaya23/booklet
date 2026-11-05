import { useEffect, useRef, useState } from 'react'
import { useBlocker, useNavigate } from 'react-router-dom'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import { ArrowLeft, Image, Trash2, X } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
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
import { cn } from '@/lib/utils'
import { useVerses, type Verse } from '@/features/verses/api/useVerses'
import VerseCombobox from '@/features/verses/components/VerseCombobox'
import { CHARACTERS_QUERY_KEY, type Character } from '../api/useCharacters'
import { useUpdateCharacter } from '../api/useUpdateCharacter'
import { useDeleteCharacter } from '../api/useDeleteCharacter'
import { useInitAvatarUpload } from '../api/useInitAvatarUpload'
import { useCompleteAvatarUpload } from '../api/useCompleteAvatarUpload'
import { usePublicFolders } from '../api/usePublicFolders'
import { useFolderDiff } from '../api/useFolderDiff'
import { useAvatarPicker } from '../hooks/useAvatarPicker'
import { uploadCharacterAvatar } from '../lib/uploadCharacterAvatar'
import { characterEditSummary } from '../lib/characterEditSummary'
import AvatarCropDialog from './AvatarCropDialog'
import { FolderChipsEdit } from './CharacterFolderChips'
import { ChooseFoldersDrawer } from './ChooseFoldersDrawer'

interface CharacterEditFormProps {
  character: Character
  /** Leave edit mode and return to the view. */
  onClose: () => void
}

export default function CharacterEditForm({ character, onClose }: CharacterEditFormProps) {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const { getToken } = useKindeAuth()

  const [name, setName] = useState(character.name)
  const [notes, setNotes] = useState(character.notes ?? '')
  // undefined = untouched, so the form is not dirty while verses are still loading.
  const [verseOverride, setVerseOverride] = useState<Verse | null | undefined>(undefined)
  const [drawerOpen, setDrawerOpen] = useState(false)
  const [discardOpen, setDiscardOpen] = useState(false)
  const [deleteOpen, setDeleteOpen] = useState(false)
  const [isSaving, setIsSaving] = useState(false)
  const leavingRef = useRef(false)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const versesQuery = useVerses()
  const publicFoldersQuery = usePublicFolders()
  const updateMutation = useUpdateCharacter()
  const deleteMutation = useDeleteCharacter()
  const initUploadMutation = useInitAvatarUpload()
  const completeUploadMutation = useCompleteAvatarUpload()
  const avatar = useAvatarPicker({ currentUrl: character.avatar_url })
  const { handlePaste } = avatar

  const {
    selectedFolders,
    selectedIds,
    invalidFolders,
    addFolder,
    removeFolder,
    dismissInvalid,
  } = useFolderDiff({ originalFolders: character.folders, publicFolders: publicFoldersQuery.data })

  const effectiveVerse =
    verseOverride !== undefined
      ? verseOverride
      : ((versesQuery.data ?? []).find((v) => v.id === character.verse_id) ?? null)
  const effectiveVerseId = verseOverride !== undefined ? (verseOverride?.id ?? null) : character.verse_id

  const summary = characterEditSummary({
    character,
    name,
    notes,
    verseId: effectiveVerseId,
    avatar: { hasNewFile: avatar.localFile !== null, cleared: avatar.avatarCleared },
    selectedFolderIds: selectedFolders.map((f) => f.id),
  })

  const blocker = useBlocker(() => !leavingRef.current && summary !== null)

  useEffect(() => {
    if (blocker.state === 'blocked') setDiscardOpen(true)
  }, [blocker.state])

  // Nothing has focus by default on a page, so listen for pastes on the document.
  useEffect(() => {
    document.addEventListener('paste', handlePaste)
    return () => document.removeEventListener('paste', handlePaste)
  }, [handlePaste])

  function handleBack() {
    if (summary !== null) {
      setDiscardOpen(true)
      return
    }
    onClose()
  }

  function handleDiscardConfirm() {
    setDiscardOpen(false)
    if (blocker.state === 'blocked') {
      blocker.proceed()
    } else {
      onClose()
    }
  }

  function handleDiscardCancel() {
    setDiscardOpen(false)
    if (blocker.state === 'blocked') blocker.reset()
  }

  async function handleSave() {
    if (!name.trim() || isSaving) return
    setIsSaving(true)
    try {
      if (avatar.localFile) {
        try {
          await uploadCharacterAvatar(
            character.id,
            avatar.localFile,
            initUploadMutation.mutateAsync,
            completeUploadMutation.mutateAsync,
          )
        } catch {
          toast.error('Failed to upload avatar')
          return
        }
      } else if (avatar.avatarCleared) {
        const res = await apiFetch(`/characters/${character.id}/avatar`, getToken, { method: 'DELETE' })
        if (!res.ok) {
          toast.error('Failed to remove avatar')
          return
        }
      }
      if (avatar.isChanged) {
        // The avatar is saved from here on; show it even if the update below fails.
        avatar.reset()
        queryClient.invalidateQueries({ queryKey: CHARACTERS_QUERY_KEY })
      }

      await updateMutation.mutateAsync({
        id: character.id,
        name: name.trim(),
        notes: notes.trim() || null,
        folder_ids: selectedFolders.map((f) => f.id),
        verse_id: effectiveVerseId,
      })

      toast.success('Character updated')
      leavingRef.current = true
      onClose()
    } catch {
      toast.error('Failed to update character')
    } finally {
      setIsSaving(false)
    }
  }

  function handleDeleteConfirm() {
    deleteMutation.mutate(character.id, {
      onSuccess: () => {
        toast.success('Character deleted')
        leavingRef.current = true
        navigate('/app/characters')
      },
      onError: () => {
        toast.error('Failed to delete character')
        setDeleteOpen(false)
      },
    })
  }

  const busy = isSaving || deleteMutation.isPending

  return (
    <div className="flex flex-col gap-6">
      <button
        type="button"
        onClick={handleBack}
        className="flex w-fit items-center gap-1 text-sm text-muted-foreground hover:text-foreground"
      >
        <ArrowLeft className="size-4" /> {character.name}
      </button>

      <div className="flex flex-col gap-1.5">
        <div className="flex items-center gap-2">
          <Input
            aria-label="Name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="Character name"
            className="h-auto py-2 text-3xl font-bold"
            disabled={busy}
          />
          <Button variant="outline" onClick={onClose} disabled={busy}>
            Cancel
          </Button>
          <Button onClick={handleSave} disabled={busy || !name.trim()}>
            Save
          </Button>
        </div>
        {summary && (
          <p className="flex items-center gap-1.5 text-sm text-muted-foreground">
            <span className="size-2 shrink-0 rounded-full bg-amber-500" />
            {summary}
          </p>
        )}
      </div>

      <div className="grid grid-cols-1 gap-6 sm:grid-cols-[1fr_208px]">
        <div className="flex flex-col gap-4">
          <div className="flex flex-col gap-1.5">
            <label className="text-sm font-medium" htmlFor="character-verse">Verse</label>
            <VerseCombobox
              value={effectiveVerse}
              onChange={setVerseOverride}
              disabled={busy}
              inputId="character-verse"
            />
          </div>
          <div className="flex flex-col gap-1.5">
            <label className="text-sm font-medium" htmlFor="character-notes">Notes</label>
            <Textarea
              id="character-notes"
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
              placeholder="Notes about this character..."
              rows={5}
              disabled={busy}
            />
          </div>
        </div>

        <div className="relative sm:order-none">
          <button
            type="button"
            aria-label="Change avatar"
            className={cn(
              'relative flex aspect-square w-full cursor-pointer items-center justify-center overflow-hidden rounded-2xl bg-muted transition-colors',
              avatar.dragActive && 'ring-2 ring-ring',
            )}
            onClick={() => fileInputRef.current?.click()}
            {...avatar.dragHandlers}
            disabled={busy}
          >
            {avatar.displayUrl ? (
              <>
                <img src={avatar.displayUrl} alt="Avatar" className="size-full object-cover" />
                <span className="absolute inset-x-2 bottom-2 rounded-lg bg-black/55 py-1 text-center text-sm font-medium text-white">
                  Change
                </span>
              </>
            ) : (
              <span className="flex flex-col items-center gap-1 text-sm text-muted-foreground">
                <Image className="size-8" />
                Click, drop or paste an image
              </span>
            )}
          </button>
          {avatar.canClear && (
            <button
              type="button"
              aria-label="Remove avatar"
              className="absolute right-2 top-2 z-10 inline-flex size-7 items-center justify-center rounded-full bg-black/55 text-white hover:bg-black/70"
              onClick={avatar.clear}
              disabled={busy}
            >
              <X className="size-4" />
            </button>
          )}
          <input
            ref={fileInputRef}
            type="file"
            accept="image/jpeg,image/png"
            className="hidden"
            onChange={avatar.handleFileChange}
          />
        </div>
      </div>

      <FolderChipsEdit
        selectedFolders={selectedFolders}
        invalidFolders={invalidFolders}
        onRemove={removeFolder}
        onDismissInvalid={dismissInvalid}
        onChoose={() => setDrawerOpen(true)}
        disabled={busy}
      />

      <div className="border-t pt-4">
        <Button
          variant="ghost"
          className="gap-1.5 text-destructive hover:text-destructive"
          onClick={() => setDeleteOpen(true)}
          disabled={busy}
        >
          <Trash2 className="size-4" />
          Delete character
        </Button>
      </div>

      <ChooseFoldersDrawer
        open={drawerOpen}
        onOpenChange={setDrawerOpen}
        selectedIds={selectedIds}
        onAdd={addFolder}
        onRemove={removeFolder}
        folders={publicFoldersQuery.data ?? null}
        isLoading={publicFoldersQuery.isLoading}
        isError={publicFoldersQuery.isError && !publicFoldersQuery.data}
        onRefresh={() => publicFoldersQuery.refetch()}
        updatedAt={publicFoldersQuery.dataUpdatedAt}
      />

      {avatar.cropSourceUrl && (
        <AvatarCropDialog
          open
          sourceUrl={avatar.cropSourceUrl}
          onCrop={avatar.confirmCrop}
          onCancel={avatar.cancelCrop}
        />
      )}

      <AlertDialog open={deleteOpen} onOpenChange={setDeleteOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete character?</AlertDialogTitle>
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

      <AlertDialog
        open={discardOpen}
        onOpenChange={(open) => {
          if (!open) handleDiscardCancel()
        }}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Discard changes?</AlertDialogTitle>
            <AlertDialogDescription>Any unsaved changes will be lost.</AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Keep editing</AlertDialogCancel>
            <AlertDialogAction onClick={handleDiscardConfirm}>Discard</AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  )
}
