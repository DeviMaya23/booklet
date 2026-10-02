import { useMutation } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'

export interface InitFileUploadResult {
  id: string
  upload_url: string
  expires_at: string
}

export function useInitFileUpload() {
  const { getToken } = useKindeAuth()

  return useMutation({
    mutationFn: async ({ mimeType, name }: { mimeType: string; name?: string }) => {
      const res = await apiFetch('/files', getToken, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ mime_type: mimeType, name: name ?? null }),
      })
      if (!res.ok) throw new Error('Failed to initiate file upload')
      return res.json() as Promise<InitFileUploadResult>
    },
  })
}
