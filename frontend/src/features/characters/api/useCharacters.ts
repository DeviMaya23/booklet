import { useQuery } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'

export interface Character {
  id: string
  name: string
  avatar_url: string | null
  notes: string | null
  is_public: boolean
  folder_ids: string[]
  created_at: string
  updated_at: string
}

export const CHARACTERS_QUERY_KEY = ['characters'] as const

export function useCharacters() {
  const { getToken } = useKindeAuth()

  return useQuery({
    queryKey: CHARACTERS_QUERY_KEY,
    queryFn: async () => {
      const res = await apiFetch('/characters', getToken)
      if (!res.ok) throw new Error('Failed to fetch characters')
      return res.json() as Promise<Character[]>
    },
  })
}
