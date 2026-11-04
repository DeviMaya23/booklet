import { render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { vi, describe, it, expect, beforeEach } from 'vitest'
import ArtpieceFormModal from './ArtpieceFormModal'
import { createWrapper } from '@/test/queryClient'
import { characterRequests, mockCharacterPickerApi } from '@/test/characterPickerApi'

vi.mock('@kinde-oss/kinde-auth-react', () => ({
  useKindeAuth: () => ({ getToken: vi.fn().mockResolvedValue('tok') }),
}))

vi.mock('@/lib/api', () => ({
  apiFetch: vi.fn(),
}))

vi.mock('@/features/artists/components/ArtistCombobox', () => ({ default: () => null }))
vi.mock('@/features/artists/components/ArtistFormModal', () => ({ default: () => null }))
vi.mock('./ArtpieceFilesInput', () => ({ default: () => null }))
vi.mock('sonner', () => ({ toast: { success: vi.fn(), error: vi.fn() } }))

beforeEach(() => {
  vi.clearAllMocks()
  mockCharacterPickerApi()
})

async function openAndType(text: string) {
  render(<ArtpieceFormModal open onOpenChange={() => {}} />, { wrapper: createWrapper() })
  await userEvent.type(await screen.findByPlaceholderText(/search characters/i), text)
}

describe('ArtpieceFormModal — character recency', () => {
  it('lists suggestions in the order the API returns them', async () => {
    await openAndType('al')

    const options = await screen.findAllByRole('listitem')
    expect(options.map((o) => o.textContent)).toEqual(['Alma', 'Alba', 'Alto'])
  })

  it('excludes characters already selected', async () => {
    await openAndType('alb')
    await userEvent.click(await screen.findByText('Alba'))
    const charactersField = screen.getByText('Characters').parentElement!
    await userEvent.type(within(charactersField).getByRole('textbox'), 'al')

    const options = await screen.findAllByRole('listitem')
    expect(options.map((o) => o.textContent)).toEqual(['Alma', 'Alto'])
  })

  it('requests the characters list with sort=recent', async () => {
    await openAndType('al')

    expect(characterRequests()).toEqual(['/characters?sort=recent'])
  })
})
