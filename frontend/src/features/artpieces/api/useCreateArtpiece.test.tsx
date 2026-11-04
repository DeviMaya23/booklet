import { renderHook, waitFor } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useCreateArtpiece } from './useCreateArtpiece'
import { CHARACTERS_QUERY_KEY } from '@/features/characters/api/useCharacters'
import { TAGS_QUERY_KEY } from '@/features/tags/api/useTags'
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
  vi.mocked(apiFetch).mockResolvedValue({ ok: true, json: async () => ({}) } as Response)
})

describe('useCreateArtpiece', () => {
  it('invalidates the characters query on success', async () => {
    const client = createTestQueryClient()
    const invalidate = vi.spyOn(client, 'invalidateQueries')
    const { result } = renderHook(() => useCreateArtpiece(), { wrapper: createWrapper(client) })

    result.current.mutate({ title: 't', characterIds: ['c1'], tagIds: [], fileIds: [] })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(invalidate).toHaveBeenCalledWith({ queryKey: CHARACTERS_QUERY_KEY })
  })

  it('invalidates the tags query on success', async () => {
    const client = createTestQueryClient()
    const invalidate = vi.spyOn(client, 'invalidateQueries')
    const { result } = renderHook(() => useCreateArtpiece(), { wrapper: createWrapper(client) })

    result.current.mutate({ title: 't', characterIds: [], tagIds: ['t1'], fileIds: [] })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(invalidate).toHaveBeenCalledWith({ queryKey: TAGS_QUERY_KEY })
  })

  it('sends tag_ids in the request body', async () => {
    const { result } = renderHook(() => useCreateArtpiece(), { wrapper: createWrapper() })

    result.current.mutate({ title: 't', characterIds: [], tagIds: ['t1', 't2'], fileIds: [] })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    const body = JSON.parse(vi.mocked(apiFetch).mock.calls[0][2]?.body as string)
    expect(body.tag_ids).toEqual(['t1', 't2'])
  })

  it('sends an empty tag_ids array when no tags are selected', async () => {
    const { result } = renderHook(() => useCreateArtpiece(), { wrapper: createWrapper() })

    result.current.mutate({ title: 't', characterIds: [], tagIds: [], fileIds: [] })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    const body = JSON.parse(vi.mocked(apiFetch).mock.calls[0][2]?.body as string)
    expect(body.tag_ids).toEqual([])
  })
})
