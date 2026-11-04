import { render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { vi, describe, it, expect, afterEach } from 'vitest'
import ArtistCharacterFilter from './ArtistCharacterFilter'
import { type Artist } from '@/features/artists/api/useArtists'
import { type Verse } from '@/features/verses/api/useVerses'

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

vi.mock('@/features/verses/api/useVerses', () => ({
  useVerses: () => ({ data: [verse], isError: false }),
}))

let mockTagsError = false
const tagRequests: unknown[] = []

vi.mock('@/features/tags/api/useTags', () => ({
  useTags: (options: unknown) => {
    tagRequests.push(options)
    return {
      data: mockTagsError ? undefined : [
        { id: 'tag-1', name: 'wip' },
        { id: 'tag-2', name: 'sketch' },
        { id: 'tag-3', name: 'outfit' },
      ],
      isError: mockTagsError,
      isPending: false,
    }
  },
}))

const noop = () => {}

const artist: Artist = {
  id: 'a1',
  name: 'Chevira',
  notes: null,
  links: [],
  created_at: '',
  updated_at: '',
}

const verse: Verse = {
  id: 'v1',
  name: 'Elowen',
  sort_order: 0,
  last_used_at: null,
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

  it('calls onVerseChange(null) when provided and clicked', async () => {
    const onArtistChange = vi.fn()
    const onCharactersChange = vi.fn()
    const onVerseChange = vi.fn()
    renderFilter({
      selectedVerse: verse,
      onVerseChange,
      onArtistChange,
      onCharactersChange,
    })
    await userEvent.click(screen.getByRole('button', { name: /1 verse/i }))
    await userEvent.click(screen.getByRole('button', { name: /clear all/i }))
    expect(onVerseChange).toHaveBeenCalledWith(null)
    expect(onArtistChange).toHaveBeenCalledWith(null)
    expect(onCharactersChange).toHaveBeenCalledWith([])
  })
})

describe('ArtistCharacterFilter — Verse field', () => {
  it('does not render a Verse field when selectedVerse/onVerseChange are omitted', async () => {
    renderFilter()
    await userEvent.click(screen.getByRole('button', { name: /filter/i }))
    expect(screen.queryByText('Verse')).not.toBeInTheDocument()
  })

  it('renders a Verse field when selectedVerse and onVerseChange are provided', async () => {
    renderFilter({ selectedVerse: null, onVerseChange: noop })
    await userEvent.click(screen.getByRole('button', { name: /filter/i }))
    expect(screen.getByText('Verse')).toBeInTheDocument()
  })

  it('shows "1 verse" in the trigger label when a verse is selected', () => {
    renderFilter({ selectedVerse: verse, onVerseChange: noop })
    expect(screen.getByRole('button', { name: '1 verse' })).toBeInTheDocument()
  })

  it('combines verse with artist and characters in the label', () => {
    renderFilter({ artist, selectedVerse: verse, onVerseChange: noop })
    expect(screen.getByRole('button', { name: '1 artist · 1 verse' })).toBeInTheDocument()
  })
})

describe('ArtistCharacterFilter — footerAction slot', () => {
  it('renders slot content when footerAction is provided', async () => {
    renderFilter({
      artist,
      footerAction: <button type="button">Save to dashboard</button>,
    })
    await userEvent.click(screen.getByRole('button', { name: /1 artist/i }))
    expect(screen.getByRole('button', { name: 'Save to dashboard' })).toBeInTheDocument()
  })

  it('does not render slot content when footerAction is omitted', async () => {
    renderFilter({ artist })
    await userEvent.click(screen.getByRole('button', { name: /1 artist/i }))
    expect(screen.queryByRole('button', { name: 'Save to dashboard' })).not.toBeInTheDocument()
  })
})

describe('ArtistCharacterFilter — Tags field', () => {
  const twoTags = [
    { id: 'tag-1', name: 'wip' },
    { id: 'tag-2', name: 'sketch' },
  ]

  it('does not render a Tags field when tags/onTagsChange are omitted', async () => {
    renderFilter()
    await userEvent.click(screen.getByRole('button', { name: /filter/i }))
    expect(screen.queryByText('Tags')).not.toBeInTheDocument()
  })

  it('does not request tags when the Tags field is not provided', () => {
    tagRequests.length = 0
    renderFilter()
    expect(tagRequests).toEqual([expect.objectContaining({ enabled: false })])
  })

  it('renders a Tags field when tags and onTagsChange are provided', async () => {
    renderFilter({ tags: [], onTagsChange: noop })
    await userEvent.click(screen.getByRole('button', { name: /filter/i }))
    expect(screen.getByText('Tags')).toBeInTheDocument()
  })

  it('requests the recency-ordered tag list', () => {
    tagRequests.length = 0
    renderFilter({ tags: [], onTagsChange: noop })
    expect(tagRequests).toEqual([expect.objectContaining({ sort: 'recent', enabled: true })])
  })

  it('selecting a suggestion calls onTagsChange with the tag appended', async () => {
    const onTagsChange = vi.fn()
    renderFilter({ tags: [], onTagsChange })
    await userEvent.click(screen.getByRole('button', { name: /filter/i }))
    await userEvent.type(screen.getByPlaceholderText('Search tags…'), 'wi')
    await userEvent.click(await screen.findByText('wip'))
    expect(onTagsChange).toHaveBeenCalledWith([{ id: 'tag-1', name: 'wip' }])
  })

  it('removing a chip calls onTagsChange with the tag removed', async () => {
    const onTagsChange = vi.fn()
    renderFilter({ tags: twoTags, onTagsChange })
    await userEvent.click(screen.getByRole('button', { name: /2 tags/i }))
    await userEvent.click(screen.getByRole('button', { name: /remove wip/i }))
    expect(onTagsChange).toHaveBeenCalledWith([{ id: 'tag-2', name: 'sketch' }])
  })

  it('excludes already-selected tags from the suggestions', async () => {
    renderFilter({ tags: [{ id: 'tag-2', name: 'sketch' }], onTagsChange: noop })
    await userEvent.click(screen.getByRole('button', { name: /1 tag/i }))
    const section = screen.getByText('Tags').parentElement!.parentElement!
    await userEvent.type(within(section).getByRole('textbox'), 's')
    expect(screen.getAllByText('sketch')).toHaveLength(1) // the chip only
  })

  it('does not offer to create an unknown tag', async () => {
    renderFilter({ tags: [], onTagsChange: noop })
    await userEvent.click(screen.getByRole('button', { name: /filter/i }))
    await userEvent.type(screen.getByPlaceholderText('Search tags…'), 'brand-new')
    expect(screen.queryByText(/create/i)).not.toBeInTheDocument()
  })

  it('shows an error message when tags fail to load', async () => {
    mockTagsError = true
    renderFilter({ tags: [], onTagsChange: noop })
    await userEvent.click(screen.getByRole('button', { name: /filter/i }))
    expect(screen.getByText('Failed to load tags.')).toBeInTheDocument()
  })
})

