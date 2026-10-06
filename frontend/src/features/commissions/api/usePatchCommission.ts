import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'
import { COMMISSIONS_QUERY_KEY } from './useCommissions'

export interface PatchCommissionParams {
  id: string
  status?: string
  paid?: boolean
  paidDate?: string | null
  lastContactedAt?: string | null
}

export function usePatchCommission() {
  const { getToken } = useKindeAuth()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ id, status, paid, paidDate, lastContactedAt }: PatchCommissionParams) => {
      const body: Record<string, unknown> = {}
      if (status !== undefined) body.status = status
      if (paid !== undefined) body.paid = paid
      if (paidDate !== undefined) body.paid_date = paidDate
      if (lastContactedAt !== undefined) body.last_contacted_at = lastContactedAt

      const res = await apiFetch(`/commissions/${id}`, getToken, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      })
      if (!res.ok) throw new Error('Failed to patch commission')
      return res.json()
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: COMMISSIONS_QUERY_KEY })
    },
  })
}
