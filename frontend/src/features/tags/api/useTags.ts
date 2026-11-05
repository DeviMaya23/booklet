import { useQuery } from '@tanstack/react-query'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { apiFetch } from '@/lib/api'

export interface Tag {
  id: string
  name: string
  // Only present when the list was requested with `withUsage`.
  artpiece_count?: number
}

export const TAGS_QUERY_KEY = ['tags'] as const

export class TagConflictError extends Error {
  constructor() {
    super('Tag name already exists')
    this.name = 'TagConflictError'
  }
}

export type TagsSort = 'recent'

export interface UseTagsOptions {
  sort?: TagsSort
  // Adds `artpiece_count` to each tag. Costs an extra join server-side, so only
  // the Tags page asks for it.
  withUsage?: boolean
  // Set false to skip the request (e.g. a caller that only sometimes shows tags).
  enabled?: boolean
}

export function tagsQueryKey(sort?: TagsSort, withUsage?: boolean) {
  const base = sort ? ([...TAGS_QUERY_KEY, sort] as const) : TAGS_QUERY_KEY
  return withUsage ? ([...base, 'usage'] as const) : base
}

export function useTags(options?: UseTagsOptions) {
  const { getToken } = useKindeAuth()
  const sort = options?.sort
  const withUsage = options?.withUsage

  return useQuery({
    queryKey: tagsQueryKey(sort, withUsage),
    enabled: options?.enabled ?? true,
    queryFn: async () => {
      const params = new URLSearchParams()
      if (sort) params.set('sort', sort)
      if (withUsage) params.set('with_usage', 'true')
      const query = params.toString()
      const path = query ? `/tags?${query}` : '/tags'
      const res = await apiFetch(path, getToken)
      if (!res.ok) throw new Error('Failed to fetch tags')
      return res.json() as Promise<Tag[]>
    },
  })
}
