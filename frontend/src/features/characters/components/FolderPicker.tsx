import { CircleHelp } from 'lucide-react'
import {
  Tooltip,
  TooltipTrigger,
  TooltipContent,
  TooltipProvider,
} from '@/components/ui/tooltip'
import TokenInput from './TokenInput'

export interface FolderItem {
  id: string
  name: string
}

interface FolderPickerProps {
  selected: FolderItem[]
  available: FolderItem[]
  onAdd: (folder: FolderItem) => void
  onRemove: (id: string) => void
  disabled: boolean
  isError?: boolean
  onRetry?: () => void
}

export function FolderPicker({
  selected,
  available,
  onAdd,
  onRemove,
  disabled,
  isError,
  onRetry,
}: FolderPickerProps) {
  function handleChange(items: FolderItem[]) {
    const removedId = selected.find((s) => !items.some((i) => i.id === s.id))?.id
    const added = items.find((i) => !selected.some((s) => s.id === i.id))
    if (removedId) onRemove(removedId)
    if (added) onAdd(added)
  }

  return (
    <div className="flex flex-col gap-1.5">
      <div className="flex items-center gap-1.5">
        <span className="text-sm font-medium">Folders</span>
        <TooltipProvider>
          <Tooltip>
            <TooltipTrigger>
              <button type="button" className="text-muted-foreground hover:text-foreground">
                <CircleHelp className="size-3.5" />
              </button>
            </TooltipTrigger>
            <TooltipContent>
              The folder list comes from Bookleaf. Only your public folders are shown.
            </TooltipContent>
          </Tooltip>
        </TooltipProvider>
      </div>
      <TokenInput
        items={selected}
        onChange={handleChange}
        suggestions={available}
        disabled={disabled}
        placeholder="Search folders…"
      />
      {isError && (
        <p className="text-sm text-muted-foreground">
          Couldn't reach folder list.{' '}
          <button
            type="button"
            className="underline hover:text-foreground"
            onClick={onRetry}
          >
            Retry
          </button>
        </p>
      )}
    </div>
  )
}
