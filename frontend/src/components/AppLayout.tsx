import { Outlet } from 'react-router-dom'
import { useMaintenanceActive } from '../lib/maintenanceStore'
import MaintenancePage from './MaintenancePage'

export default function AppLayout() {
  const maintenance = useMaintenanceActive()

  if (maintenance) return <MaintenancePage />

  return <Outlet />
}
