import { useMutation } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'

export interface UpdateFileParams {
  id: string
  name?: string | null
  notes?: string | null
}

export function useUpdateFile() {
  const { getToken } = useKindeAuth()

  return useMutation({
    mutationFn: async ({ id, name, notes }: UpdateFileParams) => {
      const res = await apiFetch(`/files/${id}`, getToken, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name: name ?? null, notes: notes ?? null }),
      })
      if (!res.ok) throw new Error('Failed to update file')
    },
  })
}
