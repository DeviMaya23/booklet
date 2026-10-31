import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { vi, describe, it, expect, beforeEach } from 'vitest'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { createMemoryRouter, RouterProvider } from 'react-router-dom'
import ArtpieceDetailView from './ArtpieceDetailView'
import { type ArtpieceDetail } from '../api/useArtpiece'

vi.mock('@kinde-oss/kinde-auth-react', () => ({
  useKindeAuth: () => ({ getToken: vi.fn().mockResolvedValue('tok') }),
}))

vi.mock('../api/useArtpiece', () => ({
  useArtpiece: () => ({ data: mockArtpiece, isLoading: false, isError: false }),
  artpieceQueryKey: (id: string) => ['artpieces', id],
}))

const mockMutateAsync = vi.fn().mockResolvedValue({})

vi.mock('../api/useUpdateArtpiece', () => ({
  useUpdateArtpiece: () => ({ mutateAsync: mockMutateAsync, isPending: false }),
}))
vi.mock('../api/useSetCover', () => ({
  useSetCover: () => ({ mutateAsync: vi.fn(), isPending: false }),
}))
vi.mock('../api/useAttachFilesToArtpiece', () => ({
  useAttachFilesToArtpiece: () => ({ mutateAsync: vi.fn(), isPending: false }),
}))
vi.mock('../api/useDownloadArtpiece', () => ({
  useDownloadArtpiece: () => ({ mutateAsync: vi.fn(), isPending: false }),
}))
vi.mock('@/features/files/api/useInitFileUpload', () => ({
  useInitFileUpload: () => ({ mutateAsync: vi.fn(), isPending: false }),
}))
vi.mock('@/features/files/api/useCompleteFileUpload', () => ({
  useCompleteFileUpload: () => ({ mutateAsync: vi.fn(), isPending: false }),
}))
vi.mock('@/features/characters/api/useCharacters', () => ({
  useCharacters: () => ({ data: [], isPending: false }),
}))
vi.mock('@/features/artists/components/ArtistCombobox', () => ({
  default: ({ value }: { value: { name: string } | null }) => (
    <div data-testid="artist-combobox">{value?.name ?? 'no artist'}</div>
  ),
}))
vi.mock('@/features/artists/components/ArtistFormModal', () => ({
  default: ({
    open,
    onOpenChange,
    onCreated,
  }: {
    open: boolean
    onOpenChange: (open: boolean) => void
    onCreated?: (artist: { id: string; name: string }) => void
  }) =>
    open ? (
      <div role="dialog" aria-label="Artist form">
        <button
          type="button"
          onClick={() => {
            onCreated?.({ id: 'a2', name: 'Bob' })
            onOpenChange(false)
          }}
        >
          Create artist
        </button>
        <button type="button" onClick={() => onOpenChange(false)}>
          Close artist form
        </button>
      </div>
    ) : null,
}))
vi.mock('./ArtpieceDetailFileGrid', () => ({
  default: () => null,
}))
vi.mock('./DeleteArtpieceDialog', () => ({
  default: () => null,
}))
vi.mock('@/features/files/components/FileViewer', () => ({
  default: () => null,
}))
vi.mock('sonner', () => ({ toast: { success: vi.fn(), error: vi.fn() } }))

beforeEach(() => {
  vi.clearAllMocks()
  mockMutateAsync.mockResolvedValue({})
})

let mockArtpiece: ArtpieceDetail

const baseArtpiece: ArtpieceDetail = {
  id: 'ap1',
  title: 'Sunrise',
  artist_id: 'a1',
  artist_name: 'Alice',
  artist_link: null,
  cover_file_id: null,
  thumbnail_url: null,
  notes: null,
  is_reference: false,
  characters: [],
  files: [],
  created_at: '2024-01-01T00:00:00Z',
  updated_at: '2024-01-01T00:00:00Z',
}

function renderWithProviders(ui: React.ReactElement) {
  const qc = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  const router = createMemoryRouter([{ path: '/', element: ui }])
  return render(
    <QueryClientProvider client={qc}>
      <RouterProvider router={router} />
    </QueryClientProvider>,
  )
}

describe('ArtpieceDetailView — artist link', () => {
  it('shows a link when artist_link is set', () => {
    mockArtpiece = { ...baseArtpiece, artist_link: 'https://artist.example.com' }
    renderWithProviders(<ArtpieceDetailView artpieceId="ap1" onClose={() => {}} onDeleted={() => {}} />)
    expect(screen.getByRole('link', { name: /open artist link/i })).toHaveAttribute('href', 'https://artist.example.com')
  })

  it('shows no link when artist_link is null', () => {
    mockArtpiece = { ...baseArtpiece, artist_link: null }
    renderWithProviders(<ArtpieceDetailView artpieceId="ap1" onClose={() => {}} onDeleted={() => {}} />)
    expect(screen.queryByRole('link')).toBeNull()
  })

  it('shows the artist name', () => {
    mockArtpiece = { ...baseArtpiece, artist_name: 'Alice' }
    renderWithProviders(<ArtpieceDetailView artpieceId="ap1" onClose={() => {}} onDeleted={() => {}} />)
    expect(screen.getByText('Alice')).toBeInTheDocument()
  })
})

