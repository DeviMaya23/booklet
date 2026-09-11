import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'
import { CHARACTERS_QUERY_KEY, type Character } from './useCharacters'

export interface UpdateCharacterInput {
  id: string
  name: string
  notes: string | null
}

export function useUpdateCharacter() {
  const { getToken } = useKindeAuth()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ id, name, notes }: UpdateCharacterInput) => {
      const res = await apiFetch(`/characters/${id}`, getToken, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name, notes, is_public: false, folder_ids: [] }),
      })
      if (!res.ok) {
        const err = new Error('Failed to update character')
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
