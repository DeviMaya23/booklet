import { useEffect } from 'react'
import { Outlet, useNavigate } from 'react-router-dom'
import { toast } from 'sonner'
import { useMaintenanceActive } from '../lib/maintenanceStore'
import { useSessionExpired, setSessionExpired } from '../lib/sessionExpiredStore'
import MaintenancePage from './MaintenancePage'

export default function AppLayout() {
  const maintenance = useMaintenanceActive()
  const sessionExpired = useSessionExpired()
  const navigate = useNavigate()

  useEffect(() => {
    if (!sessionExpired) return
    setSessionExpired(false)
    toast.error('Please log back in')
    navigate('/')
  }, [sessionExpired, navigate])

  if (maintenance) return <MaintenancePage />

  return <Outlet />
}
