import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'
import { ARTPIECES_QUERY_KEY } from '@/features/artpieces/api/useArtpieces'
import { TAGS_QUERY_KEY, TagConflictError, type Tag } from './useTags'

export function useUpdateTag() {
  const { getToken } = useKindeAuth()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ id, name }: { id: string; name: string }) => {
      const res = await apiFetch(`/tags/${id}`, getToken, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name }),
      })
      if (res.status === 409) throw new TagConflictError()
      if (!res.ok) throw new Error('Failed to rename tag')
      return res.json() as Promise<Tag>
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: TAGS_QUERY_KEY })
      queryClient.invalidateQueries({ queryKey: ARTPIECES_QUERY_KEY })
    },
  })
}
