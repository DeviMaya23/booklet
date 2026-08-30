import { Routes, Route } from 'react-router-dom'
import PublicThemeLock from './components/PublicThemeLock'
import AuthGuard from './features/auth/components/AuthGuard'
import AppLayout from './components/AppLayout'
import HomePage from './pages/HomePage'
import CallbackPage from './pages/CallbackPage'
import AppPage from './pages/AppPage'

export default function App() {
  return (
    <Routes>
      <Route element={<PublicThemeLock />}>
        <Route path="/" element={<HomePage />} />
        <Route path="/callback" element={<CallbackPage />} />
      </Route>
      <Route element={<AuthGuard />}>
        <Route element={<AppLayout />}>
          <Route path="/app" element={<AppPage />} />
        </Route>
      </Route>
    </Routes>
  )
}
