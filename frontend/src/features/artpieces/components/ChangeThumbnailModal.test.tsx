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

function makeArtpiece(id: string, title: string): ArtpieceSummary {
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
    characters: [],
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
  }
}

function makeFilter(overrides: Partial<SavedFilter> = {}): SavedFilter {
  return {
    id: 'filter-1',
    name: 'My Filter',
    thumbnail_url: null,
    filter_payload: { artist_id: 'artist-1', character_ids: [], match_mode: 'all', is_reference: false },
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
    ...overrides,
  }
}

function okJson(data: unknown) {
  return Promise.resolve({ ok: true, json: async () => data } as Response)
}

function setupFetch(artpieces: ArtpieceSummary[]) {
  vi.mocked(apiFetch).mockImplementation((url: string) => {
    if (url === '/artpieces') return okJson(artpieces)
    if (url === '/artists') return okJson([ARTIST])
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
