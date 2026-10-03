import { useCallback, useRef, useState } from 'react'
import { Loader2 } from 'lucide-react'
import { toast } from 'sonner'
import { type File as InboxFile } from '../api/useFiles'
import { useInitFileUpload } from '../api/useInitFileUpload'
import { useCompleteFileUpload } from '../api/useCompleteFileUpload'
import FileTile from './FileTile'
import FileContextMenu from './FileContextMenu'
import FileDeleteDialog from './FileDeleteDialog'
import { useDeleteFile } from '../api/useDeleteFile'
import { useBulkDeleteFiles } from '../api/useBulkDeleteFiles'

interface PlaceholderTile {
  clientId: string
}

interface FileInboxGridProps {
  files: InboxFile[]
  selection: Set<string>
  onSelectionChange: (next: Set<string>) => void
  onNewArtpiece?: (fileIds: string[]) => void
  onAddToExisting?: (fileIds: string[]) => void
}

export default function FileInboxGrid({ files, selection, onSelectionChange, onNewArtpiece, onAddToExisting }: FileInboxGridProps) {
  const [placeholders, setPlaceholders] = useState<PlaceholderTile[]>([])
  const [contextMenu, setContextMenu] = useState<{ x: number; y: number; ids: string[] } | null>(null)
  const [deleteTarget, setDeleteTarget] = useState<string[] | null>(null)
  const [dragOver, setDragOver] = useState(false)
  const lastClickedId = useRef<string | null>(null)

  const initUpload = useInitFileUpload()
  const completeUpload = useCompleteFileUpload()
  const deleteMutation = useDeleteFile()
  const bulkDeleteMutation = useBulkDeleteFiles()

  const allIds = files.map((f) => f.id)

  function handleTileClick(fileId: string, e: React.MouseEvent) {
    e.preventDefault()
    const next = new Set(selection)

    if (e.shiftKey && lastClickedId.current) {
      const fromIdx = allIds.indexOf(lastClickedId.current)
      const toIdx = allIds.indexOf(fileId)
      const [lo, hi] = fromIdx < toIdx ? [fromIdx, toIdx] : [toIdx, fromIdx]
      for (let i = lo; i <= hi; i++) next.add(allIds[i])
    } else if (e.ctrlKey || e.metaKey) {
      if (next.has(fileId)) { next.delete(fileId) } else { next.add(fileId) }
    } else {
      next.clear()
      next.add(fileId)
    }

    lastClickedId.current = fileId
    onSelectionChange(next)
  }

  function handleTileContextMenu(fileId: string, e: React.MouseEvent) {
    e.preventDefault()
    let ids: string[]
    if (selection.has(fileId)) {
      ids = Array.from(selection)
    } else {
      const next = new Set([fileId])
      onSelectionChange(next)
      ids = [fileId]
    }
    setContextMenu({ x: e.clientX, y: e.clientY, ids })
  }

  const uploadFiles = useCallback(async (rawFiles: FileList | globalThis.File[]) => {
    const list = Array.from(rawFiles) as globalThis.File[]
    const clientIds = list.map(() => crypto.randomUUID())
    setPlaceholders((prev) => [...prev, ...clientIds.map((clientId) => ({ clientId }))])

    await Promise.all(
      list.map(async (file, i) => {
        const clientId = clientIds[i]
        try {
          const result = await initUpload.mutateAsync({ mimeType: file.type || 'application/octet-stream', name: file.name })
          await fetch(result.upload_url, { method: 'PUT', body: file })
          await completeUpload.mutateAsync(result.id)
        } catch {
          toast.error(`Failed to upload ${file.name}`)
        } finally {
          setPlaceholders((prev) => prev.filter((p) => p.clientId !== clientId))
        }
      }),
    )
  }, [initUpload, completeUpload])

  function handleDrop(e: React.DragEvent) {
    e.preventDefault()
    setDragOver(false)
    if (e.dataTransfer.files.length > 0) {
      void uploadFiles(e.dataTransfer.files)
    }
  }

  function handleDeleteConfirm() {
    if (!deleteTarget) return
    const ids = deleteTarget

    if (ids.length === 1) {
      deleteMutation.mutate(ids[0], {
        onSuccess: () => { toast.success('File deleted'); setDeleteTarget(null) },
        onError: () => { toast.error('Failed to delete file'); setDeleteTarget(null) },
      })
    } else {
      bulkDeleteMutation.mutate(ids, {
        onSuccess: () => { toast.success(`${ids.length} files deleted`); onSelectionChange(new Set()); setDeleteTarget(null) },
        onError: () => { toast.error('Failed to delete files'); setDeleteTarget(null) },
      })
    }
  }

  const isPendingDelete = deleteMutation.isPending || bulkDeleteMutation.isPending

  return (
    <div
      className={`relative flex-1 rounded-lg transition-colors ${dragOver ? 'bg-muted/50' : ''}`}
      onDragOver={(e) => { e.preventDefault(); setDragOver(true) }}
      onDragLeave={() => setDragOver(false)}
      onDrop={handleDrop}
      onClick={(e) => { if (e.target === e.currentTarget) onSelectionChange(new Set()) }}
    >
      <div className="grid grid-cols-3 gap-3 sm:grid-cols-4 md:grid-cols-5 lg:grid-cols-6">
        {placeholders.map((p) => (
          <div key={p.clientId} className="aspect-square rounded-md bg-muted flex items-center justify-center">
            <Loader2 className="h-6 w-6 text-muted-foreground animate-spin" />
          </div>
        ))}
        {files.map((file) => (
          <FileTile
            key={file.id}
            file={file}
            selected={selection.has(file.id)}
            onClick={(e) => handleTileClick(file.id, e)}
            onDoubleClick={() => {/* no-op */}}
            onDeleteClick={() => setDeleteTarget([file.id])}
            onContextMenu={(e) => handleTileContextMenu(file.id, e)}
          />
        ))}
      </div>

      {contextMenu && (
        <FileContextMenu
          x={contextMenu.x}
          y={contextMenu.y}
          onNewArtpiece={() => onNewArtpiece?.(contextMenu.ids)}
          onAddToExisting={() => onAddToExisting?.(contextMenu.ids)}
          onDelete={() => setDeleteTarget(contextMenu.ids)}
          onClose={() => setContextMenu(null)}
        />
      )}

      <FileDeleteDialog
        open={deleteTarget !== null}
        count={deleteTarget?.length ?? 0}
        isPending={isPendingDelete}
        onConfirm={handleDeleteConfirm}
        onCancel={() => setDeleteTarget(null)}
      />
    </div>
  )
}
