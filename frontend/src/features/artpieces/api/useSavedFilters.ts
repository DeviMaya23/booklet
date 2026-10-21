import { useQuery } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'
import { type SavedFilter } from './useCreateSavedFilter'

export const SAVED_FILTERS_QUERY_KEY = ['saved_filters'] as const

export function useSavedFilters() {
  const { getToken } = useKindeAuth()

  return useQuery({
    queryKey: SAVED_FILTERS_QUERY_KEY,
    queryFn: async () => {
      const res = await apiFetch('/saved_filters', getToken)
      if (!res.ok) throw new Error('Failed to fetch saved filters')
      return res.json() as Promise<SavedFilter[]>
    },
  })
}
