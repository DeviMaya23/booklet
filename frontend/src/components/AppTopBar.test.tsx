import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { ThemeProvider } from '@/hooks/useTheme'
import AppTopBar from './AppTopBar'

const logout = vi.fn()
const mockUser = vi.fn()
vi.mock('@kinde-oss/kinde-auth-react', () => ({
  useKindeAuth: () => ({
    user: mockUser(),
    logout,
    getToken: vi.fn().mockResolvedValue('tok'),
  }),
}))

const mockIsMobile = vi.fn()
vi.mock('@/hooks/use-mobile', () => ({ useIsMobile: () => mockIsMobile() }))

const fullUser = {
  id: 'k1',
  givenName: 'Shatolla',
  familyName: 'Rose',
  email: 'me@example.com',
  picture: 'https://example.com/me.png',
}

function renderTopBar() {
  const qc = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  return render(
    <QueryClientProvider client={qc}>
      <ThemeProvider>
        <AppTopBar />
      </ThemeProvider>
    </QueryClientProvider>,
  )
}

beforeEach(() => {
  vi.clearAllMocks()
  mockUser.mockReturnValue(fullUser)
  mockIsMobile.mockReturnValue(false)
})

describe('AppTopBar', () => {
  it('opens the dropdown on avatar click without logging out', async () => {
    renderTopBar()

    await userEvent.click(screen.getByRole('button', { name: 'Account menu' }))

    expect(await screen.findByText('me@example.com')).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: 'Settings' })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: 'Sign out' })).toBeInTheDocument()
    expect(logout).not.toHaveBeenCalled()
  })

  it('logs out when Sign out is selected', async () => {
    renderTopBar()
    await userEvent.click(screen.getByRole('button', { name: 'Account menu' }))

    await userEvent.click(await screen.findByRole('menuitem', { name: 'Sign out' }))

    expect(logout).toHaveBeenCalledTimes(1)
  })

  it('opens the Settings modal on the Account section', async () => {
    renderTopBar()
    await userEvent.click(screen.getByRole('button', { name: 'Account menu' }))

    await userEvent.click(await screen.findByRole('menuitem', { name: 'Settings' }))

    expect(await screen.findByRole('heading', { name: 'Account' })).toBeInTheDocument()
  })

  it('reopens Settings on Account after it was left on App', async () => {
    renderTopBar()
    await userEvent.click(screen.getByRole('button', { name: 'Account menu' }))
    await userEvent.click(await screen.findByRole('menuitem', { name: 'Settings' }))
    await userEvent.click(await screen.findByRole('button', { name: 'App' }))
    await userEvent.click(screen.getByRole('button', { name: 'Close' }))

    await userEvent.click(screen.getByRole('button', { name: 'Account menu' }))
    await userEvent.click(await screen.findByRole('menuitem', { name: 'Settings' }))

    expect(await screen.findByRole('heading', { name: 'Account' })).toBeInTheDocument()
  })

  it('hides Settings on mobile but keeps email and Sign out', async () => {
    mockIsMobile.mockReturnValue(true)
    renderTopBar()

    await userEvent.click(screen.getByRole('button', { name: 'Account menu' }))

    expect(await screen.findByRole('menuitem', { name: 'Sign out' })).toBeInTheDocument()
    expect(screen.queryByRole('menuitem', { name: 'Settings' })).not.toBeInTheDocument()
    expect(screen.getByText('me@example.com')).toBeInTheDocument()
  })

  it('shows the picture when the user has one', () => {
    renderTopBar()

    expect(screen.getByRole('img', { name: 'Shatolla' })).toHaveAttribute('src', fullUser.picture)
  })

  it('falls back to the given-name initial without a picture', () => {
    mockUser.mockReturnValue({ ...fullUser, picture: undefined })
    renderTopBar()

    expect(screen.getByText('S')).toBeInTheDocument()
  })

  it('falls back to "?" when there is no picture or name', () => {
    mockUser.mockReturnValue({ id: 'k1' })
    renderTopBar()

    expect(screen.getByText('?')).toBeInTheDocument()
  })
})
