import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { vi, describe, it, expect, beforeEach } from 'vitest'
import FilesPage from './FilesPage'
import { createWrapper } from '@/test/queryClient'
import { mockCharacterPickerApi } from '@/test/characterPickerApi'

vi.mock('@kinde-oss/kinde-auth-react', () => ({
  useKindeAuth: () => ({ getToken: vi.fn().mockResolvedValue('tok') }),
}))
vi.mock('@/lib/api', () => ({ apiFetch: vi.fn() }))
vi.mock('@/features/files/api/useFiles', () => ({
  useFiles: () => ({ data: [], isLoading: false }),
}))
vi.mock('@/features/files/components/FileInboxGrid', () => ({
  default: ({ onNewArtpiece }: { onNewArtpiece: (ids: string[]) => void }) => (
    <button onClick={() => onNewArtpiece(['f1'])}>New artpiece from file</button>
  ),
}))
vi.mock('@/features/files/components/AddToExistingArtpieceModal', () => ({ default: () => null }))
vi.mock('@/features/files/components/AddFilesModal', () => ({ default: () => null }))
vi.mock('@/features/files/components/ArtpieceFilesInput', () => ({ default: () => null }))
vi.mock('@/features/artists/components/ArtistCombobox', () => ({ default: () => null }))
vi.mock('@/features/artists/components/ArtistFormModal', () => ({ default: () => null }))
vi.mock('sonner', () => ({ toast: { success: vi.fn(), error: vi.fn() } }))

beforeEach(() => {
  vi.clearAllMocks()
  mockCharacterPickerApi()
})

describe('FilesPage — Inbox artpiece modal', () => {
  it('lists characters on focus when the modal is opened from the Inbox', async () => {
    render(<FilesPage />, { wrapper: createWrapper() })

    await userEvent.click(screen.getByRole('button', { name: /new artpiece from file/i }))
    await userEvent.click(await screen.findByPlaceholderText(/search characters/i))

    const options = await screen.findAllByRole('listitem')
    expect(options.map((o) => o.textContent)).toEqual(['Alma', 'Alba', 'Alto'])
  })
})
