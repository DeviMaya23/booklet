import { useMutation } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'

export interface InitAvatarUploadInput {
  characterId: string
  mimeType: string
  fileSizeBytes: number
}

export interface InitAvatarUploadResult {
  id: string
  upload_url: string
  expires_at: string
}

export function useInitAvatarUpload() {
  const { getToken } = useKindeAuth()

  return useMutation({
    mutationFn: async ({ characterId, mimeType, fileSizeBytes }: InitAvatarUploadInput) => {
      const res = await apiFetch(`/characters/${characterId}/avatar/init`, getToken, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ mime_type: mimeType, file_size: fileSizeBytes }),
      })
      if (!res.ok) throw new Error('Failed to initiate avatar upload')
      return res.json() as Promise<InitAvatarUploadResult>
    },
  })
}
