import { useState, useEffect, useLayoutEffect, useRef } from 'react'
import { useSearchParams } from 'react-router-dom'
import { type ArtpieceSummary } from '../api/useArtpieces'
import { type Artist, useArtists } from '@/features/artists/api/useArtists'
import { useCharacters } from '@/features/characters/api/useCharacters'

export type SortMode = 'newest' | 'alpha'

const FILTER_STORAGE_KEY = 'artpieces.filterState'

interface StoredFilter {
  artistId: string | null
  characterIds: string[]
  matchMode: 'all' | 'any'
  isReference: boolean
}

function readStoredFilter(): StoredFilter | null {
  try {
    const raw = localStorage.getItem(FILTER_STORAGE_KEY)
    if (!raw) return null
    return JSON.parse(raw) as StoredFilter
  } catch {
    return null
  }
}

function writeStoredFilter(filter: StoredFilter): void {
  try {
    localStorage.setItem(FILTER_STORAGE_KEY, JSON.stringify(filter))
  } catch {
    // storage blocked (e.g. private mode)
  }
}

interface FilterParams {
  search: string
  sort: SortMode
  selectedArtist: Artist | null
  selectedCharacters: { id: string; name: string }[]
  characterMatch: 'all' | 'any'
  isReference: boolean
}

export function filterAndSortArtpieces(
  artpieces: ArtpieceSummary[],
  { search, sort, selectedArtist, selectedCharacters, characterMatch, isReference }: FilterParams,
): ArtpieceSummary[] {
  return artpieces
    .filter((a) => {
      if (isReference) {
        if (!a.is_reference) return false
      } else {
        if (a.is_reference) return false
      }
      if (search.trim()) {
        if (a.title === null) return false
        if (!a.title.toLowerCase().includes(search.toLowerCase())) return false
      }
      if (selectedArtist) {
        if (a.artist_id !== selectedArtist.id) return false
      }
      if (selectedCharacters.length > 0) {
        const ids = a.characters.map((c) => c.id)
        if (characterMatch === 'all') {
          if (!selectedCharacters.every((sc) => ids.includes(sc.id))) return false
        } else {
          if (!selectedCharacters.some((sc) => ids.includes(sc.id))) return false
        }
      }
      return true
    })
    .sort((a, b) => {
      if (sort === 'newest') {
        return new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
      }
      if (a.title === null && b.title === null) return 0
      if (a.title === null) return 1
      if (b.title === null) return -1
      return a.title.localeCompare(b.title)
    })
}

export function useArtpiecesFilter(artpieces: ArtpieceSummary[]) {
  const [searchParams, setSearchParams] = useSearchParams()
  const artistsQuery = useArtists()
  const charactersQuery = useCharacters()

  const [search, setSearch] = useState('')
  const [sort, setSort] = useState<SortMode>('newest')

  // Derive filter state from URL params + cached query data
  const artistParam = searchParams.get('artist')
  const selectedArtist: Artist | null = artistParam
    ? (artistsQuery.data?.find((a) => a.id === artistParam) ?? null)
    : null

  const charactersParam = searchParams.get('characters')
  const characterIdList = charactersParam ? charactersParam.split(',').filter(Boolean) : []
  const selectedCharacters = characterIdList
    .map((id) => charactersQuery.data?.find((c) => c.id === id))
    .filter((c): c is NonNullable<typeof c> => c != null)
    .map((c) => ({ id: c.id, name: c.name }))

  const matchParam = searchParams.get('match')
  const characterMatch: 'all' | 'any' = matchParam === 'any' ? 'any' : 'all'

  const refParam = searchParams.get('ref')
  const isReference = refParam === 'true'

  const hasFilterParams = artistParam !== null || charactersParam !== null || refParam !== null || matchParam !== null

  // pendingFilterRef always reflects the current filter state and is mutated
  // eagerly by each setter before it commits. React Router's setSearchParams
  // does not chain: the updater form always receives the params of the render
  // that created the setter, so several setters called in the same event (e.g.
  // "Clear all") would undo each other's deletions and the last navigation
  // would win. Committing the whole filter from the ref makes that last write
  // the full correct state instead.
  const pendingFilterRef = useRef<StoredFilter>({
    artistId: artistParam,
    characterIds: characterIdList,
    matchMode: characterMatch,
    isReference,
  })
  // Params not owned by the filter are carried over untouched on each commit.
  const paramsRef = useRef(searchParams)
  // Sync both refs to current URL state after each render so setters start from
  // a fresh base. useLayoutEffect (not render-time assignment) keeps the write
  // out of the render phase.
  useLayoutEffect(() => {
    pendingFilterRef.current = { artistId: artistParam, characterIds: characterIdList, matchMode: characterMatch, isReference }
    paramsRef.current = searchParams
  })

  function commitFilter(filter: StoredFilter) {
    pendingFilterRef.current = filter
    writeStoredFilter(filter)

    const next = new URLSearchParams(paramsRef.current)
    if (filter.artistId) next.set('artist', filter.artistId)
    else next.delete('artist')
    if (filter.characterIds.length > 0) next.set('characters', filter.characterIds.join(','))
    else next.delete('characters')
    if (filter.matchMode === 'any') next.set('match', 'any')
    else next.delete('match')
    if (filter.isReference) next.set('ref', 'true')
    else next.delete('ref')

    paramsRef.current = next
    setSearchParams(next, { replace: true })
  }

  // On mount: restore from localStorage when URL has no filter params, or write
  // the URL-derived state when arriving with params (e.g. from a saved-filter link).
  useEffect(() => {
    if (hasFilterParams) {
      writeStoredFilter({ artistId: artistParam, characterIds: characterIdList, matchMode: characterMatch, isReference })
    } else {
      const stored = readStoredFilter()
      if (!stored || (!stored.artistId && stored.characterIds.length === 0 && !stored.isReference)) return
      commitFilter({
        artistId: stored.artistId,
        characterIds: stored.characterIds,
        matchMode: stored.matchMode === 'any' ? 'any' : 'all',
        isReference: stored.isReference,
      })
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  function setIsReference(value: boolean) {
    commitFilter({ ...pendingFilterRef.current, isReference: value })
  }

  function setSelectedArtist(artist: Artist | null) {
    commitFilter({ ...pendingFilterRef.current, artistId: artist?.id ?? null })
  }

  function setSelectedCharacters(chars: { id: string; name: string }[]) {
    commitFilter({ ...pendingFilterRef.current, characterIds: chars.map((c) => c.id) })
  }

  function setCharacterMatch(mode: 'all' | 'any') {
    commitFilter({ ...pendingFilterRef.current, matchMode: mode })
  }

  const filtered = filterAndSortArtpieces(artpieces, {
    search,
    sort,
    selectedArtist,
    selectedCharacters,
    characterMatch,
    isReference,
  })

  return {
    search, setSearch,
    sort, setSort,
    selectedArtist, setSelectedArtist,
    selectedCharacters, setSelectedCharacters,
    characterMatch, setCharacterMatch,
    isReference, setIsReference,
    filtered,
  }
}
