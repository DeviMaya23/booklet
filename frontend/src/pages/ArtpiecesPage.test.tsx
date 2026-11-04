import { render, screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { vi, describe, it, expect, beforeEach } from 'vitest'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import ArtpiecesPage from './ArtpiecesPage'
import { type ArtpieceSummary } from '@/features/artpieces/api/useArtpieces'
import { createWrapper } from '@/test/queryClient'
import { apiFetch } from '@/lib/api'

vi.mock('@kinde-oss/kinde-auth-react', () => ({
  useKindeAuth: () => ({ getToken: vi.fn().mockResolvedValue('tok') }),
}))
vi.mock('@/lib/api', () => ({ apiFetch: vi.fn() }))
vi.mock('@/features/files/components/ArtpieceFormModal', () => ({ default: () => null }))
vi.mock('@/features/artpieces/components/ArtpieceDetailView', () => ({ default: () => null }))
vi.mock('@/features/files/components/FileViewer', () => ({ default: () => null }))
vi.mock('sonner', () => ({ toast: { success: vi.fn(), error: vi.fn() } }))

const outfit = { id: 't-outfit', name: 'outfit' }
const sketch = { id: 't-sketch', name: 'sketch' }
const allTags = [outfit, sketch]

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

function mockApi() {
  vi.mocked(apiFetch).mockImplementation(async (path: string, _t: unknown, init?: RequestInit) => {
    const ok = (body: unknown, status = 200) => ({ ok: true, status, json: async () => body }) as Response
    if (path === '/saved_filters' && init?.method === 'POST') {
      return ok({ id: 'sf-1', name: 'Mine', thumbnail_url: null, filter_payload: {}, created_at: '', updated_at: '' }, 201)
    }
    if (path === '/artpieces') return ok(artpieces)
    if (path.startsWith('/tags')) return ok(allTags)
    return ok([])
  })
}

function renderPage(initialUrl = '/app/artpieces') {
  const Wrapper = createWrapper()
  return render(
    <Wrapper>
      <MemoryRouter initialEntries={[initialUrl]}>
        <Routes>
          <Route path="/app/artpieces" element={<ArtpiecesPage />} />
        </Routes>
      </MemoryRouter>
    </Wrapper>,
  )
}

function tile(title: string) {
  return screen.queryByRole('button', { name: `View ${title}` })
}

async function openFilter(name: RegExp | string = /filter/i) {
  await userEvent.click(screen.getByRole('button', { name }))
}

function tagsInput() {
  const section = screen.getByText('Tags').parentElement!.parentElement!
  return within(section).getByRole('textbox')
}

async function pickTag(name: string) {
  const input = tagsInput()
  await waitFor(() => expect(input).toBeEnabled())
  await userEvent.type(input, name.slice(0, 3))
  await userEvent.click(await screen.findByText(name))
}

beforeEach(() => {
  vi.clearAllMocks()
  localStorage.clear()
  mockApi()
})

describe('ArtpiecesPage — tag filter', () => {
  it('selecting a tag narrows the grid to artpieces that have it', async () => {
    renderPage()
    await screen.findByRole('button', { name: 'View Plain' })

    await openFilter()
    await pickTag('sketch')

    await waitFor(() => expect(tile('Plain')).toBeNull())
    expect(tile('Outfit only')).toBeNull()
    expect(tile('Both tags')).not.toBeNull()
  })

  it('applies tags from the URL, matching All by default and Any with tag_match=any', async () => {
    const first = renderPage('/app/artpieces?tags=t-outfit,t-sketch')
    await screen.findByRole('button', { name: 'View Both tags' })
    expect(tile('Outfit only')).toBeNull()
    first.unmount()

    renderPage('/app/artpieces?tags=t-outfit,t-sketch&tag_match=any')
    await screen.findByRole('button', { name: 'View Outfit only' })
    expect(tile('Both tags')).not.toBeNull()
    expect(tile('Plain')).toBeNull()
  })

  it('ignores a deleted tag id in the URL instead of emptying the grid', async () => {
    renderPage('/app/artpieces?tags=t-outfit,t-deleted')

    await screen.findByRole('button', { name: 'View Outfit only' })
    expect(tile('Both tags')).not.toBeNull()
    expect(tile('Plain')).toBeNull()
  })

  it('shows the tags in the trigger label', async () => {
    renderPage('/app/artpieces?tags=t-outfit,t-sketch&tag_match=any')

    expect(await screen.findByRole('button', { name: '2 tags, any' })).toBeInTheDocument()
  })
})

describe('ArtpiecesPage — saving a tag filter', () => {
  it('offers "Save to dashboard" for a tag-only filter', async () => {
    renderPage('/app/artpieces?tags=t-outfit')
    await screen.findByRole('button', { name: 'View Outfit only' })

    await openFilter('1 tag, all')

    expect(screen.getByRole('button', { name: 'Save to dashboard' })).toBeInTheDocument()
  })

  it('does not offer "Save to dashboard" with no filters', async () => {
    renderPage()
    await screen.findByRole('button', { name: 'View Plain' })

    await openFilter()

    expect(screen.queryByRole('button', { name: 'Save to dashboard' })).not.toBeInTheDocument()
  })

  async function saveFilter() {
    await userEvent.click(screen.getByRole('button', { name: 'Save to dashboard' }))
    await userEvent.type(await screen.findByRole('textbox', { name: /name/i }), 'Mine')
    await userEvent.click(screen.getByRole('button', { name: 'Save' }))
    await waitFor(() => expect(screen.queryByRole('dialog')).not.toBeInTheDocument())
  }

  function savedPayload() {
    const call = vi.mocked(apiFetch).mock.calls.find((c) => c[0] === '/saved_filters' && c[2]?.method === 'POST')
    return JSON.parse(call![2]!.body as string).filter_payload
  }

  it('saves only tag ids that still exist', async () => {
    renderPage('/app/artpieces?tags=t-outfit,t-deleted&tag_match=any')
    await screen.findByRole('button', { name: 'View Outfit only' })
    await openFilter('1 tag, any')

    await saveFilter()

    expect(savedPayload()).toEqual(expect.objectContaining({ tag_ids: ['t-outfit'], tag_match_mode: 'any' }))
  })

  it('shows "Saved as" after saving and clears it when the tags change', async () => {
    renderPage('/app/artpieces?tags=t-outfit')
    await screen.findByRole('button', { name: 'View Outfit only' })
    await openFilter('1 tag, all')
    await saveFilter()

    await openFilter('1 tag, all')
    expect(await screen.findByText(/Saved as/)).toBeInTheDocument()

    await pickTag('sketch')

    await waitFor(() => expect(screen.queryByText(/Saved as/)).not.toBeInTheDocument())
    expect(screen.getByRole('button', { name: 'Save to dashboard' })).toBeInTheDocument()
  })

  it('clears "Saved as" when the tag match mode changes', async () => {
    renderPage('/app/artpieces?tags=t-outfit,t-sketch')
    await screen.findByRole('button', { name: 'View Both tags' })
    await openFilter('2 tags, all')
    await saveFilter()

    await openFilter('2 tags, all')
    expect(await screen.findByText(/Saved as/)).toBeInTheDocument()

    const section = screen.getByText('Tags').parentElement!.parentElement!
    await userEvent.click(within(section).getByRole('button', { name: 'Any' }))

    await waitFor(() => expect(screen.queryByText(/Saved as/)).not.toBeInTheDocument())
  })
})

describe('ArtpiecesPage — clearing tags', () => {
  function storedFilter() {
    return JSON.parse(localStorage.getItem('artpieces.filterState')!)
  }

  it('"Clear all" removes the tags and resets the stored tag match mode', async () => {
    renderPage('/app/artpieces?tags=t-outfit,t-sketch&tag_match=any')
    await screen.findByRole('button', { name: 'View Outfit only' })
    await openFilter('2 tags, any')

    await userEvent.click(screen.getByRole('button', { name: /clear all/i }))

    await waitFor(() => expect(storedFilter()).toEqual(expect.objectContaining({ tagIds: [], tagMatch: 'all' })))
    expect(await screen.findByRole('button', { name: 'View Plain' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Filter' })).toBeInTheDocument()
  })

  it('removing the last tag chip resets the stored tag match mode', async () => {
    renderPage('/app/artpieces?tags=t-outfit&tag_match=any')
    await screen.findByRole('button', { name: 'View Outfit only' })
    await openFilter('1 tag, any')

    await userEvent.click(screen.getByRole('button', { name: /remove outfit/i }))

    await waitFor(() => expect(storedFilter()).toEqual(expect.objectContaining({ tagIds: [], tagMatch: 'all' })))
  })
})
