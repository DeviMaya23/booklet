import { useMutation } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'

export function useDeleteSavedFilter() {
  const { getToken } = useKindeAuth()

  return useMutation({
    mutationFn: async (id: string) => {
      const res = await apiFetch(`/saved_filters/${id}`, getToken, {
        method: 'DELETE',
      })
      if (!res.ok && res.status !== 404) throw new Error('Failed to delete saved filter')
    },
  })
}
