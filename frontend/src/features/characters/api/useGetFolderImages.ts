import { useQuery } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'

export interface FolderImage {
  image_id: string
  thumbnail_url: string
}

export function useGetFolderImages(folderID: string | undefined) {
  const { getToken } = useKindeAuth()

  return useQuery({
    queryKey: ['folderImages', folderID],
    queryFn: async () => {
      const res = await apiFetch(`/folders/${folderID}/images`, getToken)
      if (!res.ok) throw new Error('Failed to fetch folder images')
      const data = (await res.json()) as { images: FolderImage[] }
      return data.images
    },
    enabled: !!folderID,
    retry: false,
  })
}
