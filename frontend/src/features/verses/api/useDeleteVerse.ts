import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'
import { VERSES_QUERY_KEY } from './useVerses'
import { CHARACTERS_QUERY_KEY } from '@/features/characters/api/useCharacters'

export function useDeleteVerse() {
  const { getToken } = useKindeAuth()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (id: string) => {
      const res = await apiFetch(`/verses/${id}`, getToken, { method: 'DELETE' })
      if (!res.ok) throw new Error('Failed to delete verse')
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: VERSES_QUERY_KEY })
      queryClient.invalidateQueries({ queryKey: CHARACTERS_QUERY_KEY })
    },
  })
}
