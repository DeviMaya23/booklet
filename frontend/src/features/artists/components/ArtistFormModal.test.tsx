import { render, screen } from '@testing-library/react'
import { vi, describe, it, expect } from 'vitest'
import ArtistFormModal from './ArtistFormModal'
import { type Artist } from '../api/useArtists'

vi.mock('../api/useCreateArtist', () => ({
  useCreateArtist: () => ({ mutate: vi.fn(), isPending: false }),
}))
vi.mock('../api/useUpdateArtist', () => ({
  useUpdateArtist: () => ({ mutate: vi.fn(), isPending: false }),
}))
vi.mock('../api/useDeleteArtist', () => ({
  useDeleteArtist: () => ({ mutate: vi.fn(), isPending: false }),
}))
vi.mock('sonner', () => ({ toast: { success: vi.fn(), error: vi.fn() } }))

const artistWithLinks: Artist = {
  id: 'a1',
  name: 'Alice',
  notes: null,
  links: [
    { id: 'l1', url: 'https://twitter.com/alice', is_primary: true },
    { id: 'l2', url: 'https://portfolio.alice.com', is_primary: false },
  ],
  created_at: '',
  updated_at: '',
}

describe('ArtistFormModal', () => {
  it('pre-populates link URLs in edit mode', () => {
    render(<ArtistFormModal open artist={artistWithLinks} onOpenChange={() => {}} />)
    expect(screen.getByDisplayValue('https://twitter.com/alice')).toBeInTheDocument()
    expect(screen.getByDisplayValue('https://portfolio.alice.com')).toBeInTheDocument()
  })

  it('shows no link inputs in create mode', () => {
    render(<ArtistFormModal open onOpenChange={() => {}} />)
    expect(screen.queryByPlaceholderText('https://...')).toBeNull()
  })
})
