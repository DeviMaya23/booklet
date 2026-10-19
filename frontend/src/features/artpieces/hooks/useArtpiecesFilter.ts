import { useState } from 'react'
import { useSearchParams } from 'react-router-dom'
import { type ArtpieceSummary } from '../api/useArtpieces'
import { type Artist, useArtists } from '@/features/artists/api/useArtists'
import { useCharacters } from '@/features/characters/api/useCharacters'

export type SortMode = 'newest' | 'alpha'

interface FilterParams {
  search: string
  sort: SortMode
  selectedArtist: Artist | null
  selectedCharacters: { id: string; name: string }[]
  characterMatch: 'all' | 'any'
}

export function filterAndSortArtpieces(
  artpieces: ArtpieceSummary[],
  { search, sort, selectedArtist, selectedCharacters, characterMatch }: FilterParams,
): ArtpieceSummary[] {
  return artpieces
    .filter((a) => {
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

  function setSelectedArtist(artist: Artist | null) {
    setSearchParams(
      (prev) => {
        const next = new URLSearchParams(prev)
        if (artist) next.set('artist', artist.id)
        else next.delete('artist')
        return next
      },
      { replace: true },
    )
  }

  function setSelectedCharacters(chars: { id: string; name: string }[]) {
    setSearchParams(
      (prev) => {
        const next = new URLSearchParams(prev)
        if (chars.length > 0) next.set('characters', chars.map((c) => c.id).join(','))
        else next.delete('characters')
        return next
      },
      { replace: true },
    )
  }

  function setCharacterMatch(mode: 'all' | 'any') {
    setSearchParams(
      (prev) => {
        const next = new URLSearchParams(prev)
        if (mode === 'any') next.set('match', 'any')
        else next.delete('match')
        return next
      },
      { replace: true },
    )
  }

  const filtered = filterAndSortArtpieces(artpieces, {
    search,
    sort,
    selectedArtist,
    selectedCharacters,
    characterMatch,
  })

  return {
    search, setSearch,
    sort, setSort,
    selectedArtist, setSelectedArtist,
    selectedCharacters, setSelectedCharacters,
    characterMatch, setCharacterMatch,
    filtered,
  }
}
