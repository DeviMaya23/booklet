import { useQuery } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'

export interface Image {
  id: string
  mime_type: string
  title: string | null
  thumbnail_url: string | null
  artist_id: string | null
  artist_name: string | null
  notes: string | null
  characters: { id: string; name: string }[]
  created_at: string
  updated_at: string
}

export const IMAGES_QUERY_KEY = ['images'] as const

export function useImages() {
  const { getToken } = useKindeAuth()

  return useQuery({
    queryKey: IMAGES_QUERY_KEY,
    queryFn: async () => {
      const res = await apiFetch('/images', getToken)
      if (!res.ok) throw new Error('Failed to fetch images')
      return res.json() as Promise<Image[]>
    },
  })
}
