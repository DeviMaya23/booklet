import { render, screen } from '@testing-library/react'
import { vi, describe, it, expect } from 'vitest'
import ArtistsList from './ArtistsList'
import { type Artist } from '../api/useArtists'

vi.mock('sonner', () => ({ toast: { success: vi.fn() } }))

const noop = () => {}

function makeArtist(overrides: Partial<Artist>): Artist {
  return {
    id: 'a1',
    name: 'Test Artist',
    notes: null,
    links: [],
    created_at: '',
    updated_at: '',
    ...overrides,
  }
}

describe('ArtistsList', () => {
  it('enables the copy link button when the artist has a primary link', () => {
    const artist = makeArtist({
      links: [{ id: 'l1', url: 'https://example.com', is_primary: true }],
    })
    render(<ArtistsList artists={[artist]} onEditClick={noop} />)
    expect(screen.getByTitle('Copy link')).not.toBeDisabled()
  })

  it('disables the copy link button when the artist has no links', () => {
    const artist = makeArtist({ links: [] })
    render(<ArtistsList artists={[artist]} onEditClick={noop} />)
    expect(screen.getByTitle('No link available')).toBeDisabled()
  })

  it('renders nothing when the artists list is empty', () => {
    const { container } = render(<ArtistsList artists={[]} onEditClick={noop} />)
    expect(container.firstChild).toBeEmptyDOMElement()
  })
})
