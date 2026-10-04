import { useEffect, useRef, useState } from 'react'
import { FileIcon, Loader2, X } from 'lucide-react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog'
import { useFiles } from '../api/useFiles'
import { useInitFileUpload } from '../api/useInitFileUpload'
import { useCompleteFileUpload } from '../api/useCompleteFileUpload'
import { useUpdateFile } from '../api/useUpdateFile'

type ThumbnailGenState = 'pending' | 'done' | 'failed' | 'not_applicable'

type DraftFile =
  | { status: 'uploading'; clientId: string; fileName: string }
  | {
      status: 'uploaded'
      clientId: string
      id: string
      name: string
      notes: string
      thumbnailUrl: string | null
      thumbnailGenState: ThumbnailGenState
      nameError: string | null
      notesError: string | null
    }

interface AddFilesModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export default function AddFilesModal({ open, onOpenChange }: AddFilesModalProps) {
  const [files, setFiles] = useState<DraftFile[]>([])
  const [dragOver, setDragOver] = useState(false)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const { data: inboxFiles = [] } = useFiles()
  const initUpload = useInitFileUpload()
  const completeUpload = useCompleteFileUpload()
  const updateFile = useUpdateFile()

  // Sync thumbnailUrl and thumbnailGenState from the useFiles poll
  useEffect(() => {
    setFiles((prev) =>
      prev.map((f) => {
        if (f.status !== 'uploaded') return f
        const match = inboxFiles.find((inf) => inf.id === f.id)
        if (!match) return f
        if (match.thumbnail_url === f.thumbnailUrl && match.thumbnail_gen_state === f.thumbnailGenState) return f
        return { ...f, thumbnailUrl: match.thumbnail_url, thumbnailGenState: match.thumbnail_gen_state }
      }),
    )
  }, [inboxFiles])

  // Reset file list when modal closes
  useEffect(() => {
    if (!open) setFiles([])
  }, [open])

  async function uploadFiles(rawFiles: FileList | File[]) {
    const list = Array.from(rawFiles)
    const clientIds = list.map(() => crypto.randomUUID())

    setFiles((prev) => [
      ...prev,
      ...list.map((f, i) => ({
        status: 'uploading' as const,
        clientId: clientIds[i],
        fileName: f.name,
      })),
    ])

    await Promise.all(
      list.map(async (file, i) => {
        const clientId = clientIds[i]
        try {
          const result = await initUpload.mutateAsync({
            mimeType: file.type || 'application/octet-stream',
            name: file.name,
          })
          await fetch(result.upload_url, { method: 'PUT', body: file })
          await completeUpload.mutateAsync(result.id)
          setFiles((prev) =>
            prev.map((f) =>
              f.clientId === clientId
                ? {
                    status: 'uploaded' as const,
                    clientId,
                    id: result.id,
                    name: file.name,
                    notes: '',
                    thumbnailUrl: null,
                    thumbnailGenState: 'pending' as ThumbnailGenState,
                    nameError: null,
                    notesError: null,
                  }
                : f,
            ),
          )
        } catch {
          toast.error(`Failed to upload ${file.name}`)
          setFiles((prev) => prev.filter((f) => f.clientId !== clientId))
        }
      }),
    )
  }

  function handleDrop(e: React.DragEvent) {
    e.preventDefault()
    setDragOver(false)
    if (e.dataTransfer.files.length > 0) void uploadFiles(e.dataTransfer.files)
  }

  function handleFileInputChange(e: React.ChangeEvent<HTMLInputElement>) {
    if (e.target.files && e.target.files.length > 0) {
      void uploadFiles(e.target.files)
      e.target.value = ''
    }
  }

  function removeFile(clientId: string) {
    setFiles((prev) => prev.filter((f) => f.clientId !== clientId))
  }

  function updateLocalName(clientId: string, name: string) {
    setFiles((prev) =>
      prev.map((f) => (f.clientId === clientId && f.status === 'uploaded' ? { ...f, name } : f)),
    )
  }

  function updateLocalNotes(clientId: string, notes: string) {
    setFiles((prev) =>
      prev.map((f) => (f.clientId === clientId && f.status === 'uploaded' ? { ...f, notes } : f)),
    )
  }

