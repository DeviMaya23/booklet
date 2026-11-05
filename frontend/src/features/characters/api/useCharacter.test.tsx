import { renderHook, waitFor } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useCharacter } from './useCharacter'
import { CHARACTERS_QUERY_KEY, type Character } from './useCharacters'
import { createTestQueryClient, createWrapper } from '@/test/queryClient'

vi.mock('@kinde-oss/kinde-auth-react', () => ({
  useKindeAuth: vi.fn().mockReturnValue({ getToken: vi.fn().mockResolvedValue('tok') }),
}))

vi.mock('@/lib/api', () => ({
  apiFetch: vi.fn(),
}))

import { apiFetch } from '@/lib/api'

const tanelle: Character = {
  id: 'c-1',
  name: 'Tanelle',
  avatar_url: null,
  notes: 'Apothecary from Dyffrynne',
  is_public: false,
  verse_id: null,
  folders: [],
  created_at: '',
  updated_at: '',
}

beforeEach(() => {
  vi.clearAllMocks()
})

describe('useCharacter', () => {
  it('fetches GET /characters/:id and returns the character', async () => {
    vi.mocked(apiFetch).mockResolvedValue({ ok: true, json: async () => tanelle } as Response)

    const { result } = renderHook(() => useCharacter('c-1'), { wrapper: createWrapper() })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(apiFetch).toHaveBeenCalledWith('/characters/c-1', expect.any(Function))
    expect(result.current.data).toEqual(tanelle)
  })

  it('errors when the character is not found', async () => {
    vi.mocked(apiFetch).mockResolvedValue({ ok: false, status: 404, json: async () => ({}) } as Response)

    const { result } = renderHook(() => useCharacter('missing'), { wrapper: createWrapper() })

    await waitFor(() => expect(result.current.isError).toBe(true))
    expect(result.current.error?.message).toBe('Failed to fetch character')
  })

  it('does not fetch without an id', () => {
    renderHook(() => useCharacter(undefined), { wrapper: createWrapper() })

    expect(apiFetch).not.toHaveBeenCalled()
  })

  it('is refetched when CHARACTERS_QUERY_KEY is invalidated', async () => {
    vi.mocked(apiFetch).mockResolvedValue({ ok: true, json: async () => tanelle } as Response)
    const client = createTestQueryClient()
    const { result } = renderHook(() => useCharacter('c-1'), { wrapper: createWrapper(client) })
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(apiFetch).toHaveBeenCalledTimes(1)

    await client.invalidateQueries({ queryKey: CHARACTERS_QUERY_KEY })

    await waitFor(() => expect(apiFetch).toHaveBeenCalledTimes(2))
  })
})
