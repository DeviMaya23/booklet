import { render, screen } from '@testing-library/react'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import { describe, it, expect, afterEach, vi } from 'vitest'
import { setMaintenanceActive } from '../lib/maintenanceStore'
import { setSessionExpired } from '../lib/sessionExpiredStore'
import AppLayout from './AppLayout'

vi.mock('sonner', () => ({ toast: { error: vi.fn() } }))

afterEach(() => {
  setMaintenanceActive(false)
  setSessionExpired(false)
})

function renderAppLayout(initialPath = '/app') {
  return render(
    <MemoryRouter initialEntries={[initialPath]}>
      <Routes>
        <Route path="/" element={<div>Login page</div>} />
        <Route element={<AppLayout />}>
          <Route path="/app" element={<div>App content</div>} />
        </Route>
      </Routes>
    </MemoryRouter>,
  )
}

describe('AppLayout', () => {
  it('renders outlet content when maintenance is inactive', () => {
    setMaintenanceActive(false)

    renderAppLayout()

    expect(screen.getByText('App content')).toBeInTheDocument()
    expect(screen.queryByTestId('maintenance-page')).not.toBeInTheDocument()
  })

  it('renders maintenance page when maintenance is active', () => {
    setMaintenanceActive(true)

    renderAppLayout()

    expect(screen.getByTestId('maintenance-page')).toBeInTheDocument()
    expect(screen.queryByText('App content')).not.toBeInTheDocument()
  })

  it('redirects to / and shows toast when session is expired', async () => {
    const { toast } = await import('sonner')
    setSessionExpired(true)

    renderAppLayout()

    expect(toast.error).toHaveBeenCalledWith('Please log back in')
    expect(toast.error).toHaveBeenCalledTimes(1)
    expect(screen.getByText('Login page')).toBeInTheDocument()
    expect(screen.queryByText('App content')).not.toBeInTheDocument()
  })

  it('resets the session-expired flag after redirect', async () => {
    const { useSessionExpired } = await import('../lib/sessionExpiredStore')
    setSessionExpired(true)

    renderAppLayout()

    const { renderHook } = await import('@testing-library/react')
    const { result } = renderHook(() => useSessionExpired())
    expect(result.current).toBe(false)
  })
})
