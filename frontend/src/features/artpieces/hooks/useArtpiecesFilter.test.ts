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
  { id: 'char-a', name: 'Alpha', avatar_url: null, notes: null, is_public: false, verse_id: null, folders: [], created_at: '', updated_at: '' },
  { id: 'char-b', name: 'Beta', avatar_url: null, notes: null, is_public: false, verse_id: null, folders: [], created_at: '', updated_at: '' },
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
// Params handed to the hook on its last render. React Router's setSearchParams
// closes over these, so the updater form never sees the result of an earlier
// call in the same event — the mock has to behave the same way or it hides
// setters that undo each other.
let renderedSearchParams = mockSearchParams
type SearchParamsInit = URLSearchParams | ((prev: URLSearchParams) => URLSearchParams)
const setSearchParamsMock = vi.fn((next: SearchParamsInit) => {
  mockSearchParams = typeof next === 'function'
    ? next(new URLSearchParams(renderedSearchParams))
    : new URLSearchParams(next)
})

vi.mock('react-router-dom', () => ({
  useSearchParams: () => {
    renderedSearchParams = mockSearchParams
    return [mockSearchParams, setSearchParamsMock]
  },
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
    is_reference: false,
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
  isReference: false,
}

let localStore: Record<string, string> = {}
const localStorageMock = {
  getItem: vi.fn((key: string) => localStore[key] ?? null),
  setItem: vi.fn((key: string, value: string) => { localStore[key] = value }),
  removeItem: vi.fn((key: string) => { delete localStore[key] }),
  clear: vi.fn(() => { localStore = {} }),
}
vi.stubGlobal('localStorage', localStorageMock)

beforeEach(() => {
  localStore = {}
  mockSearchParams = new URLSearchParams()
  renderedSearchParams = mockSearchParams
  setSearchParamsMock.mockClear()
  localStorageMock.getItem.mockClear()
  localStorageMock.setItem.mockClear()
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

describe('filterAndSortArtpieces — isReference filter', () => {
  const apNormal = makeArtpiece({ id: 'ap-normal', title: 'Normal', is_reference: false })
  const apRef = makeArtpiece({ id: 'ap-ref', title: 'Ref', is_reference: true })

  it('toggle off (default) hides references', () => {
    const result = filterAndSortArtpieces([apNormal, apRef], { ...defaults, isReference: false })
    expect(result).toEqual([apNormal])
  })

  it('toggle on shows only references', () => {
    const result = filterAndSortArtpieces([apNormal, apRef], { ...defaults, isReference: true })
    expect(result).toEqual([apRef])
  })

  it('other filters still apply when toggle is on', () => {
    const apRefAlice = makeArtpiece({ id: 'ap-ref-alice', title: 'Ref Alice', is_reference: true, artist_id: artist1.id })
    const apRefBob = makeArtpiece({ id: 'ap-ref-bob', title: 'Ref Bob', is_reference: true, artist_id: artist2.id })
    const result = filterAndSortArtpieces([apRefAlice, apRefBob], {
      ...defaults,
      isReference: true,
      selectedArtist: artist1,
    })
    expect(result).toEqual([apRefAlice])
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

  it('absent ref param defaults isReference to false', () => {
    const { result } = renderHook(() => useArtpiecesFilter([]))
    expect(result.current.isReference).toBe(false)
  })

  it('ref=true restores isReference to true', () => {
    mockSearchParams = new URLSearchParams({ ref: 'true' })
    const { result } = renderHook(() => useArtpiecesFilter([]))
    expect(result.current.isReference).toBe(true)
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

  it('setIsReference(true) sets ref param', () => {
    const { result } = renderHook(() => useArtpiecesFilter([]))
    act(() => result.current.setIsReference(true))
    expect(mockSearchParams.get('ref')).toBe('true')
  })

  it('setIsReference(false) removes ref param', () => {
    mockSearchParams = new URLSearchParams({ ref: 'true' })
    const { result } = renderHook(() => useArtpiecesFilter([]))
    act(() => result.current.setIsReference(false))
    expect(mockSearchParams.has('ref')).toBe(false)
  })
})

// --- localStorage persistence tests ---

const FILTER_KEY = 'artpieces.filterState'

describe('useArtpiecesFilter — mount with no URL params restores from localStorage', () => {
  it('calls setSearchParams to restore stored artist when URL has no filter params', () => {
    localStore[FILTER_KEY] = JSON.stringify({ artistId: 'artist-1', characterIds: [], matchMode: 'all', isReference: false })
    renderHook(() => useArtpiecesFilter([]))
    expect(setSearchParamsMock).toHaveBeenCalled()
    expect(mockSearchParams.get('artist')).toBe('artist-1')
  })

  it('calls setSearchParams to restore stored characters when URL has no filter params', () => {
    localStore[FILTER_KEY] = JSON.stringify({ artistId: null, characterIds: ['char-a', 'char-b'], matchMode: 'any', isReference: false })
    renderHook(() => useArtpiecesFilter([]))
    expect(mockSearchParams.get('characters')).toBe('char-a,char-b')
    expect(mockSearchParams.get('match')).toBe('any')
  })

  it('calls setSearchParams to restore isReference when URL has no filter params', () => {
    localStore[FILTER_KEY] = JSON.stringify({ artistId: null, characterIds: [], matchMode: 'all', isReference: true })
    renderHook(() => useArtpiecesFilter([]))
    expect(mockSearchParams.get('ref')).toBe('true')
  })
})

describe('useArtpiecesFilter — mount with URL params writes to localStorage', () => {
  it('writes URL-derived state to localStorage when arriving with filter params', () => {
    mockSearchParams = new URLSearchParams({ artist: 'artist-1', characters: 'char-a', match: 'any' })
    renderHook(() => useArtpiecesFilter([]))
    expect(localStorageMock.setItem).toHaveBeenCalledWith(
      FILTER_KEY,
      JSON.stringify({ artistId: 'artist-1', characterIds: ['char-a'], matchMode: 'any', isReference: false }),
    )
  })

  it('writes isReference state to localStorage when ref param is present', () => {
    mockSearchParams = new URLSearchParams({ ref: 'true' })
    renderHook(() => useArtpiecesFilter([]))
    expect(localStorageMock.setItem).toHaveBeenCalledWith(
      FILTER_KEY,
      JSON.stringify({ artistId: null, characterIds: [], matchMode: 'all', isReference: true }),
    )
  })
})

describe('useArtpiecesFilter — setters write to localStorage', () => {
  it('setSelectedArtist writes full filter state to localStorage', () => {
    const { result } = renderHook(() => useArtpiecesFilter([]))
    act(() => result.current.setSelectedArtist(artist1))
    expect(localStorageMock.setItem).toHaveBeenCalledWith(
      FILTER_KEY,
      JSON.stringify({ artistId: 'artist-1', characterIds: [], matchMode: 'all', isReference: false }),
    )
  })

  it('setSelectedCharacters writes full filter state to localStorage', () => {
    const { result } = renderHook(() => useArtpiecesFilter([]))
    act(() => result.current.setSelectedCharacters([charA, charB]))
    expect(localStorageMock.setItem).toHaveBeenCalledWith(
      FILTER_KEY,
      JSON.stringify({ artistId: null, characterIds: ['char-a', 'char-b'], matchMode: 'all', isReference: false }),
    )
  })

  it('setCharacterMatch writes full filter state to localStorage', () => {
    const { result } = renderHook(() => useArtpiecesFilter([]))
    act(() => result.current.setCharacterMatch('any'))
    expect(localStorageMock.setItem).toHaveBeenCalledWith(
      FILTER_KEY,
      JSON.stringify({ artistId: null, characterIds: [], matchMode: 'any', isReference: false }),
    )
  })

  it('setIsReference writes full filter state to localStorage', () => {
    const { result } = renderHook(() => useArtpiecesFilter([]))
    act(() => result.current.setIsReference(true))
    expect(localStorageMock.setItem).toHaveBeenCalledWith(
      FILTER_KEY,
      JSON.stringify({ artistId: null, characterIds: [], matchMode: 'all', isReference: true }),
    )
  })
})

describe('useArtpiecesFilter — mount with no params and no localStorage', () => {
  it('does not call setSearchParams when localStorage is empty and URL has no filter params', () => {
    renderHook(() => useArtpiecesFilter([]))
    expect(setSearchParamsMock).not.toHaveBeenCalled()
  })
})

describe('useArtpiecesFilter — clearing filter writes empty state to localStorage', () => {
  it('setSelectedArtist(null) with no other filters writes empty state', () => {
    mockSearchParams = new URLSearchParams({ artist: 'artist-1' })
    const { result } = renderHook(() => useArtpiecesFilter([]))
    act(() => result.current.setSelectedArtist(null))
    expect(localStorageMock.setItem).toHaveBeenCalledWith(
      FILTER_KEY,
      JSON.stringify({ artistId: null, characterIds: [], matchMode: 'all', isReference: false }),
    )
  })
})

describe('useArtpiecesFilter — several setters in one event', () => {
  // "Clear all" calls setSelectedArtist, setSelectedCharacters and setIsReference
  // in a single click handler; each one has to survive the others.
  it('"Clear all" removes the artist, characters and ref params together', () => {
    mockSearchParams = new URLSearchParams({ artist: 'artist-1', characters: 'char-a,char-b', ref: 'true' })
    const { result } = renderHook(() => useArtpiecesFilter([]))
    act(() => {
      result.current.setSelectedArtist(null)
      result.current.setSelectedCharacters([])
      result.current.setIsReference(false)
    })
    expect(mockSearchParams.toString()).toBe('')
  })

  it('"Clear all" writes empty state to localStorage when several filters were active', () => {
    mockSearchParams = new URLSearchParams({ artist: 'artist-1', characters: 'char-a', ref: 'true' })
    const { result } = renderHook(() => useArtpiecesFilter([]))
    act(() => {
      result.current.setSelectedArtist(null)
      result.current.setSelectedCharacters([])
      result.current.setIsReference(false)
    })
    expect(localStorageMock.setItem).toHaveBeenLastCalledWith(
      FILTER_KEY,
      JSON.stringify({ artistId: null, characterIds: [], matchMode: 'all', isReference: false }),
    )
  })

  it('two setters in one event both take effect', () => {
    const { result } = renderHook(() => useArtpiecesFilter([]))
    act(() => {
      result.current.setSelectedArtist(artist1)
      result.current.setSelectedCharacters([charA])
    })
    expect(mockSearchParams.get('artist')).toBe('artist-1')
    expect(mockSearchParams.get('characters')).toBe('char-a')
  })
})

describe('useArtpiecesFilter — non-filter params', () => {
  it('preserves unrelated search params when a filter changes', () => {
    mockSearchParams = new URLSearchParams({ view: 'grid' })
    const { result } = renderHook(() => useArtpiecesFilter([]))
    act(() => result.current.setSelectedArtist(artist1))
    expect(mockSearchParams.get('view')).toBe('grid')
    expect(mockSearchParams.get('artist')).toBe('artist-1')
  })
})
