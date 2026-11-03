import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'
import { CHARACTERS_QUERY_KEY } from '@/features/characters/api/useCharacters'

export interface CreateArtpieceParams {
  title: string
  notes?: string | null
  artistId?: string | null
  characterIds: string[]
  fileIds: string[]
  isReference?: boolean
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
  const queryClient = useQueryClient()

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
          is_reference: params.isReference ?? false,
        }),
      })
      if (!res.ok) throw new Error('Failed to create artpiece')
      return res.json() as Promise<Artpiece>
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: CHARACTERS_QUERY_KEY })
    },
  })
}
