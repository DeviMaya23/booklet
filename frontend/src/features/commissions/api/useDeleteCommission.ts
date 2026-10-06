import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'
import { COMMISSIONS_QUERY_KEY } from './useCommissions'

export function useDeleteCommission() {
  const { getToken } = useKindeAuth()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (id: string) => {
      const res = await apiFetch(`/commissions/${id}`, getToken, { method: 'DELETE' })
      if (!res.ok) throw new Error('Failed to delete commission')
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: COMMISSIONS_QUERY_KEY })
    },
  })
}
