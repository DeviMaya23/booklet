import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { MemoryRouter } from 'react-router-dom'
import React from 'react'
import { vi, describe, it, expect, beforeEach } from 'vitest'
import ChangeThumbnailModal from './ChangeThumbnailModal'
import { type SavedFilter } from '../api/useCreateSavedFilter'
import { type ArtpieceSummary } from '../api/useArtpieces'
import { apiFetch } from '@/lib/api'

vi.mock('@kinde-oss/kinde-auth-react', () => ({
  useKindeAuth: vi.fn().mockReturnValue({ getToken: vi.fn().mockResolvedValue('tok') }),
}))

vi.mock('@/lib/api', () => ({
  apiFetch: vi.fn(),
}))

const ARTIST = { id: 'artist-1', name: 'Alice', notes: null, links: [], created_at: '', updated_at: '' }

function makeArtpiece(id: string, title: string, characters: { id: string; name: string }[] = []): ArtpieceSummary {
  return {
    id,
    title,
    artist_id: 'artist-1',
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
    characters,
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
  }
}

function makeFilter(overrides: Partial<SavedFilter> = {}): SavedFilter {
  return {
    id: 'filter-1',
    name: 'My Filter',
    thumbnail_url: null,
    filter_payload: { artist_id: 'artist-1', character_ids: [], match_mode: 'all', is_reference: false, verse_id: null },
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
    ...overrides,
  }
}

function okJson(data: unknown) {
  return Promise.resolve({ ok: true, json: async () => data } as Response)
}

type CharacterFixture = { id: string; name: string; verse_id: string | null }

function setupFetch(artpieces: ArtpieceSummary[], characters: CharacterFixture[] = [], verses: { id: string; name: string }[] = []) {
  vi.mocked(apiFetch).mockImplementation((url: string) => {
    if (url === '/artpieces') return okJson(artpieces)
    if (url === '/artists') return okJson([ARTIST])
    if (url === '/characters') {
      return okJson(characters.map((c) => ({
        id: c.id,
        name: c.name,
        avatar_url: null,
        notes: null,
        is_public: false,
        verse_id: c.verse_id,
        folders: [],
        created_at: '',
        updated_at: '',
      })))
    }
    if (url === '/verses') {
      return okJson(verses.map((v) => ({ id: v.id, name: v.name, sort_order: 0, last_used_at: null, created_at: '', updated_at: '' })))
    }
    return okJson({})
  })
}

function wrapper({ children }: { children: React.ReactNode }) {
  const qc = new QueryClient({
    defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
  })
  return (
    <QueryClientProvider client={qc}>
      <MemoryRouter>{children}</MemoryRouter>
    </QueryClientProvider>
  )
}

const onOpenChange = vi.fn()
const onSaved = vi.fn()

beforeEach(() => {
  vi.clearAllMocks()
})

describe('ChangeThumbnailModal', () => {
  it('renders up to 10 filtered artpieces', async () => {
    setupFetch(Array.from({ length: 15 }, (_, i) => makeArtpiece(`a-${i}`, `Art ${i}`)))
    render(<ChangeThumbnailModal filter={makeFilter()} open onOpenChange={onOpenChange} onSaved={onSaved} />, { wrapper })
    const buttons = await screen.findAllByRole('button', { name: /art \d/i })
    expect(buttons).toHaveLength(10)
  })

  it('pre-selects the first artpiece', async () => {
    setupFetch([makeArtpiece('a-1', 'First'), makeArtpiece('a-2', 'Second')])
    render(<ChangeThumbnailModal filter={makeFilter()} open onOpenChange={onOpenChange} onSaved={onSaved} />, { wrapper })
    const firstBtn = await screen.findByRole('button', { name: 'First' })
    expect(firstBtn).toHaveAttribute('aria-pressed', 'true')
    expect(screen.getByRole('button', { name: 'Second' })).toHaveAttribute('aria-pressed', 'false')
  })

  it('renders empty state and Save is disabled when no artpieces match the filter', async () => {
    setupFetch([])
    render(<ChangeThumbnailModal filter={makeFilter()} open onOpenChange={onOpenChange} onSaved={onSaved} />, { wrapper })
    expect(await screen.findByText(/no artpieces match this filter/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /save/i })).toBeDisabled()
  })

  it('clicking Cancel calls onOpenChange(false) without sending a PATCH', async () => {
    const user = userEvent.setup()
    setupFetch([makeArtpiece('a-1', 'Art')])
    render(<ChangeThumbnailModal filter={makeFilter()} open onOpenChange={onOpenChange} onSaved={onSaved} />, { wrapper })
    await screen.findByRole('button', { name: 'Art' })
    await user.click(screen.getByRole('button', { name: /cancel/i }))
    expect(onOpenChange).toHaveBeenCalledWith(false)
    expect(vi.mocked(apiFetch)).not.toHaveBeenCalledWith(
      expect.stringContaining('/saved_filters'),
      expect.anything(),
      expect.objectContaining({ method: 'PATCH' }),
    )
  })
})

