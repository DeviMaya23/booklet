import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'
import { FILES_QUERY_KEY } from './useFiles'

export function useCompleteFileUpload() {
  const { getToken } = useKindeAuth()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (pendingId: string) => {
      const res = await apiFetch(`/files/${pendingId}/complete`, getToken, {
        method: 'POST',
      })
      if (!res.ok) throw new Error('Failed to complete file upload')
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: FILES_QUERY_KEY })
    },
  })
}
