import { setMaintenanceActive } from './maintenanceStore'

const BASE_URL = import.meta.env.VITE_API_BASE_URL ?? ''
const MAINTENANCE_BYPASS_STORAGE_KEY = 'booklet-maintenance-bypass'

export async function apiFetch(
  path: string,
  getToken: () => Promise<string | undefined>,
  options?: RequestInit,
): Promise<Response> {
  const token = await getToken()
  const bypassToken = localStorage.getItem(MAINTENANCE_BYPASS_STORAGE_KEY)
  const res = await fetch(`${BASE_URL}${path}`, {
    ...options,
    headers: {
      ...options?.headers,
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...(bypassToken ? { 'X-Booklet-Bypass': bypassToken } : {}),
    },
  })

  setMaintenanceActive(res.headers.get('X-Booklet-Maintenance') === 'true')

  return res
}
