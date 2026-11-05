import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { vi, describe, it, expect, beforeEach } from 'vitest'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { createMemoryRouter, RouterProvider } from 'react-router-dom'
import ArtpieceDetailView from './ArtpieceDetailView'
import { type ArtpieceDetail, type ArtpieceFile } from '../api/useArtpiece'

vi.mock('@kinde-oss/kinde-auth-react', () => ({
  useKindeAuth: () => ({ getToken: vi.fn().mockResolvedValue('tok') }),
}))
vi.mock('../api/useArtpiece', () => ({
  useArtpiece: () => ({ data: mockArtpiece, isLoading: false, isError: false }),
  artpieceQueryKey: (id: string) => ['artpieces', id],
}))
vi.mock('../api/useUpdateArtpiece', () => ({
  useUpdateArtpiece: () => ({ mutateAsync: vi.fn(), isPending: false }),
}))
const mockSetCover = vi.fn()
vi.mock('../api/useSetCover', () => ({ useSetCover: () => ({ mutateAsync: mockSetCover, isPending: false }) }))
vi.mock('../api/useAttachFilesToArtpiece', () => ({ useAttachFilesToArtpiece: () => ({ mutateAsync: vi.fn(), isPending: false }) }))
vi.mock('../api/useDownloadArtpiece', () => ({ useDownloadArtpiece: () => ({ mutate: vi.fn(), isPending: false }) }))
vi.mock('@/features/files/api/useInitFileUpload', () => ({ useInitFileUpload: () => ({ mutateAsync: vi.fn(), isPending: false }) }))
vi.mock('@/features/files/api/useCompleteFileUpload', () => ({ useCompleteFileUpload: () => ({ mutateAsync: vi.fn(), isPending: false }) }))
vi.mock('@/features/characters/api/useCharacters', () => ({ useCharacters: () => ({ data: [], isPending: false }) }))
vi.mock('@/features/artists/components/ArtistCombobox', () => ({ default: () => null }))
vi.mock('@/features/artists/components/ArtistFormModal', () => ({ default: () => null }))
vi.mock('@/features/tags/components/TagPicker', () => ({ default: () => null }))
vi.mock('./ArtpieceDetailFileGrid', () => ({
  default: ({
    files,
    mode,
    onSetCover,
    onOpen,
  }: {
    files: ArtpieceFile[]
    mode: string
    onSetCover?: (id: string) => void
    onOpen?: (id: string) => void
  }) =>
    mode === 'view' ? (
      <div>
        {files.filter((f) => f.is_visible).map((f) => (
          <button key={f.id} type="button" onClick={() => onOpen?.(f.id)}>
            open {f.id}
          </button>
        ))}
      </div>
    ) : mode === 'edit' ? (
      <div>
        {files.filter((f) => f.is_visible).map((f) => (
          <button key={f.id} type="button" onClick={() => onSetCover?.(f.id)}>
            cover {f.id}
          </button>
        ))}
      </div>
    ) : null,
}))
vi.mock('./DeleteArtpieceDialog', () => ({ default: () => null }))
vi.mock('@/features/files/components/FileViewer', () => ({
  default: ({ files, initialIndex }: { files: { id: string }[]; initialIndex: number }) => (
    <div data-testid="viewer">{files[initialIndex]?.id}</div>
  ),
}))
vi.mock('./ArtpieceFileDetailsModal', () => ({
  default: ({ open, files }: { open: boolean; files: ArtpieceFile[] }) =>
    open ? <div role="dialog" aria-label="File details">{files.length} rows</div> : null,
}))
vi.mock('sonner', () => ({ toast: { success: vi.fn(), error: vi.fn() } }))

let mockArtpiece: ArtpieceDetail

function makeFile(id: string, is_visible: boolean): ArtpieceFile {
  return {
    id,
    file_url: `https://cdn/${id}`,
    thumbnail_url: null,
    mime_type: 'image/png',
    name: id,
    notes: null,
    is_visible,
    extension: 'png',
    width: null,
    height: null,
  }
}

