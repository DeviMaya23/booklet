import { useState } from 'react'
import { Link } from 'react-router-dom'
import { ImageOff, MoreHorizontal } from 'lucide-react'
import { type LucideIcon } from 'lucide-react'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'

export interface ThumbnailTileMenuItem {
  label: string
  icon?: LucideIcon
  variant?: 'destructive'
  onClick: () => void
}

interface ThumbnailTileProps {
  imageUrl: string | null
  label: string
  sublabel?: string
  href: string
  menuItems?: ThumbnailTileMenuItem[]
}

export default function ThumbnailTile({ imageUrl, label, sublabel, href, menuItems }: ThumbnailTileProps) {
  const [imgFailed, setImgFailed] = useState(false)
  const hasMenu = menuItems != null && menuItems.length > 0

  return (
    <figure className="group/tile select-none">
      <div className="relative aspect-square overflow-hidden rounded-lg bg-muted transition-transform duration-150 group-hover/tile:scale-[1.02] group-hover/tile:shadow-md">
        <Link to={href} className="block size-full" aria-label={label}>
          {imageUrl && !imgFailed ? (
            <img
              src={imageUrl}
              alt={label}
              className="size-full object-cover"
              onError={() => setImgFailed(true)}
            />
          ) : (
            <div className="flex size-full items-center justify-center">
              <ImageOff className="size-8 text-muted-foreground" />
            </div>
          )}
        </Link>

        {hasMenu && (
          <div className="absolute right-1 top-1">
            <DropdownMenu>
              <DropdownMenuTrigger
                className="inline-flex size-6 items-center justify-center rounded-[min(var(--radius-md),10px)] bg-background/80 text-foreground opacity-0 transition-opacity duration-150 hover:bg-background focus-visible:outline-none group-hover/tile:opacity-100"
                aria-label="More options"
              >
                <MoreHorizontal className="size-3" />
              </DropdownMenuTrigger>
              <DropdownMenuContent side="bottom" align="end">
                {menuItems!.map((item) => {
                  const Icon = item.icon
                  return (
                    <DropdownMenuItem
                      key={item.label}
                      variant={item.variant}
                      onClick={item.onClick}
                    >
                      {Icon && <Icon />}
                      {item.label}
                    </DropdownMenuItem>
                  )
                })}
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        )}
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
