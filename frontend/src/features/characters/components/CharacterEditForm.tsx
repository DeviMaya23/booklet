import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { ArrowLeft, Trash2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { useVerses, type Verse } from '@/features/verses/api/useVerses'
import VerseCombobox from '@/features/verses/components/VerseCombobox'
import { type Character } from '../api/useCharacters'
import { usePublicFolders } from '../api/usePublicFolders'
import { useFolderDiff } from '../api/useFolderDiff'
import { useAvatarPicker } from '../hooks/useAvatarPicker'
import { useCharacterSave } from '../hooks/useCharacterSave'
import { useCharacterDelete } from '../hooks/useCharacterDelete'
import { useUnsavedChangesGuard } from '../hooks/useUnsavedChangesGuard'
import { characterEditSummary } from '../lib/characterEditSummary'
import AvatarCropDialog from './AvatarCropDialog'
import CharacterAvatarEditor from './CharacterAvatarEditor'
import DeleteCharacterDialog from './DeleteCharacterDialog'
import DiscardChangesDialog from './DiscardChangesDialog'
import { FolderChipsEdit } from './CharacterFolderChips'
import { ChooseFoldersDrawer } from './ChooseFoldersDrawer'

interface CharacterEditFormProps {
  character: Character
  /** Leave edit mode and return to the view. */
  onClose: () => void
}

export default function CharacterEditForm({ character, onClose }: CharacterEditFormProps) {
  const navigate = useNavigate()

  const [name, setName] = useState(character.name)
  const [notes, setNotes] = useState(character.notes ?? '')
  // undefined = untouched, so the form is not dirty while verses are still loading.
  const [verseOverride, setVerseOverride] = useState<Verse | null | undefined>(undefined)
  const [drawerOpen, setDrawerOpen] = useState(false)

  const versesQuery = useVerses()
  const publicFoldersQuery = usePublicFolders()
  const avatar = useAvatarPicker({ currentUrl: character.avatar_url, listenForPaste: true })
  const { selectedFolders, selectedIds, invalidFolders, addFolder, removeFolder, dismissInvalid } =
    useFolderDiff({ originalFolders: character.folders, publicFolders: publicFoldersQuery.data })

  const effectiveVerse =
    verseOverride !== undefined
      ? verseOverride
      : ((versesQuery.data ?? []).find((v) => v.id === character.verse_id) ?? null)
  const effectiveVerseId = verseOverride !== undefined ? (verseOverride?.id ?? null) : character.verse_id

  // Bookleaf is the source of truth for which folders exist, and the save drops stale
  // ones on its own, so their absence must not make the form look edited.
  const liveFolderIds = publicFoldersQuery.data ? new Set(publicFoldersQuery.data.map((f) => f.id)) : null
  const staleFolderIds = liveFolderIds
    ? character.folders.filter((f) => !liveFolderIds.has(f.id)).map((f) => f.id)
    : []

  const summary = characterEditSummary({
    character,
    name,
    notes,
    verseId: effectiveVerseId,
    avatar: { hasNewFile: avatar.localFile !== null, cleared: avatar.avatarCleared },
    selectedFolderIds: selectedFolders.map((f) => f.id),
    staleFolderIds,
  })

  const guard = useUnsavedChangesGuard({ isDirty: summary !== null, onLeave: onClose })
  const { isSaving, save } = useCharacterSave({ character, avatar, onSaved: onClose })
  const remove = useCharacterDelete({
    characterId: character.id,
    onDeleted: () => {
      guard.allowLeaving()
      navigate('/app/characters')
    },
  })

  const busy = isSaving || remove.isPending

  return (
    <div className="flex flex-col gap-6">
      <button
        type="button"
        onClick={guard.requestLeave}
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
          <Button
            onClick={() =>
              save({
                name,
                notes,
                verseId: effectiveVerseId,
                folderIds: selectedFolders.map((f) => f.id),
              })
            }
            disabled={busy || !name.trim()}
          >
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
        <CharacterAvatarEditor avatar={avatar} disabled={busy} />
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
          onClick={remove.askToDelete}
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

      <DeleteCharacterDialog
        open={remove.confirmOpen}
        pending={remove.isPending}
        onCancel={remove.cancel}
        onConfirm={remove.confirm}
      />
      <DiscardChangesDialog
        open={guard.promptOpen}
        onKeepEditing={guard.keepEditing}
        onDiscard={guard.discard}
      />
    </div>
  )
}
