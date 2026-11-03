import { useState } from 'react'
import { XIcon } from 'lucide-react'
import { Dialog, DialogClose, DialogContent, DialogTitle } from '@/components/ui/dialog'
import { cn } from '@/lib/utils'
import AccountSection from './AccountSection'
import AppSection from './AppSection'

type SettingsSection = 'account' | 'app'

const SECTIONS: { id: SettingsSection; label: string }[] = [
  { id: 'account', label: 'Account' },
  { id: 'app', label: 'App' },
]

interface SettingsModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export default function SettingsModal({ open, onOpenChange }: SettingsModalProps) {
  const [activeSection, setActiveSection] = useState<SettingsSection>('account')

  const activeLabel = SECTIONS.find((section) => section.id === activeSection)?.label

  return (
    <Dialog
      open={open}
      onOpenChange={onOpenChange}
      onOpenChangeComplete={(isOpen) => {
        if (!isOpen) setActiveSection('account')
      }}
    >
      <DialogContent className="gap-0 overflow-hidden p-0 sm:max-w-3xl">
        <DialogTitle className="sr-only">Settings</DialogTitle>
        <div className="flex min-h-[580px]">
          <nav
            aria-label="Settings sections"
            className="flex w-44 shrink-0 flex-col gap-1 border-r border-border bg-sidebar p-3"
          >
            <p className="px-2 pb-3 font-serif text-base font-semibold text-foreground">Settings</p>
            {SECTIONS.map((section) => (
              <button
                key={section.id}
                type="button"
                onClick={() => setActiveSection(section.id)}
                aria-current={activeSection === section.id ? 'true' : undefined}
                className={cn(
                  'rounded-md px-2 py-1.5 text-left text-sm',
                  activeSection === section.id
                    ? 'bg-accent font-medium text-accent-foreground'
                    : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground',
                )}
              >
                {section.label}
              </button>
            ))}
          </nav>
          <div className="flex min-w-0 flex-1 flex-col">
            <div className="flex items-center justify-between border-b border-border px-6 py-4">
              <h2 className="font-serif text-base font-semibold text-foreground">{activeLabel}</h2>
              <DialogClose
                render={
                  <button
                    type="button"
                    className="flex size-[26px] shrink-0 items-center justify-center rounded-full bg-accent text-muted-foreground hover:bg-muted"
                  />
                }
              >
                <XIcon className="size-3.5" />
                <span className="sr-only">Close</span>
              </DialogClose>
            </div>
            <div className="min-h-0 flex-1 overflow-y-auto p-6">
              {activeSection === 'account' && <AccountSection />}
              {activeSection === 'app' && <AppSection />}
            </div>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  )
}
