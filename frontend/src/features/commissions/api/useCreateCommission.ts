import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'
import { COMMISSIONS_QUERY_KEY } from './useCommissions'

export interface CreateCommissionParams {
  title?: string | null
  artistId?: string | null
  status?: string
  price?: number | null
  paid?: boolean
  paidDate?: string | null
  finishDate?: string | null
  notes?: string | null
  artpieceIds?: string[]
}

export function useCreateCommission() {
  const { getToken } = useKindeAuth()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (params: CreateCommissionParams) => {
      const body: Record<string, unknown> = {
        status: params.status ?? 'waitlist',
        paid: params.paid ?? false,
        artpiece_ids: params.artpieceIds ?? [],
      }
      if (params.title !== undefined) body.title = params.title
      if (params.artistId !== undefined) body.artist_id = params.artistId
      if (params.price !== undefined) body.price = params.price
      if (params.paidDate !== undefined) body.paid_date = params.paidDate
      if (params.finishDate !== undefined) body.finish_date = params.finishDate
      if (params.notes !== undefined) body.notes = params.notes

      const res = await apiFetch('/commissions', getToken, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      })
      if (!res.ok) throw new Error('Failed to create commission')
      return res.json()
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: COMMISSIONS_QUERY_KEY })
    },
  })
}
