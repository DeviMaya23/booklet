import { useQuery } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'

export interface Artist {
  id: string
  name: string
  notes: string | null
  artist_link: string | null
  created_at: string
  updated_at: string
}

export const ARTISTS_QUERY_KEY = ['artists'] as const

export function useArtists() {
  const { getToken } = useKindeAuth()

  return useQuery({
    queryKey: ARTISTS_QUERY_KEY,
    queryFn: async () => {
      const res = await apiFetch('/artists', getToken)
      if (!res.ok) throw new Error('Failed to fetch artists')
      return res.json() as Promise<Artist[]>
    },
  })
}
