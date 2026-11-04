import { describe, it, expect, vi, beforeEach } from 'vitest'
import { renderHook, act } from '@testing-library/react'
import { filterAndSortArtpieces, useArtpiecesFilter } from './useArtpiecesFilter'
import { type ArtpieceSummary } from '../api/useArtpieces'
import { type Artist } from '@/features/artists/api/useArtists'
import { type Character } from '@/features/characters/api/useCharacters'
import { type Verse } from '@/features/verses/api/useVerses'
import { type Tag } from '@/features/tags/api/useTags'

// --- Mocks ---

const mockArtists: Artist[] = [
  { id: 'artist-1', name: 'Alice', notes: null, links: [], created_at: '', updated_at: '' },
  { id: 'artist-2', name: 'Bob', notes: null, links: [], created_at: '', updated_at: '' },
]

const mockCharacters: Character[] = [
  { id: 'char-a', name: 'Alpha', avatar_url: null, notes: null, is_public: false, verse_id: null, folders: [], created_at: '', updated_at: '' },
  { id: 'char-b', name: 'Beta', avatar_url: null, notes: null, is_public: false, verse_id: null, folders: [], created_at: '', updated_at: '' },
]

const mockVerses: Verse[] = [
  { id: 'verse-1', name: 'Elowen', sort_order: 0, last_used_at: null, created_at: '', updated_at: '' },
  { id: 'verse-2', name: 'Thalassa', sort_order: 1, last_used_at: null, created_at: '', updated_at: '' },
]

const mockTags: Tag[] = [
  { id: 'tag-a', name: 'outfit' },
  { id: 'tag-b', name: 'sketch' },
]

let mockArtistData: Artist[] | undefined = mockArtists
let mockTagData: Tag[] | undefined = mockTags
let mockTagError = false
let mockCharacterData: Character[] | undefined = mockCharacters
let mockVerseData: Verse[] | undefined = mockVerses

vi.mock('@/features/artists/api/useArtists', () => ({
  useArtists: () => ({ data: mockArtistData }),
}))

vi.mock('@/features/characters/api/useCharacters', () => ({
  useCharacters: () => ({ data: mockCharacterData }),
}))

vi.mock('@/features/verses/api/useVerses', () => ({
  useVerses: () => ({ data: mockVerseData }),
}))