describe('ChangeThumbnailModal — verse filter', () => {
  it('excludes artpieces with no character in the saved filter verse', async () => {
    setupFetch(
      [
        makeArtpiece('a-in-verse', 'In verse', [{ id: 'char-1', name: 'Alpha' }]),
        makeArtpiece('a-out-verse', 'Out of verse', [{ id: 'char-2', name: 'Beta' }]),
      ],
      [
        { id: 'char-1', name: 'Alpha', verse_id: 'verse-1' },
        { id: 'char-2', name: 'Beta', verse_id: 'verse-2' },
      ],
      [{ id: 'verse-1', name: 'Elowen' }],
    )
    render(
      <ChangeThumbnailModal
        filter={makeFilter({ filter_payload: { artist_id: null, character_ids: [], match_mode: 'all', is_reference: false, verse_id: 'verse-1' } })}
        open
        onOpenChange={onOpenChange}
        onSaved={onSaved}
      />,
      { wrapper },
    )
    expect(await screen.findByRole('button', { name: 'In verse' })).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'Out of verse' })).not.toBeInTheDocument()
  })

  it('includes "verse" in the filter chip label when verse_id is set', async () => {
    setupFetch(
      [makeArtpiece('a-1', 'Art', [{ id: 'char-1', name: 'Alpha' }])],
      [{ id: 'char-1', name: 'Alpha', verse_id: 'verse-1' }],
      [{ id: 'verse-1', name: 'Elowen' }],
    )
    render(
      <ChangeThumbnailModal
        filter={makeFilter({ filter_payload: { artist_id: null, character_ids: [], match_mode: 'all', is_reference: false, verse_id: 'verse-1' } })}
        open
        onOpenChange={onOpenChange}
        onSaved={onSaved}
      />,
      { wrapper },
    )
    expect(await screen.findByText(/verse/i)).toBeInTheDocument()
  })

  it('shows a loading spinner (not the empty state) while characters are still loading', async () => {
    const artpieces = [makeArtpiece('a-in-verse', 'In verse', [{ id: 'char-1', name: 'Alpha' }])]
    let resolveCharacters!: (value: Response) => void
    vi.mocked(apiFetch).mockImplementation((url: string) => {
      if (url === '/artpieces') return okJson(artpieces)
      if (url === '/artists') return okJson([ARTIST])
      if (url === '/characters') return new Promise((resolve) => { resolveCharacters = resolve })
      if (url === '/verses') return okJson([{ id: 'verse-1', name: 'Elowen', sort_order: 0, last_used_at: null, created_at: '', updated_at: '' }])
      return okJson({})
    })

    render(
      <ChangeThumbnailModal
        filter={makeFilter({ filter_payload: { artist_id: null, character_ids: [], match_mode: 'all', is_reference: false, verse_id: 'verse-1' } })}
        open
        onOpenChange={onOpenChange}
        onSaved={onSaved}
      />,
      { wrapper },
    )

    // Artpieces have resolved by now, but characters haven't — the modal
    // must still show a spinner, not "No artpieces match this filter".
    await screen.findByText('Filter:')
    expect(screen.queryByText(/no artpieces match this filter/i)).not.toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'In verse' })).not.toBeInTheDocument()

    resolveCharacters(await okJson([{ id: 'char-1', name: 'Alpha', avatar_url: null, notes: null, is_public: false, verse_id: 'verse-1', folders: [], created_at: '', updated_at: '' }]))
    expect(await screen.findByRole('button', { name: 'In verse' })).toBeInTheDocument()
  })
})
