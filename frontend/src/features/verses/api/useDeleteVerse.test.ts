import { renderHook, act, waitFor } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import React from 'react'
import { useDeleteVerse } from './useDeleteVerse'

vi.mock('@kinde-oss/kinde-auth-react', () => ({
  useKindeAuth: vi.fn().mockReturnValue({ getToken: vi.fn().mockResolvedValue('tok') }),
}))

vi.mock('@/lib/api', () => ({
  apiFetch: vi.fn(),
}))

import { apiFetch } from '@/lib/api'

function makeWrapper() {
  const qc = new QueryClient({ defaultOptions: { mutations: { retry: false } } })
  const invalidateSpy = vi.spyOn(qc, 'invalidateQueries')
  function wrapper({ children }: { children: React.ReactNode }) {
    return React.createElement(QueryClientProvider, { client: qc }, children)
  }
  return { wrapper, invalidateSpy }
}

beforeEach(() => {
  vi.clearAllMocks()
})

describe('useDeleteVerse', () => {
  it('DELETEs /verses/:id on success', async () => {
    vi.mocked(apiFetch).mockResolvedValue({ ok: true } as Response)
    const { wrapper } = makeWrapper()

    const { result } = renderHook(() => useDeleteVerse(), { wrapper })
    act(() => { result.current.mutate('v-1') })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(apiFetch).toHaveBeenCalledWith('/verses/v-1', expect.any(Function), expect.objectContaining({ method: 'DELETE' }))
  })

  it('invalidates both the verses and characters query caches on success', async () => {
    vi.mocked(apiFetch).mockResolvedValue({ ok: true } as Response)
    const { wrapper, invalidateSpy } = makeWrapper()

    const { result } = renderHook(() => useDeleteVerse(), { wrapper })
    act(() => { result.current.mutate('v-1') })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(invalidateSpy).toHaveBeenCalledWith({ queryKey: ['verses'] })
    expect(invalidateSpy).toHaveBeenCalledWith({ queryKey: ['characters'] })
  })

  it('enters error state when the API returns a non-ok response', async () => {
    vi.mocked(apiFetch).mockResolvedValue({ ok: false } as Response)
    const { wrapper } = makeWrapper()

    const { result } = renderHook(() => useDeleteVerse(), { wrapper })
    act(() => { result.current.mutate('v-1') })

    await waitFor(() => expect(result.current.isError).toBe(true))
    expect((result.current.error as Error).message).toBe('Failed to delete verse')
  })
})
