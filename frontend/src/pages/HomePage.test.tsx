import { render, screen } from '@testing-library/react'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import { vi, describe, it, expect } from 'vitest'
import HomePage from './HomePage'

vi.mock('@kinde-oss/kinde-auth-react', () => ({
  useKindeAuth: vi.fn(),
}))

import { useKindeAuth } from '@kinde-oss/kinde-auth-react'

function renderHomePage() {
  return render(
    <MemoryRouter initialEntries={['/']}>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/app" element={<div>App shell</div>} />
      </Routes>
    </MemoryRouter>,
  )
}

describe('HomePage', () => {
  it('renders login button when unauthenticated', () => {
    vi.mocked(useKindeAuth).mockReturnValue({
      isAuthenticated: false,
      isLoading: false,
      login: vi.fn(),
    } as unknown as ReturnType<typeof useKindeAuth>)

    renderHomePage()

    expect(screen.getByRole('button', { name: /log in/i })).toBeInTheDocument()
  })

  it('redirects to /app when already authenticated', () => {
    vi.mocked(useKindeAuth).mockReturnValue({
      isAuthenticated: true,
      isLoading: false,
      login: vi.fn(),
    } as unknown as ReturnType<typeof useKindeAuth>)

    renderHomePage()

    expect(screen.getByText('App shell')).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: /log in/i })).not.toBeInTheDocument()
  })
})
