import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'
import { ARTISTS_QUERY_KEY, type Artist } from './useArtists'

export interface CreateArtistInput {
  name: string
  notes?: string
  artist_link?: string
}

export function useCreateArtist() {
  const { getToken } = useKindeAuth()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (input: CreateArtistInput) => {
      const res = await apiFetch('/artists', getToken, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(input),
      })
      if (res.status === 409) {
        const err = new Error('An artist with this name already exists')
        ;(err as Error & { status: number }).status = 409
        throw err
      }
      if (!res.ok) throw new Error('Failed to create artist')
      return res.json() as Promise<Artist>
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ARTISTS_QUERY_KEY })
    },
  })
}
