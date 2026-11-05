import { useRef, useState } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Switch } from '@/components/ui/switch'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog'
import { useUpdateFile } from '@/features/files/api/useUpdateFile'
import FileThumb from '@/features/files/components/FileThumb'
import { type ArtpieceFile } from '../api/useArtpiece'

const SAVE_ERROR = "Couldn't save. Try again."

interface ArtpieceFileDetailsModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  artpieceId: string
  files: ArtpieceFile[]
}

export default function ArtpieceFileDetailsModal({
  open,
  onOpenChange,
  artpieceId,
  files,
}: ArtpieceFileDetailsModalProps) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-xl">
        <DialogHeader>
          <DialogTitle>Edit file details</DialogTitle>
        </DialogHeader>

        {files.length === 0 ? (
          <p className="text-sm text-muted-foreground">No files attached.</p>
        ) : (
          <div className="flex max-h-96 flex-col gap-2 overflow-y-auto pr-1">
            {files.map((file) => (
              <FileDetailsRow key={file.id} file={file} artpieceId={artpieceId} />
            ))}
          </div>
        )}

        <DialogFooter>
          <Button variant="outline" type="button" onClick={() => onOpenChange(false)}>
            Close
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

function FileDetailsRow({ file, artpieceId }: { file: ArtpieceFile; artpieceId: string }) {
  const updateFile = useUpdateFile()

  const [name, setName] = useState(file.name ?? '')
  const [notes, setNotes] = useState(file.notes ?? '')
  const [isVisible, setIsVisible] = useState(file.is_visible)
  const [nameError, setNameError] = useState<string | null>(null)
  const [notesError, setNotesError] = useState<string | null>(null)
  const [visibilityError, setVisibilityError] = useState<string | null>(null)

  // Last successfully saved values, so blurring one field never sends the other
  // field's unsaved draft to the server.
  const lastSavedName = useRef(file.name ?? '')
  const lastSavedNotes = useRef(file.notes ?? '')

  function handleFieldBlur(field: 'name' | 'notes') {
    const setError = field === 'name' ? setNameError : setNotesError
    const draft = field === 'name' ? name : notes
    const saved = field === 'name' ? lastSavedName.current : lastSavedNotes.current
    // Nothing to persist unless the value differs from the last successful save
    // (a failed save leaves this stale, so blurring again retries).
    if (draft === saved) return

    updateFile.mutate(
      {
        id: file.id,
        artpieceId,
        name: field === 'name' ? draft : lastSavedName.current,
        notes: field === 'notes' ? draft : lastSavedNotes.current,
      },
      {
        onSuccess: () => {
          if (field === 'name') lastSavedName.current = draft
          else lastSavedNotes.current = draft
          setError(null)
        },
        onError: () => setError(SAVE_ERROR),
      },
    )
  }

  function handleToggle(next: boolean) {
    const previous = isVisible
    setIsVisible(next)
    updateFile.mutate(
      { id: file.id, artpieceId, isVisible: next },
      {
        onSuccess: () => {
          setVisibilityError(null)
        },
        onError: () => {
          setIsVisible(previous)
          setVisibilityError(SAVE_ERROR)
        },
      },
    )
  }

  const label = file.name ?? file.extension ?? 'file'

  return (
    <div className="flex gap-3 rounded-md border p-2">
      <div className={`flex h-16 w-16 shrink-0 overflow-hidden rounded-md ${isVisible ? '' : 'opacity-50'}`}>
        <FileThumb thumbnailUrl={file.thumbnail_url} extension={file.extension} />
      </div>
      <div className="flex flex-1 flex-col gap-1.5">
        <div className="flex flex-col gap-0.5">
          <Input
            value={name}
            placeholder="Name"
            aria-label={`Name for ${label}`}
            onChange={(e) => setName(e.target.value)}
            onBlur={() => handleFieldBlur('name')}
            className={nameError ? 'border-destructive focus-visible:border-destructive' : ''}
          />
          {nameError && <p className="text-xs text-destructive">{nameError}</p>}
        </div>
        <div className="flex flex-col gap-0.5">
          <Input
            value={notes}
            placeholder="Notes"
            aria-label={`Notes for ${label}`}
            onChange={(e) => setNotes(e.target.value)}
            onBlur={() => handleFieldBlur('notes')}
            className={notesError ? 'border-destructive focus-visible:border-destructive' : ''}
          />
          {notesError && <p className="text-xs text-destructive">{notesError}</p>}
        </div>
        <div className="flex flex-col gap-0.5">
          <div className="flex items-center gap-2">
            <Switch
              checked={isVisible}
              onCheckedChange={handleToggle}
              aria-label={`Visible: ${label}`}
            />
            <span className="text-xs text-muted-foreground">{isVisible ? 'Visible' : 'Hidden'}</span>
          </div>
          {visibilityError && <p className="text-xs text-destructive">{visibilityError}</p>}
        </div>
      </div>
    </div>
  )
}
