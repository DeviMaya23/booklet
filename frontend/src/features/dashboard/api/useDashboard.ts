import { useQuery } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'

export interface DashboardArtpiece {
  id: string
  title: string | null
  artist_name: string | null
  thumbnail_url: string | null
}

export interface DashboardCommission {
  id: string
  title: string | null
  artist_name: string | null
  artist_link: string | null
  status: 'waitlist' | 'wip'
  paid: boolean
  last_contacted_at: string | null
  created_at: string
}

export interface DashboardHousekeepingItem {
  id: string
  title: string | null
}

export interface DashboardHousekeeping {
  commissions_no_artist: DashboardHousekeepingItem[]
  artpieces_no_artist: DashboardHousekeepingItem[]
  done_no_artpieces: DashboardHousekeepingItem[]
  artpieces_no_files: DashboardHousekeepingItem[]
}

export interface DashboardData {
  recent_artpieces: DashboardArtpiece[]
  in_progress: DashboardCommission[]
  housekeeping: DashboardHousekeeping
}

export const DASHBOARD_QUERY_KEY = ['dashboard'] as const

export function useDashboard() {
  const { getToken } = useKindeAuth()

  return useQuery({
    queryKey: DASHBOARD_QUERY_KEY,
    queryFn: async () => {
      const res = await apiFetch('/dashboard', getToken)
      if (!res.ok) throw new Error('Failed to fetch dashboard')
      return res.json() as Promise<DashboardData>
    },
  })
}
