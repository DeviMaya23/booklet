import { useRef, useState } from 'react'
import { toast } from 'sonner'
import { Image } from 'lucide-react'
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
import { CHARACTERS_QUERY_KEY } from '../api/useCharacters'
import { useCreateCharacter } from '../api/useCreateCharacter'
import { useInitAvatarUpload } from '../api/useInitAvatarUpload'
import { useCompleteAvatarUpload } from '../api/useCompleteAvatarUpload'
import { usePublicFolders } from '../api/usePublicFolders'
import { useFolderDiff } from '../api/useFolderDiff'
import { useAvatarPicker } from '../hooks/useAvatarPicker'
import { uploadCharacterAvatar } from '../lib/uploadCharacterAvatar'
import { FolderList } from './FolderList'
import { ChooseFoldersPanel } from './ChooseFoldersPanel'
import { type Verse } from '@/features/verses/api/useVerses'
import VerseCombobox from '@/features/verses/components/VerseCombobox'

interface CharacterFormModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

const NO_FOLDERS: { id: string; name: string }[] = []

export default function CharacterFormModal({ open, onOpenChange }: CharacterFormModalProps) {
  const [name, setName] = useState('')
  const [notes, setNotes] = useState('')
  const [verse, setVerse] = useState<Verse | null>(null)
  const [discardOpen, setDiscardOpen] = useState(false)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [panelOpen, setPanelOpen] = useState(false)

  const fileInputRef = useRef<HTMLInputElement>(null)
  const { getToken } = useKindeAuth()
  const queryClient = useQueryClient()

  const createMutation = useCreateCharacter()
  const initUploadMutation = useInitAvatarUpload()
  const completeUploadMutation = useCompleteAvatarUpload()
  const publicFoldersQuery = usePublicFolders()
  const avatar = useAvatarPicker()

  const {
    selectedFolders,
    selectedIds,
    invalidFolders,
    addFolder,
    removeFolder,
    dismissInvalid,
    reset: resetFolderDiff,
  } = useFolderDiff({ originalFolders: NO_FOLDERS, publicFolders: publicFoldersQuery.data })

  function resetForm() {
    setName('')
    setNotes('')
    setVerse(null)
    avatar.reset()
    resetFolderDiff()
    setPanelOpen(false)
  }

  function isDirty(): boolean {
    return (
      name.trim() !== '' ||
      notes.trim() !== '' ||
      verse !== null ||
      avatar.localFile !== null ||
      selectedFolders.length > 0
    )
  }

  function handleOpenChange(nextOpen: boolean) {
    if (!nextOpen) resetForm()
    onOpenChange(nextOpen)
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!name.trim()) return

    setIsSubmitting(true)
    try {
      await handleCreateSubmit()
    } catch {
      toast.error('Failed to create character')
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
      verse_id: verse?.id ?? null,
    })

    if (avatar.localFile) {
      try {
        await uploadCharacterAvatar(
          created.id,
          avatar.localFile,
          initUploadMutation.mutateAsync,
          completeUploadMutation.mutateAsync,
        )
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
          onPaste={avatar.handlePaste}
        >
          <div className="flex min-h-0">
            {/* Left column: the form */}
            <div className={cn('flex flex-col gap-4', panelOpen ? 'w-[400px] shrink-0 pr-6' : 'w-full')}>
              <DialogHeader>
                <DialogTitle>New character</DialogTitle>
              </DialogHeader>

              <form onSubmit={handleSubmit} className="flex flex-col gap-4">
                {/* Avatar + Name row */}
                <div className="flex items-start gap-4">
                  <div className="relative shrink-0">
                    <button
                      type="button"
                      className={cn(
                        'relative flex size-20 cursor-pointer items-center justify-center overflow-hidden rounded-lg bg-muted transition-colors',
                        avatar.dragActive && 'ring-2 ring-ring',
                      )}
                      onClick={() => fileInputRef.current?.click()}
                      {...avatar.dragHandlers}
                      disabled={isSubmitting}
                    >
                      {avatar.displayUrl ? (
                        <>
                          <img
                            src={avatar.displayUrl}
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
                    <input
                      ref={fileInputRef}
                      type="file"
                      accept="image/jpeg,image/png"
                      className="hidden"
                      onChange={avatar.handleFileChange}
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
                      disabled={isSubmitting}
                    />
                  </div>
                </div>

                <div className="flex flex-col gap-1.5">
                  <label className="text-sm font-medium" htmlFor="character-verse">Verse</label>
                  <VerseCombobox
                    value={verse}
                    onChange={setVerse}
                    disabled={isSubmitting}
                    inputId="character-verse"
                  />
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
                    disabled={isSubmitting}
                  />
                </div>

                <FolderList
                  selectedFolders={selectedFolders}
                  invalidFolders={invalidFolders}
                  onRemove={removeFolder}
                  onDismissInvalid={dismissInvalid}
                  onOpenPanel={() => setPanelOpen(true)}
                  disabled={isSubmitting}
                />
              </form>

              <DialogFooter>
                <Button
                  variant="outline"
                  type="button"
                  onClick={() => handleOpenChange(false)}
                  disabled={isSubmitting}
                >
                  Cancel
                </Button>
                <Button
                  type="submit"
                  onClick={handleSubmit}
                  disabled={isSubmitting || !name.trim()}
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
                  updatedAt={publicFoldersQuery.dataUpdatedAt}
                />
              </div>
            )}
          </div>
        </DialogContent>
      </Dialog>

      {avatar.cropSourceUrl && (
        <AvatarCropDialog
          open={!!avatar.cropSourceUrl}
          sourceUrl={avatar.cropSourceUrl}
          onCrop={avatar.confirmCrop}
          onCancel={avatar.cancelCrop}
        />
      )}

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
