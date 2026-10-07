import { render, screen, fireEvent } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import React from 'react'
import FileViewer, { type ViewerFile } from './FileViewer'

vi.mock('@kinde-oss/kinde-auth-react', () => ({
  useKindeAuth: vi.fn().mockReturnValue({ getToken: vi.fn().mockResolvedValue('tok') }),
}))

vi.mock('@/lib/api', () => ({
  apiFetch: vi.fn().mockResolvedValue({
    ok: true,
    json: async () => ({ download_url: 'https://cdn.example.com/file?sig=x' }),
  }),
}))

vi.mock('sonner', () => ({ toast: { error: vi.fn() } }))

function wrapper({ children }: { children: React.ReactNode }) {
  const qc = new QueryClient({ defaultOptions: { mutations: { retry: false } } })
  return React.createElement(QueryClientProvider, { client: qc }, children)
}

function makeFile(overrides: Partial<ViewerFile> = {}): ViewerFile {
  return {
    id: 'f1',
    name: 'photo.jpg',
    mimeType: 'image/jpeg',
    thumbnailUrl: 'https://cdn.example.com/thumb.jpg',
    previewUrl: 'https://cdn.example.com/file.jpg',
    ...overrides,
  }
}

beforeEach(() => {
  vi.clearAllMocks()
})

describe('FileViewer — keyboard navigation', () => {
  it('ArrowRight advances to the next file', () => {
    const files = [makeFile({ id: 'f1', name: 'first.jpg' }), makeFile({ id: 'f2', name: 'second.jpg' })]
    render(<FileViewer files={files} initialIndex={0} onClose={vi.fn()} />, { wrapper })

    expect(screen.getByText('1 of 2')).toBeInTheDocument()
    fireEvent.keyDown(document, { key: 'ArrowRight' })
    expect(screen.getByText('2 of 2')).toBeInTheDocument()
  })

  it('ArrowLeft goes back to the previous file', () => {
    const files = [makeFile({ id: 'f1', name: 'first.jpg' }), makeFile({ id: 'f2', name: 'second.jpg' })]
    render(<FileViewer files={files} initialIndex={1} onClose={vi.fn()} />, { wrapper })

    expect(screen.getByText('2 of 2')).toBeInTheDocument()
    fireEvent.keyDown(document, { key: 'ArrowLeft' })
    expect(screen.getByText('1 of 2')).toBeInTheDocument()
  })

  it('Escape calls onClose', () => {
    const onClose = vi.fn()
    render(<FileViewer files={[makeFile()]} initialIndex={0} onClose={onClose} />, { wrapper })

    fireEvent.keyDown(document, { key: 'Escape' })
    expect(onClose).toHaveBeenCalledOnce()
  })
})

describe('FileViewer — no-preview fallback', () => {
  it('renders the fallback card for files with no thumbnail', () => {
    const file = makeFile({ mimeType: 'image/vnd.adobe.photoshop', thumbnailUrl: null, name: 'artwork.psd' })
    render(<FileViewer files={[file]} initialIndex={0} onClose={vi.fn()} />, { wrapper })

    expect(screen.getAllByText('artwork.psd')).toHaveLength(2) // top bar + fallback card
    expect(screen.getByText(/Photoshop document · no preview available/i)).toBeInTheDocument()
    expect(screen.queryByRole('img')).not.toBeInTheDocument()
  })
})
