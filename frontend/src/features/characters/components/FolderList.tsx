import { AlertTriangle, ExternalLink, FolderOpen, PanelRight, X } from 'lucide-react'
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip'
import { Button } from '@/components/ui/button'
import { CircleHelp } from 'lucide-react'
import type { PublicFolder } from '../api/usePublicFolders'

interface FolderListProps {
  selectedFolders: PublicFolder[]
  invalidFolders: { id: string; name: string }[]
  onRemove: (id: string) => void
  onDismissInvalid: (id: string) => void
  onOpenPanel: () => void
  disabled?: boolean
}

export function FolderList({
  selectedFolders,
  invalidFolders,
  onRemove,
  onDismissInvalid,
  onOpenPanel,
  disabled,
}: FolderListProps) {
  const bookleafHost = import.meta.env.VITE_BOOKLEAF_HOST ?? ''

  return (
    <div className="flex flex-col gap-1.5">
      <div className="flex items-center gap-1.5">
        <span className="text-sm font-medium">Folders</span>
        <TooltipProvider>
          <Tooltip>
            <TooltipTrigger>
              <button type="button" className="text-muted-foreground hover:text-foreground">
                <CircleHelp className="size-3.5" />
                <span className="sr-only">About folders</span>
              </button>
            </TooltipTrigger>
            <TooltipContent>
              The folder list comes from Bookleaf. Only your public folders are shown.
            </TooltipContent>
          </Tooltip>
        </TooltipProvider>
      </div>

      {(selectedFolders.length > 0 || invalidFolders.length > 0) && (
        <div className="flex flex-col gap-1">
          {selectedFolders.map((folder) => (
            <div
              key={folder.id}
              className="flex items-center gap-2 rounded-md border px-3 py-2 text-sm"
            >
              <FolderOpen className="size-4 shrink-0 text-muted-foreground" />
              <span className="flex-1 truncate">{folder.name}</span>
              {folder.token && (
                <a
                  href={`${bookleafHost}/share/${folder.token}`}
                  target="_blank"
                  rel="noreferrer"
                  className="text-muted-foreground hover:text-foreground"
                  onClick={(e) => e.stopPropagation()}
                  aria-label={`Open ${folder.name} in Bookleaf`}
                >
                  <ExternalLink className="size-3.5" />
                </a>
              )}
              <button
                type="button"
                className="text-muted-foreground hover:text-foreground"
                onClick={() => onRemove(folder.id)}
                disabled={disabled}
                aria-label={`Remove ${folder.name}`}
              >
                <X className="size-3.5" />
              </button>
            </div>
          ))}

          {invalidFolders.map((folder) => (
            <div
              key={folder.id}
              className="flex items-center gap-2 rounded-md border border-orange-200 bg-orange-50 px-3 py-2 text-sm dark:border-orange-900 dark:bg-orange-950/30"
            >
              <AlertTriangle className="size-4 shrink-0 text-orange-500" />
              <div className="min-w-0 flex-1">
                <p className="truncate font-medium">{folder.name}</p>
                <p className="text-xs text-muted-foreground">No longer public in Bookleaf.</p>
              </div>
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={() => onDismissInvalid(folder.id)}
                disabled={disabled}
              >
                Remove
              </Button>
            </div>
          ))}
        </div>
      )}

      <Button
        type="button"
        variant="outline"
        size="sm"
        className="self-start gap-1.5"
        onClick={onOpenPanel}
        disabled={disabled}
      >
        <PanelRight className="size-4" />
        Choose folders
      </Button>
    </div>
  )
}
