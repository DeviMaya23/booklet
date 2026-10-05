import { describe, it, expect } from 'vitest'
import { filterAndSortArtpieces } from './useArtpiecesFilter'
import { type ArtpieceSummary } from '../api/useArtpieces'

const artist1 = { id: 'artist-1', name: 'Alice', notes: null, artist_link: null, created_at: '', updated_at: '' }
const artist2 = { id: 'artist-2', name: 'Bob', notes: null, artist_link: null, created_at: '', updated_at: '' }
const charA = { id: 'char-a', name: 'Alpha' }
const charB = { id: 'char-b', name: 'Beta' }

function makeArtpiece(overrides: Partial<ArtpieceSummary>): ArtpieceSummary {
  return {
    id: 'ap-default',
    title: 'Default',
    artist_id: null,
    artist_name: null,
    cover_file_id: null,
    thumbnail_url: null,
    notes: null,
    characters: [],
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
    ...overrides,
  }
}

const apAlice = makeArtpiece({ id: 'ap-alice', title: 'Sunrise', artist_id: artist1.id, artist_name: artist1.name })
const apBob = makeArtpiece({ id: 'ap-bob', title: 'Moonrise', artist_id: artist2.id, artist_name: artist2.name })
const apCharAB = makeArtpiece({ id: 'ap-ab', title: 'Both', characters: [charA, charB] })
const apCharA = makeArtpiece({ id: 'ap-a', title: 'Alpha only', characters: [charA] })
const apNoTitle = makeArtpiece({ id: 'ap-notitle', title: null })

const defaults = {
  search: '',
  sort: 'newest' as const,
  selectedArtist: null,
  selectedCharacters: [],
  characterMatch: 'all' as const,
}

describe('filterAndSortArtpieces — search', () => {
  it('returns all artpieces when search is empty', () => {
    const result = filterAndSortArtpieces([apAlice, apBob], defaults)
    expect(result).toHaveLength(2)
  })

  it('filters by title case-insensitively', () => {
    const result = filterAndSortArtpieces([apAlice, apBob], { ...defaults, search: 'sunrise' })
    expect(result).toEqual([apAlice])
  })

  it('excludes artpieces with null title when search is non-empty', () => {
    const result = filterAndSortArtpieces([apAlice, apNoTitle], { ...defaults, search: 'any' })
    expect(result).not.toContain(apNoTitle)
  })
})

describe('filterAndSortArtpieces — artist filter', () => {
  it('shows only artpieces matching the selected artist', () => {
    const result = filterAndSortArtpieces([apAlice, apBob], { ...defaults, selectedArtist: artist1 })
    expect(result).toEqual([apAlice])
  })
})

describe('filterAndSortArtpieces — character filter', () => {
  it('All mode: requires every selected character to be present', () => {
    const result = filterAndSortArtpieces([apCharAB, apCharA], {
      ...defaults,
      selectedCharacters: [charA, charB],
      characterMatch: 'all',
    })
    expect(result).toEqual([apCharAB])
  })

  it('Any mode: requires at least one selected character to be present', () => {
    const result = filterAndSortArtpieces([apCharAB, apCharA], {
      ...defaults,
      selectedCharacters: [charA, charB],
      characterMatch: 'any',
    })
    expect(result).toHaveLength(2)
  })

  it('no character chips: character match mode has no effect', () => {
    const result = filterAndSortArtpieces([apAlice, apBob], {
      ...defaults,
      selectedCharacters: [],
      characterMatch: 'all',
    })
    expect(result).toHaveLength(2)
  })
})

describe('filterAndSortArtpieces — combined filters', () => {
  it('artist and character filters are AND-combined', () => {
    const apAliceWithChar = makeArtpiece({ id: 'ap-alice-char', title: 'X', artist_id: artist1.id, characters: [charA] })
    const apAliceNoChar = makeArtpiece({ id: 'ap-alice-nc', title: 'Y', artist_id: artist1.id, characters: [] })
    const result = filterAndSortArtpieces([apAliceWithChar, apAliceNoChar, apBob], {
      ...defaults,
      selectedArtist: artist1,
      selectedCharacters: [charA],
      characterMatch: 'all',
    })
    expect(result).toEqual([apAliceWithChar])
  })
})

describe('filterAndSortArtpieces — sort', () => {
  const older = makeArtpiece({ id: 'ap-old', title: 'Zebra', created_at: '2023-01-01T00:00:00Z' })
  const newer = makeArtpiece({ id: 'ap-new', title: 'Apple', created_at: '2024-06-01T00:00:00Z' })

  it('newest first is the default sort', () => {
    const result = filterAndSortArtpieces([older, newer], { ...defaults, sort: 'newest' })
    expect(result[0]).toEqual(newer)
  })

  it('alphabetical sorts by title ascending, null titles last', () => {
    const result = filterAndSortArtpieces([older, newer, apNoTitle], { ...defaults, sort: 'alpha' })
    expect(result[0]).toEqual(newer)
    expect(result[result.length - 1]).toEqual(apNoTitle)
  })
})
