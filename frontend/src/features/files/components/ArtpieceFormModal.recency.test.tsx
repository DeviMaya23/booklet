import { render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { vi, describe, it, expect, beforeEach } from 'vitest'
import ArtpieceFormModal from './ArtpieceFormModal'
import { createWrapper } from '@/test/queryClient'
import { characterRequests, manyCharacters, mockCharacterPickerApi } from '@/test/characterPickerApi'

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

describe('ArtpieceFormModal — characters on focus', () => {
  async function openAndFocus() {
    render(<ArtpieceFormModal open onOpenChange={() => {}} />, { wrapper: createWrapper() })
    await userEvent.click(await screen.findByPlaceholderText(/search characters/i))
  }

  it('lists characters in recency order on focus', async () => {
    await openAndFocus()

    const options = await screen.findAllByRole('listitem')
    expect(options.map((o) => o.textContent)).toEqual(['Alma', 'Alba', 'Alto'])
  })

  it('caps the focus list at 10 with a hint row', async () => {
    const chars = manyCharacters(25)
    mockCharacterPickerApi(chars)
    await openAndFocus()

    const options = await screen.findAllByRole('listitem')
    expect(options.slice(0, 10).map((o) => o.textContent)).toEqual(chars.slice(0, 10).map((c) => c.name))
    expect(options).toHaveLength(11)
    expect(options[10]).toHaveTextContent(/type to search/i)
  })

  it('typing finds a character beyond the first 10', async () => {
    mockCharacterPickerApi(manyCharacters(25))
    await openAndType('Char22')

    expect(await screen.findByText('Char22')).toBeInTheDocument()
  })
})
