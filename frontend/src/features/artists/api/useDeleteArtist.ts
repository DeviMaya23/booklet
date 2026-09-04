import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'
import { ARTISTS_QUERY_KEY } from './useArtists'

export function useDeleteArtist() {
  const { getToken } = useKindeAuth()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (id: string) => {
      const res = await apiFetch(`/artists/${id}`, getToken, { method: 'DELETE' })
      if (!res.ok) throw new Error('Failed to delete artist')
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ARTISTS_QUERY_KEY })
    },
  })
}
