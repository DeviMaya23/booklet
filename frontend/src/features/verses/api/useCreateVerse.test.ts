import { renderHook, act, waitFor } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import React from 'react'
import { useCreateVerse } from './useCreateVerse'

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

describe('useCreateVerse', () => {
  it('POSTs to /verses and returns the created verse', async () => {
    const created = { id: 'v-1', name: 'Solaris', sort_order: 0, last_used_at: null, created_at: '2026-01-01T00:00:00Z', updated_at: '2026-01-01T00:00:00Z' }
    vi.mocked(apiFetch).mockResolvedValue({
      ok: true,
      json: async () => created,
    } as Response)

    const { result } = renderHook(() => useCreateVerse(), { wrapper })
    act(() => { result.current.mutate({ name: 'Solaris' }) })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(apiFetch).toHaveBeenCalledWith('/verses', expect.any(Function), expect.objectContaining({
      method: 'POST',
      body: JSON.stringify({ name: 'Solaris' }),
    }))
    expect(result.current.data).toEqual(created)
  })

  it('attaches the HTTP status to the thrown error on failure', async () => {
    vi.mocked(apiFetch).mockResolvedValue({ ok: false, status: 409 } as Response)

    const { result } = renderHook(() => useCreateVerse(), { wrapper })
    act(() => { result.current.mutate({ name: 'Solaris' }) })

    await waitFor(() => expect(result.current.isError).toBe(true))
    expect((result.current.error as Error & { status: number }).status).toBe(409)
  })
})
