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
  sublabel?: string
  onDeleteClick: () => void
  onClick?: () => void
  onDoubleClick?: () => void
}

export default function ResourceCard({ imageUrl, label, sublabel, onDeleteClick, onClick, onDoubleClick }: ResourceCardProps) {
  return (
    <figure className="group/card select-none">
      <div
        className={`relative aspect-square overflow-hidden rounded-lg bg-muted transition-transform duration-150 group-hover/card:scale-[1.02] group-hover/card:shadow-md${onClick || onDoubleClick ? ' cursor-pointer' : ''}`}
        onClick={onClick}
        onDoubleClick={onDoubleClick}
      >
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

        <div className="absolute right-1 top-1" onClick={(e) => e.stopPropagation()}>
          <DropdownMenu>
            <DropdownMenuTrigger className="inline-flex size-6 items-center justify-center rounded-[min(var(--radius-md),10px)] bg-background/80 text-foreground opacity-0 transition-opacity duration-150 hover:bg-background focus-visible:outline-none group-hover/card:opacity-100">
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

      <figcaption className="mt-1.5 px-0.5">
        <p className="truncate text-sm font-medium leading-tight">{label}</p>
        {sublabel && (
          <p className="truncate text-xs text-muted-foreground">{sublabel}</p>
        )}
      </figcaption>
    </figure>
  )
}
