import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'
import { VERSES_QUERY_KEY, type Verse } from './useVerses'

export interface CreateVerseInput {
  name: string
}

export function useCreateVerse() {
  const { getToken } = useKindeAuth()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (input: CreateVerseInput) => {
      const res = await apiFetch('/verses', getToken, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(input),
      })
      if (!res.ok) {
        const err = new Error('Failed to create verse')
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