  function handleFieldBlur(file: Extract<DraftFile, { status: 'uploaded' }>, field: 'name' | 'notes') {
    const errorKey = field === 'name' ? 'nameError' : 'notesError'
    updateFile.mutate(
      { id: file.id, name: file.name, notes: file.notes },
      {
        onSuccess: () =>
          setFiles((prev) =>
            prev.map((f) =>
              f.clientId === file.clientId && f.status === 'uploaded' ? { ...f, [errorKey]: null } : f,
            ),
          ),
        onError: () =>
          setFiles((prev) =>
            prev.map((f) =>
              f.clientId === file.clientId && f.status === 'uploaded'
                ? { ...f, [errorKey]: "Couldn't save. Try again." }
                : f,
            ),
          ),
      },
    )
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="flex max-h-[calc(100dvh-2rem)] max-w-md flex-col overflow-hidden">
        <DialogHeader>
          <DialogTitle>Add files to dump</DialogTitle>
        </DialogHeader>

        <div className="flex min-h-0 flex-1 flex-col gap-3 overflow-y-auto pr-1">
          {/* Drop zone */}
          <button
            type="button"
            className={`flex w-full cursor-pointer items-center justify-center rounded-lg border border-dashed border-border py-8 text-sm text-muted-foreground transition-colors hover:bg-muted/50 ${dragOver ? 'ring-2 ring-ring bg-muted/50' : ''}`}
            onClick={() => fileInputRef.current?.click()}
            onDragOver={(e) => { e.preventDefault(); setDragOver(true) }}
            onDragLeave={() => setDragOver(false)}
            onDrop={handleDrop}
          >
            {dragOver ? 'Drop to add' : 'Add files here (drag and drop, file picker)'}
          </button>
          <input
            ref={fileInputRef}
            type="file"
            multiple
            className="hidden"
            onChange={handleFileInputChange}
          />

          {/* File rows */}
          {files.length > 0 && (
            <div className="flex max-h-96 flex-col gap-2 overflow-y-auto pr-1">
              {files.map((file) => {
                if (file.status === 'uploading') {
                  return (
                    <div key={file.clientId} className="flex gap-3 rounded-md border p-2">
                      <div className="flex h-16 w-16 shrink-0 items-center justify-center rounded-md bg-muted">
                        <Loader2 className="h-5 w-5 animate-spin text-muted-foreground" />
                      </div>
                      <div className="flex flex-1 flex-col gap-1.5">
                        <Input value={file.fileName} disabled placeholder="Name" />
                        <Input value="" disabled placeholder="Notes" />
                      </div>
                    </div>
                  )
                }

                const thumbnailPending = file.thumbnailGenState === 'pending'
                return (
                  <div key={file.clientId} className="flex gap-3 rounded-md border p-2">
                    <div className="flex h-16 w-16 shrink-0 items-center justify-center overflow-hidden rounded-md bg-muted">
                      {file.thumbnailUrl ? (
                        <img src={file.thumbnailUrl} alt="" className="h-full w-full object-cover" />
                      ) : thumbnailPending ? (
                        <Loader2 className="h-5 w-5 animate-spin text-muted-foreground" />
                      ) : (
                        <FileIcon className="h-5 w-5 text-muted-foreground" />
                      )}
                    </div>
                    <div className="flex flex-1 flex-col gap-1.5">
                      <div className="flex gap-1.5">
                        <div className="flex flex-1 flex-col gap-0.5">
                          <Input
                            value={file.name}
                            placeholder="Name"
                            disabled={thumbnailPending}
                            onChange={(e) => updateLocalName(file.clientId, e.target.value)}
                            onBlur={() => handleFieldBlur(file, 'name')}
                            className={file.nameError ? 'border-destructive focus-visible:border-destructive' : ''}
                          />
                          {file.nameError && (
                            <p className="text-xs text-destructive">{file.nameError}</p>
                          )}
                        </div>
                        <button
                          type="button"
                          className="inline-flex h-9 w-9 shrink-0 items-center justify-center rounded-md border border-input hover:bg-accent"
                          onClick={() => removeFile(file.clientId)}
                          aria-label="Remove file"
                        >
                          <X className="h-4 w-4" />
                        </button>
                      </div>
                      <div className="flex flex-col gap-0.5">
                        <Input
                          value={file.notes}
                          placeholder="Notes"
                          disabled={thumbnailPending}
                          onChange={(e) => updateLocalNotes(file.clientId, e.target.value)}
                          onBlur={() => handleFieldBlur(file, 'notes')}
                          className={file.notesError ? 'border-destructive focus-visible:border-destructive' : ''}
                        />
                        {file.notesError && (
                          <p className="text-xs text-destructive">{file.notesError}</p>
                        )}
                      </div>
                    </div>
                  </div>
                )
              })}
            </div>
          )}
        </div>

        <DialogFooter>
          <Button variant="outline" type="button" onClick={() => onOpenChange(false)}>
            Close
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
