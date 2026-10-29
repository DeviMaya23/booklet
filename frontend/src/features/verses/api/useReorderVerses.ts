import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'
import { VERSES_QUERY_KEY } from './useVerses'

export function useReorderVerses() {
  const { getToken } = useKindeAuth()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (ids: string[]) => {
      const res = await apiFetch('/verses/reorder', getToken, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ ids }),
      })
      if (!res.ok) throw new Error('Failed to reorder verses')
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: VERSES_QUERY_KEY })
    },
  })
}
