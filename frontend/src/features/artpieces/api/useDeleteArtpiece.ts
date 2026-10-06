import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'
import { ARTPIECES_QUERY_KEY } from './useArtpieces'

export function useDeleteArtpiece() {
  const { getToken } = useKindeAuth()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ id, deleteFiles }: { id: string; deleteFiles: boolean }) => {
      const url = deleteFiles ? `/artpieces/${id}?delete_files=true` : `/artpieces/${id}`
      const res = await apiFetch(url, getToken, { method: 'DELETE' })
      if (!res.ok) throw new Error('Failed to delete artpiece')
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ARTPIECES_QUERY_KEY })
    },
  })
}
