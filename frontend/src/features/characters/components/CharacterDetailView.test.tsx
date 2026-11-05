import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mockApiRoutes, callsTo } from '@/test/mockApi'
import {
  makeCharacter,
  makeVerse,
  foldersResponse,
  renderCharacterDetail,
} from '@/test/characterFixtures'

vi.mock('@kinde-oss/kinde-auth-react', () => ({
  useKindeAuth: vi.fn().mockReturnValue({ getToken: vi.fn().mockResolvedValue('tok') }),
}))
vi.mock('@/lib/api', () => ({ apiFetch: vi.fn() }))
vi.mock('sonner', () => ({ toast: { success: vi.fn(), error: vi.fn() } }))

const tanelle = makeCharacter({
  avatar_url: 'https://img/tanelle.jpg',
  notes: 'Apothecary from Dyffrynne',
  verse_id: 'v-1',
})

beforeEach(() => {
  vi.clearAllMocks()
})

describe('CharacterDetailView — view mode', () => {
  it('shows the avatar, name, verse chip and notes', async () => {
    mockApiRoutes({
      'GET /characters/c-1': { body: tanelle },
      'GET /verses': { body: [makeVerse()] },
    })
    renderCharacterDetail()

    expect(await screen.findByRole('heading', { name: 'Tanelle' })).toBeInTheDocument()
    expect(screen.getByAltText('Tanelle')).toHaveAttribute('src', 'https://img/tanelle.jpg')
    expect(await screen.findByText('Coffee Shop AU')).toBeInTheDocument()
    expect(screen.getByText('Apothecary from Dyffrynne')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /edit/i })).toBeInTheDocument()
  })

  it('shows no verse chip when the character has no verse', async () => {
    mockApiRoutes({
      'GET /characters/c-1': { body: { ...tanelle, verse_id: null } },
      'GET /verses': { body: [makeVerse()] },
    })
    renderCharacterDetail()

    await screen.findByRole('heading', { name: 'Tanelle' })
    expect(screen.queryByText('Coffee Shop AU')).not.toBeInTheDocument()
  })

  it('shows no notes text when the character has none', async () => {
    mockApiRoutes({ 'GET /characters/c-1': { body: { ...tanelle, notes: null } } })
    renderCharacterDetail()

    await screen.findByRole('heading', { name: 'Tanelle' })
    expect(screen.queryByText('Apothecary from Dyffrynne')).not.toBeInTheDocument()
  })

  it('shows a placeholder when the character has no avatar', async () => {
    mockApiRoutes({ 'GET /characters/c-1': { body: { ...tanelle, avatar_url: null } } })
    renderCharacterDetail()

    expect(await screen.findByLabelText('No avatar')).toBeInTheDocument()
    expect(screen.queryByAltText('Tanelle')).not.toBeInTheDocument()
  })

  it('shows a loading state with the back link while the character loads', () => {
    mockApiRoutes({ 'GET /characters/c-1': { body: tanelle } })
    renderCharacterDetail()

    expect(screen.getByText('Loading…')).toBeInTheDocument()
    expect(screen.getByRole('link', { name: /characters/i })).toBeInTheDocument()
  })

  it('shows an error state with a back link when the character cannot be loaded', async () => {
    mockApiRoutes({ 'GET /characters/c-1': { ok: false, status: 404, body: {} } })
    renderCharacterDetail()

    expect(await screen.findByText('Failed to load character.')).toBeInTheDocument()
    expect(screen.getByRole('link', { name: /characters/i })).toHaveAttribute('href', '/app/characters')
  })

  it('returns to the characters list from the back link', async () => {
    mockApiRoutes({ 'GET /characters/c-1': { body: tanelle } })
    renderCharacterDetail()

    await userEvent.click(await screen.findByRole('link', { name: /characters/i }))

    expect(await screen.findByText('Characters list')).toBeInTheDocument()
  })
})

describe('CharacterDetailView — folders row', () => {
  const withFolders = {
    ...tanelle,
    folders: [
      { id: 'f-a', name: 'Plant' },
      { id: 'f-b', name: 'Cottagecore' },
    ],
  }

  it('links attached folders to Bookleaf when the live list has their tokens', async () => {
    mockApiRoutes({
      'GET /characters/c-1': { body: withFolders },
      'GET /folders': {
        body: foldersResponse([
          { id: 'f-a', name: 'Plant', token: 'tok-a' },
          { id: 'f-b', name: 'Cottagecore', token: 'tok-b' },
        ]),
      },
    })
    renderCharacterDetail()

    const plant = await screen.findByRole('link', { name: 'Plant' })
    expect(plant).toHaveAttribute('href', expect.stringMatching(/\/share\/tok-a$/))
    expect(plant).toHaveAttribute('target', '_blank')
    expect(screen.getByRole('link', { name: 'Cottagecore' })).toBeInTheDocument()
  })

  it('shows plain folder names, with no error, when Bookleaf is unreachable', async () => {
    mockApiRoutes({
      'GET /characters/c-1': { body: withFolders },
      'GET /folders': { ok: false, status: 502, body: {} },
    })
    renderCharacterDetail()

    expect(await screen.findByText('Plant')).toBeInTheDocument()
    expect(screen.getByText('Cottagecore')).toBeInTheDocument()
    await waitFor(() => expect(callsTo('GET', '/folders').length).toBeGreaterThan(0))
    expect(screen.queryByRole('link', { name: 'Plant' })).not.toBeInTheDocument()
    expect(screen.queryByText(/couldn't reach bookleaf/i)).not.toBeInTheDocument()
  })

  it('says "No folders attached" when the character has none', async () => {
    mockApiRoutes({ 'GET /characters/c-1': { body: { ...tanelle, folders: [] } } })
    renderCharacterDetail()

    expect(await screen.findByText('No folders attached')).toBeInTheDocument()
  })
})
