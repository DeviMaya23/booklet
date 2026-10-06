import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'
import { artpieceQueryKey } from './useArtpiece'

export function useSetCover() {
  const { getToken } = useKindeAuth()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ artpieceId, fileId }: { artpieceId: string; fileId: string }) => {
      const res = await apiFetch(`/artpieces/${artpieceId}/cover`, getToken, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ file_id: fileId }),
      })
      if (!res.ok) throw new Error('Failed to set cover')
      return res.json()
    },
    onSuccess: (_data, { artpieceId }) => {
      queryClient.invalidateQueries({ queryKey: artpieceQueryKey(artpieceId) })
    },
  })
}
