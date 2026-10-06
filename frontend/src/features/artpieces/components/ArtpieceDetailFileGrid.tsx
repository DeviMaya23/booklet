import { FileIcon, MoreHorizontal, Star } from 'lucide-react'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { type ArtpieceFile } from '../api/useArtpiece'

interface ViewModeProps {
  mode: 'view'
}

interface EditModeProps {
  mode: 'edit'
  locallyRemovedIds: Set<string>
  pendingCoverFileId: string | null
  onSetCover: (fileId: string) => void
  onRemove: (fileId: string) => void
}

type ArtpieceDetailFileGridProps = {
  files: ArtpieceFile[]
  coverFileId: string | null
} & (ViewModeProps | EditModeProps)

export default function ArtpieceDetailFileGrid(props: ArtpieceDetailFileGridProps) {
  const { files, coverFileId, mode } = props

  const visibleFiles =
    mode === 'edit'
      ? files.filter((f) => !props.locallyRemovedIds.has(f.id))
      : files

  const effectiveCoverId =
    mode === 'edit' ? (props.pendingCoverFileId ?? coverFileId) : coverFileId

  if (visibleFiles.length === 0) {
    return <p className="text-sm text-muted-foreground">No files attached.</p>
  }

  return (
    <div className="grid grid-cols-3 gap-2 sm:grid-cols-4 md:grid-cols-5">
      {visibleFiles.map((file) => {
        const isCover = file.id === effectiveCoverId
        return (
          <div key={file.id} className="relative aspect-square">
            <div className="relative h-full w-full overflow-hidden rounded-md bg-muted">
              {file.thumbnail_url ? (
                <img src={file.thumbnail_url} alt="" className="h-full w-full object-cover" />
              ) : (
                <div className="flex h-full w-full items-center justify-center">
                  <FileIcon className="size-6 text-muted-foreground" />
                </div>
              )}

              {isCover && (
                <div className="absolute bottom-1 right-1 rounded-full bg-black/50 p-0.5">
                  <Star className="size-3 fill-white text-white" />
                </div>
              )}

              {mode === 'edit' && (
                <div className="absolute right-1 top-1" onClick={(e) => e.stopPropagation()}>
                  <DropdownMenu>
                    <DropdownMenuTrigger className="inline-flex size-6 items-center justify-center rounded-[min(var(--radius-md),10px)] bg-black/30 text-white hover:bg-black/50 focus-visible:outline-none">
                      <MoreHorizontal className="size-3" />
                    </DropdownMenuTrigger>
                    <DropdownMenuContent side="bottom" align="start">
                      <DropdownMenuItem
                        disabled={isCover}
                        onClick={() => props.onSetCover(file.id)}
                      >
                        Set as cover
                      </DropdownMenuItem>
                      <DropdownMenuItem
                        variant="destructive"
                        onClick={() => props.onRemove(file.id)}
                      >
                        Remove from artpiece
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </div>
              )}
            </div>
          </div>
        )
      })}
    </div>
  )
}
