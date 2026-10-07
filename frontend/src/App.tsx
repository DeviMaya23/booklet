import { Routes, Route, Navigate } from 'react-router-dom'
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

export default function App() {
  return (
    <Routes>
      <Route element={<PublicThemeLock />}>
        <Route path="/" element={<HomePage />} />
        <Route path="/callback" element={<CallbackPage />} />
      </Route>
      <Route element={<AuthGuard />}>
        <Route element={<AppLayout />}>
          <Route element={<AppShell />}>
            <Route index path="/app" element={<Navigate to="/app/characters" replace />} />
            <Route path="/app/characters" element={<CharactersPage />} />
            <Route path="/app/artists" element={<ArtistsPage />} />
            <Route path="/app/files" element={<FilesPage />} />
            <Route path="/app/artpieces" element={<ArtpiecesPage />} />
            <Route path="/app/artpieces/:id" element={<ArtpiecesPage />} />
            <Route path="/app/commissions" element={<CommissionsPage />} />
          </Route>
        </Route>
      </Route>
    </Routes>
  )
}
