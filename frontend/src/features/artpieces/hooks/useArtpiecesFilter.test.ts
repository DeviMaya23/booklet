import { describe, it, expect, vi, beforeEach } from 'vitest'
import { renderHook, act } from '@testing-library/react'
import { filterAndSortArtpieces, useArtpiecesFilter } from './useArtpiecesFilter'
import { type ArtpieceSummary } from '../api/useArtpieces'
import { type Artist } from '@/features/artists/api/useArtists'
import { type Character } from '@/features/characters/api/useCharacters'

// --- Mocks ---

const mockArtists: Artist[] = [
  { id: 'artist-1', name: 'Alice', notes: null, links: [], created_at: '', updated_at: '' },
  { id: 'artist-2', name: 'Bob', notes: null, links: [], created_at: '', updated_at: '' },
]

const mockCharacters: Character[] = [
  { id: 'char-a', name: 'Alpha', avatar_url: null, notes: null, is_public: false, folders: [], created_at: '', updated_at: '' },
  { id: 'char-b', name: 'Beta', avatar_url: null, notes: null, is_public: false, folders: [], created_at: '', updated_at: '' },
]

let mockArtistData: Artist[] | undefined = mockArtists
let mockCharacterData: Character[] | undefined = mockCharacters

vi.mock('@/features/artists/api/useArtists', () => ({
  useArtists: () => ({ data: mockArtistData }),
}))

vi.mock('@/features/characters/api/useCharacters', () => ({
  useCharacters: () => ({ data: mockCharacterData }),
}))

let mockSearchParams = new URLSearchParams()
const setSearchParamsMock = vi.fn((updater: (prev: URLSearchParams) => URLSearchParams) => {
  mockSearchParams = updater(mockSearchParams)
})

vi.mock('react-router-dom', () => ({
  useSearchParams: () => [mockSearchParams, setSearchParamsMock],
}))

// --- Test data ---

const artist1 = mockArtists[0]
const artist2 = mockArtists[1]
const charA = { id: 'char-a', name: 'Alpha' }
const charB = { id: 'char-b', name: 'Beta' }

function makeArtpiece(overrides: Partial<ArtpieceSummary>): ArtpieceSummary {
  return {
    id: 'ap-default',
    title: 'Default',
    artist_id: null,
    artist_name: null,
    artist_link: null,
    cover_file_id: null,
    commission_id: null,
    thumbnail_url: null,
    cover_file_url: null,
    cover_file_mime_type: null,
    cover_file_name: null,
    notes: null,
    characters: [],
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
    ...overrides,
  }
}

const apAlice = makeArtpiece({ id: 'ap-alice', title: 'Sunrise', artist_id: artist1.id, artist_name: artist1.name })
const apBob = makeArtpiece({ id: 'ap-bob', title: 'Moonrise', artist_id: artist2.id, artist_name: artist2.name })
const apNoTitle = makeArtpiece({ id: 'ap-notitle', title: null })

const defaults = {
  search: '',
  sort: 'newest' as const,
  selectedArtist: null,
  selectedCharacters: [],
  characterMatch: 'all' as const,
}

beforeEach(() => {
  mockSearchParams = new URLSearchParams()
  setSearchParamsMock.mockClear()
  mockArtistData = mockArtists
  mockCharacterData = mockCharacters
})

// --- Pure function tests ---

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
  const apCharAB = makeArtpiece({ id: 'ap-ab', title: 'Both', characters: [charA, charB] })
  const apCharA = makeArtpiece({ id: 'ap-a', title: 'Alpha only', characters: [charA] })

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

// --- Hook tests ---

describe('useArtpiecesFilter — URL param hydration', () => {
  it('initialises selectedArtist from artist URL param', () => {
    mockSearchParams = new URLSearchParams({ artist: 'artist-1' })
    const { result } = renderHook(() => useArtpiecesFilter([apAlice, apBob]))
    expect(result.current.selectedArtist?.id).toBe('artist-1')
    expect(result.current.selectedArtist?.name).toBe('Alice')
  })

  it('initialises selectedCharacters from characters URL param', () => {
    mockSearchParams = new URLSearchParams({ characters: 'char-a,char-b' })
    const { result } = renderHook(() => useArtpiecesFilter([]))
    expect(result.current.selectedCharacters).toHaveLength(2)
    expect(result.current.selectedCharacters[0].id).toBe('char-a')
    expect(result.current.selectedCharacters[1].id).toBe('char-b')
  })

  it('initialises characterMatch from match URL param', () => {
    mockSearchParams = new URLSearchParams({ match: 'any' })
    const { result } = renderHook(() => useArtpiecesFilter([]))
    expect(result.current.characterMatch).toBe('any')
  })

  it('defaults to null artist and empty characters when params absent', () => {
    const { result } = renderHook(() => useArtpiecesFilter([]))
    expect(result.current.selectedArtist).toBeNull()
    expect(result.current.selectedCharacters).toHaveLength(0)
    expect(result.current.characterMatch).toBe('all')
  })
})

describe('useArtpiecesFilter — filter changes update URL', () => {
  it('setSelectedArtist sets artist param', () => {
    const { result } = renderHook(() => useArtpiecesFilter([]))
    act(() => result.current.setSelectedArtist(artist1))
    expect(setSearchParamsMock).toHaveBeenCalledOnce()
    expect(mockSearchParams.get('artist')).toBe('artist-1')
  })

  it('setSelectedArtist(null) removes artist param', () => {
    mockSearchParams = new URLSearchParams({ artist: 'artist-1' })
    const { result } = renderHook(() => useArtpiecesFilter([]))
    act(() => result.current.setSelectedArtist(null))
    expect(mockSearchParams.has('artist')).toBe(false)
  })

  it('setSelectedCharacters sets characters param', () => {
    const { result } = renderHook(() => useArtpiecesFilter([]))
    act(() => result.current.setSelectedCharacters([charA, charB]))
    expect(mockSearchParams.get('characters')).toBe('char-a,char-b')
  })

  it('setSelectedCharacters([]) removes characters param', () => {
    mockSearchParams = new URLSearchParams({ characters: 'char-a' })
    const { result } = renderHook(() => useArtpiecesFilter([]))
    act(() => result.current.setSelectedCharacters([]))
    expect(mockSearchParams.has('characters')).toBe(false)
  })

  it('setCharacterMatch("any") sets match param', () => {
    const { result } = renderHook(() => useArtpiecesFilter([]))
    act(() => result.current.setCharacterMatch('any'))
    expect(mockSearchParams.get('match')).toBe('any')
  })

  it('setCharacterMatch("all") removes match param', () => {
    mockSearchParams = new URLSearchParams({ match: 'any' })
    const { result } = renderHook(() => useArtpiecesFilter([]))
    act(() => result.current.setCharacterMatch('all'))
    expect(mockSearchParams.has('match')).toBe(false)
  })
})
