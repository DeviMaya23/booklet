import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import ArtistCharacterFilter from './ArtistCharacterFilter'
import { createWrapper } from '@/test/queryClient'
import { characterRequests, mockCharacterPickerApi } from '@/test/characterPickerApi'

vi.mock('@kinde-oss/kinde-auth-react', () => ({
  useKindeAuth: vi.fn().mockReturnValue({ getToken: vi.fn().mockResolvedValue('tok') }),
}))

vi.mock('@/lib/api', () => ({
  apiFetch: vi.fn(),
}))

beforeEach(() => {
  vi.clearAllMocks()
  mockCharacterPickerApi()
})

async function openAndType(selected: { id: string; name: string }[] = []) {
  render(
    <ArtistCharacterFilter
      artist={null}
      onArtistChange={() => {}}
      characters={selected}
      onCharactersChange={() => {}}
    />,
    { wrapper: createWrapper() },
  )
  await userEvent.click(screen.getByRole('button', { name: /filter|character/i }))
  // TokenInput is the only plain textbox in the panel (the artist field is a combobox).
  const input = await screen.findByRole('textbox')
  await waitFor(() => expect(input).toBeEnabled())
  await userEvent.type(input, 'al')
}

describe('ArtistCharacterFilter — character recency', () => {
  it('lists suggestions in the order the API returns them', async () => {
    await openAndType()

    const options = await screen.findAllByRole('listitem')
    expect(options.map((o) => o.textContent)).toEqual(['Alma', 'Alba', 'Alto'])
  })

  it('excludes already-selected characters', async () => {
    await openAndType([{ id: 'c-b', name: 'Alba' }])

    const options = await screen.findAllByRole('listitem')
    expect(options.map((o) => o.textContent)).toEqual(['Alma', 'Alto'])
  })

  it('requests the characters list with sort=recent', async () => {
    await openAndType()

    expect(characterRequests()).toEqual(['/characters?sort=recent'])
  })
})
