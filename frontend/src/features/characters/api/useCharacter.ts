import { useQuery } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'
import { CHARACTERS_QUERY_KEY, type Character } from './useCharacters'

export function characterQueryKey(id: string) {
  return [...CHARACTERS_QUERY_KEY, 'detail', id] as const
}

export function useCharacter(id: string | undefined) {
  const { getToken } = useKindeAuth()

  return useQuery({
    queryKey: characterQueryKey(id ?? ''),
    queryFn: async () => {
      const res = await apiFetch(`/characters/${id}`, getToken)
      if (!res.ok) throw new Error('Failed to fetch character')
      return res.json() as Promise<Character>
    },
    enabled: !!id,
  })
}
