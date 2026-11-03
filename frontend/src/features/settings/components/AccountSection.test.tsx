import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import AccountSection from './AccountSection'

const logout = vi.fn()
vi.mock('@kinde-oss/kinde-auth-react', () => ({
  useKindeAuth: () => ({
    user: { id: 'k1', givenName: 'Shatolla', familyName: 'Rose', email: 'me@example.com' },
    logout,
    getToken: vi.fn().mockResolvedValue('tok'),
  }),
}))

vi.mock('sonner', () => ({ toast: { error: vi.fn() } }))
vi.mock('@/lib/api', () => ({ apiFetch: vi.fn() }))

import { toast } from 'sonner'
import { apiFetch } from '@/lib/api'

function renderSection() {
  const qc = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })
  return render(
    <QueryClientProvider client={qc}>
      <AccountSection />
    </QueryClientProvider>,
  )
}

beforeEach(() => {
  vi.clearAllMocks()
})

describe('AccountSection', () => {
  it('shows the full name and email', () => {
    renderSection()

    expect(screen.getByText('Shatolla Rose')).toBeInTheDocument()
    expect(screen.getByText('me@example.com')).toBeInTheDocument()
  })

  it('logs out when Log out is clicked', async () => {
    renderSection()

    await userEvent.click(screen.getByRole('button', { name: 'Log out' }))

    expect(logout).toHaveBeenCalledTimes(1)
  })

  it('reveals the confirmation prompt with a disabled confirm button', async () => {
    renderSection()

    await userEvent.click(screen.getByRole('button', { name: 'Delete account…' }))

    expect(screen.getByLabelText('Account email confirmation')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Delete account' })).toBeDisabled()
  })

  it('enables confirm only once the typed email matches, ignoring case and whitespace', async () => {
    renderSection()
    await userEvent.click(screen.getByRole('button', { name: 'Delete account…' }))
    const input = screen.getByLabelText('Account email confirmation')

    await userEvent.type(input, 'other@example.com')
    expect(screen.getByRole('button', { name: 'Delete account' })).toBeDisabled()

    await userEvent.clear(input)
    await userEvent.type(input, ' ME@example.com ')
    expect(screen.getByRole('button', { name: 'Delete account' })).toBeEnabled()
  })

  it('Cancel hides the prompt, clears the input and shows the opening button again', async () => {
    renderSection()
    await userEvent.click(screen.getByRole('button', { name: 'Delete account…' }))
    await userEvent.type(screen.getByLabelText('Account email confirmation'), 'me@example.com')

    await userEvent.click(screen.getByRole('button', { name: 'Cancel' }))

    expect(screen.queryByLabelText('Account email confirmation')).not.toBeInTheDocument()
    await userEvent.click(screen.getByRole('button', { name: 'Delete account…' }))
    expect(screen.getByLabelText('Account email confirmation')).toHaveValue('')
  })

  it('calls DELETE /me and logs out on success', async () => {
    vi.mocked(apiFetch).mockResolvedValue({ ok: true, status: 202 } as Response)
    renderSection()
    await userEvent.click(screen.getByRole('button', { name: 'Delete account…' }))
    await userEvent.type(screen.getByLabelText('Account email confirmation'), 'me@example.com')

    await userEvent.click(screen.getByRole('button', { name: 'Delete account' }))

    await waitFor(() => expect(logout).toHaveBeenCalledTimes(1))
    expect(apiFetch).toHaveBeenCalledWith('/me', expect.any(Function), { method: 'DELETE' })
    expect(toast.error).not.toHaveBeenCalled()
  })

  it('shows an error toast, stays signed in and keeps the prompt when deletion fails', async () => {
    vi.mocked(apiFetch).mockResolvedValue({ ok: false, status: 502 } as Response)
    renderSection()
    await userEvent.click(screen.getByRole('button', { name: 'Delete account…' }))
    await userEvent.type(screen.getByLabelText('Account email confirmation'), 'me@example.com')

    await userEvent.click(screen.getByRole('button', { name: 'Delete account' }))

    await waitFor(() => expect(toast.error).toHaveBeenCalledWith('Failed to delete account'))
    expect(logout).not.toHaveBeenCalled()
    expect(screen.getByLabelText('Account email confirmation')).toBeInTheDocument()
  })
})
