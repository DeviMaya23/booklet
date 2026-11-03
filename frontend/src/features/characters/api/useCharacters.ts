import { useQuery } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'

export interface Character {
  id: string
  name: string
  avatar_url: string | null
  notes: string | null
  is_public: boolean
  verse_id: string | null
  folders: { id: string; name: string }[]
  created_at: string
  updated_at: string
}

export const CHARACTERS_QUERY_KEY = ['characters'] as const

export type CharactersSort = 'recent'

export interface UseCharactersOptions {
  sort?: CharactersSort
}

export function charactersQueryKey(sort?: CharactersSort) {
  return sort ? ([...CHARACTERS_QUERY_KEY, sort] as const) : CHARACTERS_QUERY_KEY
}

export function useCharacters(options?: UseCharactersOptions) {
  const { getToken } = useKindeAuth()
  const sort = options?.sort

  return useQuery({
    queryKey: charactersQueryKey(sort),
    queryFn: async () => {
      const path = sort ? `/characters?sort=${sort}` : '/characters'
      const res = await apiFetch(path, getToken)
      if (!res.ok) throw new Error('Failed to fetch characters')
      return res.json() as Promise<Character[]>
    },
  })
}
