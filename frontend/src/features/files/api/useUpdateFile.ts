import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'
import { FILES_QUERY_KEY } from './useFiles'

export interface UpdateFileParams {
  id: string
  name?: string | null
  notes?: string | null
}

export function useUpdateFile() {
  const { getToken } = useKindeAuth()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ id, name, notes }: UpdateFileParams) => {
      const res = await apiFetch(`/files/${id}`, getToken, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name: name ?? null, notes: notes ?? null }),
      })
      if (!res.ok) throw new Error('Failed to update file')
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: FILES_QUERY_KEY })
    },
  })
}
