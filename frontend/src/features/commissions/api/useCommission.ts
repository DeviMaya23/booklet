import { useQuery } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'
import { type Commission } from './useCommissions'

export interface CommissionArtpieceSummary {
  id: string
  thumbnail_url: string | null
}

export interface CommissionDetail extends Commission {
  artpieces: CommissionArtpieceSummary[]
}

export function commissionQueryKey(id: string) {
  return ['commissions', id] as const
}

export function useCommission(id: string | null) {
  const { getToken } = useKindeAuth()

  return useQuery({
    queryKey: commissionQueryKey(id ?? ''),
    queryFn: async () => {
      const res = await apiFetch(`/commissions/${id}`, getToken)
      if (!res.ok) throw new Error('Failed to fetch commission')
      const data = await res.json() as CommissionDetail
      return { ...data, artpieces: data.artpieces ?? [] }
    },
    enabled: id !== null,
  })
}
