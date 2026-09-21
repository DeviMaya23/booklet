import { useMutation } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'

export interface InitImageUploadInput {
  mimeType: string
  title?: string
  artistId?: string
  notes?: string
  characterIds?: string[]
}

export interface InitImageUploadResult {
  id: string
  upload_url: string
  expires_at: string
}

export function useInitImageUpload() {
  const { getToken } = useKindeAuth()

  return useMutation({
    mutationFn: async (input: InitImageUploadInput) => {
      const res = await apiFetch('/images', getToken, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          mime_type: input.mimeType,
          ...(input.title ? { title: input.title } : {}),
          ...(input.artistId ? { artist_id: input.artistId } : {}),
          ...(input.notes ? { notes: input.notes } : {}),
          ...(input.characterIds?.length ? { character_ids: input.characterIds } : {}),
        }),
      })
      if (!res.ok) throw new Error('Failed to initiate image upload')
      return res.json() as Promise<InitImageUploadResult>
    },
  })
}
