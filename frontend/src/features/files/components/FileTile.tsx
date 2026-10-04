import { FileIcon, Trash2 } from 'lucide-react'
import { Skeleton } from '@/components/ui/skeleton'
import { type File } from '../api/useFiles'

interface FileTileProps {
  file: File
  selected: boolean
  onClick: (e: React.MouseEvent) => void
  onDoubleClick: () => void
  onDeleteClick: () => void
  onContextMenu: (e: React.MouseEvent) => void
}

export default function FileTile({ file, selected, onClick, onDoubleClick, onDeleteClick, onContextMenu }: FileTileProps) {
  return (
    <div className="flex flex-col gap-1">
      <div
        className={`group relative aspect-square rounded-md overflow-hidden border-2 cursor-pointer select-none ${
          selected ? 'border-blue-500' : 'border-transparent'
        }`}
        onClick={onClick}
        onDoubleClick={onDoubleClick}
        onContextMenu={onContextMenu}
      >
        {file.thumbnail_url ? (
          <img
            src={file.thumbnail_url}
            alt=""
            className="w-full h-full object-cover"
            draggable={false}
          />
        ) : file.thumbnail_gen_state === 'pending' ? (
          <Skeleton className="w-full h-full rounded-none" />
        ) : (
          <div className="flex h-full w-full items-center justify-center bg-muted">
            <FileIcon size={24} className="text-muted-foreground" />
          </div>
        )}

        <div
          className="absolute top-1 right-1"
          onClick={(e) => e.stopPropagation()}
        >
          <button
            type="button"
            className="p-1 rounded bg-black/40 hover:bg-black/60 text-white opacity-0 group-hover:opacity-100 focus:opacity-100"
            aria-label="Delete file"
            onClick={onDeleteClick}
          >
            <Trash2 size={14} />
          </button>
        </div>
      </div>

      {file.name && (
        <p className="truncate text-xs text-muted-foreground px-0.5">{file.name}</p>
      )}
    </div>
  )
}
