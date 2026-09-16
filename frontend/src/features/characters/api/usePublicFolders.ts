import { useQuery } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'

export interface PublicFolder {
  id: string
  name: string
}

const PUBLIC_FOLDERS_QUERY_KEY = ['publicFolders'] as const

export function usePublicFolders() {
  const { getToken } = useKindeAuth()

  return useQuery({
    queryKey: PUBLIC_FOLDERS_QUERY_KEY,
    queryFn: async () => {
      const res = await apiFetch('/folders', getToken)
      if (!res.ok) throw new Error('Failed to fetch folders')
      const data = (await res.json()) as { folder_list: { folder_id: string; folder_name: string }[] }
      return data.folder_list.map(
        (f): PublicFolder => ({ id: f.folder_id, name: f.folder_name }),
      )
    },
    retry: false,
    staleTime: 15 * 60 * 1000,
  })
}
