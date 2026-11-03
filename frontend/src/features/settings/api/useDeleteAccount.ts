import { useMutation } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'

export function useDeleteAccount() {
  const { getToken } = useKindeAuth()

  return useMutation({
    mutationFn: async () => {
      const res = await apiFetch('/me', getToken, { method: 'DELETE' })
      if (!res.ok) throw new Error('Failed to delete account')
    },
  })
}
