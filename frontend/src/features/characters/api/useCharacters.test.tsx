import { renderHook, waitFor } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useCharacters, CHARACTERS_QUERY_KEY } from './useCharacters'
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
  vi.mocked(apiFetch).mockResolvedValue({ ok: true, json: async () => [] } as Response)
})

describe('useCharacters', () => {
  it('requests /characters with no sort param by default', async () => {
    const { result } = renderHook(() => useCharacters(), { wrapper: createWrapper() })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))

    expect(apiFetch).toHaveBeenCalledWith('/characters', expect.any(Function))
  })

  it('requests /characters?sort=recent for the recent variant', async () => {
    const { result } = renderHook(() => useCharacters({ sort: 'recent' }), { wrapper: createWrapper() })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))

    expect(apiFetch).toHaveBeenCalledWith('/characters?sort=recent', expect.any(Function))
  })

  it('invalidating CHARACTERS_QUERY_KEY refetches both variants', async () => {
    const client = createTestQueryClient()
    const wrapper = createWrapper(client)
    const { result } = renderHook(
      () => ({ plain: useCharacters(), recent: useCharacters({ sort: 'recent' }) }),
      { wrapper },
    )
    await waitFor(() => {
      expect(result.current.plain.isSuccess).toBe(true)
      expect(result.current.recent.isSuccess).toBe(true)
    })
    expect(apiFetch).toHaveBeenCalledTimes(2)

    await client.invalidateQueries({ queryKey: CHARACTERS_QUERY_KEY })

    await waitFor(() => expect(apiFetch).toHaveBeenCalledTimes(4))
  })
})
