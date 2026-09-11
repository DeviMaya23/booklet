import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'
import { CHARACTERS_QUERY_KEY, type Character } from './useCharacters'

export interface CreateCharacterInput {
  name: string
  notes?: string
}

export function useCreateCharacter() {
  const { getToken } = useKindeAuth()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (input: CreateCharacterInput) => {
      const res = await apiFetch('/characters', getToken, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(input),
      })
      if (!res.ok) {
        const err = new Error('Failed to create character')
        ;(err as Error & { status: number }).status = res.status
        throw err
      }
      return res.json() as Promise<Character>
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: CHARACTERS_QUERY_KEY })
    },
  })
}
