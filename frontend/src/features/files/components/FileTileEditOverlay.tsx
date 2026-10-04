import { useEffect, useRef, useState } from 'react'
import { FileIcon, X } from 'lucide-react'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { type File } from '../api/useFiles'
import { useUpdateFile } from '../api/useUpdateFile'

interface FileTileEditOverlayProps {
  file: File
  onClose: () => void
}

export default function FileTileEditOverlay({ file, onClose }: FileTileEditOverlayProps) {
  const [name, setName] = useState(file.name ?? '')
  const [notes, setNotes] = useState(file.notes ?? '')
  const [nameError, setNameError] = useState<string | null>(null)
  const [notesError, setNotesError] = useState<string | null>(null)

  // Track last successfully saved values so blurring one field doesn't overwrite
  // the other field's in-progress draft value on the server.
  const lastSavedName = useRef(file.name ?? '')
  const lastSavedNotes = useRef(file.notes ?? '')

  const updateFile = useUpdateFile()

  useEffect(() => {
    function handleKeyDown(e: KeyboardEvent) {
      if (e.key === 'Escape') onClose()
    }
    document.addEventListener('keydown', handleKeyDown)
    return () => document.removeEventListener('keydown', handleKeyDown)
  }, [onClose])

  function handleFieldBlur(field: 'name' | 'notes') {
    const setError = field === 'name' ? setNameError : setNotesError
    const nameToSend = field === 'name' ? name : lastSavedName.current
    const notesToSend = field === 'notes' ? notes : lastSavedNotes.current

    updateFile.mutate(
      { id: file.id, name: nameToSend, notes: notesToSend },
      {
        onSuccess: () => {
          if (field === 'name') lastSavedName.current = name
          else lastSavedNotes.current = notes
          setError(null)
        },
        onError: () => setError("Couldn't save. Try again."),
      },
    )
  }

  return (
    <>
      {/* Backdrop */}
      <div className="fixed inset-0 z-40" onClick={onClose} />

      {/* Panel */}
      <div
        role="dialog"
        aria-modal="false"
        aria-labelledby="file-edit-title"
        className="fixed right-0 top-0 z-50 flex h-full w-80 flex-col gap-4 border-l bg-background p-4 shadow-xl"
      >
        <div className="flex items-center justify-between">
          <span id="file-edit-title" className="text-sm font-medium">Edit file</span>
          <Button variant="ghost" size="icon" onClick={onClose} aria-label="Close">
            <X className="h-4 w-4" />
          </Button>
        </div>

        {/* Thumbnail */}
        <div className="flex h-48 w-full items-center justify-center overflow-hidden rounded-md bg-muted">
          {file.thumbnail_url ? (
            <img src={file.thumbnail_url} alt="" className="h-full w-full object-cover" />
          ) : (
            <FileIcon className="h-8 w-8 text-muted-foreground" />
          )}
        </div>

        {/* Name */}
        <div className="flex flex-col gap-1">
          <label htmlFor="file-edit-name" className="text-sm font-medium">Name</label>
          <Input
            id="file-edit-name"
            value={name}
            placeholder="File name"
            onChange={(e) => setName(e.target.value)}
            onBlur={() => handleFieldBlur('name')}
            className={nameError ? 'border-destructive focus-visible:border-destructive' : ''}
          />
          {nameError && <p className="text-xs text-destructive">{nameError}</p>}
        </div>

        {/* Notes */}
        <div className="flex flex-col gap-1">
          <label htmlFor="file-edit-notes" className="text-sm font-medium">Notes</label>
          <Input
            id="file-edit-notes"
            value={notes}
            placeholder="Notes"
            onChange={(e) => setNotes(e.target.value)}
            onBlur={() => handleFieldBlur('notes')}
            className={notesError ? 'border-destructive focus-visible:border-destructive' : ''}
          />
          {notesError && <p className="text-xs text-destructive">{notesError}</p>}
        </div>
      </div>
    </>
  )
}
