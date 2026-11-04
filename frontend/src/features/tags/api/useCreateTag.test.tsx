import { renderHook, waitFor } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useCreateTag, TagConflictError } from './useCreateTag'
import { TAGS_QUERY_KEY } from './useTags'
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
})

describe('useCreateTag', () => {
  it('posts the name and returns the created tag', async () => {
    vi.mocked(apiFetch).mockResolvedValue({ ok: true, status: 201, json: async () => ({ id: 't1', name: 'wip' }) } as Response)
    const { result } = renderHook(() => useCreateTag(), { wrapper: createWrapper() })

    result.current.mutate('wip')

    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data).toEqual({ id: 't1', name: 'wip' })
    const [path, , init] = vi.mocked(apiFetch).mock.calls[0]
    expect(path).toBe('/tags')
    expect(init?.method).toBe('POST')
    expect(JSON.parse(init?.body as string)).toEqual({ name: 'wip' })
  })

  it('invalidates the tags query on success', async () => {
    vi.mocked(apiFetch).mockResolvedValue({ ok: true, status: 201, json: async () => ({ id: 't1', name: 'wip' }) } as Response)
    const client = createTestQueryClient()
    const invalidate = vi.spyOn(client, 'invalidateQueries')
    const { result } = renderHook(() => useCreateTag(), { wrapper: createWrapper(client) })

    result.current.mutate('wip')

    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(invalidate).toHaveBeenCalledWith({ queryKey: TAGS_QUERY_KEY })
  })

  it('rejects with TagConflictError on 409', async () => {
    vi.mocked(apiFetch).mockResolvedValue({ ok: false, status: 409 } as Response)
    const { result } = renderHook(() => useCreateTag(), { wrapper: createWrapper() })

    result.current.mutate('wip')

    await waitFor(() => expect(result.current.isError).toBe(true))
    expect(result.current.error).toBeInstanceOf(TagConflictError)
  })

  it('rejects with a generic error on other failures', async () => {
    vi.mocked(apiFetch).mockResolvedValue({ ok: false, status: 500 } as Response)
    const { result } = renderHook(() => useCreateTag(), { wrapper: createWrapper() })

    result.current.mutate('wip')

    await waitFor(() => expect(result.current.isError).toBe(true))
    expect(result.current.error).not.toBeInstanceOf(TagConflictError)
    expect(result.current.error?.message).toBe('Failed to create tag')
  })
})
