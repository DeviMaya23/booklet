import { useEffect } from 'react'
import { Outlet, useNavigate } from 'react-router-dom'
import { toast } from 'sonner'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { useMaintenanceActive } from '../lib/maintenanceStore'
import { useSessionExpired, setSessionExpired } from '../lib/sessionExpiredStore'
import MaintenancePage from './MaintenancePage'

export default function AppLayout() {
  const maintenance = useMaintenanceActive()
  const sessionExpired = useSessionExpired()
  const { isAuthenticated, isLoading } = useKindeAuth()
  const navigate = useNavigate()

  useEffect(() => {
    if (!sessionExpired) return
    if (isLoading || !isAuthenticated) return
    setSessionExpired(false)
    toast.error('Please log back in')
    navigate('/')
  }, [sessionExpired, isLoading, isAuthenticated, navigate])

  if (maintenance) return <MaintenancePage />

  return <Outlet />
}
