import { renderHook, waitFor } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useTags, tagsQueryKey, TAGS_QUERY_KEY } from './useTags'
import { createTestQueryClient, createWrapper } from '@/test/queryClient'

vi.mock('@kinde-oss/kinde-auth-react', () => ({
  useKindeAuth: vi.fn().mockReturnValue({ getToken: vi.fn().mockResolvedValue('tok') }),
}))

vi.mock('@/lib/api', () => ({
  apiFetch: vi.fn(),
}))

import { apiFetch } from '@/lib/api'

beforeEach(() => {
  vi.clearAllMocks()
  vi.mocked(apiFetch).mockResolvedValue({ ok: true, json: async () => [{ id: 't1', name: 'wip' }] } as Response)
})

describe('useTags', () => {
  it('requests /tags with no sort param by default', async () => {
    const { result } = renderHook(() => useTags(), { wrapper: createWrapper() })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))

    expect(apiFetch).toHaveBeenCalledWith('/tags', expect.any(Function))
    expect(result.current.data).toEqual([{ id: 't1', name: 'wip' }])
  })

  it('requests /tags?sort=recent for the recent variant', async () => {
    const { result } = renderHook(() => useTags({ sort: 'recent' }), { wrapper: createWrapper() })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))

    expect(apiFetch).toHaveBeenCalledWith('/tags?sort=recent', expect.any(Function))
  })

  it('requests with_usage=true, combined with sort', async () => {
    const { result } = renderHook(() => useTags({ sort: 'recent', withUsage: true }), { wrapper: createWrapper() })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))

    expect(apiFetch).toHaveBeenCalledWith('/tags?sort=recent&with_usage=true', expect.any(Function))
  })

  it('requests only with_usage=true when no sort is given', async () => {
    const { result } = renderHook(() => useTags({ withUsage: true }), { wrapper: createWrapper() })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))

    expect(apiFetch).toHaveBeenCalledWith('/tags?with_usage=true', expect.any(Function))
  })

  it('keeps usage and non-usage lists in separate cache entries', async () => {
    const client = createTestQueryClient()
    const plain = renderHook(() => useTags(), { wrapper: createWrapper(client) })
    await waitFor(() => expect(plain.result.current.isSuccess).toBe(true))
    expect(apiFetch).toHaveBeenCalledTimes(1)

    const usage = renderHook(() => useTags({ withUsage: true }), { wrapper: createWrapper(client) })
    await waitFor(() => expect(usage.result.current.isSuccess).toBe(true))

    expect(apiFetch).toHaveBeenCalledTimes(2)
    expect(tagsQueryKey()).toEqual(TAGS_QUERY_KEY)
    expect(tagsQueryKey('recent', true)).toEqual([...TAGS_QUERY_KEY, 'recent', 'usage'])
  })

  it('does not request tags when disabled', async () => {
    const { result } = renderHook(() => useTags({ sort: 'recent', enabled: false }), { wrapper: createWrapper() })

    await new Promise((r) => setTimeout(r, 20))

    expect(apiFetch).not.toHaveBeenCalled()
    expect(result.current.fetchStatus).toBe('idle')
  })

  it('invalidating TAGS_QUERY_KEY refetches the recent variant', async () => {
    const client = createTestQueryClient()
    const { result } = renderHook(() => useTags({ sort: 'recent' }), { wrapper: createWrapper(client) })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(apiFetch).toHaveBeenCalledTimes(1)

    await client.invalidateQueries({ queryKey: TAGS_QUERY_KEY })

    await waitFor(() => expect(apiFetch).toHaveBeenCalledTimes(2))
  })

  it('surfaces a failed request as an error', async () => {
    vi.mocked(apiFetch).mockResolvedValue({ ok: false, status: 500 } as Response)
    const { result } = renderHook(() => useTags(), { wrapper: createWrapper() })

    await waitFor(() => expect(result.current.isError).toBe(true))
  })
})