describe('ArtpieceDetailView — is_reference checkbox', () => {
  it('initializes checkbox unchecked when artpiece.is_reference is false', async () => {
    const user = userEvent.setup()
    mockArtpiece = { ...baseArtpiece, is_reference: false }
    renderWithProviders(<ArtpieceDetailView artpieceId="ap1" onClose={() => {}} onDeleted={() => {}} />)
    await user.click(screen.getByRole('button', { name: /edit/i }))
    expect(screen.getByRole('checkbox', { name: /this is a reference/i })).not.toBeChecked()
  })

  it('initializes checkbox checked when artpiece.is_reference is true', async () => {
    const user = userEvent.setup()
    mockArtpiece = { ...baseArtpiece, is_reference: true }
    renderWithProviders(<ArtpieceDetailView artpieceId="ap1" onClose={() => {}} onDeleted={() => {}} />)
    await user.click(screen.getByRole('button', { name: /edit/i }))
    expect(screen.getByRole('checkbox', { name: /this is a reference/i })).toBeChecked()
  })

  it('passes isReference: true to mutateAsync when checkbox is checked on save', async () => {
    const user = userEvent.setup()
    mockArtpiece = { ...baseArtpiece, is_reference: false }
    renderWithProviders(<ArtpieceDetailView artpieceId="ap1" onClose={() => {}} onDeleted={() => {}} />)
    await user.click(screen.getByRole('button', { name: /edit/i }))
    await user.click(screen.getByRole('checkbox', { name: /this is a reference/i }))
    await user.click(screen.getByRole('button', { name: /^save$/i }))
    expect(mockMutateAsync).toHaveBeenCalledWith(expect.objectContaining({ isReference: true }))
  })

  it('passes isReference: false to mutateAsync when checkbox is unchecked on save', async () => {
    const user = userEvent.setup()
    mockArtpiece = { ...baseArtpiece, is_reference: true }
    renderWithProviders(<ArtpieceDetailView artpieceId="ap1" onClose={() => {}} onDeleted={() => {}} />)
    await user.click(screen.getByRole('button', { name: /edit/i }))
    await user.click(screen.getByRole('checkbox', { name: /this is a reference/i }))
    await user.click(screen.getByRole('button', { name: /^save$/i }))
    expect(mockMutateAsync).toHaveBeenCalledWith(expect.objectContaining({ isReference: false }))
  })

  it('toggling checkbox alone makes unsaved changes summary non-null', async () => {
    const user = userEvent.setup()
    mockArtpiece = { ...baseArtpiece, is_reference: false }
    renderWithProviders(<ArtpieceDetailView artpieceId="ap1" onClose={() => {}} onDeleted={() => {}} />)
    await user.click(screen.getByRole('button', { name: /edit/i }))
    await user.click(screen.getByRole('checkbox', { name: /this is a reference/i }))
    expect(screen.getByText(/reference changed/i)).toBeInTheDocument()
  })
})

describe('ArtpieceDetailView — add artist from edit mode', () => {
  async function enterEdit() {
    const user = userEvent.setup()
    mockArtpiece = { ...baseArtpiece }
    renderWithProviders(<ArtpieceDetailView artpieceId="ap1" onClose={() => {}} onDeleted={() => {}} />)
    await user.click(screen.getByRole('button', { name: /edit/i }))
    return user
  }

  it('opens the artist-creation modal and keeps unsaved form values', async () => {
    const user = await enterEdit()
    const title = screen.getByDisplayValue('Sunrise')
    await user.clear(title)
    await user.type(title, 'Sunset')
    await user.click(screen.getByRole('button', { name: /add artist/i }))
    expect(screen.getByRole('dialog', { name: /artist form/i })).toBeInTheDocument()
    await user.click(screen.getByRole('button', { name: /close artist form/i }))
    expect(screen.getByDisplayValue('Sunset')).toBeInTheDocument()
  })

  it('selects the created artist and reports it as an unsaved change', async () => {
    const user = await enterEdit()
    await user.click(screen.getByRole('button', { name: /add artist/i }))
    await user.click(screen.getByRole('button', { name: /create artist/i }))
    expect(screen.getByTestId('artist-combobox')).toHaveTextContent('Bob')
    expect(screen.getByText(/artist changed/i)).toBeInTheDocument()
    expect(mockMutateAsync).not.toHaveBeenCalled()
  })

  it('saves the created artist id', async () => {
    const user = await enterEdit()
    await user.click(screen.getByRole('button', { name: /add artist/i }))
    await user.click(screen.getByRole('button', { name: /create artist/i }))
    await user.click(screen.getByRole('button', { name: /^save$/i }))
    expect(mockMutateAsync).toHaveBeenCalledWith(expect.objectContaining({ artistId: 'a2' }))
  })

  it('keeps the selected artist when the modal is closed without creating', async () => {
    const user = await enterEdit()
    await user.click(screen.getByRole('button', { name: /add artist/i }))
    await user.click(screen.getByRole('button', { name: /close artist form/i }))
    expect(screen.getByTestId('artist-combobox')).toHaveTextContent('no artist')
    expect(screen.queryByText(/unsaved changes/i)).toBeNull()
  })

  it('disables the add button while saving', async () => {
    mockMutateAsync.mockReturnValue(new Promise(() => {}))
    const user = await enterEdit()
    await user.click(screen.getByRole('button', { name: /^save$/i }))
    expect(screen.getByRole('button', { name: /add artist/i })).toBeDisabled()
  })
})
