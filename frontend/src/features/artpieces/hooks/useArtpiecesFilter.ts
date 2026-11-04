import { useState, useEffect, useLayoutEffect, useRef } from 'react'
import { useSearchParams } from 'react-router-dom'
import { type ArtpieceSummary } from '../api/useArtpieces'
import { type Artist, useArtists } from '@/features/artists/api/useArtists'
import { useCharacters } from '@/features/characters/api/useCharacters'
import { useTags } from '@/features/tags/api/useTags'
import { type Verse, useVerses } from '@/features/verses/api/useVerses'

export type SortMode = 'newest' | 'alpha'

const FILTER_STORAGE_KEY = 'artpieces.filterState'

interface StoredFilter {
  artistId: string | null
  characterIds: string[]
  matchMode: 'all' | 'any'
  isReference: boolean
  verseId: string | null
  tagIds: string[]
  tagMatch: 'all' | 'any'
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
  selectedTags?: { id: string; name: string }[]
  tagMatch?: 'all' | 'any'
  isReference: boolean
  selectedVerse?: Verse | null
  allCharacters?: { id: string; verse_id: string | null }[]
  charactersPending?: boolean
}

export function filterAndSortArtpieces(
  artpieces: ArtpieceSummary[],
  { search, sort, selectedArtist, selectedCharacters, characterMatch, selectedTags = [], tagMatch = 'all', isReference, selectedVerse, allCharacters, charactersPending }: FilterParams,
): ArtpieceSummary[] {
  // allCharacters is still loading: membership of the verse is unknown, not
  // empty, so the verse filter must not exclude anything until it resolves.
  const verseCharacterIds = selectedVerse && !charactersPending
    ? new Set((allCharacters ?? []).filter((c) => c.verse_id === selectedVerse.id).map((c) => c.id))
    : null

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
      if (selectedTags.length > 0) {
        const ids = a.tags.map((t) => t.id)
        if (tagMatch === 'all') {
          if (!selectedTags.every((st) => ids.includes(st.id))) return false
        } else {
          if (!selectedTags.some((st) => ids.includes(st.id))) return false
        }
      }
      if (verseCharacterIds) {
        if (!a.characters.some((c) => verseCharacterIds.has(c.id))) return false
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
  const versesQuery = useVerses()
  const tagsQuery = useTags({ sort: 'recent' })

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

  const verseParam = searchParams.get('verse')
  const selectedVerse: Verse | null = verseParam
    ? (versesQuery.data?.find((v) => v.id === verseParam) ?? null)
    : null

  // Tag ids that no longer resolve (deleted tags) are dropped here, and the
  // tag filter is not applied until the tag list has loaded. If the list
  // failed to load nothing can be resolved, so the ids are applied as given
  // rather than silently showing every artpiece under a "tagged" URL.
  const tagsParam = searchParams.get('tags')
  const tagIdList = tagsParam ? tagsParam.split(',').filter(Boolean) : []
  const selectedTags = tagsQuery.isError
    ? tagIdList.map((id) => ({ id, name: '' }))
    : tagIdList
        .map((id) => tagsQuery.data?.find((t) => t.id === id))
        .filter((t): t is NonNullable<typeof t> => t != null)
        .map((t) => ({ id: t.id, name: t.name }))

  const tagMatchParam = searchParams.get('tag_match')
  const tagMatch: 'all' | 'any' = tagMatchParam === 'any' ? 'any' : 'all'

  const hasFilterParams =
    artistParam !== null ||
    charactersParam !== null ||
    refParam !== null ||
    matchParam !== null ||
    verseParam !== null ||
    tagsParam !== null ||
    tagMatchParam !== null

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
    verseId: verseParam,
    tagIds: tagIdList,
    tagMatch,
  })
  // Params not owned by the filter are carried over untouched on each commit.
  const paramsRef = useRef(searchParams)
  // Sync both refs to current URL state after each render so setters start from
  // a fresh base. useLayoutEffect (not render-time assignment) keeps the write
  // out of the render phase.
  useLayoutEffect(() => {
    pendingFilterRef.current = { artistId: artistParam, characterIds: characterIdList, matchMode: characterMatch, isReference, verseId: verseParam, tagIds: tagIdList, tagMatch }
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
    if (filter.verseId) next.set('verse', filter.verseId)
    else next.delete('verse')
    if (filter.tagIds.length > 0) next.set('tags', filter.tagIds.join(','))
    else next.delete('tags')
    if (filter.tagMatch === 'any') next.set('tag_match', 'any')
    else next.delete('tag_match')

    paramsRef.current = next
    setSearchParams(next, { replace: true })
  }

  // On mount: restore from localStorage when URL has no filter params, or write
  // the URL-derived state when arriving with params (e.g. from a saved-filter link).
  useEffect(() => {
    if (hasFilterParams) {
      writeStoredFilter({ artistId: artistParam, characterIds: characterIdList, matchMode: characterMatch, isReference, verseId: verseParam, tagIds: tagIdList, tagMatch })
    } else {
      const stored = readStoredFilter()
      // State stored before tags existed has no tag fields.
      const storedTagIds = Array.isArray(stored?.tagIds) ? stored.tagIds : []
      if (!stored || (!stored.artistId && stored.characterIds.length === 0 && !stored.isReference && !stored.verseId && storedTagIds.length === 0)) return
      commitFilter({
        artistId: stored.artistId,
        characterIds: stored.characterIds,
        matchMode: stored.matchMode === 'any' ? 'any' : 'all',
        isReference: stored.isReference,
        verseId: stored.verseId ?? null,
        tagIds: storedTagIds,
        tagMatch: stored.tagMatch === 'any' ? 'any' : 'all',
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

  function setSelectedTags(tags: { id: string; name: string }[]) {
    // Emptying the tags also resets the match mode, so no stale tag_match=any
    // survives to silently apply Any the next time tags are added.
    commitFilter({
      ...pendingFilterRef.current,
      tagIds: tags.map((t) => t.id),
      tagMatch: tags.length === 0 ? 'all' : pendingFilterRef.current.tagMatch,
    })
  }

  function setTagMatch(mode: 'all' | 'any') {
    commitFilter({ ...pendingFilterRef.current, tagMatch: mode })
  }

  function setSelectedVerse(verse: Verse | null) {
    commitFilter({ ...pendingFilterRef.current, verseId: verse?.id ?? null })
  }

  const filtered = filterAndSortArtpieces(artpieces, {
    search,
    sort,
    selectedArtist,
    selectedCharacters,
    characterMatch,
    selectedTags,
    tagMatch,
    isReference,
    selectedVerse,
    allCharacters: charactersQuery.data,
    charactersPending: charactersQuery.isPending,
  })

  return {
    search, setSearch,
    sort, setSort,
    selectedArtist, setSelectedArtist,
    selectedCharacters, setSelectedCharacters,
    characterMatch, setCharacterMatch,
    selectedTags, setSelectedTags,
    tagMatch, setTagMatch,
    isReference, setIsReference,
    selectedVerse, setSelectedVerse,
    filtered,
  }
}
