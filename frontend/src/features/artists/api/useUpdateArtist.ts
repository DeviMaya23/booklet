import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'
import { ARTISTS_QUERY_KEY, type Artist } from './useArtists'

export interface UpdateArtistInput {
  id: string
  name: string
  notes: string | null
  artist_link: string | null
}

export function useUpdateArtist() {
  const { getToken } = useKindeAuth()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ id, ...body }: UpdateArtistInput) => {
      const res = await apiFetch(`/artists/${id}`, getToken, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      })
      if (!res.ok) {
        const err = new Error('Failed to update artist')
        ;(err as Error & { status: number }).status = res.status
        throw err
      }
      return res.json() as Promise<Artist>
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ARTISTS_QUERY_KEY })
    },
  })
}
