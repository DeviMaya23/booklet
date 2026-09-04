import { Routes, Route, Navigate } from 'react-router-dom'
import PublicThemeLock from './components/PublicThemeLock'
import AuthGuard from './features/auth/components/AuthGuard'
import AppLayout from './components/AppLayout'
import AppShell from './components/AppShell'
import HomePage from './pages/HomePage'
import CallbackPage from './pages/CallbackPage'
import CharactersPage from './pages/CharactersPage'
import ImagesPage from './pages/ImagesPage'
import ArtistsPage from './pages/ArtistsPage'

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
            <Route path="/app/images" element={<ImagesPage />} />
            <Route path="/app/artists" element={<ArtistsPage />} />
          </Route>
        </Route>
      </Route>
    </Routes>
  )
}
