import { useState } from 'react'
import { type ArtpieceSummary } from '../api/useArtpieces'
import { type Artist } from '@/features/artists/api/useArtists'

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
  const [search, setSearch] = useState('')
  const [sort, setSort] = useState<SortMode>('newest')
  const [selectedArtist, setSelectedArtist] = useState<Artist | null>(null)
  const [selectedCharacters, setSelectedCharacters] = useState<{ id: string; name: string }[]>([])
  const [characterMatch, setCharacterMatch] = useState<'all' | 'any'>('all')

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
