import { Sheet, SheetContent, SheetTitle } from '@/components/ui/sheet'
import { ChooseFoldersPanel } from './ChooseFoldersPanel'

type PanelProps = Parameters<typeof ChooseFoldersPanel>[0]

interface ChooseFoldersDrawerProps extends Omit<PanelProps, 'onClose'> {
  open: boolean
  onOpenChange: (open: boolean) => void
}

/** The Choose folders panel in a right-side drawer over the detail page. */
export function ChooseFoldersDrawer({ open, onOpenChange, ...panelProps }: ChooseFoldersDrawerProps) {
  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent showCloseButton={false} className="w-full gap-0 p-0 data-[side=right]:sm:max-w-xl">
        <SheetTitle className="sr-only">Choose folders</SheetTitle>
        <div className="flex min-h-0 flex-1 flex-col p-6">
          <ChooseFoldersPanel {...panelProps} onClose={() => onOpenChange(false)} />
        </div>
      </SheetContent>
    </Sheet>
  )
}
