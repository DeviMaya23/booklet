import { Info, ImageIcon } from 'lucide-react'
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip'

interface ArtpieceTileProps {
  title: string | null
  artistName: string | null
  imageUrl: string | null
  showDetails: boolean
  onClick: () => void
  onInfoClick: () => void
}

export default function ArtpieceTile({
  title,
  artistName,
  imageUrl,
  showDetails,
  onClick,
  onInfoClick,
}: ArtpieceTileProps) {
  return (
    <figure className="group/tile select-none">
      <div className="relative aspect-square overflow-hidden rounded-lg bg-muted transition-transform duration-150 group-hover/tile:scale-[1.02] group-hover/tile:shadow-md">
        {imageUrl ? (
          <img
            src={imageUrl}
            alt={title ?? 'Untitled'}
            className="h-full w-full object-cover"
          />
        ) : (
          <div
            role="img"
            aria-label="No cover image"
            className="flex h-full w-full items-center justify-center"
          >
            <ImageIcon className="h-8 w-8 text-muted-foreground" aria-hidden />
          </div>
        )}

        {/* Transparent full-area button — keyboard and click access to lightbox */}
        <button
          type="button"
          aria-label={`View ${title ?? 'Untitled'}`}
          className="absolute inset-0 cursor-pointer"
          onClick={onClick}
        />

        <Tooltip>
          <TooltipTrigger
            render={
              <button
                type="button"
                aria-label="Open details"
                className="absolute right-2 top-2 z-10 flex h-7 w-7 items-center justify-center rounded-full bg-background/80 text-foreground opacity-0 transition-opacity duration-150 hover:bg-background group-hover/tile:opacity-100"
                onClick={onInfoClick}
              />
            }
          >
            <Info className="h-4 w-4" />
          </TooltipTrigger>
          <TooltipContent side="bottom">Open details</TooltipContent>
        </Tooltip>
      </div>

      {showDetails && (
        <figcaption className="mt-1.5 px-0.5">
          <p className="truncate text-sm font-medium leading-tight">
            {title ?? 'Untitled'}
          </p>
          {artistName && (
            <p className="truncate text-xs text-muted-foreground">
              {artistName}
            </p>
          )}
        </figcaption>
      )}
    </figure>
  )
}