function artpieceWith(files: ArtpieceFile[]): ArtpieceDetail {
  return {
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
    files,
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
  }
}

function renderView() {
  const qc = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  const router = createMemoryRouter([
    { path: '/', element: <ArtpieceDetailView artpieceId="ap1" onClose={() => {}} onDeleted={() => {}} /> },
  ])
  return render(
    <QueryClientProvider client={qc}>
      <RouterProvider router={router} />
    </QueryClientProvider>,
  )
}

beforeEach(() => vi.clearAllMocks())

describe('ArtpieceDetailView — hidden files header', () => {
  it('shows the visible count and hidden hint when files are hidden', () => {
    mockArtpiece = artpieceWith([makeFile('a', true), makeFile('b', false), makeFile('c', false)])
    renderView()
    expect(screen.getByText('1')).toBeInTheDocument()
    expect(screen.getByText(/2 hidden/)).toBeInTheDocument()
  })

  it('shows no hidden hint when no files are hidden', () => {
    mockArtpiece = artpieceWith([makeFile('a', true)])
    renderView()
    expect(screen.queryByText(/hidden/)).not.toBeInTheDocument()
  })
})

describe('ArtpieceDetailView — Edit details entry point', () => {
  it('opens the modal with every attached file from view mode', async () => {
    mockArtpiece = artpieceWith([makeFile('a', true), makeFile('b', false)])
    renderView()
    await userEvent.click(screen.getByRole('button', { name: 'Edit details' }))
    expect(screen.getByRole('dialog', { name: 'File details' })).toHaveTextContent('2 rows')
  })

  it('opens the modal from edit mode', async () => {
    mockArtpiece = artpieceWith([makeFile('a', true), makeFile('b', false)])
    renderView()
    await userEvent.click(screen.getByRole('button', { name: 'Edit' }))
    await userEvent.click(screen.getByRole('button', { name: 'Edit details' }))
    expect(screen.getByRole('dialog', { name: 'File details' })).toHaveTextContent('2 rows')
  })
})

describe('ArtpieceDetailView — pending cover vs hidden files', () => {
  it('does not set a pending cover that was hidden via the details modal', async () => {
    mockArtpiece = artpieceWith([makeFile('a', true), makeFile('b', true)])
    renderView()
    await userEvent.click(screen.getByRole('button', { name: /^edit$/i }))
    await userEvent.click(screen.getByRole('button', { name: 'cover b' }))

    // file b gets hidden through the modal; the next render sees the updated data
    mockArtpiece = artpieceWith([makeFile('a', true), makeFile('b', false)])
    await userEvent.click(screen.getByRole('button', { name: 'Edit details' }))
    await userEvent.click(screen.getByRole('button', { name: 'Save' }))

    expect(mockSetCover).not.toHaveBeenCalled()
  })

  it('still sets a pending cover that stays visible', async () => {
    mockArtpiece = artpieceWith([makeFile('a', true), makeFile('b', true)])
    renderView()
    await userEvent.click(screen.getByRole('button', { name: /^edit$/i }))
    await userEvent.click(screen.getByRole('button', { name: 'cover b' }))
    await userEvent.click(screen.getByRole('button', { name: 'Save' }))

    expect(mockSetCover).toHaveBeenCalledWith({ artpieceId: 'ap1', fileId: 'b' })
  })
})

describe('ArtpieceDetailView — viewer with hidden files', () => {
  it('opens the viewer on the clicked file even when a hidden file precedes it', async () => {
    mockArtpiece = artpieceWith([makeFile('hidden', false), makeFile('a', true), makeFile('b', true)])
    renderView()
    await userEvent.click(screen.getByRole('button', { name: 'open b' }))
    expect(screen.getByTestId('viewer')).toHaveTextContent('b')
  })
})

