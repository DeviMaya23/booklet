import { render, screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import ArtistCharacterFilter from './ArtistCharacterFilter'
import { createWrapper } from '@/test/queryClient'
import { characterRequests, mockCharacterPickerApi } from '@/test/characterPickerApi'
import { apiFetch } from '@/lib/api'

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

describe('ArtistCharacterFilter — tag recency', () => {
  const recencyTags = [
    { id: 't-a', name: 'sketch' },
    { id: 't-b', name: 'skirt' },
    { id: 't-c', name: 'skull' },
  ]

  function mockTagsApi() {
    vi.mocked(apiFetch).mockImplementation(async (path: string) => {
      const body = path.startsWith('/tags') ? recencyTags : []
      return { ok: true, json: async () => body } as Response
    })
  }

  async function openTags(selected: { id: string; name: string }[] = []) {
    render(
      <ArtistCharacterFilter
        artist={null}
        onArtistChange={() => {}}
        characters={[]}
        onCharactersChange={() => {}}
        tags={selected}
        onTagsChange={() => {}}
      />,
      { wrapper: createWrapper() },
    )
    await userEvent.click(screen.getByRole('button', { name: /filter|tag/i }))
    // The placeholder disappears once a chip exists, so locate the input by section.
    const section = (await screen.findByText('Tags')).parentElement!.parentElement!
    const input = within(section).getByRole('textbox')
    await waitFor(() => expect(input).toBeEnabled())
    await userEvent.type(input, 'sk')
  }

  it('lists suggestions in the order the API returns them', async () => {
    mockTagsApi()
    await openTags()

    const options = await screen.findAllByRole('listitem')
    expect(options.map((o) => o.textContent)).toEqual(['sketch', 'skirt', 'skull'])
  })

  it('excludes already-selected tags', async () => {
    mockTagsApi()
    await openTags([{ id: 't-b', name: 'skirt' }])

    const options = await screen.findAllByRole('listitem')
    expect(options.map((o) => o.textContent)).toEqual(['sketch', 'skull'])
  })

  it('requests the tags list with sort=recent', async () => {
    mockTagsApi()
    await openTags()

    const tagPaths = vi.mocked(apiFetch).mock.calls.map((c) => c[0] as string).filter((p) => p.startsWith('/tags'))
    expect(tagPaths).toEqual(['/tags?sort=recent'])
  })
})
