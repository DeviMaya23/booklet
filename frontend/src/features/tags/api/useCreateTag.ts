import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'
import { TAGS_QUERY_KEY, type Tag } from './useTags'

export class TagConflictError extends Error {
  constructor() {
    super('Tag name already exists')
    this.name = 'TagConflictError'
  }
}

export function useCreateTag() {
  const { getToken } = useKindeAuth()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (name: string) => {
      const res = await apiFetch('/tags', getToken, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name }),
      })
      if (res.status === 409) throw new TagConflictError()
      if (!res.ok) throw new Error('Failed to create tag')
      return res.json() as Promise<Tag>
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: TAGS_QUERY_KEY })
    },
  })
}
