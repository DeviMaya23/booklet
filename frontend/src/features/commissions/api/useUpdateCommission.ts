import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'
import { COMMISSIONS_QUERY_KEY } from './useCommissions'

export interface UpdateCommissionParams {
  id: string
  title?: string | null
  artistId?: string | null
  status: string
  price?: number | null
  paid: boolean
  paidDate?: string | null
  finishDate?: string | null
  notes?: string | null
}

export function useUpdateCommission() {
  const { getToken } = useKindeAuth()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ id, ...params }: UpdateCommissionParams) => {
      const body: Record<string, unknown> = {
        status: params.status,
        paid: params.paid,
        title: params.title ?? null,
        artist_id: params.artistId ?? null,
        price: params.price ?? null,
        paid_date: params.paidDate ?? null,
        finish_date: params.finishDate ?? null,
        notes: params.notes ?? null,
      }

      const res = await apiFetch(`/commissions/${id}`, getToken, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      })
      if (!res.ok) throw new Error('Failed to update commission')
      return res.json()
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: COMMISSIONS_QUERY_KEY })
    },
  })
}
