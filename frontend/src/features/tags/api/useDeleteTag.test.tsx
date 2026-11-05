import { renderHook, waitFor } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useDeleteTag } from './useDeleteTag'
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

describe('useDeleteTag', () => {
  it('DELETEs the tag and invalidates tag and artpiece caches', async () => {
    vi.mocked(apiFetch).mockResolvedValue({ ok: true, status: 204 } as Response)
    const client = createTestQueryClient()
    const invalidate = vi.spyOn(client, 'invalidateQueries')
    const { result } = renderHook(() => useDeleteTag(), { wrapper: createWrapper(client) })

    result.current.mutate('t1')

    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    const [path, , init] = vi.mocked(apiFetch).mock.calls[0]
    expect(path).toBe('/tags/t1')
    expect(init?.method).toBe('DELETE')
    expect(invalidate).toHaveBeenCalledWith({ queryKey: TAGS_QUERY_KEY })
    expect(invalidate).toHaveBeenCalledWith({ queryKey: ARTPIECES_QUERY_KEY })
  })

  it('rejects when the request fails', async () => {
    vi.mocked(apiFetch).mockResolvedValue({ ok: false, status: 500 } as Response)
    const { result } = renderHook(() => useDeleteTag(), { wrapper: createWrapper() })

    result.current.mutate('t1')

    await waitFor(() => expect(result.current.isError).toBe(true))
    expect(result.current.error?.message).toBe('Failed to delete tag')
  })
})
