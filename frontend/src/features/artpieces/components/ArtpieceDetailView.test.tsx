import { render, screen } from '@testing-library/react'
import { vi, describe, it, expect } from 'vitest'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import ArtpieceDetailView from './ArtpieceDetailView'
import { type ArtpieceDetail } from '../api/useArtpiece'

vi.mock('@kinde-oss/kinde-auth-react', () => ({
  useKindeAuth: () => ({ getToken: vi.fn().mockResolvedValue('tok') }),
}))

vi.mock('../api/useArtpiece', () => ({
  useArtpiece: (_id: string) => ({ data: mockArtpiece, isLoading: false, isError: false }),
  artpieceQueryKey: (id: string) => ['artpieces', id],
}))

vi.mock('../api/useUpdateArtpiece', () => ({
  useUpdateArtpiece: () => ({ mutateAsync: vi.fn(), isPending: false }),
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
  default: () => null,
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
  characters: [],
  files: [],
  created_at: '2024-01-01T00:00:00Z',
  updated_at: '2024-01-01T00:00:00Z',
}

function wrapper({ children }: { children: React.ReactNode }) {
  const qc = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  return <QueryClientProvider client={qc}>{children}</QueryClientProvider>
}

describe('ArtpieceDetailView — artist link', () => {
  it('shows a link when artist_link is set', () => {
    mockArtpiece = { ...baseArtpiece, artist_link: 'https://artist.example.com' }
    render(<ArtpieceDetailView artpieceId="ap1" onClose={() => {}} onDeleted={() => {}} />, { wrapper })
    expect(screen.getByRole('link', { name: /open artist link/i })).toHaveAttribute('href', 'https://artist.example.com')
  })

  it('shows no link when artist_link is null', () => {
    mockArtpiece = { ...baseArtpiece, artist_link: null }
    render(<ArtpieceDetailView artpieceId="ap1" onClose={() => {}} onDeleted={() => {}} />, { wrapper })
    expect(screen.queryByRole('link')).toBeNull()
  })

  it('shows the artist name', () => {
    mockArtpiece = { ...baseArtpiece, artist_name: 'Alice' }
    render(<ArtpieceDetailView artpieceId="ap1" onClose={() => {}} onDeleted={() => {}} />, { wrapper })
    expect(screen.getByText('Alice')).toBeInTheDocument()
  })
})
