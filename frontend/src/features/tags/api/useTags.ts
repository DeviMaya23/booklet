import { useQuery } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'

export interface Tag {
  id: string
  name: string
}

export const TAGS_QUERY_KEY = ['tags'] as const

export type TagsSort = 'recent'

export interface UseTagsOptions {
  sort?: TagsSort
  // Set false to skip the request (e.g. a caller that only sometimes shows tags).
  enabled?: boolean
}

export function tagsQueryKey(sort?: TagsSort) {
  return sort ? ([...TAGS_QUERY_KEY, sort] as const) : TAGS_QUERY_KEY
}

export function useTags(options?: UseTagsOptions) {
  const { getToken } = useKindeAuth()
  const sort = options?.sort

  return useQuery({
    queryKey: tagsQueryKey(sort),
    enabled: options?.enabled ?? true,
    queryFn: async () => {
      const path = sort ? `/tags?sort=${sort}` : '/tags'
      const res = await apiFetch(path, getToken)
      if (!res.ok) throw new Error('Failed to fetch tags')
      return res.json() as Promise<Tag[]>
    },
  })
}
