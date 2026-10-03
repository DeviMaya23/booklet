import { useMutation } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'

export interface CreateArtpieceParams {
  title: string
  notes?: string | null
  artistId?: string | null
  characterIds: string[]
  fileIds: string[]
}

export interface Artpiece {
  id: string
  title: string | null
  notes: string | null
  artist_id: string | null
  created_at: string
  updated_at: string
}

export function useCreateArtpiece() {
  const { getToken } = useKindeAuth()

  return useMutation({
    mutationFn: async (params: CreateArtpieceParams) => {
      const res = await apiFetch('/artpieces', getToken, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          title: params.title,
          notes: params.notes ?? null,
          artist_id: params.artistId ?? null,
          character_ids: params.characterIds,
          file_ids: params.fileIds,
        }),
      })
      if (!res.ok) throw new Error('Failed to create artpiece')
      return res.json() as Promise<Artpiece>
    },
  })
}
