import { useState } from 'react'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { useIsMobile } from '@/hooks/use-mobile'
import SettingsModal from '@/features/settings/components/SettingsModal'

export default function AppTopBar() {
  const { user, logout } = useKindeAuth()
  const isMobile = useIsMobile()
  const [settingsOpen, setSettingsOpen] = useState(false)

  return (
    <header className="flex h-12 items-center gap-2 border-b border-border px-4">
      <span className="text-xl font-bold font-title">Booklet</span>
      <div className="ml-auto">
        <DropdownMenu>
          <DropdownMenuTrigger
            className="size-8 overflow-hidden rounded-full focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
            aria-label="Account menu"
          >
            {user?.picture ? (
              <img
                src={user.picture}
                alt={user.givenName ?? 'User'}
                className="size-full object-cover"
              />
            ) : (
              <span className="flex size-full items-center justify-center bg-muted text-xs font-medium text-muted-foreground">
                {user?.givenName?.[0]?.toUpperCase() ?? '?'}
              </span>
            )}
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-[220px]">
            {user?.email && (
              <div className="truncate px-2 py-1.5 text-xs text-muted-foreground">{user.email}</div>
            )}
            {!isMobile && (
              <DropdownMenuItem onClick={() => setSettingsOpen(true)}>Settings</DropdownMenuItem>
            )}
            <DropdownMenuSeparator />
            <DropdownMenuItem variant="destructive" onClick={() => logout()}>
              Sign out
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
      <SettingsModal open={settingsOpen} onOpenChange={setSettingsOpen} />
    </header>
  )
}
