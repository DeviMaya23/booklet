import { render, screen } from '@testing-library/react'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import { describe, it, expect, afterEach } from 'vitest'
import { setMaintenanceActive } from '../lib/maintenanceStore'
import AppLayout from './AppLayout'

afterEach(() => {
  setMaintenanceActive(false)
})

function renderAppLayout() {
  return render(
    <MemoryRouter initialEntries={['/app']}>
      <Routes>
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
})
