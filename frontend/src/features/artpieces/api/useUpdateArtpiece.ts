import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'
import { ARTPIECES_QUERY_KEY } from './useArtpieces'
import { artpieceQueryKey } from './useArtpiece'

export interface UpdateArtpieceParams {
  id: string
  title: string | null
  notes: string | null
  artistId: string | null
  characterIds: string[]
}

export function useUpdateArtpiece() {
  const { getToken } = useKindeAuth()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ id, title, notes, artistId, characterIds }: UpdateArtpieceParams) => {
      const res = await apiFetch(`/artpieces/${id}`, getToken, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          title,
          notes,
          artist_id: artistId,
          character_ids: characterIds,
        }),
      })
      if (!res.ok) throw new Error('Failed to update artpiece')
      return res.json()
    },
    onSuccess: (_data, { id }) => {
      queryClient.invalidateQueries({ queryKey: artpieceQueryKey(id) })
      queryClient.invalidateQueries({ queryKey: ARTPIECES_QUERY_KEY })
    },
  })
}
