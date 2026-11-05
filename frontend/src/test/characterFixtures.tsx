import { render } from '@testing-library/react'
import { createMemoryRouter, RouterProvider } from 'react-router-dom'
import CharacterDetailPage from '@/pages/CharacterDetailPage'
import type { Character } from '@/features/characters/api/useCharacters'
import type { ArtpieceSummary } from '@/features/artpieces/api/useArtpieces'
import type { Verse } from '@/features/verses/api/useVerses'
import { createWrapper } from './queryClient'

export function makeCharacter(overrides: Partial<Character> = {}): Character {
  return {
    id: 'c-1',
    name: 'Tanelle',
    avatar_url: null,
    notes: null,
    is_public: false,
    verse_id: null,
    folders: [],
    created_at: '',
    updated_at: '',
    ...overrides,
  }
}

export function makeVerse(overrides: Partial<Verse> = {}): Verse {
  return {
    id: 'v-1',
    name: 'Coffee Shop AU',
    sort_order: 0,
    last_used_at: null,
    created_at: '',
    updated_at: '',
    ...overrides,
  }
}

export function makeArtpiece(overrides: Partial<ArtpieceSummary> & { id: string }): ArtpieceSummary {
  return {
    title: overrides.id,
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
    created_at: '',
    updated_at: '',
    ...overrides,
  }
}

/** Shape returned by GET /folders. */
export function foldersResponse(folders: { id: string; name: string; token: string }[]) {
  return {
    folder_list: folders.map((f) => ({ folder_id: f.id, folder_name: f.name, token: f.token })),
  }
}

/**
 * Renders the character detail route inside a data router (needed by useBlocker),
 * with stub pages for the destinations the page can navigate to.
 */
export function renderCharacterDetail(characterId = 'c-1') {
  const router = createMemoryRouter(
    [
      { path: '/app/characters', element: <div>Characters list</div> },
      { path: '/app/characters/:id', element: <CharacterDetailPage /> },
      { path: '/app/artpieces/:id', element: <div>Artpiece detail</div> },
      { path: '/app/other', element: <div>Other page</div> },
    ],
    { initialEntries: [`/app/characters/${characterId}`] },
  )
  const Wrapper = createWrapper()
  const utils = render(
    <Wrapper>
      <RouterProvider router={router} />
    </Wrapper>,
  )
  return { router, ...utils }
}
