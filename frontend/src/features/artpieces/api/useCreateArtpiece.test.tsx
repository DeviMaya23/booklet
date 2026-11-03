import { renderHook, waitFor } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useCreateArtpiece } from './useCreateArtpiece'
import { CHARACTERS_QUERY_KEY } from '@/features/characters/api/useCharacters'
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

    result.current.mutate({ title: 't', characterIds: ['c1'], fileIds: [] })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(invalidate).toHaveBeenCalledWith({ queryKey: CHARACTERS_QUERY_KEY })
  })
})