vi.mock('@/features/tags/api/useTags', () => ({
  useTags: () => ({ data: mockTagError ? undefined : mockTagData, isPending: mockTagData === undefined && !mockTagError, isError: mockTagError }),
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
const verse1 = mockVerses[0]

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
    tags: [],
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
  selectedVerse: null,
  allCharacters: mockCharacters,
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
  mockVerseData = mockVerses
  mockTagData = mockTags
  mockTagError = false
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

const tagA = { id: 'tag-a', name: 'outfit' }
const tagB = { id: 'tag-b', name: 'sketch' }

describe('filterAndSortArtpieces — tag filter', () => {
  const apTagAB = makeArtpiece({ id: 'ap-tag-ab', title: 'Both tags', tags: [tagA, tagB] })
  const apTagA = makeArtpiece({ id: 'ap-tag-a', title: 'Tag A only', tags: [tagA] })
  const apUntagged = makeArtpiece({ id: 'ap-untagged', title: 'Untagged' })

  it('All mode: requires every selected tag to be present', () => {
    const result = filterAndSortArtpieces([apTagAB, apTagA, apUntagged], {
      ...defaults,
      selectedTags: [tagA, tagB],
      tagMatch: 'all',
    })
    expect(result).toEqual([apTagAB])
  })

  it('Any mode: requires at least one selected tag to be present', () => {
    const result = filterAndSortArtpieces([apTagAB, apTagA, apUntagged], {
      ...defaults,
      selectedTags: [tagA, tagB],
      tagMatch: 'any',
    })
    expect(result.map((a) => a.id).sort()).toEqual(['ap-tag-a', 'ap-tag-ab'])
  })

  it('single tag: shows only artpieces that have it', () => {
    const result = filterAndSortArtpieces([apTagAB, apTagA, apUntagged], { ...defaults, selectedTags: [tagB] })
    expect(result).toEqual([apTagAB])
  })

  it('defaults to All when tagMatch is not provided', () => {
    const result = filterAndSortArtpieces([apTagAB, apTagA], { ...defaults, selectedTags: [tagA, tagB] })
    expect(result).toEqual([apTagAB])
  })

  it('no tags selected: no tag filtering', () => {
    const result = filterAndSortArtpieces([apTagAB, apTagA, apUntagged], { ...defaults, selectedTags: [], tagMatch: 'all' })
    expect(result).toHaveLength(3)
  })

  it('combines with the artist filter', () => {
    const byAlice = makeArtpiece({ id: 'ap-alice-tagged', title: 'Alice tagged', artist_id: artist1.id, tags: [tagA] })
    const byBobTagged = makeArtpiece({ id: 'ap-bob-tagged', title: 'Bob tagged', artist_id: artist2.id, tags: [tagA] })
    const result = filterAndSortArtpieces([byAlice, byBobTagged], { ...defaults, selectedArtist: artist1, selectedTags: [tagA] })
    expect(result).toEqual([byAlice])
  })

  it('character Any and tag All are applied independently', () => {
    const match = makeArtpiece({ id: 'ap-match', title: 'Match', characters: [charA], tags: [tagA, tagB] })
    const wrongTags = makeArtpiece({ id: 'ap-wrong-tags', title: 'Wrong tags', characters: [charA], tags: [tagA] })
    const wrongChars = makeArtpiece({ id: 'ap-wrong-chars', title: 'Wrong chars', characters: [], tags: [tagA, tagB] })
    const result = filterAndSortArtpieces([match, wrongTags, wrongChars], {
      ...defaults,
      selectedCharacters: [charA, charB],
      characterMatch: 'any',
      selectedTags: [tagA, tagB],
      tagMatch: 'all',
    })
    expect(result).toEqual([match])
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

describe('filterAndSortArtpieces — verse filter', () => {
  const charactersWithVerses: Character[] = [
    { ...mockCharacters[0], verse_id: 'verse-1' },
    { ...mockCharacters[1], verse_id: 'verse-2' },
  ]

  it('matches an artpiece with a character in the selected verse', () => {
    const apInVerse = makeArtpiece({ id: 'ap-in-verse', title: 'In verse', characters: [charA] })
    const result = filterAndSortArtpieces([apInVerse], {
      ...defaults,
      allCharacters: charactersWithVerses,
      selectedVerse: verse1,
    })
    expect(result).toEqual([apInVerse])
  })

  it('excludes an artpiece whose characters are all verse-less or in other verses', () => {
    const apOtherVerse = makeArtpiece({ id: 'ap-other-verse', title: 'Other verse', characters: [charB] })
    const apNoVerse = makeArtpiece({ id: 'ap-no-verse', title: 'No verse', characters: [charA] })
    const charactersNoneInVerse1: Character[] = [
      { ...mockCharacters[0], verse_id: null },
      { ...mockCharacters[1], verse_id: 'verse-2' },
    ]
    const result = filterAndSortArtpieces([apOtherVerse, apNoVerse], {
      ...defaults,
      allCharacters: charactersNoneInVerse1,
      selectedVerse: verse1,
    })
    expect(result).toHaveLength(0)
  })

  it('matches an artpiece with characters from multiple verses, including the selected one', () => {
    const apMixed = makeArtpiece({ id: 'ap-mixed', title: 'Mixed', characters: [charA, charB] })
    const result = filterAndSortArtpieces([apMixed], {
      ...defaults,
      allCharacters: charactersWithVerses,
      selectedVerse: verse1,
    })
    expect(result).toEqual([apMixed])
  })

  it('combines with artist filter (AND)', () => {
    const apAliceInVerse = makeArtpiece({ id: 'ap-alice-verse', title: 'Alice verse', artist_id: artist1.id, characters: [charA] })
    const apBobInVerse = makeArtpiece({ id: 'ap-bob-verse', title: 'Bob verse', artist_id: artist2.id, characters: [charA] })
    const result = filterAndSortArtpieces([apAliceInVerse, apBobInVerse], {
      ...defaults,
      allCharacters: charactersWithVerses,
      selectedVerse: verse1,
      selectedArtist: artist1,
    })
    expect(result).toEqual([apAliceInVerse])
  })

  it('combines with character filter (AND)', () => {
    const apBoth = makeArtpiece({ id: 'ap-both', title: 'Both', characters: [charA, charB] })
    const apVerseOnly = makeArtpiece({ id: 'ap-verse-only', title: 'Verse only', characters: [charA] })
    const result = filterAndSortArtpieces([apBoth, apVerseOnly], {
      ...defaults,
      allCharacters: charactersWithVerses,
      selectedVerse: verse1,
      selectedCharacters: [charB],
      characterMatch: 'all',
    })
    expect(result).toEqual([apBoth])
  })

  it('no verse selected: verse filter has no effect', () => {
    const result = filterAndSortArtpieces([apAlice, apBob], {
      ...defaults,
      allCharacters: charactersWithVerses,
      selectedVerse: null,
    })
    expect(result).toHaveLength(2)
  })

  it('does not exclude artpieces while characters are still loading', () => {
    const apOtherVerse = makeArtpiece({ id: 'ap-other-verse', title: 'Other verse', characters: [charB] })
    const result = filterAndSortArtpieces([apOtherVerse], {
      ...defaults,
      allCharacters: [],
      selectedVerse: verse1,
      charactersPending: true,
    })
    expect(result).toEqual([apOtherVerse])
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

  it('initialises selectedVerse from verse URL param', () => {
    mockSearchParams = new URLSearchParams({ verse: 'verse-1' })
    const { result } = renderHook(() => useArtpiecesFilter([]))
    expect(result.current.selectedVerse?.id).toBe('verse-1')
    expect(result.current.selectedVerse?.name).toBe('Elowen')
  })

  it('defaults selectedVerse to null when verse param absent', () => {
    const { result } = renderHook(() => useArtpiecesFilter([]))
    expect(result.current.selectedVerse).toBeNull()
  })

  it('resolves selectedVerse to null when the stored verse id does not exist in current verses', () => {
    mockSearchParams = new URLSearchParams({ verse: 'verse-deleted' })
    const { result } = renderHook(() => useArtpiecesFilter([]))
    expect(result.current.selectedVerse).toBeNull()
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

  it('setSelectedVerse sets verse param', () => {
    const { result } = renderHook(() => useArtpiecesFilter([]))
    act(() => result.current.setSelectedVerse(verse1))
    expect(mockSearchParams.get('verse')).toBe('verse-1')
  })

  it('setSelectedVerse(null) removes verse param', () => {
    mockSearchParams = new URLSearchParams({ verse: 'verse-1' })
    const { result } = renderHook(() => useArtpiecesFilter([]))
    act(() => result.current.setSelectedVerse(null))
    expect(mockSearchParams.has('verse')).toBe(false)
  })
})

// --- localStorage persistence tests ---

const FILTER_KEY = 'artpieces.filterState'

describe('useArtpiecesFilter — mount with no URL params restores from localStorage', () => {
  it('calls setSearchParams to restore stored artist when URL has no filter params', () => {
    localStore[FILTER_KEY] = JSON.stringify({ artistId: 'artist-1', characterIds: [], matchMode: 'all', isReference: false, verseId: null })
    renderHook(() => useArtpiecesFilter([]))
    expect(setSearchParamsMock).toHaveBeenCalled()
    expect(mockSearchParams.get('artist')).toBe('artist-1')
  })

  it('calls setSearchParams to restore stored characters when URL has no filter params', () => {
    localStore[FILTER_KEY] = JSON.stringify({ artistId: null, characterIds: ['char-a', 'char-b'], matchMode: 'any', isReference: false, verseId: null })
    renderHook(() => useArtpiecesFilter([]))
    expect(mockSearchParams.get('characters')).toBe('char-a,char-b')
    expect(mockSearchParams.get('match')).toBe('any')
  })

  it('calls setSearchParams to restore isReference when URL has no filter params', () => {
    localStore[FILTER_KEY] = JSON.stringify({ artistId: null, characterIds: [], matchMode: 'all', isReference: true, verseId: null })
    renderHook(() => useArtpiecesFilter([]))
    expect(mockSearchParams.get('ref')).toBe('true')
  })

  it('calls setSearchParams to restore stored verse when URL has no filter params', () => {
    localStore[FILTER_KEY] = JSON.stringify({ artistId: null, characterIds: [], matchMode: 'all', isReference: false, verseId: 'verse-1' })
    renderHook(() => useArtpiecesFilter([]))
    expect(mockSearchParams.get('verse')).toBe('verse-1')
  })
})

describe('useArtpiecesFilter — mount with URL params writes to localStorage', () => {
  it('writes URL-derived state to localStorage when arriving with filter params', () => {
    mockSearchParams = new URLSearchParams({ artist: 'artist-1', characters: 'char-a', match: 'any' })
    renderHook(() => useArtpiecesFilter([]))
    expect(localStorageMock.setItem).toHaveBeenCalledWith(
      FILTER_KEY,
      JSON.stringify({ artistId: 'artist-1', characterIds: ['char-a'], matchMode: 'any', isReference: false, verseId: null, tagIds: [], tagMatch: 'all' }),
    )
  })

  it('writes isReference state to localStorage when ref param is present', () => {
    mockSearchParams = new URLSearchParams({ ref: 'true' })
    renderHook(() => useArtpiecesFilter([]))
    expect(localStorageMock.setItem).toHaveBeenCalledWith(
      FILTER_KEY,
      JSON.stringify({ artistId: null, characterIds: [], matchMode: 'all', isReference: true, verseId: null, tagIds: [], tagMatch: 'all' }),
    )
  })

  it('writes verse state to localStorage when verse param is present', () => {
    mockSearchParams = new URLSearchParams({ verse: 'verse-1' })
    renderHook(() => useArtpiecesFilter([]))
    expect(localStorageMock.setItem).toHaveBeenCalledWith(
      FILTER_KEY,
      JSON.stringify({ artistId: null, characterIds: [], matchMode: 'all', isReference: false, verseId: 'verse-1', tagIds: [], tagMatch: 'all' }),
    )
  })
})

describe('useArtpiecesFilter — setters write to localStorage', () => {
  it('setSelectedArtist writes full filter state to localStorage', () => {
    const { result } = renderHook(() => useArtpiecesFilter([]))
    act(() => result.current.setSelectedArtist(artist1))
    expect(localStorageMock.setItem).toHaveBeenCalledWith(
      FILTER_KEY,
      JSON.stringify({ artistId: 'artist-1', characterIds: [], matchMode: 'all', isReference: false, verseId: null, tagIds: [], tagMatch: 'all' }),
    )
  })

  it('setSelectedCharacters writes full filter state to localStorage', () => {
    const { result } = renderHook(() => useArtpiecesFilter([]))
    act(() => result.current.setSelectedCharacters([charA, charB]))
    expect(localStorageMock.setItem).toHaveBeenCalledWith(
      FILTER_KEY,
      JSON.stringify({ artistId: null, characterIds: ['char-a', 'char-b'], matchMode: 'all', isReference: false, verseId: null, tagIds: [], tagMatch: 'all' }),
    )
  })

  it('setCharacterMatch writes full filter state to localStorage', () => {
    const { result } = renderHook(() => useArtpiecesFilter([]))
    act(() => result.current.setCharacterMatch('any'))
    expect(localStorageMock.setItem).toHaveBeenCalledWith(
      FILTER_KEY,
      JSON.stringify({ artistId: null, characterIds: [], matchMode: 'any', isReference: false, verseId: null, tagIds: [], tagMatch: 'all' }),
    )
  })

  it('setIsReference writes full filter state to localStorage', () => {
    const { result } = renderHook(() => useArtpiecesFilter([]))
    act(() => result.current.setIsReference(true))
    expect(localStorageMock.setItem).toHaveBeenCalledWith(
      FILTER_KEY,
      JSON.stringify({ artistId: null, characterIds: [], matchMode: 'all', isReference: true, verseId: null, tagIds: [], tagMatch: 'all' }),
    )
  })

  it('setSelectedVerse writes full filter state to localStorage', () => {
    const { result } = renderHook(() => useArtpiecesFilter([]))
    act(() => result.current.setSelectedVerse(verse1))
    expect(localStorageMock.setItem).toHaveBeenCalledWith(
      FILTER_KEY,
      JSON.stringify({ artistId: null, characterIds: [], matchMode: 'all', isReference: false, verseId: 'verse-1', tagIds: [], tagMatch: 'all' }),
    )
  })

  it('setSelectedVerse(null) removes verseId from stored filter state', () => {
    mockSearchParams = new URLSearchParams({ verse: 'verse-1' })
    const { result } = renderHook(() => useArtpiecesFilter([]))
    act(() => result.current.setSelectedVerse(null))
    expect(localStorageMock.setItem).toHaveBeenCalledWith(
      FILTER_KEY,
      JSON.stringify({ artistId: null, characterIds: [], matchMode: 'all', isReference: false, verseId: null, tagIds: [], tagMatch: 'all' }),
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
      JSON.stringify({ artistId: null, characterIds: [], matchMode: 'all', isReference: false, verseId: null, tagIds: [], tagMatch: 'all' }),
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
      JSON.stringify({ artistId: null, characterIds: [], matchMode: 'all', isReference: false, verseId: null, tagIds: [], tagMatch: 'all' }),
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

// --- Tag filter state ---

describe('useArtpiecesFilter — tag URL hydration', () => {
  it('initialises selectedTags from the tags URL param', () => {
    mockSearchParams = new URLSearchParams({ tags: 'tag-a,tag-b' })
    const { result } = renderHook(() => useArtpiecesFilter([]))
    expect(result.current.selectedTags).toEqual([tagA, tagB])
  })

  it('initialises tagMatch from the tag_match URL param', () => {
    mockSearchParams = new URLSearchParams({ tags: 'tag-a,tag-b', tag_match: 'any' })
    const { result } = renderHook(() => useArtpiecesFilter([]))
    expect(result.current.tagMatch).toBe('any')
  })

  it('defaults to no tags and match All when params are absent', () => {
    const { result } = renderHook(() => useArtpiecesFilter([]))
    expect(result.current.selectedTags).toEqual([])
    expect(result.current.tagMatch).toBe('all')
  })

  it('tag_match does not affect the character match mode', () => {
    mockSearchParams = new URLSearchParams({ tag_match: 'any' })
    const { result } = renderHook(() => useArtpiecesFilter([]))
    expect(result.current.characterMatch).toBe('all')
  })

  it('filters the artpiece list by the URL tags', () => {
    const tagged = makeArtpiece({ id: 'ap-tagged', title: 'Tagged', tags: [tagA] })
    const other = makeArtpiece({ id: 'ap-other', title: 'Other' })
    mockSearchParams = new URLSearchParams({ tags: 'tag-a' })
    const { result } = renderHook(() => useArtpiecesFilter([tagged, other]))
    expect(result.current.filtered).toEqual([tagged])
  })
})

describe('useArtpiecesFilter — deleted ids resolve to not selected', () => {
  const both = makeArtpiece({ id: 'ap-both', title: 'Both', tags: [tagA, tagB], characters: [charA, charB] })
  const onlyA = makeArtpiece({ id: 'ap-only-a', title: 'Only A', tags: [tagA], characters: [charA] })
  const plain = makeArtpiece({ id: 'ap-plain', title: 'Plain' })

  it('a deleted tag in All mode is dropped, leaving the remaining tag applied', () => {
    mockSearchParams = new URLSearchParams({ tags: 'tag-a,tag-deleted' })
    const { result } = renderHook(() => useArtpiecesFilter([both, onlyA, plain]))
    expect(result.current.selectedTags).toEqual([tagA])
    expect(result.current.filtered.map((a) => a.id).sort()).toEqual(['ap-both', 'ap-only-a'])
  })

  it('when every selected tag is deleted no tag filter is applied', () => {
    mockSearchParams = new URLSearchParams({ tags: 'tag-gone-1,tag-gone-2' })
    const { result } = renderHook(() => useArtpiecesFilter([both, onlyA, plain]))
    expect(result.current.selectedTags).toEqual([])
    expect(result.current.filtered).toHaveLength(3)
  })

  it('a deleted character in All mode is dropped, leaving the remaining character applied', () => {
    mockSearchParams = new URLSearchParams({ characters: 'char-a,char-deleted' })
    const { result } = renderHook(() => useArtpiecesFilter([both, onlyA, plain]))
    expect(result.current.selectedCharacters).toEqual([charA])
    expect(result.current.filtered.map((a) => a.id).sort()).toEqual(['ap-both', 'ap-only-a'])
  })

  it('does not apply the tag filter while the tag list is still loading', () => {
    mockTagData = undefined
    mockSearchParams = new URLSearchParams({ tags: 'tag-a' })
    const { result } = renderHook(() => useArtpiecesFilter([both, onlyA, plain]))
    expect(result.current.selectedTags).toEqual([])
    expect(result.current.filtered).toHaveLength(3)
  })

  it('keeps unresolved tag ids in the URL when another filter changes', () => {
    mockTagData = undefined
    mockSearchParams = new URLSearchParams({ tags: 'tag-a' })
    const { result } = renderHook(() => useArtpiecesFilter([]))
    act(() => result.current.setSelectedArtist(artist1))
    expect(mockSearchParams.get('tags')).toBe('tag-a')
    expect(mockSearchParams.get('artist')).toBe('artist-1')
  })

  it('editing the tags prunes a deleted id', () => {
    mockSearchParams = new URLSearchParams({ tags: 'tag-a,tag-deleted' })
    const { result } = renderHook(() => useArtpiecesFilter([]))
    act(() => result.current.setSelectedTags([tagA, tagB]))
    expect(mockSearchParams.get('tags')).toBe('tag-a,tag-b')
  })
})

describe('useArtpiecesFilter — tag setters update the URL', () => {
  it('setSelectedTags sets the tags param', () => {
    const { result } = renderHook(() => useArtpiecesFilter([]))
    act(() => result.current.setSelectedTags([tagA, tagB]))
    expect(mockSearchParams.get('tags')).toBe('tag-a,tag-b')
  })

  it('setSelectedTags([]) removes the tags param', () => {
    mockSearchParams = new URLSearchParams({ tags: 'tag-a' })
    const { result } = renderHook(() => useArtpiecesFilter([]))
    act(() => result.current.setSelectedTags([]))
    expect(mockSearchParams.has('tags')).toBe(false)
  })

  it('setTagMatch("any") sets tag_match and leaves match alone', () => {
    const { result } = renderHook(() => useArtpiecesFilter([]))
    act(() => result.current.setTagMatch('any'))
    expect(mockSearchParams.get('tag_match')).toBe('any')
    expect(mockSearchParams.has('match')).toBe(false)
  })

  it('setTagMatch("all") removes tag_match', () => {
    mockSearchParams = new URLSearchParams({ tags: 'tag-a,tag-b', tag_match: 'any' })
    const { result } = renderHook(() => useArtpiecesFilter([]))
    act(() => result.current.setTagMatch('all'))
    expect(mockSearchParams.has('tag_match')).toBe(false)
  })

  it('"Clear all" removes the tags and tag_match params together with the others', () => {
    mockSearchParams = new URLSearchParams({ artist: 'artist-1', tags: 'tag-a,tag-b', tag_match: 'any' })
    const { result } = renderHook(() => useArtpiecesFilter([]))
    act(() => {
      result.current.setSelectedArtist(null)
      result.current.setSelectedTags([])
      result.current.setTagMatch('all')
    })
    expect(mockSearchParams.toString()).toBe('')
  })
})

describe('useArtpiecesFilter — tag state in localStorage', () => {
  it('setSelectedTags writes the full filter state including tags', () => {
    const { result } = renderHook(() => useArtpiecesFilter([]))
    act(() => result.current.setSelectedTags([tagA]))
    expect(localStorageMock.setItem).toHaveBeenLastCalledWith(
      FILTER_KEY,
      JSON.stringify({ artistId: null, characterIds: [], matchMode: 'all', isReference: false, verseId: null, tagIds: ['tag-a'], tagMatch: 'all' }),
    )
  })

  it('setTagMatch writes the tag match mode', () => {
    const { result } = renderHook(() => useArtpiecesFilter([]))
    act(() => result.current.setTagMatch('any'))
    expect(localStorageMock.setItem).toHaveBeenLastCalledWith(
      FILTER_KEY,
      JSON.stringify({ artistId: null, characterIds: [], matchMode: 'all', isReference: false, verseId: null, tagIds: [], tagMatch: 'any' }),
    )
  })

  it('clearing the last tag writes the empty state', () => {
    mockSearchParams = new URLSearchParams({ tags: 'tag-a' })
    const { result } = renderHook(() => useArtpiecesFilter([]))
    act(() => result.current.setSelectedTags([]))
    expect(localStorageMock.setItem).toHaveBeenLastCalledWith(
      FILTER_KEY,
      JSON.stringify({ artistId: null, characterIds: [], matchMode: 'all', isReference: false, verseId: null, tagIds: [], tagMatch: 'all' }),
    )
  })

  it('writes URL-derived tag state when arriving with tag params', () => {
    mockSearchParams = new URLSearchParams({ tags: 'tag-a,tag-b', tag_match: 'any' })
    renderHook(() => useArtpiecesFilter([]))
    expect(localStorageMock.setItem).toHaveBeenCalledWith(
      FILTER_KEY,
      JSON.stringify({ artistId: null, characterIds: [], matchMode: 'all', isReference: false, verseId: null, tagIds: ['tag-a', 'tag-b'], tagMatch: 'any' }),
    )
  })

  it('restores stored tags and match mode when the URL has no filter params', () => {
    localStore[FILTER_KEY] = JSON.stringify({ artistId: null, characterIds: [], matchMode: 'all', isReference: false, verseId: null, tagIds: ['tag-a', 'tag-b'], tagMatch: 'any' })
    renderHook(() => useArtpiecesFilter([]))
    expect(mockSearchParams.get('tags')).toBe('tag-a,tag-b')
    expect(mockSearchParams.get('tag_match')).toBe('any')
  })

  it('a tag-only stored filter is restorable, not treated as empty', () => {
    localStore[FILTER_KEY] = JSON.stringify({ artistId: null, characterIds: [], matchMode: 'all', isReference: false, verseId: null, tagIds: ['tag-a'], tagMatch: 'all' })
    renderHook(() => useArtpiecesFilter([]))
    expect(setSearchParamsMock).toHaveBeenCalled()
    expect(mockSearchParams.get('tags')).toBe('tag-a')
  })

  it('state stored before tags existed restores with no tags and match All', () => {
    localStore[FILTER_KEY] = JSON.stringify({ artistId: 'artist-1', characterIds: [], matchMode: 'all', isReference: false, verseId: null })
    const { result } = renderHook(() => useArtpiecesFilter([]))
    expect(mockSearchParams.get('artist')).toBe('artist-1')
    expect(mockSearchParams.has('tags')).toBe(false)
    expect(mockSearchParams.has('tag_match')).toBe(false)
    expect(result.current.selectedTags).toEqual([])
    expect(result.current.tagMatch).toBe('all')
  })

  it('an empty stored tag state with nothing else does not touch the URL', () => {
    localStore[FILTER_KEY] = JSON.stringify({ artistId: null, characterIds: [], matchMode: 'all', isReference: false, verseId: null, tagIds: [], tagMatch: 'all' })
    renderHook(() => useArtpiecesFilter([]))
    expect(setSearchParamsMock).not.toHaveBeenCalled()
  })
})

describe('useArtpiecesFilter — clearing tags resets the tag match mode', () => {
  it('setSelectedTags([]) removes tag_match as well as tags', () => {
    mockSearchParams = new URLSearchParams({ tags: 'tag-a,tag-b', tag_match: 'any' })
    const { result } = renderHook(() => useArtpiecesFilter([]))
    act(() => result.current.setSelectedTags([]))
    expect(mockSearchParams.has('tags')).toBe(false)
    expect(mockSearchParams.has('tag_match')).toBe(false)
  })

  it('setSelectedTags([]) stores the default match mode', () => {
    mockSearchParams = new URLSearchParams({ tags: 'tag-a,tag-b', tag_match: 'any' })
    const { result } = renderHook(() => useArtpiecesFilter([]))
    act(() => result.current.setSelectedTags([]))
    expect(localStorageMock.setItem).toHaveBeenLastCalledWith(
      FILTER_KEY,
      JSON.stringify({ artistId: null, characterIds: [], matchMode: 'all', isReference: false, verseId: null, tagIds: [], tagMatch: 'all' }),
    )
  })

  it('"Clear all" as the component calls it (no setTagMatch) leaves no tag params behind', () => {
    mockSearchParams = new URLSearchParams({ artist: 'artist-1', tags: 'tag-a,tag-b', tag_match: 'any' })
    const { result } = renderHook(() => useArtpiecesFilter([]))
    act(() => {
      result.current.setSelectedArtist(null)
      result.current.setSelectedTags([])
    })
    expect(mockSearchParams.toString()).toBe('')
  })

  it('removing one of several tags keeps the match mode', () => {
    mockSearchParams = new URLSearchParams({ tags: 'tag-a,tag-b', tag_match: 'any' })
    const { result } = renderHook(() => useArtpiecesFilter([]))
    act(() => result.current.setSelectedTags([tagA]))
    expect(mockSearchParams.get('tags')).toBe('tag-a')
    expect(mockSearchParams.get('tag_match')).toBe('any')
  })
})

describe('useArtpiecesFilter — tag list failed to load', () => {
  const tagged = makeArtpiece({ id: 'ap-tagged', title: 'Tagged', tags: [tagA] })
  const plain = makeArtpiece({ id: 'ap-plain', title: 'Plain' })

  it('applies the URL tag ids as given instead of showing every artpiece', () => {
    mockTagError = true
    mockSearchParams = new URLSearchParams({ tags: 'tag-a' })
    const { result } = renderHook(() => useArtpiecesFilter([tagged, plain]))
    expect(result.current.selectedTags).toEqual([{ id: 'tag-a', name: '' }])
    expect(result.current.filtered).toEqual([tagged])
  })

  it('applies no tag filter when the URL has no tags', () => {
    mockTagError = true
    const { result } = renderHook(() => useArtpiecesFilter([tagged, plain]))
    expect(result.current.selectedTags).toEqual([])
    expect(result.current.filtered).toHaveLength(2)
  })
})
