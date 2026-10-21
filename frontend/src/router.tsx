import { createBrowserRouter, Navigate } from 'react-router-dom'
import PublicThemeLock from './components/PublicThemeLock'
import AuthGuard from './features/auth/components/AuthGuard'
import AppLayout from './components/AppLayout'
import AppShell from './components/AppShell'
import HomePage from './pages/HomePage'
import CallbackPage from './pages/CallbackPage'
import CharactersPage from './pages/CharactersPage'
import ArtistsPage from './pages/ArtistsPage'
import FilesPage from './pages/FilesPage'
import ArtpiecesPage from './pages/ArtpiecesPage'
import CommissionsPage from './pages/CommissionsPage'
import DashboardPage from './pages/DashboardPage'

export const router = createBrowserRouter([
  {
    element: <PublicThemeLock />,
    children: [
      { path: '/', element: <HomePage /> },
      { path: '/callback', element: <CallbackPage /> },
    ],
  },
  {
    element: <AuthGuard />,
    children: [
      {
        element: <AppLayout />,
        children: [
          {
            element: <AppShell />,
            children: [
              { path: '/app', element: <Navigate to="/app/dashboard" replace /> },
              { path: '/app/dashboard', element: <DashboardPage /> },
              { path: '/app/characters', element: <CharactersPage /> },
              { path: '/app/artists', element: <ArtistsPage /> },
              { path: '/app/files', element: <FilesPage /> },
              { path: '/app/artpieces', element: <ArtpiecesPage /> },
              { path: '/app/artpieces/:id', element: <ArtpiecesPage /> },
              { path: '/app/commissions', element: <CommissionsPage /> },
            ],
          },
        ],
      },
    ],
  },
])
