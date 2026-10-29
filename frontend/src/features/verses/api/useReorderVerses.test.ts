import { renderHook, act, waitFor } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import React from 'react'
import { useReorderVerses } from './useReorderVerses'

vi.mock('@kinde-oss/kinde-auth-react', () => ({
  useKindeAuth: vi.fn().mockReturnValue({ getToken: vi.fn().mockResolvedValue('tok') }),
}))

vi.mock('@/lib/api', () => ({
  apiFetch: vi.fn(),
}))

import { apiFetch } from '@/lib/api'

function wrapper({ children }: { children: React.ReactNode }) {
  const qc = new QueryClient({ defaultOptions: { mutations: { retry: false } } })
  return React.createElement(QueryClientProvider, { client: qc }, children)
}

beforeEach(() => {
  vi.clearAllMocks()
})

describe('useReorderVerses', () => {
  it('POSTs the id array wrapped in { ids } to /verses/reorder', async () => {
    vi.mocked(apiFetch).mockResolvedValue({ ok: true } as Response)

    const { result } = renderHook(() => useReorderVerses(), { wrapper })
    act(() => { result.current.mutate(['v-2', 'v-1', 'v-3']) })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(apiFetch).toHaveBeenCalledWith('/verses/reorder', expect.any(Function), expect.objectContaining({
      method: 'POST',
      body: JSON.stringify({ ids: ['v-2', 'v-1', 'v-3'] }),
    }))
  })

  it('enters error state when the API returns a non-ok response', async () => {
    vi.mocked(apiFetch).mockResolvedValue({ ok: false } as Response)

    const { result } = renderHook(() => useReorderVerses(), { wrapper })
    act(() => { result.current.mutate(['v-1']) })

    await waitFor(() => expect(result.current.isError).toBe(true))
    expect((result.current.error as Error).message).toBe('Failed to reorder verses')
  })
})
