import { vi, describe, it, expect, beforeEach } from 'vitest'
import { renderHook } from '@testing-library/react'
import { apiFetch } from './api'
import { useMaintenanceActive, setMaintenanceActive } from './maintenanceStore'

describe('apiFetch', () => {
  beforeEach(() => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response()))
    localStorage.clear()
    setMaintenanceActive(false)
  })

  it('attaches Authorization header with bearer token', async () => {
    const getToken = vi.fn().mockResolvedValue('test-token')

    await apiFetch('/images', getToken)

    expect(fetch).toHaveBeenCalledWith(
      expect.stringContaining('/images'),
      expect.objectContaining({
        headers: expect.objectContaining({
          Authorization: 'Bearer test-token',
        }),
      }),
    )
  })

  it('sends request without Authorization header when token is undefined', async () => {
    const getToken = vi.fn().mockResolvedValue(undefined)

    await apiFetch('/images', getToken)

    const [, options] = vi.mocked(fetch).mock.calls[0]
    expect((options?.headers as Record<string, string>)?.Authorization).toBeUndefined()
  })

  it('attaches X-Booklet-Bypass header when a bypass token is stored', async () => {
    localStorage.setItem('booklet-maintenance-bypass', 'bypass-token')
    const getToken = vi.fn().mockResolvedValue('test-token')

    await apiFetch('/images', getToken)

    expect(fetch).toHaveBeenCalledWith(
      expect.stringContaining('/images'),
      expect.objectContaining({
        headers: expect.objectContaining({
          'X-Booklet-Bypass': 'bypass-token',
        }),
      }),
    )
  })

  it('omits X-Booklet-Bypass header when no bypass token is stored', async () => {
    const getToken = vi.fn().mockResolvedValue('test-token')

    await apiFetch('/images', getToken)

    const [, options] = vi.mocked(fetch).mock.calls[0]
    expect((options?.headers as Record<string, string>)?.['X-Booklet-Bypass']).toBeUndefined()
  })

  it('sets maintenance state to true when response includes X-Booklet-Maintenance: true', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response(null, {
      headers: { 'X-Booklet-Maintenance': 'true' },
    })))
    const getToken = vi.fn().mockResolvedValue('test-token')

    await apiFetch('/images', getToken)

    const { result } = renderHook(() => useMaintenanceActive())
    expect(result.current).toBe(true)
  })

  it('sets maintenance state to false when response does not include X-Booklet-Maintenance', async () => {
    setMaintenanceActive(true)
    const getToken = vi.fn().mockResolvedValue('test-token')

    await apiFetch('/images', getToken)

    const { result } = renderHook(() => useMaintenanceActive())
    expect(result.current).toBe(false)
  })
})
