import { renderHook, waitFor } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import React from 'react'
import { useDashboard } from './useDashboard'

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

const mockDashboardData = {
  recent_artpieces: [{ id: 'a1', title: 'Chibi', artist_name: 'Jill', thumbnail_url: null }],
  in_progress: [{ id: 'c1', title: 'Commission 1', artist_name: 'Jill', artist_link: null, status: 'wip', paid: false, last_contacted_at: null, created_at: new Date().toISOString() }],
  housekeeping: {
    commissions_no_artist: [],
    artpieces_no_artist: [],
    done_no_artpieces: [],
    artpieces_no_files: [],
  },
}

describe('useDashboard', () => {
  it('returns dashboard data on successful fetch', async () => {
    vi.mocked(apiFetch).mockResolvedValue({
      ok: true,
      json: async () => mockDashboardData,
    } as Response)

    const { result } = renderHook(() => useDashboard(), { wrapper })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(apiFetch).toHaveBeenCalledWith('/dashboard', expect.any(Function))
    expect(result.current.data?.recent_artpieces).toHaveLength(1)
    expect(result.current.data?.recent_artpieces[0].id).toBe('a1')
    expect(result.current.data?.in_progress[0].status).toBe('wip')
  })

  it('sets isError when the API returns a non-ok response', async () => {
    vi.mocked(apiFetch).mockResolvedValue({ ok: false } as Response)

    const { result } = renderHook(() => useDashboard(), { wrapper })

    await waitFor(() => expect(result.current.isError).toBe(true))
    expect((result.current.error as Error).message).toBe('Failed to fetch dashboard')
  })
})
