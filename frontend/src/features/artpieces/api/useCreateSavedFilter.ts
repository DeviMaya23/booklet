import { useMutation } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'

export interface CreateSavedFilterParams {
  name: string
  thumbnailArtpieceId?: string | null
  filterPayload: {
    artist_id: string | null
    character_ids: string[]
    match_mode: 'all' | 'any'
  }
}

export interface SavedFilter {
  id: string
  name: string
  thumbnail_url: string | null
  filter_payload: {
    artist_id: string | null
    character_ids: string[]
    match_mode: 'all' | 'any'
  }
  created_at: string
  updated_at: string
}

export interface DuplicateNameError {
  type: 'duplicate_name'
}

export function useCreateSavedFilter() {
  const { getToken } = useKindeAuth()

  return useMutation({
    mutationFn: async (params: CreateSavedFilterParams) => {
      const res = await apiFetch('/saved_filters', getToken, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          name: params.name,
          thumbnail_artpiece_id: params.thumbnailArtpieceId ?? null,
          filter_payload: params.filterPayload,
        }),
      })

      if (res.status === 422) {
        const body = await res.json() as { errors?: { field: string; message: string }[] }
        const nameError = body.errors?.find((e) => e.field === 'name')
        if (nameError) {
          const err: DuplicateNameError & Error = Object.assign(new Error(nameError.message), { type: 'duplicate_name' as const })
          throw err
        }
        throw new Error('Validation error')
      }

      if (!res.ok) throw new Error('Failed to create saved filter')
      return res.json() as Promise<SavedFilter>
    },
  })
}
