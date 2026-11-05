import { screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mockApiRoutes } from '@/test/mockApi'
import { makeCharacter, makeArtpiece, renderCharacterDetail } from '@/test/characterFixtures'

vi.mock('@kinde-oss/kinde-auth-react', () => ({
  useKindeAuth: vi.fn().mockReturnValue({ getToken: vi.fn().mockResolvedValue('tok') }),
}))
vi.mock('@/lib/api', () => ({ apiFetch: vi.fn() }))
vi.mock('sonner', () => ({ toast: { success: vi.fn(), error: vi.fn() } }))

const me = { id: 'c-1', name: 'Tanelle' }
const someoneElse = { id: 'c-2', name: 'Other' }

const artpieces = [
  makeArtpiece({ id: 'a1', title: 'Apothecary', characters: [me] }),
  makeArtpiece({ id: 'a2', title: 'Market day', characters: [me, someoneElse] }),
  makeArtpiece({ id: 'a3', title: 'Pose sheet', characters: [me], is_reference: true }),
  makeArtpiece({ id: 'a4', title: 'Not mine', characters: [someoneElse] }),
  makeArtpiece({ id: 'a5', title: 'Other reference', characters: [someoneElse], is_reference: true }),
]

function mockPage(list = artpieces) {
  mockApiRoutes({
    'GET /characters/c-1': { body: makeCharacter() },
    'GET /artpieces': { body: list },
  })
}

beforeEach(() => {
  vi.clearAllMocks()
})

describe('CharacterArtpieceTabs', () => {
  it('splits the character\'s artpieces by reference flag with counts', async () => {
    mockPage()
    renderCharacterDetail()

    expect(await screen.findByRole('tab', { name: /artpieces\s*2/i })).toBeInTheDocument()
    expect(screen.getByRole('tab', { name: /references\s*1/i })).toBeInTheDocument()
    expect(await screen.findByRole('button', { name: 'View Apothecary' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'View Market day' })).toBeInTheDocument()
  })

  it('does not list other characters\' artpieces or references under Artpieces', async () => {
    mockPage()
    renderCharacterDetail()

    await screen.findByRole('button', { name: 'View Apothecary' })
    expect(screen.queryByRole('button', { name: 'View Not mine' })).not.toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'View Pose sheet' })).not.toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'View Other reference' })).not.toBeInTheDocument()
  })

  it('shows only reference artpieces under References', async () => {
    mockPage()
    renderCharacterDetail()

    await userEvent.click(await screen.findByRole('tab', { name: /references/i }))

    expect(await screen.findByRole('button', { name: 'View Pose sheet' })).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'View Apothecary' })).not.toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'View Other reference' })).not.toBeInTheDocument()
  })

  it('shows neutral empty states that name the character', async () => {
    mockPage([])
    renderCharacterDetail()

    expect(await screen.findByText('No artpieces with Tanelle yet')).toBeInTheDocument()
    expect(screen.getByText(/add this character in an artpiece's characters field/i)).toBeInTheDocument()
    expect(screen.queryByText(/\b(her|his|she|he)\b/i)).not.toBeInTheDocument()

    await userEvent.click(screen.getByRole('tab', { name: /references/i }))

    expect(await screen.findByText('No references with Tanelle yet')).toBeInTheDocument()
  })

  it('opens the file viewer when a tile is clicked', async () => {
    mockPage()
    renderCharacterDetail()

    await userEvent.click(await screen.findByRole('button', { name: 'View Market day' }))

    expect(await screen.findByRole('dialog', { name: 'File viewer' })).toBeInTheDocument()
    expect(screen.getByText('2 of 2')).toBeInTheDocument()
  })

  it('opens the artpiece detail from a tile\'s info control', async () => {
    mockPage()
    renderCharacterDetail()
    await screen.findByRole('button', { name: 'View Apothecary' })

    await userEvent.click(screen.getAllByRole('button', { name: 'Open details' })[0])

    expect(await screen.findByText('Artpiece detail')).toBeInTheDocument()
  })
})
