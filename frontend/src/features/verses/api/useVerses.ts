import { useQuery } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'

export interface Verse {
  id: string
  name: string
  sort_order: number
  last_used_at: string | null
  created_at: string
  updated_at: string
}

export const VERSES_QUERY_KEY = ['verses'] as const

export function useVerses() {
  const { getToken } = useKindeAuth()

  return useQuery({
    queryKey: VERSES_QUERY_KEY,
    queryFn: async () => {
      const res = await apiFetch('/verses', getToken)
      if (!res.ok) throw new Error('Failed to fetch verses')
      return res.json() as Promise<Verse[]>
    },
  })
}
