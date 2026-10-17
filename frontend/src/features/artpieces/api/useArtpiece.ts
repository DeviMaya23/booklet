import { useQuery } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'
import { type ThumbnailGenState } from '@/features/files/api/useFiles'
import { type ArtpieceCharacter } from './useArtpieces'

export interface ArtpieceFile {
  id: string
  file_url: string
  thumbnail_url: string | null
  thumbnail_gen_state: ThumbnailGenState
  mime_type: string
  name: string | null
  extension: string | null
  width: number | null
  height: number | null
}

export interface ArtpieceDetail {
  id: string
  title: string | null
  artist_id: string | null
  artist_name: string | null
  artist_link: string | null
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
    refetchInterval: (query) => {
      const data = query.state.data
      return data?.files.some((f) => f.thumbnail_gen_state === 'pending') ? 2000 : false
    },
  })
}
