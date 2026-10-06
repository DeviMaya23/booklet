import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'
import { COMMISSIONS_QUERY_KEY } from './useCommissions'
import { commissionQueryKey } from './useCommission'

export function useReplaceArtpieces() {
  const { getToken } = useKindeAuth()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ id, artpieceIds }: { id: string; artpieceIds: string[] }) => {
      const res = await apiFetch(`/commissions/${id}/artpieces`, getToken, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ artpiece_ids: artpieceIds }),
      })
      if (!res.ok) throw new Error('Failed to replace artpieces')
      return res.json()
    },
    onSuccess: (_data, { id }) => {
      queryClient.invalidateQueries({ queryKey: COMMISSIONS_QUERY_KEY })
      queryClient.invalidateQueries({ queryKey: commissionQueryKey(id) })
    },
  })
}
