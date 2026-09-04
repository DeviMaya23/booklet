import { useKindeAuth } from '@kinde-oss/kinde-auth-react'

export default function AppTopBar() {
  const { user, logout } = useKindeAuth()

  return (
    <header className="flex h-12 items-center gap-2 border-b border-border px-4">
      <span className="text-xl font-bold font-title">Booklet</span>
      <div className="ml-auto">
        <button
          onClick={() => logout()}
          className="size-8 overflow-hidden rounded-full focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          aria-label="Log out"
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
        </button>
      </div>
    </header>
  )
}
