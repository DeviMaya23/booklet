import { renderHook, waitFor } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useUpdateTag } from './useUpdateTag'
import { TagConflictError } from './useTags'
import { TAGS_QUERY_KEY } from './useTags'
import { ARTPIECES_QUERY_KEY } from '@/features/artpieces/api/useArtpieces'
import { createTestQueryClient, createWrapper } from '@/test/queryClient'

vi.mock('@kinde-oss/kinde-auth-react', () => ({
  useKindeAuth: vi.fn().mockReturnValue({ getToken: vi.fn().mockResolvedValue('tok') }),
}))
vi.mock('@/lib/api', () => ({ apiFetch: vi.fn() }))

import { apiFetch } from '@/lib/api'

beforeEach(() => {
  vi.clearAllMocks()
})

describe('useUpdateTag', () => {
  it('PATCHes the new name and invalidates tag and artpiece caches', async () => {
    vi.mocked(apiFetch).mockResolvedValue({ ok: true, status: 200, json: async () => ({ id: 't1', name: 'outfits' }) } as Response)
    const client = createTestQueryClient()
    const invalidate = vi.spyOn(client, 'invalidateQueries')
    const { result } = renderHook(() => useUpdateTag(), { wrapper: createWrapper(client) })

    result.current.mutate({ id: 't1', name: 'outfits' })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    const [path, , init] = vi.mocked(apiFetch).mock.calls[0]
    expect(path).toBe('/tags/t1')
    expect(init?.method).toBe('PATCH')
    expect(JSON.parse(init?.body as string)).toEqual({ name: 'outfits' })
    expect(invalidate).toHaveBeenCalledWith({ queryKey: TAGS_QUERY_KEY })
    expect(invalidate).toHaveBeenCalledWith({ queryKey: ARTPIECES_QUERY_KEY })
  })

  it('rejects with TagConflictError on 409', async () => {
    vi.mocked(apiFetch).mockResolvedValue({ ok: false, status: 409 } as Response)
    const { result } = renderHook(() => useUpdateTag(), { wrapper: createWrapper() })

    result.current.mutate({ id: 't1', name: 'x' })

    await waitFor(() => expect(result.current.isError).toBe(true))
    expect(result.current.error).toBeInstanceOf(TagConflictError)
  })

  it('rejects with a generic error on other failures', async () => {
    vi.mocked(apiFetch).mockResolvedValue({ ok: false, status: 500 } as Response)
    const { result } = renderHook(() => useUpdateTag(), { wrapper: createWrapper() })

    result.current.mutate({ id: 't1', name: 'x' })

    await waitFor(() => expect(result.current.isError).toBe(true))
    expect(result.current.error?.message).toBe('Failed to rename tag')
  })
})
