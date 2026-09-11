import { useMutation } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'

export interface CompleteAvatarUploadInput {
  characterId: string
  uploadId: string
}

export function useCompleteAvatarUpload() {
  const { getToken } = useKindeAuth()

  return useMutation({
    mutationFn: async ({ characterId, uploadId }: CompleteAvatarUploadInput) => {
      const res = await apiFetch(
        `/characters/${characterId}/avatar/${uploadId}/complete`,
        getToken,
        { method: 'POST' },
      )
      if (!res.ok) throw new Error('Failed to complete avatar upload')
    },
  })
}
