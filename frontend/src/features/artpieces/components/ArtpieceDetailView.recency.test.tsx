import { render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { vi, describe, it, expect, beforeEach } from 'vitest'
import { createMemoryRouter, RouterProvider } from 'react-router-dom'
import ArtpieceDetailView from './ArtpieceDetailView'
import { type ArtpieceDetail } from '../api/useArtpiece'
import { createWrapper } from '@/test/queryClient'
import { characterRequests, mockCharacterPickerApi } from '@/test/characterPickerApi'

vi.mock('@kinde-oss/kinde-auth-react', () => ({
  useKindeAuth: () => ({ getToken: vi.fn().mockResolvedValue('tok') }),
}))

vi.mock('@/lib/api', () => ({
  apiFetch: vi.fn(),
}))

vi.mock('../api/useArtpiece', () => ({
  useArtpiece: () => ({ data: artpiece, isLoading: false, isError: false }),
  artpieceQueryKey: (id: string) => ['artpieces', id],
}))
vi.mock('@/features/artists/components/ArtistCombobox', () => ({ default: () => null }))
vi.mock('@/features/artists/components/ArtistFormModal', () => ({ default: () => null }))
vi.mock('./ArtpieceDetailFileGrid', () => ({ default: () => null }))
vi.mock('./DeleteArtpieceDialog', () => ({ default: () => null }))
vi.mock('@/features/files/components/FileViewer', () => ({ default: () => null }))
vi.mock('sonner', () => ({ toast: { success: vi.fn(), error: vi.fn() } }))

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

let artpiece: ArtpieceDetail

beforeEach(() => {
  vi.clearAllMocks()
  artpiece = baseArtpiece
  mockCharacterPickerApi()
})

async function enterEditAndType(text: string) {
  const router = createMemoryRouter([
    { path: '/', element: <ArtpieceDetailView artpieceId="ap1" onClose={() => {}} onDeleted={() => {}} /> },
  ])
  render(<RouterProvider router={router} />, { wrapper: createWrapper() })
  await userEvent.click(screen.getByRole('button', { name: /^edit$/i }))
  const charactersField = screen.getByText('Characters').parentElement!
  await userEvent.type(within(charactersField).getByRole('textbox'), text)
}

describe('ArtpieceDetailView — character recency in edit mode', () => {
  it('lists suggestions in the order the API returns them', async () => {
    await enterEditAndType('al')

    const options = await screen.findAllByRole('listitem')
    expect(options.map((o) => o.textContent)).toEqual(['Alma', 'Alba', 'Alto'])
  })

  it('excludes characters already on the artpiece', async () => {
    artpiece = { ...baseArtpiece, characters: [{ id: 'c-b', name: 'Alba' }] as ArtpieceDetail['characters'] }

    await enterEditAndType('al')

    const options = await screen.findAllByRole('listitem')
    expect(options.map((o) => o.textContent)).toEqual(['Alma', 'Alto'])
  })

  it('requests the characters list with sort=recent', async () => {
    await enterEditAndType('al')

    expect(characterRequests()).toEqual(['/characters?sort=recent'])
  })
})
