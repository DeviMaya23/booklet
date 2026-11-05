import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'
import { artpieceQueryKey } from '@/features/artpieces/api/useArtpiece'
import { FILES_QUERY_KEY } from './useFiles'

export interface UpdateFileParams {
  id: string
  name?: string | null
  notes?: string | null
  isVisible?: boolean
  // When the file is attached to an artpiece, its cached detail is refreshed on success
  // (cover reassignment, name/notes/visibility shown there).
  artpieceId?: string
}

export function useUpdateFile() {
  const { getToken } = useKindeAuth()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ id, name, notes, isVisible }: UpdateFileParams) => {
      const res = await apiFetch(`/files/${id}`, getToken, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          // Only send what the caller provided. Empty inputs clear the column (null)
          // rather than storing "".
          ...(name !== undefined && { name: name || null }),
          ...(notes !== undefined && { notes: notes || null }),
          ...(isVisible !== undefined && { is_visible: isVisible }),
        }),
      })
      if (!res.ok) throw new Error('Failed to update file')
    },
    // Invalidation lives here, not in per-call mutate callbacks, because those are
    // skipped when the calling component unmounts before the request returns.
    onSuccess: (_data, { artpieceId }) => {
      void queryClient.invalidateQueries({ queryKey: FILES_QUERY_KEY })
      if (artpieceId) {
        void queryClient.invalidateQueries({ queryKey: artpieceQueryKey(artpieceId) })
      }
    },
  })
}
