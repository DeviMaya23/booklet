import { renderHook, waitFor } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import React from 'react'
import { useVerses } from './useVerses'

vi.mock('@kinde-oss/kinde-auth-react', () => ({
  useKindeAuth: vi.fn().mockReturnValue({ getToken: vi.fn().mockResolvedValue('tok') }),
}))

vi.mock('@/lib/api', () => ({
  apiFetch: vi.fn(),
}))

import { apiFetch } from '@/lib/api'

function wrapper({ children }: { children: React.ReactNode }) {
  const qc = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  return React.createElement(QueryClientProvider, { client: qc }, children)
}

beforeEach(() => {
  vi.clearAllMocks()
})

describe('useVerses', () => {
  it('returns the parsed verse list on success', async () => {
    const verses = [
      { id: 'v-1', name: 'Solaris', sort_order: 0, last_used_at: null, created_at: '2026-01-01T00:00:00Z', updated_at: '2026-01-01T00:00:00Z' },
    ]
    vi.mocked(apiFetch).mockResolvedValue({
      ok: true,
      json: async () => verses,
    } as Response)

    const { result } = renderHook(() => useVerses(), { wrapper })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(apiFetch).toHaveBeenCalledWith('/verses', expect.any(Function))
    expect(result.current.data).toEqual(verses)
  })

  it('enters error state when the API returns a non-ok response', async () => {
    vi.mocked(apiFetch).mockResolvedValue({ ok: false } as Response)

    const { result } = renderHook(() => useVerses(), { wrapper })

    await waitFor(() => expect(result.current.isError).toBe(true))
    expect((result.current.error as Error).message).toBe('Failed to fetch verses')
  })
})
