import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'
import { FILES_QUERY_KEY } from './useFiles'

export function useBulkDeleteFiles() {
  const { getToken } = useKindeAuth()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (ids: string[]) => {
      const res = await apiFetch('/files', getToken, {
        method: 'DELETE',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ ids }),
      })
      if (!res.ok) throw new Error('Failed to bulk delete files')
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: FILES_QUERY_KEY })
    },
  })
}
