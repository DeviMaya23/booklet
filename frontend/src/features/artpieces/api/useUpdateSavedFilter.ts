import { useMutation } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'
import { type SavedFilter } from './useCreateSavedFilter'

export interface UpdateSavedFilterParams {
  id: string
  thumbnailArtpieceId: string
}

export function useUpdateSavedFilter() {
  const { getToken } = useKindeAuth()

  return useMutation({
    mutationFn: async (params: UpdateSavedFilterParams) => {
      const res = await apiFetch(`/saved_filters/${params.id}`, getToken, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ thumbnail_artpiece_id: params.thumbnailArtpieceId }),
      })
      if (!res.ok) throw new Error('Failed to update saved filter')
      return res.json() as Promise<SavedFilter>
    },
  })
}
