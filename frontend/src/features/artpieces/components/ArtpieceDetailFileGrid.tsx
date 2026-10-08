import { FileIcon, Loader2, Plus, Star, X } from 'lucide-react'
import { type ArtpieceFile } from '../api/useArtpiece'
import { mimeTypeLabel } from '../lib/mimeTypeLabel'

interface ViewModeProps {
  mode: 'view'
}

interface EditModeProps {
  mode: 'edit'
  locallyRemovedIds: Set<string>
  pendingCoverFileId: string | null
  onSetCover: (fileId: string) => void
  onRemove: (fileId: string) => void
  onUndoRemove: (fileId: string) => void
  onAddFiles: () => void
  dragHandlers: {
    dragOver: boolean
    onDragOver: (e: React.DragEvent) => void
    onDragLeave: () => void
    onDrop: (e: React.DragEvent) => void
  }
  uploadPlaceholders: { clientId: string }[]
}

type ArtpieceDetailFileGridProps = {
  files: ArtpieceFile[]
  coverFileId: string | null
  onOpen?: (index: number) => void
} & (ViewModeProps | EditModeProps)

export default function ArtpieceDetailFileGrid(props: ArtpieceDetailFileGridProps) {
  const { files, coverFileId, mode, onOpen } = props

  const effectiveCoverId =
    mode === 'edit' ? (props.pendingCoverFileId ?? coverFileId) : coverFileId

  if (files.length === 0 && mode === 'view') {
    return <p className="text-sm text-muted-foreground">No files attached.</p>
  }

  return (
    <div className="grid grid-cols-3 gap-2 sm:grid-cols-4 md:grid-cols-5">
      {mode === 'edit' && props.uploadPlaceholders.map((p) => (
        <div key={p.clientId} className="aspect-square rounded-md bg-muted flex items-center justify-center">
          <Loader2 className="size-6 text-muted-foreground animate-spin" />
        </div>
      ))}
      {files.map((file, fileIndex) => {
        const isCover = file.id === effectiveCoverId
        const isRemoved = mode === 'edit' && props.locallyRemovedIds.has(file.id)
        const label = mimeTypeLabel(file.mime_type)

        return (
          <div key={file.id} className="relative aspect-square">
            <div
              className={`relative h-full w-full overflow-hidden rounded-md bg-muted transition-opacity ${isRemoved ? 'opacity-40' : ''} ${mode === 'view' && onOpen ? 'cursor-pointer' : ''}`}
              onClick={mode === 'view' && onOpen ? () => onOpen(fileIndex) : undefined}
            >
              {file.thumbnail_url ? (
                <img src={file.thumbnail_url} alt="" className="h-full w-full object-cover" />
              ) : (
                <div className="flex h-full w-full items-center justify-center">
                  <FileIcon className="size-6 text-muted-foreground" />
                </div>
              )}

              {/* Cover chip */}
              {isCover && !isRemoved && (
                <div className="absolute left-1 top-1 flex items-center gap-0.5 rounded-full bg-black/60 px-1.5 py-0.5 text-[10px] font-medium text-white">
                  <Star className="size-2.5 fill-white" />
                  Cover
                </div>
              )}

              {/* Type label */}
              {label && (
                <div className="absolute bottom-1 left-1 rounded bg-black/60 px-1 py-0.5 text-[10px] font-medium text-white">
                  {label}
                </div>
              )}

              {/* Edit controls */}
              {mode === 'edit' && !isRemoved && (
                <div className="absolute right-1 top-1 flex gap-1" onClick={(e) => e.stopPropagation()}>
                  {!isCover && (
                    <button
                      type="button"
                      title="Set as cover"
                      className="inline-flex size-6 items-center justify-center rounded-[min(var(--radius-md),10px)] bg-black/40 text-white hover:bg-black/60 focus-visible:outline-none"
                      onClick={() => props.onSetCover(file.id)}
                    >
                      <Star className="size-3" />
                    </button>
                  )}
                  <button
                    type="button"
                    title="Remove"
                    className="inline-flex size-6 items-center justify-center rounded-[min(var(--radius-md),10px)] bg-black/40 text-white hover:bg-black/60 focus-visible:outline-none"
                    onClick={() => props.onRemove(file.id)}
                  >
                    <X className="size-3" />
                  </button>
                </div>
              )}
            </div>

            {/* Removed overlay */}
            {mode === 'edit' && isRemoved && (
              <div className="absolute inset-0 flex flex-col items-center justify-center gap-1 rounded-md">
                <p className="text-center text-[10px] font-medium text-foreground">Will be removed</p>
                <button
                  type="button"
                  className="rounded border border-border bg-background px-2 py-0.5 text-[10px] font-medium hover:bg-muted focus-visible:outline-none"
                  onClick={() => props.onUndoRemove(file.id)}
                >
                  Undo
                </button>
              </div>
            )}
          </div>
        )
      })}

      {/* Add files cell — last in grid, edit mode only */}
      {mode === 'edit' && (
        <div className="aspect-square">
          <button
            type="button"
            className={`flex h-full w-full flex-col items-center justify-center gap-1 rounded-md border border-dashed border-border text-muted-foreground transition-colors hover:bg-muted/50 focus-visible:outline-none ${props.dragHandlers.dragOver ? 'bg-muted/50 ring-2 ring-ring' : ''}`}
            onClick={props.onAddFiles}
            onDragOver={props.dragHandlers.onDragOver}
            onDragLeave={props.dragHandlers.onDragLeave}
            onDrop={props.dragHandlers.onDrop}
          >
            <>
              <Plus className="size-4" />
              <span className="text-center text-[10px]">
                Add files<br />or drop here
              </span>
            </>
          </button>
        </div>
      )}
    </div>
  )
}
