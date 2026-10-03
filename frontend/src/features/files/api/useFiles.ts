import { useQuery } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'

export interface File {
  id: string
  file_url: string | null
  mime_type: string
  thumbnail_url: string | null
  thumbnail_gen_state: 'pending' | 'done' | 'failed' | 'not_applicable'
  artpiece_id: string | null
  name: string | null
  notes: string | null
  created_at: string
  updated_at: string
}

export const FILES_QUERY_KEY = ['files', 'unassigned'] as const

export function useFiles() {
  const { getToken } = useKindeAuth()

  return useQuery({
    queryKey: FILES_QUERY_KEY,
    queryFn: async () => {
      const res = await apiFetch('/files?unassigned=true', getToken)
      if (!res.ok) throw new Error('Failed to fetch files')
      return res.json() as Promise<File[]>
    },
    refetchInterval: (query) => {
      const data = query.state.data
      return Array.isArray(data) && data.some((f) => f.thumbnail_gen_state === 'pending') ? 2000 : false
    },
  })
}
