import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'
import { VERSES_QUERY_KEY, type Verse } from './useVerses'

export interface UpdateVerseInput {
  id: string
  name: string
}

export function useUpdateVerse() {
  const { getToken } = useKindeAuth()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ id, ...body }: UpdateVerseInput) => {
      const res = await apiFetch(`/verses/${id}`, getToken, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      })
      if (!res.ok) {
        const err = new Error('Failed to update verse')
        ;(err as Error & { status: number }).status = res.status
        throw err
      }
      return res.json() as Promise<Verse>
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: VERSES_QUERY_KEY })
    },
  })
}
