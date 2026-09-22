import { useMutation } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'

export function useCompleteImageUpload() {
  const { getToken } = useKindeAuth()

  return useMutation({
    mutationFn: async (pendingId: string) => {
      const res = await apiFetch(`/images/${pendingId}/complete`, getToken, {
        method: 'POST',
      })
      if (!res.ok) throw new Error('Failed to complete image upload')
    },
  })
}
