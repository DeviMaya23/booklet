import { ImageOff, MoreHorizontal, Trash2 } from 'lucide-react'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'

interface ResourceCardProps {
  imageUrl?: string | null
  label: string
  onDeleteClick: () => void
}

export default function ResourceCard({ imageUrl, label, onDeleteClick }: ResourceCardProps) {
  return (
    <div className="relative aspect-square overflow-hidden rounded-lg bg-muted">
      {imageUrl ? (
        <img
          src={imageUrl}
          alt={label}
          className="size-full object-cover"
        />
      ) : (
        <div className="flex size-full items-center justify-center">
          <ImageOff className="size-8 text-muted-foreground" />
        </div>
      )}

      <div className="absolute inset-x-0 bottom-0 bg-gradient-to-t from-black/60 to-transparent px-2 pb-2 pt-6">
        <span className="truncate text-xs font-medium text-white">{label}</span>
      </div>

      <div className="absolute right-1 top-1">
        <DropdownMenu>
          <DropdownMenuTrigger className="inline-flex size-6 items-center justify-center rounded-[min(var(--radius-md),10px)] bg-black/30 text-white hover:bg-black/50 focus-visible:outline-none">
            <MoreHorizontal className="size-3" />
          </DropdownMenuTrigger>
          <DropdownMenuContent side="bottom" align="end">
            <DropdownMenuItem variant="destructive" onClick={onDeleteClick}>
              <Trash2 />
              Delete
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </div>
  )
}
