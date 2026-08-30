import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { Navigate } from 'react-router-dom'

export default function HomePage() {
  const { isAuthenticated, isLoading, login } = useKindeAuth()

  if (!isLoading && isAuthenticated) return <Navigate to="/app" replace />

  return (
    <div className="flex min-h-screen items-center justify-center bg-background text-foreground">
      <button
        onClick={() => login()}
        className="rounded-md bg-foreground px-4 py-2 text-sm font-medium text-background"
      >
        Log in
      </button>
    </div>
  )
}
