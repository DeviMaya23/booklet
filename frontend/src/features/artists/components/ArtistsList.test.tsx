import { render, screen } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import ArtistsList from './ArtistsList'
import { type Artist } from '../api/useArtists'

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
  it('renders nothing when the artists list is empty', () => {
    const { container } = render(<ArtistsList artists={[]} onEditClick={noop} />)
    expect(container.firstChild).toBeNull()
  })

  it('renders primary link chip first regardless of array order from API', () => {
    const artist = makeArtist({
      links: [
        { id: 'l1', url: 'https://x.com/artist', is_primary: false },
        { id: 'l2', url: 'https://bsky.app/artist', is_primary: true },
      ],
    })
    render(<ArtistsList artists={[artist]} onEditClick={noop} />)
    const chips = screen.getAllByRole('link')
    expect(chips[0]).toHaveAttribute('href', 'https://bsky.app/artist')
    expect(chips[1]).toHaveAttribute('href', 'https://x.com/artist')
  })

  it('displays hostname only in chip label, not the full URL', () => {
    const artist = makeArtist({
      links: [{ id: 'l1', url: 'https://bsky.app/@handle', is_primary: true }],
    })
    render(<ArtistsList artists={[artist]} onEditClick={noop} />)
    expect(screen.getByText('bsky.app')).toBeInTheDocument()
    expect(screen.queryByText('https://bsky.app/@handle')).not.toBeInTheDocument()
  })

  it('shows at most 3 chips and a +N badge when there are more than 3 links', () => {
    const artist = makeArtist({
      links: [
        { id: 'l1', url: 'https://bsky.app/a', is_primary: true },
        { id: 'l2', url: 'https://x.com/a', is_primary: false },
        { id: 'l3', url: 'https://ko-fi.com/a', is_primary: false },
        { id: 'l4', url: 'https://pixiv.net/a', is_primary: false },
      ],
    })
    render(<ArtistsList artists={[artist]} onEditClick={noop} />)
    expect(screen.getAllByRole('link')).toHaveLength(3)
    expect(screen.getByText('+1')).toBeInTheDocument()
  })

  it('shows "No links" when the artist has no links', () => {
    const artist = makeArtist({ links: [] })
    render(<ArtistsList artists={[artist]} onEditClick={noop} />)
    expect(screen.getByText('No links')).toBeInTheDocument()
  })

  it('shows "—" in the Notes cell when notes is null', () => {
    const artist = makeArtist({ notes: null })
    render(<ArtistsList artists={[artist]} onEditClick={noop} />)
    expect(screen.getByText('—')).toBeInTheDocument()
  })
})
