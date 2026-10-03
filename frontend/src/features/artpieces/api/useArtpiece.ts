import { useQuery } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'
import { type ArtpieceCharacter } from './useArtpieces'

export interface ArtpieceFile {
  id: string
  file_url: string
  thumbnail_url: string | null
}

export interface ArtpieceDetail {
  id: string
  title: string | null
  artist_id: string | null
  artist_name: string | null
  cover_file_id: string | null
  thumbnail_url: string | null
  notes: string | null
  characters: ArtpieceCharacter[]
  files: ArtpieceFile[]
  created_at: string
  updated_at: string
}

export function artpieceQueryKey(id: string) {
  return ['artpieces', id] as const
}

export function useArtpiece(id: string | null) {
  const { getToken } = useKindeAuth()

  return useQuery({
    queryKey: artpieceQueryKey(id ?? ''),
    queryFn: async () => {
      const res = await apiFetch(`/artpieces/${id}`, getToken)
      if (!res.ok) throw new Error('Failed to fetch artpiece')
      return res.json() as Promise<ArtpieceDetail>
    },
    enabled: id !== null,
  })
}
