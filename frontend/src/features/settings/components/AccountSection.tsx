import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { getFullName } from '@/features/auth/lib/profile'
import { useDeleteAccount } from '../api/useDeleteAccount'
import { useDeleteAccountConfirmation } from '../useDeleteAccountConfirmation'

export default function AccountSection() {
  const { user, logout } = useKindeAuth()
  const queryClient = useQueryClient()
  const deleteAccount = useDeleteAccount()
  const confirmation = useDeleteAccountConfirmation(user?.email)

  function handleDelete() {
    deleteAccount.mutate(undefined, {
      onSuccess: () => {
        queryClient.clear()
        logout()
      },
      onError: () => {
        toast.error('Failed to delete account')
      },
    })
  }

  return (
    <div>
      <div className="mb-6">
        <p className="mb-2.5 text-[10px] font-semibold uppercase tracking-wider text-muted-foreground">
          Profile
        </p>
        <div className="flex items-center gap-4">
          <div className="size-[52px] shrink-0 overflow-hidden rounded-full">
            {user?.picture ? (
              <img
                src={user.picture}
                alt={user.givenName ?? 'User'}
                className="size-full object-cover"
              />
            ) : (
              <span className="flex size-full items-center justify-center bg-muted text-base font-medium text-muted-foreground">
                {user?.givenName?.[0]?.toUpperCase() ?? '?'}
              </span>
            )}
          </div>
          <div className="flex flex-1 flex-col gap-1.5">
            <div className="font-medium">{user ? getFullName(user) : ''}</div>
            <div className="rounded-md border bg-input px-2.5 py-1.5 text-sm text-muted-foreground">
              {user?.email ?? ''}
            </div>
          </div>
        </div>
      </div>

      <div className="border-b border-border pb-5">
        <Button variant="outline" onClick={() => logout()}>
          Log out
        </Button>
      </div>

      <div className="pt-5">
        <p className="mb-1.5 text-sm font-medium text-destructive">Delete account</p>
        <p className="mb-3.5 text-xs leading-relaxed text-muted-foreground">
          Permanently deletes your account and all associated data. This cannot be undone.
        </p>
        {!confirmation.confirming ? (
          <Button variant="destructive" onClick={confirmation.start}>
            Delete account…
          </Button>
        ) : (
          <div className="flex flex-col gap-3">
            <p className="text-sm text-muted-foreground">
              Type <span className="font-medium text-foreground">{user?.email}</span> to confirm.
            </p>
            <Input
              placeholder="Type your account email"
              aria-label="Account email confirmation"
              value={confirmation.typedEmail}
              onChange={(e) => confirmation.updateTypedEmail(e.target.value)}
              autoFocus
            />
            <div className="flex gap-2">
              <Button variant="outline" onClick={confirmation.cancel}>
                Cancel
              </Button>
              <Button
                variant="destructive"
                disabled={!confirmation.emailMatches || deleteAccount.isPending}
                onClick={handleDelete}
              >
                Delete account
              </Button>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
