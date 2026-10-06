import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { vi, describe, it, expect, afterEach } from 'vitest'
import ArtistCharacterFilter from './ArtistCharacterFilter'
import { type Artist } from '@/features/artists/api/useArtists'

vi.mock('@/features/artists/api/useArtists', () => ({
  useArtists: () => ({ data: [], isError: false }),
}))

let mockCharactersError = false

vi.mock('@/features/characters/api/useCharacters', () => ({
  useCharacters: () => ({
    data: mockCharactersError ? undefined : [],
    isError: mockCharactersError,
    isPending: false,
  }),
}))

const noop = () => {}

const artist: Artist = {
  id: 'a1',
  name: 'Chevira',
  notes: null,
  artist_link: null,
  created_at: '',
  updated_at: '',
}

const twoChars = [
  { id: 'c1', name: 'Amitie' },
  { id: 'c2', name: 'Tanelle' },
]

function renderFilter(props: Partial<Parameters<typeof ArtistCharacterFilter>[0]> = {}) {
  return render(
    <ArtistCharacterFilter
      artist={null}
      onArtistChange={noop}
      characters={[]}
      onCharactersChange={noop}
      {...props}
    />,
  )
}

describe('ArtistCharacterFilter — trigger label', () => {
  it('shows "Filter" when no filters are active', () => {
    renderFilter()
    expect(screen.getByRole('button', { name: /filter/i })).toBeInTheDocument()
  })

  it('shows "1 artist" when an artist is selected', () => {
    renderFilter({ artist })
    expect(screen.getByRole('button', { name: '1 artist' })).toBeInTheDocument()
  })

  it('shows character count without match suffix when characterMatch is omitted', () => {
    renderFilter({ characters: [{ id: 'c1', name: 'Amitie' }] })
    expect(screen.getByRole('button', { name: '1 character' })).toBeInTheDocument()
  })

  it('shows match suffix when characterMatch is provided', () => {
    renderFilter({ characters: twoChars, characterMatch: 'any', onCharacterMatchChange: noop })
    expect(screen.getByRole('button', { name: '2 characters, any' })).toBeInTheDocument()
  })

  it('combines artist and characters in the label', () => {
    renderFilter({ artist, characters: twoChars, characterMatch: 'all', onCharacterMatchChange: noop })
    expect(screen.getByRole('button', { name: '1 artist · 2 characters, all' })).toBeInTheDocument()
  })
})

describe('ArtistCharacterFilter — panel', () => {
  it('shows Artist and Characters fields after trigger click', async () => {
    renderFilter()
    await userEvent.click(screen.getByRole('button', { name: /filter/i }))
    expect(screen.getByText('Artist')).toBeInTheDocument()
    expect(screen.getByText('Characters')).toBeInTheDocument()
  })

  it('shows error message when characters fail to load', async () => {
    mockCharactersError = true
    renderFilter()
    await userEvent.click(screen.getByRole('button', { name: /filter/i }))
    expect(screen.getByText('Failed to load characters.')).toBeInTheDocument()
  })
})

afterEach(() => {
  mockCharactersError = false
})

describe('ArtistCharacterFilter — match toggle', () => {
  it('is hidden when characterMatch prop is omitted', async () => {
    renderFilter({ characters: twoChars })
    await userEvent.click(screen.getByRole('button', { name: /2 characters/i }))
    expect(screen.queryByRole('button', { name: 'Any' })).not.toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'All' })).not.toBeInTheDocument()
  })

  it('is hidden when fewer than 2 characters are selected', async () => {
    renderFilter({
      characters: [{ id: 'c1', name: 'Amitie' }],
      characterMatch: 'any',
      onCharacterMatchChange: noop,
    })
    await userEvent.click(screen.getByRole('button', { name: /1 character/i }))
    expect(screen.queryByRole('button', { name: 'Any' })).not.toBeInTheDocument()
  })

  it('is visible when 2+ characters are selected and characterMatch is provided', async () => {
    renderFilter({ characters: twoChars, characterMatch: 'any', onCharacterMatchChange: noop })
    await userEvent.click(screen.getByRole('button', { name: /2 characters/i }))
    expect(screen.getByRole('button', { name: 'Any' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'All' })).toBeInTheDocument()
  })
})

describe('ArtistCharacterFilter — Clear all', () => {
  it('calls both onArtistChange(null) and onCharactersChange([]) when clicked', async () => {
    const onArtistChange = vi.fn()
    const onCharactersChange = vi.fn()
    renderFilter({ artist, characters: twoChars, onArtistChange, onCharactersChange })
    await userEvent.click(screen.getByRole('button', { name: /1 artist/i }))
    await userEvent.click(screen.getByRole('button', { name: /clear all/i }))
    expect(onArtistChange).toHaveBeenCalledWith(null)
    expect(onCharactersChange).toHaveBeenCalledWith([])
  })
})
