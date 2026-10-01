import { MoreHorizontal } from 'lucide-react'
import { Skeleton } from '@/components/ui/skeleton'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
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
      ) : (
        <Skeleton className="w-full h-full rounded-none" />
      )}

      <div
        className="absolute top-1 right-1"
        onClick={(e) => e.stopPropagation()}
      >
        <DropdownMenu>
          <DropdownMenuTrigger
            className="p-1 rounded bg-black/40 hover:bg-black/60 text-white opacity-0 group-hover:opacity-100 focus:opacity-100"
            aria-label="File actions"
          >
            <MoreHorizontal size={14} />
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem onClick={() => {/* no-op */}}>
              View Detail
            </DropdownMenuItem>
            <DropdownMenuItem onClick={onDeleteClick} variant="destructive">
              Delete
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </div>
  )
}
