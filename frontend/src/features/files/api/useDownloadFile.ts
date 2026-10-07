import { useMutation } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'

export function useDownloadFile() {
  const { getToken } = useKindeAuth()

  return useMutation({
    mutationFn: async (fileId: string) => {
      const res = await apiFetch(`/files/${fileId}/download`, getToken)
      if (!res.ok) throw new Error('Failed to get download URL')
      const { download_url } = (await res.json()) as { download_url: string }
      window.open(download_url, '_blank')
    },
  })
}