describe('ArtistCharacterFilter — tag match toggle', () => {
  const oneTag = [{ id: 'tag-1', name: 'wip' }]
  const twoTags = [...oneTag, { id: 'tag-2', name: 'sketch' }]

  it('is hidden when tagMatch is omitted', async () => {
    renderFilter({ tags: twoTags, onTagsChange: noop })
    await userEvent.click(screen.getByRole('button', { name: /2 tags/i }))
    expect(screen.queryByRole('button', { name: 'Any' })).not.toBeInTheDocument()
  })

  it('is hidden with fewer than two tags', async () => {
    renderFilter({ tags: oneTag, onTagsChange: noop, tagMatch: 'all', onTagMatchChange: noop })
    await userEvent.click(screen.getByRole('button', { name: /1 tag/i }))
    expect(screen.queryByRole('button', { name: 'Any' })).not.toBeInTheDocument()
  })

  it('is visible with two or more tags', async () => {
    renderFilter({ tags: twoTags, onTagsChange: noop, tagMatch: 'all', onTagMatchChange: noop })
    await userEvent.click(screen.getByRole('button', { name: /2 tags/i }))
    expect(screen.getByRole('button', { name: 'Any' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'All' })).toBeInTheDocument()
  })

  it('clicking Any calls onTagMatchChange and leaves the character mode alone', async () => {
    const onTagMatchChange = vi.fn()
    const onCharacterMatchChange = vi.fn()
    renderFilter({ tags: twoTags, onTagsChange: noop, tagMatch: 'all', onTagMatchChange, onCharacterMatchChange, characterMatch: 'all' })
    await userEvent.click(screen.getByRole('button', { name: /2 tags/i }))
    await userEvent.click(screen.getByRole('button', { name: 'Any' }))
    expect(onTagMatchChange).toHaveBeenCalledWith('any')
    expect(onCharacterMatchChange).not.toHaveBeenCalled()
  })
})

describe('ArtistCharacterFilter — tag summary label', () => {
  const oneTag = [{ id: 'tag-1', name: 'wip' }]
  const twoTags = [...oneTag, { id: 'tag-2', name: 'sketch' }]

  it('shows "1 tag" for a single tag with nothing else active', () => {
    renderFilter({ tags: oneTag, onTagsChange: noop })
    expect(screen.getByRole('button', { name: '1 tag' })).toBeInTheDocument()
  })

  it('shows the match mode when tagMatch is provided', () => {
    renderFilter({ tags: twoTags, onTagsChange: noop, tagMatch: 'any', onTagMatchChange: noop })
    expect(screen.getByRole('button', { name: '2 tags, any' })).toBeInTheDocument()
  })

  it('combines tags with artist and characters', () => {
    renderFilter({
      artist,
      characters: twoChars,
      characterMatch: 'all',
      onCharacterMatchChange: noop,
      tags: oneTag,
      onTagsChange: noop,
      tagMatch: 'all',
      onTagMatchChange: noop,
    })
    expect(screen.getByRole('button', { name: '1 artist · 2 characters, all · 1 tag, all' })).toBeInTheDocument()
  })
})

describe('ArtistCharacterFilter — Clear all with tags', () => {
  it('calls onTagsChange([]) when provided', async () => {
    const onTagsChange = vi.fn()
    renderFilter({ artist, tags: [{ id: 'tag-1', name: 'wip' }], onTagsChange })
    await userEvent.click(screen.getByRole('button', { name: /1 artist/i }))
    await userEvent.click(screen.getByRole('button', { name: /clear all/i }))
    expect(onTagsChange).toHaveBeenCalledWith([])
  })

  it('shows Clear all for a tag-only filter', async () => {
    renderFilter({ tags: [{ id: 'tag-1', name: 'wip' }], onTagsChange: noop })
    await userEvent.click(screen.getByRole('button', { name: '1 tag' }))
    expect(screen.getByRole('button', { name: /clear all/i })).toBeInTheDocument()
  })

  it('does not call onTagsChange when the tag props are absent', async () => {
    const onArtistChange = vi.fn()
    renderFilter({ artist, onArtistChange })
    await userEvent.click(screen.getByRole('button', { name: /1 artist/i }))
    await userEvent.click(screen.getByRole('button', { name: /clear all/i }))
    expect(onArtistChange).toHaveBeenCalledWith(null)
  })
})
