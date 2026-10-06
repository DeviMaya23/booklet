import { useQuery } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'

export interface Commission {
  id: string
  title: string | null
  artist_id: string | null
  artist_name: string | null
  artist_link: string | null
  status: 'waitlist' | 'wip' | 'done'
  price: number | null
  paid: boolean
  paid_date: string | null
  finish_date: string | null
  last_contacted_at: string | null
  notes: string | null
  created_at: string
  updated_at: string
}

export const COMMISSIONS_QUERY_KEY = ['commissions'] as const

export function useCommissions() {
  const { getToken } = useKindeAuth()

  return useQuery({
    queryKey: COMMISSIONS_QUERY_KEY,
    queryFn: async () => {
      const res = await apiFetch('/commissions', getToken)
      if (!res.ok) throw new Error('Failed to fetch commissions')
      return res.json() as Promise<Commission[]>
    },
  })
}
