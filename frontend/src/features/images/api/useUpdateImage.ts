import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'
import { IMAGES_QUERY_KEY, type Image } from './useImages'

export interface UpdateImageInput {
  id: string
  title: string | null
  notes: string | null
  artistId: string | null
  characterIds: string[]
}

export function useUpdateImage() {
  const { getToken } = useKindeAuth()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ id, title, notes, artistId, characterIds }: UpdateImageInput) => {
      const res = await apiFetch(`/images/${id}`, getToken, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          title,
          notes,
          artist_id: artistId,
          character_ids: characterIds,
        }),
      })
      if (!res.ok) throw new Error('Failed to update image')
      return res.json() as Promise<Image>
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: IMAGES_QUERY_KEY })
    },
  })
}
