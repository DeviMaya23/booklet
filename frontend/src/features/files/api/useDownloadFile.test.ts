import { renderHook, act, waitFor } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import React from 'react'
import { useDownloadFile } from './useDownloadFile'

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

describe('useDownloadFile', () => {
  it('opens the download URL returned by the API', async () => {
    vi.mocked(apiFetch).mockResolvedValue({
      ok: true,
      json: async () => ({ download_url: 'https://cdn.example.com/file?sig=abc' }),
    } as Response)
    const openSpy = vi.spyOn(window, 'open').mockImplementation(() => null)

    const { result } = renderHook(() => useDownloadFile(), { wrapper })
    act(() => { result.current.mutate('file-123') })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(apiFetch).toHaveBeenCalledWith('/files/file-123/download', expect.any(Function))
    expect(openSpy).toHaveBeenCalledWith('https://cdn.example.com/file?sig=abc', '_blank')
  })

  it('throws when the API returns a non-ok response', async () => {
    vi.mocked(apiFetch).mockResolvedValue({ ok: false } as Response)

    const { result } = renderHook(() => useDownloadFile(), { wrapper })
    act(() => { result.current.mutate('file-456') })

    await waitFor(() => expect(result.current.isError).toBe(true))
    expect(result.current.error).toBeInstanceOf(Error)
    expect((result.current.error as Error).message).toBe('Failed to get download URL')
  })
})
