import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'
import { ARTPIECES_QUERY_KEY } from '@/features/artpieces/api/useArtpieces'
import { TAGS_QUERY_KEY } from './useTags'

export function useDeleteTag() {
  const { getToken } = useKindeAuth()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (id: string) => {
      const res = await apiFetch(`/tags/${id}`, getToken, { method: 'DELETE' })
      if (!res.ok) throw new Error('Failed to delete tag')
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: TAGS_QUERY_KEY })
      queryClient.invalidateQueries({ queryKey: ARTPIECES_QUERY_KEY })
    },
    // A 404 means the list on screen is stale; refresh it either way.
    onError: () => {
      queryClient.invalidateQueries({ queryKey: TAGS_QUERY_KEY })
    },
  })
}
