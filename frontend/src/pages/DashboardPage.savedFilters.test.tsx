import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { vi, describe, it, expect, beforeEach } from 'vitest'
import { MemoryRouter, Route, Routes, useLocation } from 'react-router-dom'
import DashboardPage from './DashboardPage'
import ArtpiecesPage from './ArtpiecesPage'
import { type ArtpieceSummary } from '@/features/artpieces/api/useArtpieces'
import { createWrapper } from '@/test/queryClient'
import { apiFetch } from '@/lib/api'

vi.mock('@kinde-oss/kinde-auth-react', () => ({
  useKindeAuth: () => ({ getToken: vi.fn().mockResolvedValue('tok') }),
}))
vi.mock('@/lib/api', () => ({ apiFetch: vi.fn() }))
vi.mock('@/features/files/components/AddFilesModal', () => ({ default: () => null }))
vi.mock('@/features/files/components/ArtpieceFormModal', () => ({ default: () => null }))
vi.mock('@/features/artpieces/components/ArtpieceDetailView', () => ({ default: () => null }))
vi.mock('@/features/files/components/FileViewer', () => ({ default: () => null }))
vi.mock('sonner', () => ({ toast: { success: vi.fn(), error: vi.fn() } }))

const outfit = { id: 't-outfit', name: 'outfit' }
const sketch = { id: 't-sketch', name: 'sketch' }

function makeArtpiece(id: string, title: string, tags: { id: string; name: string }[] = []): ArtpieceSummary {
  return {
    id,
    title,
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
    tags,
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
  }
}

const artpieces = [
  makeArtpiece('ap-both', 'Both tags', [outfit, sketch]),
  makeArtpiece('ap-outfit', 'Outfit only', [outfit]),
  makeArtpiece('ap-plain', 'Plain'),
]

type Payload = Record<string, unknown>

function savedFilter(name: string, payload: Payload) {
  return {
    id: `sf-${name}`,
    name,
    thumbnail_url: null,
    filter_payload: { artist_id: null, character_ids: [], match_mode: 'all', is_reference: false, verse_id: null, ...payload },
    created_at: '',
    updated_at: '',
  }
}

function mockApi(filters: ReturnType<typeof savedFilter>[]) {
  vi.mocked(apiFetch).mockImplementation(async (path: string) => {
    const ok = (body: unknown) => ({ ok: true, status: 200, json: async () => body }) as Response
    if (path === '/dashboard') {
      return ok({
        recent_artpieces: [],
        in_progress: [],
        housekeeping: { commissions_no_artist: [], artpieces_no_artist: [], done_no_artpieces: [], artpieces_no_files: [] },
      })
    }
    if (path === '/saved_filters') return ok(filters)
    if (path === '/artpieces') return ok(artpieces)
    if (path.startsWith('/tags')) return ok([outfit, sketch])
    return ok([])
  })
}

function LocationProbe() {
  const location = useLocation()
  return <output aria-label="location">{location.pathname + location.search}</output>
}

function renderDashboard() {
  const Wrapper = createWrapper()
  return render(
    <Wrapper>
      <MemoryRouter initialEntries={['/app/dashboard']}>
        <LocationProbe />
        <Routes>
          <Route path="/app/dashboard" element={<DashboardPage />} />
          <Route path="/app/artpieces" element={<ArtpiecesPage />} />
        </Routes>
      </MemoryRouter>
    </Wrapper>,
  )
}

async function openTile(name: string) {
  await userEvent.click(await screen.findByRole('link', { name }))
}

function currentParams() {
  const loc = screen.getByLabelText('location').textContent!
  return new URL(loc, 'http://x').searchParams
}

beforeEach(() => {
  vi.clearAllMocks()
  localStorage.clear()
})

describe('DashboardPage — saved filter tiles with tags', () => {
  it('navigates with tags and tag_match=any and applies them in the gallery', async () => {
    mockApi([savedFilter('Any tags', { tag_ids: ['t-outfit', 't-sketch'], tag_match_mode: 'any' })])
    renderDashboard()

    await openTile('Any tags')

    await screen.findByRole('button', { name: 'View Outfit only' })
    expect(currentParams().get('tags')).toBe('t-outfit,t-sketch')
    expect(currentParams().get('tag_match')).toBe('any')
    expect(screen.getByRole('button', { name: 'View Both tags' })).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'View Plain' })).not.toBeInTheDocument()
  })

  it('navigates with tags only for All mode', async () => {
    mockApi([savedFilter('All tags', { tag_ids: ['t-outfit', 't-sketch'], tag_match_mode: 'all' })])
    renderDashboard()

    await openTile('All tags')

    await screen.findByRole('button', { name: 'View Both tags' })
    expect(currentParams().get('tags')).toBe('t-outfit,t-sketch')
    expect(currentParams().has('tag_match')).toBe(false)
    expect(screen.queryByRole('button', { name: 'View Outfit only' })).not.toBeInTheDocument()
  })

  it('adds no tag params for a filter saved before tags existed', async () => {
    mockApi([savedFilter('Old filter', {})])
    renderDashboard()

    await openTile('Old filter')

    await screen.findByRole('button', { name: 'View Plain' })
    expect(currentParams().has('tags')).toBe(false)
    expect(currentParams().has('tag_match')).toBe(false)
  })

  it('opens the gallery with a deleted tag ignored and the rest applied', async () => {
    mockApi([savedFilter('Stale tag', { tag_ids: ['t-outfit', 't-deleted'], tag_match_mode: 'all' })])
    renderDashboard()

    await openTile('Stale tag')

    await screen.findByRole('button', { name: 'View Outfit only' })
    expect(screen.getByRole('button', { name: 'View Both tags' })).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'View Plain' })).not.toBeInTheDocument()
  })
})
