import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { describe, it, expect, vi } from 'vitest'
import { ThemeProvider } from '@/hooks/useTheme'
import SettingsModal from './SettingsModal'

vi.mock('@kinde-oss/kinde-auth-react', () => ({
  useKindeAuth: () => ({
    user: { id: 'k1', givenName: 'Shatolla', familyName: 'Rose', email: 'me@example.com' },
    logout: vi.fn(),
    getToken: vi.fn().mockResolvedValue('tok'),
  }),
}))

function renderModal(onOpenChange = vi.fn()) {
  const qc = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  render(
    <QueryClientProvider client={qc}>
      <ThemeProvider>
        <SettingsModal open onOpenChange={onOpenChange} />
      </ThemeProvider>
    </QueryClientProvider>,
  )
  return onOpenChange
}

describe('SettingsModal', () => {
  it('opens on the Account section', () => {
    renderModal()

    expect(screen.getByRole('heading', { name: 'Account' })).toBeInTheDocument()
    expect(screen.getByText('me@example.com')).toBeInTheDocument()
  })

  it('switches to the App section from the left navigation', async () => {
    renderModal()

    await userEvent.click(screen.getByRole('button', { name: 'App' }))

    expect(screen.getByRole('heading', { name: 'App' })).toBeInTheDocument()
    expect(screen.getByRole('radio', { name: 'Parchment theme' })).toBeInTheDocument()
  })

  it('offers only the Account and App sections', () => {
    renderModal()

    const nav = screen.getByRole('navigation', { name: 'Settings sections' })
    const labels = Array.from(nav.querySelectorAll('button')).map((b) => b.textContent)

    expect(labels).toEqual(['Account', 'App'])
  })

  it('asks to close when the close button is clicked', async () => {
    const onOpenChange = renderModal()

    await userEvent.click(screen.getByRole('button', { name: 'Close' }))

    expect(onOpenChange).toHaveBeenCalledWith(false, expect.anything())
  })
})
