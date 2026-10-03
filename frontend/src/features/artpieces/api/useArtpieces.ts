import { useQuery } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'

export interface ArtpieceCharacter {
  id: string
  name: string
}

export interface ArtpieceSummary {
  id: string
  title: string | null
  artist_id: string | null
  artist_name: string | null
  cover_file_id: string | null
  thumbnail_url: string | null
  notes: string | null
  characters: ArtpieceCharacter[]
  created_at: string
  updated_at: string
}

export const ARTPIECES_QUERY_KEY = ['artpieces'] as const

export function useArtpieces() {
  const { getToken } = useKindeAuth()

  return useQuery({
    queryKey: ARTPIECES_QUERY_KEY,
    queryFn: async () => {
      const res = await apiFetch('/artpieces', getToken)
      if (!res.ok) throw new Error('Failed to fetch artpieces')
      return res.json() as Promise<ArtpieceSummary[]>
    },
  })
}
