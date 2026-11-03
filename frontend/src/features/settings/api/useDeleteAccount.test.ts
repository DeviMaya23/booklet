import { renderHook, waitFor } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import React from 'react'
import { useDeleteAccount } from './useDeleteAccount'

vi.mock('@kinde-oss/kinde-auth-react', () => ({
  useKindeAuth: vi.fn().mockReturnValue({ getToken: vi.fn().mockResolvedValue('tok') }),
}))

vi.mock('@/lib/api', () => ({
  apiFetch: vi.fn(),
}))

import { apiFetch } from '@/lib/api'

function wrapper({ children }: { children: React.ReactNode }) {
  const qc = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })
  return React.createElement(QueryClientProvider, { client: qc }, children)
}

beforeEach(() => {
  vi.clearAllMocks()
})

describe('useDeleteAccount', () => {
  it('sends DELETE /me and resolves on 202', async () => {
    vi.mocked(apiFetch).mockResolvedValue({ ok: true, status: 202 } as Response)

    const { result } = renderHook(() => useDeleteAccount(), { wrapper })
    result.current.mutate()

    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(apiFetch).toHaveBeenCalledWith('/me', expect.any(Function), { method: 'DELETE' })
  })

  it('rejects with "Failed to delete account" on a non-2xx response', async () => {
    vi.mocked(apiFetch).mockResolvedValue({ ok: false, status: 502 } as Response)

    const { result } = renderHook(() => useDeleteAccount(), { wrapper })
    result.current.mutate()

    await waitFor(() => expect(result.current.isError).toBe(true))
    expect(result.current.error?.message).toBe('Failed to delete account')
  })
})
