import { renderHook, waitFor } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useMergeTag } from './useMergeTag'
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

describe('useMergeTag', () => {
  it('POSTs target_id to the source merge path and invalidates tag and artpiece caches', async () => {
    vi.mocked(apiFetch).mockResolvedValue({ ok: true, status: 204 } as Response)
    const client = createTestQueryClient()
    const invalidate = vi.spyOn(client, 'invalidateQueries')
    const { result } = renderHook(() => useMergeTag(), { wrapper: createWrapper(client) })

    result.current.mutate({ sourceId: 'src', targetId: 'tgt' })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    const [path, , init] = vi.mocked(apiFetch).mock.calls[0]
    expect(path).toBe('/tags/src/merge')
    expect(init?.method).toBe('POST')
    expect(JSON.parse(init?.body as string)).toEqual({ target_id: 'tgt' })
    expect(invalidate).toHaveBeenCalledWith({ queryKey: TAGS_QUERY_KEY })
    expect(invalidate).toHaveBeenCalledWith({ queryKey: ARTPIECES_QUERY_KEY })
  })

  it('maps 404 to a "no longer exists" error and still refreshes the tag list', async () => {
    vi.mocked(apiFetch).mockResolvedValue({ ok: false, status: 404 } as Response)
    const client = createTestQueryClient()
    const invalidate = vi.spyOn(client, 'invalidateQueries')
    const { result } = renderHook(() => useMergeTag(), { wrapper: createWrapper(client) })

    result.current.mutate({ sourceId: 'src', targetId: 'tgt' })

    await waitFor(() => expect(result.current.isError).toBe(true))
    expect(result.current.error?.message).toBe('Tag no longer exists')
    expect(invalidate).toHaveBeenCalledWith({ queryKey: TAGS_QUERY_KEY })
  })

  it('maps other failures to a generic error', async () => {
    vi.mocked(apiFetch).mockResolvedValue({ ok: false, status: 422 } as Response)
    const { result } = renderHook(() => useMergeTag(), { wrapper: createWrapper() })

    result.current.mutate({ sourceId: 'src', targetId: 'tgt' })

    await waitFor(() => expect(result.current.isError).toBe(true))
    expect(result.current.error?.message).toBe('Failed to merge tags')
  })
})
