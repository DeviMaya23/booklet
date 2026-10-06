import { useMutation } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'
import { sanitizeFilename } from '../lib/sanitizeFilename'

interface DownloadArtpieceParams {
  id: string
  title: string | null
}

export function useDownloadArtpiece() {
  const { getToken } = useKindeAuth()

  return useMutation({
    mutationFn: async ({ id, title }: DownloadArtpieceParams) => {
      const res = await apiFetch(`/artpieces/${id}/download`, getToken)
      if (!res.ok) throw new Error('Failed to download artpiece files')
      const blob = await res.blob()
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = title ? `${sanitizeFilename(title)}.zip` : 'artpiece.zip'
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)
      URL.revokeObjectURL(url)
    },
  })
}
