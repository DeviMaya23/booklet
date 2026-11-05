import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { FolderPreview } from './FolderPreview'

vi.mock('@kinde-oss/kinde-auth-react', () => ({
  useKindeAuth: vi.fn().mockReturnValue({ getToken: vi.fn().mockResolvedValue('tok') }),
}))

vi.mock('@/lib/api', () => ({
  apiFetch: vi.fn(),
}))

import { apiFetch } from '@/lib/api'

const mockApiFetch = vi.mocked(apiFetch)

function makeImages(count: number) {
  return Array.from({ length: count }, (_, i) => ({
    image_id: `img-${i}`,
    thumbnail_url: `https://cdn.example.com/thumb-${i}.jpg`,
  }))
}

function renderPreview(shareUrl?: string) {
  const qc = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  return render(
    <QueryClientProvider client={qc}>
      <FolderPreview folderID="f1" shareUrl={shareUrl} />
    </QueryClientProvider>,
  )
}

beforeEach(() => {
  vi.clearAllMocks()
})

describe('FolderPreview', () => {
  it('shows skeleton and loading label while fetching', () => {
    mockApiFetch.mockReturnValue(new Promise(() => {}))
    renderPreview()
    expect(screen.getByText('Loading images…')).toBeInTheDocument()
  })

  it('renders images after successful load', async () => {
    mockApiFetch.mockResolvedValue({
      ok: true,
      json: async () => ({ images: makeImages(5) }),
    } as Response)
    renderPreview()
    await waitFor(() => expect(screen.getAllByAltText('thumbnail')).toHaveLength(5))
    expect(screen.queryByText('Show more')).not.toBeInTheDocument()
  })

  it('shows the first 24 images with a caption, then 24 more on Show more', async () => {
    mockApiFetch.mockResolvedValue({
      ok: true,
      json: async () => ({ images: makeImages(30) }),
    } as Response)
    renderPreview()
    await waitFor(() => expect(screen.getAllByAltText('thumbnail')).toHaveLength(24))
    expect(screen.getByText('Showing the first 24')).toBeInTheDocument()
    await userEvent.click(screen.getByRole('button', { name: /show more/i }))
    expect(screen.getAllByAltText('thumbnail')).toHaveLength(30)
    expect(screen.queryByRole('button', { name: /show more/i })).not.toBeInTheDocument()
    expect(screen.getByText('Showing all 30')).toBeInTheDocument()
  })

  it('links to the folder in Bookleaf in a new tab when a share URL is given', async () => {
    mockApiFetch.mockResolvedValue({
      ok: true,
      json: async () => ({ images: makeImages(3) }),
    } as Response)
    renderPreview('https://bookleaf.test/share/tok-1')
    const link = await screen.findByRole('link', { name: /view all in bookleaf/i })
    expect(link).toHaveAttribute('href', 'https://bookleaf.test/share/tok-1')
    expect(link).toHaveAttribute('target', '_blank')
  })

  it('shows error message and Retry button when fetch fails', async () => {
    mockApiFetch.mockResolvedValue({ ok: false } as Response)
    renderPreview()
    await waitFor(() =>
      expect(screen.getByText(/couldn't load the images for this folder/i)).toBeInTheDocument(),
    )
    expect(screen.getByRole('button', { name: /retry/i })).toBeInTheDocument()
  })

  it('calls the API again when Retry is clicked', async () => {
    mockApiFetch.mockResolvedValue({ ok: false } as Response)
    renderPreview()
    await waitFor(() =>
      expect(screen.getByRole('button', { name: /retry/i })).toBeInTheDocument(),
    )
    await userEvent.click(screen.getByRole('button', { name: /retry/i }))
    expect(mockApiFetch).toHaveBeenCalledTimes(2)
  })
})
