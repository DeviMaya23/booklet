import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'
import { ARTPIECES_QUERY_KEY } from '@/features/artpieces/api/useArtpieces'
import { TAGS_QUERY_KEY, type Tag } from './useTags'

export interface MergeTarget {
  source: Tag
  target: Tag
}

export function useMergeTag() {
  const { getToken } = useKindeAuth()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ sourceId, targetId }: { sourceId: string; targetId: string }) => {
      const res = await apiFetch(`/tags/${sourceId}/merge`, getToken, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ target_id: targetId }),
      })
      if (res.status === 404) throw new Error('Tag no longer exists')
      if (!res.ok) throw new Error('Failed to merge tags')
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
