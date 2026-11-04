import { render, screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { vi, describe, it, expect, beforeEach } from 'vitest'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { createMemoryRouter, Link, RouterProvider } from 'react-router-dom'
import ArtpieceDetailView from './ArtpieceDetailView'
import { type ArtpieceDetail } from '../api/useArtpiece'
import { apiFetch } from '@/lib/api'

vi.mock('@kinde-oss/kinde-auth-react', () => ({
  useKindeAuth: () => ({ getToken: vi.fn().mockResolvedValue('tok') }),
}))
vi.mock('@/lib/api', () => ({ apiFetch: vi.fn() }))

vi.mock('../api/useArtpiece', () => ({
  useArtpiece: () => ({ data: mockArtpiece, isLoading: false, isError: false }),
  artpieceQueryKey: (id: string) => ['artpieces', id],
}))

const mockMutateAsync = vi.fn().mockResolvedValue({})
vi.mock('../api/useUpdateArtpiece', () => ({
  useUpdateArtpiece: () => ({ mutateAsync: mockMutateAsync, isPending: false }),
}))
vi.mock('../api/useSetCover', () => ({ useSetCover: () => ({ mutateAsync: vi.fn(), isPending: false }) }))
vi.mock('../api/useAttachFilesToArtpiece', () => ({ useAttachFilesToArtpiece: () => ({ mutateAsync: vi.fn(), isPending: false }) }))
vi.mock('../api/useDownloadArtpiece', () => ({ useDownloadArtpiece: () => ({ mutateAsync: vi.fn(), isPending: false }) }))
vi.mock('@/features/files/api/useInitFileUpload', () => ({ useInitFileUpload: () => ({ mutateAsync: vi.fn(), isPending: false }) }))
vi.mock('@/features/files/api/useCompleteFileUpload', () => ({ useCompleteFileUpload: () => ({ mutateAsync: vi.fn(), isPending: false }) }))
vi.mock('@/features/characters/api/useCharacters', () => ({ useCharacters: () => ({ data: [], isPending: false }) }))
vi.mock('@/features/artists/components/ArtistCombobox', () => ({ default: () => null }))
vi.mock('@/features/artists/components/ArtistFormModal', () => ({ default: () => null }))
vi.mock('./ArtpieceDetailFileGrid', () => ({ default: () => null }))
vi.mock('./DeleteArtpieceDialog', () => ({ default: () => null }))
vi.mock('@/features/files/components/FileViewer', () => ({ default: () => null }))
vi.mock('sonner', () => ({ toast: { success: vi.fn(), error: vi.fn() } }))

const outfit = { id: 't-outfit', name: 'outfit' }
const sketch = { id: 't-sketch', name: 'sketch' }
const wip = { id: 't-wip', name: 'wip' }
const allTags = [wip, outfit, sketch]

let mockArtpiece: ArtpieceDetail

const baseArtpiece: ArtpieceDetail = {
  id: 'ap1',
  title: 'Sunrise',
  artist_id: null,
  artist_name: null,
  artist_link: null,
  cover_file_id: null,
  thumbnail_url: null,
  notes: null,
  is_reference: false,
  characters: [],
  tags: [],
  files: [],
  created_at: '2024-01-01T00:00:00Z',
  updated_at: '2024-01-01T00:00:00Z',
}

beforeEach(() => {
  vi.clearAllMocks()
  mockMutateAsync.mockResolvedValue({})
  vi.mocked(apiFetch).mockImplementation(async (path: string, _t: unknown, init?: RequestInit) => {
    if (path === '/tags' && init?.method === 'POST') {
      return { ok: true, status: 201, json: async () => ({ id: 't-new', name: 'fresh' }) } as Response
    }
    if (path.startsWith('/tags')) return { ok: true, status: 200, json: async () => allTags } as Response
    return { ok: true, status: 200, json: async () => [] } as Response
  })
})

function renderView(artpiece: ArtpieceDetail) {
  mockArtpiece = artpiece
  const qc = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  const router = createMemoryRouter([
    {
      path: '/',
      element: (
        <>
          <Link to="/elsewhere">Leave</Link>
          <ArtpieceDetailView artpieceId="ap1" onClose={() => {}} onDeleted={() => {}} />
        </>
      ),
    },
    { path: '/elsewhere', element: <p>Elsewhere page</p> },
  ])
  return render(
    <QueryClientProvider client={qc}>
      <RouterProvider router={router} />
    </QueryClientProvider>,
  )
}

async function enterEdit() {
  await userEvent.click(screen.getByRole('button', { name: /^edit$/i }))
  await waitFor(() => expect(apiFetch).toHaveBeenCalledWith('/tags?sort=recent', expect.any(Function)))
}

function tagsField() {
  return within(screen.getByText('Tags').parentElement!)
}

async function pickTag(name: string) {
  await userEvent.click(tagsField().getByRole('textbox'))
  await userEvent.click(await screen.findByText(name))
}

describe('ArtpieceDetailView — tags in view mode', () => {
  it('shows each tag as a chip', () => {
    renderView({ ...baseArtpiece, tags: [outfit, sketch] })

    expect(screen.getByText('Tags')).toBeInTheDocument()
    expect(screen.getByText('outfit')).toBeInTheDocument()
    expect(screen.getByText('sketch')).toBeInTheDocument()
  })

  it('shows a dash when the artpiece has no tags', () => {
    renderView({ ...baseArtpiece, tags: [] })

    const section = screen.getByText('Tags').parentElement!
    expect(within(section).getByText('—')).toBeInTheDocument()
  })
})

describe('ArtpieceDetailView — tags in edit mode', () => {
  it('initialises the field with the artpiece tags', async () => {
    renderView({ ...baseArtpiece, tags: [outfit, sketch] })
    await enterEdit()

    expect(tagsField().getByText('outfit')).toBeInTheDocument()
    expect(tagsField().getByText('sketch')).toBeInTheDocument()
  })

  it('adding an existing tag counts as a pending change', async () => {
    renderView({ ...baseArtpiece, tags: [outfit] })
    await enterEdit()

    await pickTag('sketch')

    expect(screen.getByText(/tags changed/i)).toBeInTheDocument()
    expect(mockMutateAsync).not.toHaveBeenCalled()
  })

  it('an inline-created tag becomes a chip and a pending change without updating the artpiece', async () => {
    renderView({ ...baseArtpiece, tags: [] })
    await enterEdit()

    await userEvent.type(tagsField().getByRole('textbox'), 'fresh{Enter}')

    await waitFor(() => expect(tagsField().getByText('fresh')).toBeInTheDocument())
    expect(screen.getByText(/tags changed/i)).toBeInTheDocument()
    expect(mockMutateAsync).not.toHaveBeenCalled()
  })

  it('removing a chip counts as a pending change', async () => {
    renderView({ ...baseArtpiece, tags: [outfit] })
    await enterEdit()

    await userEvent.click(tagsField().getByRole('button', { name: /remove outfit/i }))

    expect(screen.getByText(/tags changed/i)).toBeInTheDocument()
  })

  it('removing a tag and adding it back is not a change', async () => {
    renderView({ ...baseArtpiece, tags: [outfit] })
    await enterEdit()

    await userEvent.click(tagsField().getByRole('button', { name: /remove outfit/i }))
    await pickTag('outfit')

    expect(screen.queryByText(/tags changed/i)).not.toBeInTheDocument()
  })

  it('excludes already-selected tags from suggestions and requests sort=recent', async () => {
    renderView({ ...baseArtpiece, tags: [outfit] })
    await enterEdit()

    await userEvent.click(tagsField().getByRole('textbox'))

    expect(await screen.findByText('sketch')).toBeInTheDocument()
    expect(screen.getAllByText('outfit')).toHaveLength(1) // chip only
  })
})

describe('ArtpieceDetailView — saving tags', () => {
  async function save() {
    await userEvent.click(screen.getByRole('button', { name: /^save$/i }))
  }

  it('sends the resulting tag ids when tags changed', async () => {
    renderView({ ...baseArtpiece, tags: [sketch] })
    await enterEdit()
    await userEvent.click(tagsField().getByRole('button', { name: /remove sketch/i }))
    await pickTag('wip')

    await save()

    expect(mockMutateAsync).toHaveBeenCalledWith(expect.objectContaining({ tagIds: ['t-wip'] }))
  })

  it('sends an empty tagIds array when every tag was removed', async () => {
    renderView({ ...baseArtpiece, tags: [sketch] })
    await enterEdit()
    await userEvent.click(tagsField().getByRole('button', { name: /remove sketch/i }))

    await save()

    expect(mockMutateAsync).toHaveBeenCalledWith(expect.objectContaining({ tagIds: [] }))
  })

  it('sends the existing tag ids when only the title changed', async () => {
    renderView({ ...baseArtpiece, tags: [outfit, sketch] })
    await enterEdit()
    const title = screen.getByDisplayValue('Sunrise')
    await userEvent.clear(title)
    await userEvent.type(title, 'Sunset')

    await save()

    expect(mockMutateAsync).toHaveBeenCalledWith(
      expect.objectContaining({ title: 'Sunset', tagIds: ['t-outfit', 't-sketch'] }),
    )
  })

  it('Cancel discards tag edits and view mode shows the original tags', async () => {
    renderView({ ...baseArtpiece, tags: [outfit] })
    await enterEdit()
    await userEvent.click(tagsField().getByRole('button', { name: /remove outfit/i }))

    await userEvent.click(screen.getByRole('button', { name: /^cancel$/i }))

    expect(screen.getByText('outfit')).toBeInTheDocument()
    expect(mockMutateAsync).not.toHaveBeenCalled()
  })
})

describe('ArtpieceDetailView — route blocker', () => {
  it('treats a tag-only edit as unsaved when navigating away', async () => {
    renderView({ ...baseArtpiece, tags: [outfit] })
    await enterEdit()
    await userEvent.click(tagsField().getByRole('button', { name: /remove outfit/i }))

    await userEvent.click(screen.getByRole('link', { name: 'Leave' }))

    expect(await screen.findByText('Discard changes?')).toBeInTheDocument()
    expect(screen.queryByText('Elsewhere page')).not.toBeInTheDocument()
  })

  it('does not block navigation when nothing changed', async () => {
    renderView({ ...baseArtpiece, tags: [outfit] })
    await enterEdit()

    await userEvent.click(screen.getByRole('link', { name: 'Leave' }))

    expect(await screen.findByText('Elsewhere page')).toBeInTheDocument()
  })
})
