import { useCallback, useEffect, useRef, useState } from 'react'
import { Loader2, X } from 'lucide-react'
import { toast } from 'sonner'
import { Input } from '@/components/ui/input'
import { useInitFileUpload } from '../api/useInitFileUpload'
import { useCompleteFileUpload } from '../api/useCompleteFileUpload'
import { useUpdateFile } from '../api/useUpdateFile'
import { validateFileName } from '../lib/files'
import { uploadFile, filterBySize } from '../lib/upload'
import FileThumb from './FileThumb'

type DraftFile =
  | { status: 'uploading'; clientId: string; fileName: string }
  | { status: 'uploaded'; clientId: string; id: string; name: string; notes: string; localThumbnailUrl: string | null; extension: string | null; nameError: string | null }

interface ArtpieceFilesInputProps {
  initialFiles?: { id: string; name: string | null; notes: string | null; thumbnailUrl?: string | null; extension?: string | null }[]
  onFileIdsChange: (ids: string[]) => void
  onHasUploading: (uploading: boolean) => void
  disabled?: boolean
}

export default function ArtpieceFilesInput({
  initialFiles,
  onFileIdsChange,
  onHasUploading,
  disabled,
}: ArtpieceFilesInputProps) {
  const [files, setFiles] = useState<DraftFile[]>(() =>
    (initialFiles ?? []).map((f) => ({
      status: 'uploaded' as const,
      clientId: f.id,
      id: f.id,
      name: f.name ?? '',
      notes: f.notes ?? '',
      localThumbnailUrl: f.thumbnailUrl ?? null,
      extension: f.extension ?? null,
      nameError: null,
    })),
  )
  const [dragOver, setDragOver] = useState(false)
  const fileInputRef = useRef<HTMLInputElement>(null)
  const objectUrlsRef = useRef<string[]>([])

  const initUpload = useInitFileUpload()
  const completeUpload = useCompleteFileUpload()
  const updateFile = useUpdateFile()

  useEffect(() => {
    const urls = objectUrlsRef.current
    return () => {
      urls.forEach((url) => URL.revokeObjectURL(url))
    }
  }, [])

  // Notify parent of file IDs and uploading state
  useEffect(() => {
    const ids = files.filter((f): f is Extract<DraftFile, { status: 'uploaded' }> => f.status === 'uploaded').map((f) => f.id)
    onFileIdsChange(ids)
    onHasUploading(files.some((f) => f.status === 'uploading'))
  }, [files, onFileIdsChange, onHasUploading])

  const uploadFiles = useCallback(
    async (rawFiles: FileList | File[]) => {
      const all = Array.from(rawFiles)
      const { accepted: list, rejected } = filterBySize(all)
      for (const { file, limitBytes, categoryName } of rejected) {
        toast.error(`${file.name} exceeds the ${Math.round(limitBytes / (1024 * 1024))} MB limit for ${categoryName} files`)
      }
      if (list.length === 0) return
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
            const { fileId, name: strippedName, extension, thumbnailBlob } = await uploadFile(file, initUpload.mutateAsync, completeUpload.mutateAsync)
            const localThumbnailUrl = thumbnailBlob ? URL.createObjectURL(thumbnailBlob) : null
            if (localThumbnailUrl) objectUrlsRef.current.push(localThumbnailUrl)
            setFiles((prev) =>
              prev.map((f) =>
                f.clientId === clientId
                  ? {
                      status: 'uploaded' as const,
                      clientId,
                      id: fileId,
                      name: strippedName,
                      notes: '',
                      localThumbnailUrl,
                      extension,
                      nameError: null,
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
    },
    [initUpload, completeUpload],
  )

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

  function handleNameBlur(file: Extract<DraftFile, { status: 'uploaded' }>) {
    const validationError = validateFileName(file.name)
    if (validationError) {
      setFiles((prev) =>
        prev.map((f) =>
          f.clientId === file.clientId && f.status === 'uploaded' ? { ...f, nameError: validationError } : f,
        ),
      )
      return
    }
    setFiles((prev) =>
      prev.map((f) =>
        f.clientId === file.clientId && f.status === 'uploaded' ? { ...f, nameError: null } : f,
      ),
    )
    updateFile.mutate({ id: file.id, name: file.name, notes: file.notes })
  }

  function handleNotesBlur(file: Extract<DraftFile, { status: 'uploaded' }>) {
    updateFile.mutate({ id: file.id, name: file.name, notes: file.notes })
  }

  return (
    <div className="flex flex-col gap-2">
      {/* Sticky drop zone */}
      <button
        type="button"
        className={`flex w-full cursor-pointer items-center justify-center rounded-lg border border-dashed border-border py-3 text-sm text-muted-foreground transition-colors hover:bg-muted/50 ${dragOver ? 'ring-2 ring-ring bg-muted/50' : ''}`}
        onClick={() => fileInputRef.current?.click()}
        onDragOver={(e) => { e.preventDefault(); setDragOver(true) }}
        onDragLeave={() => setDragOver(false)}
        onDrop={handleDrop}
        disabled={disabled}
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
        <div className="flex max-h-64 flex-col gap-2 overflow-y-auto pr-1">
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

            return (
              <div key={file.clientId} className="flex gap-3 rounded-md border p-2">
                <div className="flex h-16 w-16 shrink-0 overflow-hidden rounded-md">
                  <FileThumb
                    thumbnailUrl={file.localThumbnailUrl}
                    extension={file.extension}
                  />
                </div>
                <div className="flex flex-1 flex-col gap-1.5">
                  <div className="flex gap-1.5">
                    <div className="flex flex-1 flex-col gap-0.5">
                      <Input
                        value={file.name}
                        placeholder="Name"
                        disabled={disabled}
                        onChange={(e) => updateLocalName(file.clientId, e.target.value)}
                        onBlur={() => handleNameBlur(file)}
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
                      disabled={disabled}
                      aria-label="Remove file"
                    >
                      <X className="h-4 w-4" />
                    </button>
                  </div>
                  <Input
                    value={file.notes}
                    placeholder="Notes"
                    disabled={disabled}
                    onChange={(e) => updateLocalNotes(file.clientId, e.target.value)}
                    onBlur={() => handleNotesBlur(file)}
                  />
                </div>
              </div>
            )
          })}
        </div>
      )}
    </div>
  )
}
