import { useMutation } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'

export interface InitFileUploadResult {
  id: string
  upload_url: string
  thumbnail_upload_url: string | null
  expires_at: string
}

export function useInitFileUpload() {
  const { getToken } = useKindeAuth()

  return useMutation({
    mutationFn: async ({ mimeType, name, extension, fileSizeBytes }: { mimeType: string; name?: string; extension?: string; fileSizeBytes: number }) => {
      const res = await apiFetch('/files', getToken, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ mime_type: mimeType, name: name ?? null, extension: extension ?? null, file_size: fileSizeBytes }),
      })
      if (!res.ok) throw new Error('Failed to initiate file upload')
      return res.json() as Promise<InitFileUploadResult>
    },
  })
}
