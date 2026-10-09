import { render, screen } from '@testing-library/react'
import { vi, describe, it, expect } from 'vitest'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import CommissionsTable from './CommissionsTable'
import { type Commission } from '../api/useCommissions'

vi.mock('@kinde-oss/kinde-auth-react', () => ({
  useKindeAuth: () => ({ getToken: vi.fn().mockResolvedValue('tok') }),
}))

vi.mock('../api/usePatchCommission', () => ({
  usePatchCommission: () => ({ mutateAsync: vi.fn(), isPending: false }),
}))

vi.mock('sonner', () => ({ toast: { success: vi.fn(), error: vi.fn() } }))

const baseCommission: Commission = {
  id: 'c1',
  title: 'Test Commission',
  artist_id: 'a1',
  artist_name: 'Alice',
  artist_link: null,
  status: 'waitlist',
  price: null,
  paid: false,
  paid_date: null,
  finish_date: null,
  last_contacted_at: null,
  notes: null,
  artpieces: [],
  created_at: '2024-01-01T00:00:00Z',
  updated_at: '2024-01-01T00:00:00Z',
}

function wrapper({ children }: { children: React.ReactNode }) {
  const qc = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  return <QueryClientProvider client={qc}>{children}</QueryClientProvider>
}

describe('CommissionsTable — artist link', () => {
  it('shows a link when the commission has an artist_link', () => {
    const commission = { ...baseCommission, artist_link: 'https://artist.example.com' }
    render(
      <CommissionsTable commissions={[commission]} view="active" onEdit={() => {}} onDelete={() => {}} />,
      { wrapper },
    )
    expect(screen.getByRole('link', { name: /open artist link/i })).toHaveAttribute(
      'href',
      'https://artist.example.com',
    )
  })

  it('shows no artist link when artist_link is null', () => {
    const commission = { ...baseCommission, artist_link: null }
    render(
      <CommissionsTable commissions={[commission]} view="active" onEdit={() => {}} onDelete={() => {}} />,
      { wrapper },
    )
    expect(screen.queryByRole('link', { name: /open artist link/i })).toBeNull()
  })
})
