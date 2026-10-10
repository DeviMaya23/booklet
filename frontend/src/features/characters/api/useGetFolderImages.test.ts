import { renderHook, waitFor } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import React from 'react'
import { useGetFolderImages } from './useGetFolderImages'

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

describe('useGetFolderImages', () => {
  it('fetches and returns image list when folderID is provided', async () => {
    vi.mocked(apiFetch).mockResolvedValue({
      ok: true,
      json: async () => ({
        images: [
          { image_id: 'img-1', thumbnail_url: 'https://cdn.example.com/thumb1.jpg' },
        ],
      }),
    } as Response)

    const { result } = renderHook(() => useGetFolderImages('folder-123'), { wrapper })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))

    expect(apiFetch).toHaveBeenCalledWith('/folders/folder-123/images', expect.any(Function))
    expect(result.current.data).toEqual([
      { image_id: 'img-1', thumbnail_url: 'https://cdn.example.com/thumb1.jpg' },
    ])
  })

  it('does not fetch when folderID is undefined', () => {
    renderHook(() => useGetFolderImages(undefined), { wrapper })

    expect(apiFetch).not.toHaveBeenCalled()
  })

  it('does not fetch when folderID is an empty string', () => {
    renderHook(() => useGetFolderImages(''), { wrapper })

    expect(apiFetch).not.toHaveBeenCalled()
  })
})
