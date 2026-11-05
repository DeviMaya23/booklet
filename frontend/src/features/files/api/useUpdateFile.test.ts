import { renderHook, act, waitFor } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import React from 'react'
import { useUpdateFile } from './useUpdateFile'

const qc = new QueryClient({ defaultOptions: { mutations: { retry: false } } })

vi.mock('@kinde-oss/kinde-auth-react', () => ({
  useKindeAuth: vi.fn().mockReturnValue({ getToken: vi.fn().mockResolvedValue('tok') }),
}))

vi.mock('@/lib/api', () => ({
  apiFetch: vi.fn(),
}))

import { apiFetch } from '@/lib/api'

function wrapper({ children }: { children: React.ReactNode }) {
  return React.createElement(QueryClientProvider, { client: qc }, children)
}

function sentBody(): Record<string, unknown> {
  const init = vi.mocked(apiFetch).mock.calls[0][2] as RequestInit
  return JSON.parse(init.body as string)
}

beforeEach(() => {
  vi.clearAllMocks()
  vi.restoreAllMocks()
  vi.mocked(apiFetch).mockResolvedValue({ ok: true } as Response)
})

describe('useUpdateFile', () => {
  it('omits is_visible from the body when not provided', async () => {
    const { result } = renderHook(() => useUpdateFile(), { wrapper })
    act(() => { result.current.mutate({ id: 'f1', name: 'a', notes: 'b' }) })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(sentBody()).toEqual({ name: 'a', notes: 'b' })
  })

  it('omits name and notes from the body when not provided', async () => {
    const { result } = renderHook(() => useUpdateFile(), { wrapper })
    act(() => { result.current.mutate({ id: 'f1', name: 'a' }) })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(sentBody()).toEqual({ name: 'a' })
  })

  it('includes is_visible in the body when provided, including false', async () => {
    const { result } = renderHook(() => useUpdateFile(), { wrapper })
    act(() => { result.current.mutate({ id: 'f1', isVisible: false }) })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(sentBody()).toEqual({ is_visible: false })
  })

  it('sends null instead of an empty string when name or notes are cleared', async () => {
    const { result } = renderHook(() => useUpdateFile(), { wrapper })
    act(() => { result.current.mutate({ id: 'f1', name: '', notes: '' }) })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(sentBody()).toEqual({ name: null, notes: null })
  })

  it('invalidates the artpiece query on success when artpieceId is given', async () => {
    const invalidate = vi.spyOn(qc, 'invalidateQueries')
    const { result } = renderHook(() => useUpdateFile(), { wrapper })
    act(() => { result.current.mutate({ id: 'f1', artpieceId: 'ap1', isVisible: false }) })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(invalidate).toHaveBeenCalledWith({ queryKey: ['artpieces', 'ap1'] })
  })

  it('does not invalidate any artpiece query when artpieceId is omitted', async () => {
    const invalidate = vi.spyOn(qc, 'invalidateQueries')
    const { result } = renderHook(() => useUpdateFile(), { wrapper })
    act(() => { result.current.mutate({ id: 'f1', name: 'a' }) })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(invalidate).not.toHaveBeenCalledWith({ queryKey: ['artpieces', 'ap1'] })
  })
})
