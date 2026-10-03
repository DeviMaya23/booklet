import { useMutation } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'

export interface AttachFilesToArtpieceParams {
  artpieceId: string
  fileIds: string[]
}

export function useAttachFilesToArtpiece() {
  const { getToken } = useKindeAuth()

  return useMutation({
    mutationFn: async ({ artpieceId, fileIds }: AttachFilesToArtpieceParams) => {
      const res = await apiFetch(`/artpieces/${artpieceId}/files`, getToken, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ file_ids: fileIds }),
      })
      if (!res.ok) throw new Error('Failed to attach files')
      return res.json()
    },
  })
}
