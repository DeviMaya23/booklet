import { renderHook, waitFor, act } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useUpdateCharacter } from './useUpdateCharacter'
import { characterQueryKey } from './useCharacter'
import { CHARACTERS_QUERY_KEY } from './useCharacters'
import { createTestQueryClient, createWrapper } from '@/test/queryClient'
import { makeCharacter } from '@/test/characterFixtures'

vi.mock('@kinde-oss/kinde-auth-react', () => ({
  useKindeAuth: vi.fn().mockReturnValue({ getToken: vi.fn().mockResolvedValue('tok') }),
}))
vi.mock('@/lib/api', () => ({ apiFetch: vi.fn() }))

import { apiFetch } from '@/lib/api'

beforeEach(() => {
  vi.clearAllMocks()
})

const input = { id: 'c-1', name: 'New name', notes: null, folder_ids: [], verse_id: null }

describe('useUpdateCharacter', () => {
  it('writes the updated character into the detail cache as soon as the update resolves', async () => {
    const updated = makeCharacter({ name: 'New name' })
    vi.mocked(apiFetch).mockResolvedValue({ ok: true, json: async () => updated } as Response)
    const client = createTestQueryClient()
    client.setQueryData(characterQueryKey('c-1'), makeCharacter({ name: 'Old name' }))
    const { result } = renderHook(() => useUpdateCharacter(), { wrapper: createWrapper(client) })

    await act(async () => {
      await result.current.mutateAsync(input)
    })

    expect(client.getQueryData(characterQueryKey('c-1'))).toEqual(updated)
  })

  it('still invalidates the characters queries so the list refetches', async () => {
    vi.mocked(apiFetch).mockResolvedValue({ ok: true, json: async () => makeCharacter() } as Response)
    const client = createTestQueryClient()
    client.setQueryData(CHARACTERS_QUERY_KEY, [makeCharacter()])
    const { result } = renderHook(() => useUpdateCharacter(), { wrapper: createWrapper(client) })

    await act(async () => {
      await result.current.mutateAsync(input)
    })

    await waitFor(() => expect(client.getQueryState(CHARACTERS_QUERY_KEY)?.isInvalidated).toBe(true))
  })

  it('leaves the cache alone when the update fails', async () => {
    vi.mocked(apiFetch).mockResolvedValue({ ok: false, status: 500, json: async () => ({}) } as Response)
    const client = createTestQueryClient()
    const original = makeCharacter({ name: 'Old name' })
    client.setQueryData(characterQueryKey('c-1'), original)
    const { result } = renderHook(() => useUpdateCharacter(), { wrapper: createWrapper(client) })

    await act(async () => {
      await expect(result.current.mutateAsync(input)).rejects.toThrow('Failed to update character')
    })

    expect(client.getQueryData(characterQueryKey('c-1'))).toEqual(original)
  })
})
