import { act, fireEvent, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { toast } from 'sonner'
import { FILE_SIZE_LIMITS } from '@/features/files/lib/fileCategory'
import { mockApiRoutes, apiCalls, callsTo } from '@/test/mockApi'
import {
  makeCharacter,
  makeVerse,
  foldersResponse,
  renderCharacterDetail,
} from '@/test/characterFixtures'
import type { Character } from '../api/useCharacters'

vi.mock('@kinde-oss/kinde-auth-react', () => ({
  useKindeAuth: vi.fn().mockReturnValue({ getToken: vi.fn().mockResolvedValue('tok') }),
}))
vi.mock('@/lib/api', () => ({ apiFetch: vi.fn() }))
vi.mock('sonner', () => ({ toast: { success: vi.fn(), error: vi.fn() } }))
// react-easy-crop needs real layout; stand in a dialog that confirms with a fixed file.
vi.mock('./AvatarCropDialog', () => ({
  default: ({ onCrop }: { onCrop: (f: File) => void }) => (
    <button onClick={() => onCrop(new File(['jpg'], 'crop.jpg', { type: 'image/jpeg' }))}>
      Confirm crop
    </button>
  ),
}))

const folderA = { id: 'f-a', name: 'Plant', token: 'tok-a' }
const folderB = { id: 'f-b', name: 'Cottagecore', token: 'tok-b' }
const folderC = { id: 'f-c', name: 'Tavern', token: 'tok-c' }

const baseCharacter = makeCharacter({
  avatar_url: 'https://img/tanelle.jpg',
  notes: 'Apothecary from Dyffrynne',
  verse_id: 'v-1',
  folders: [
    { id: 'f-a', name: 'Plant' },
    { id: 'f-b', name: 'Cottagecore' },
  ],
})

let fetchMock: ReturnType<typeof vi.fn>
let current: Character

type Routes = Parameters<typeof mockApiRoutes>[0]

function setupApi(character: Character = baseCharacter, overrides: Routes = {}) {
  current = character
  mockApiRoutes({
    'GET /characters/c-1': () => ({ body: current }),
    'GET /verses': { body: [makeVerse(), makeVerse({ id: 'v-2', name: 'Dragon Academy' })] },
    'GET /folders': { body: foldersResponse([folderA, folderB, folderC]) },
    'PUT /characters/c-1': ({ body }) => {
      const b = body as { name: string; notes: string | null; verse_id: string | null }
      current = { ...current, name: b.name, notes: b.notes, verse_id: b.verse_id }
      return { body: current }
    },
    ...overrides,
  })
}

async function openEditor(character: Character = baseCharacter, overrides: Routes = {}) {
  setupApi(character, overrides)
  const utils = renderCharacterDetail()
  await userEvent.click(await screen.findByRole('button', { name: /edit/i }))
  await screen.findByLabelText('Name')
  return utils
}

async function renameTo(name: string) {
  const input = screen.getByLabelText('Name')
  await userEvent.clear(input)
  if (name) await userEvent.type(input, name)
}

async function pickAvatar() {
  const input = document.querySelector('input[type="file"]') as HTMLInputElement
  await userEvent.upload(input, new File(['png'], 'src.png', { type: 'image/png' }))
  await userEvent.click(await screen.findByText('Confirm crop'))
}

beforeEach(() => {
  vi.clearAllMocks()
  URL.createObjectURL = vi.fn(() => 'blob:mock')
  URL.revokeObjectURL = vi.fn()
  fetchMock = vi.fn().mockResolvedValue({ ok: true } as Response)
  vi.stubGlobal('fetch', fetchMock)
})

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('CharacterEditForm — header and fields', () => {
  it('opens pre-filled from the character', async () => {
    await openEditor()

    expect(screen.getByLabelText('Name')).toHaveValue('Tanelle')
    expect(screen.getByLabelText('Notes')).toHaveValue('Apothecary from Dyffrynne')
    await waitFor(() => expect(screen.getByLabelText('Verse')).toHaveValue('Coffee Shop AU'))
    expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Save' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Tanelle' })).toBeInTheDocument()
    expect(screen.queryByText(/unsaved changes/i)).not.toBeInTheDocument()
  })

  it('disables Save while the name is blank', async () => {
    await openEditor()

    await renameTo('   ')

    expect(screen.getByRole('button', { name: 'Save' })).toBeDisabled()
    expect(callsTo('PUT', '/characters/c-1')).toHaveLength(0)
  })

  it('Cancel discards edits and returns to the view with the original values', async () => {
    await openEditor()
    await renameTo('Renamed')

    await userEvent.click(screen.getByRole('button', { name: 'Cancel' }))

    expect(await screen.findByRole('heading', { name: 'Tanelle' })).toBeInTheDocument()
    expect(screen.queryByText('Discard changes?')).not.toBeInTheDocument()
    expect(callsTo('PUT', '/characters/c-1')).toHaveLength(0)
  })

  it('sends verse_id null after clearing the verse', async () => {
    await openEditor()
    await waitFor(() => expect(screen.getByLabelText('Verse')).toHaveValue('Coffee Shop AU'))

    await userEvent.click(screen.getByRole('button', { name: 'Clear verse' }))
    expect(screen.getByText('Unsaved changes: verse changed')).toBeInTheDocument()
    await userEvent.click(screen.getByRole('button', { name: 'Save' }))

    await waitFor(() => expect(callsTo('PUT', '/characters/c-1')).toHaveLength(1))
    expect(callsTo('PUT', '/characters/c-1')[0].body).toMatchObject({ verse_id: null })
  })

  it('sends notes null after emptying the notes', async () => {
    await openEditor()

    await userEvent.clear(screen.getByLabelText('Notes'))
    await userEvent.click(screen.getByRole('button', { name: 'Save' }))

    await waitFor(() => expect(callsTo('PUT', '/characters/c-1')).toHaveLength(1))
    expect(callsTo('PUT', '/characters/c-1')[0].body).toMatchObject({ notes: null })
  })
})

describe('CharacterEditForm — avatar', () => {
  it('previews a cropped replacement and reports it as unsaved', async () => {
    await openEditor()

    await pickAvatar()

    expect(screen.getByAltText('Avatar')).toHaveAttribute('src', 'blob:mock')
    expect(screen.getByText('Unsaved changes: avatar changed')).toBeInTheDocument()
    expect(callsTo('POST', '/characters/c-1/avatar/init')).toHaveLength(0)
  })

  it('rejects an oversized image with a toast and no crop dialog', async () => {
    await openEditor()
    const big = new File(['x'], 'big.png', { type: 'image/png' })
    Object.defineProperty(big, 'size', { value: FILE_SIZE_LIMITS['image'] + 1 })

    await userEvent.upload(document.querySelector('input[type="file"]') as HTMLInputElement, big)

    expect(toast.error).toHaveBeenCalledWith(expect.stringContaining('MB limit'))
    expect(screen.queryByText('Confirm crop')).not.toBeInTheDocument()
  })

  it('removes the avatar locally without any request before Save', async () => {
    await openEditor()

    await userEvent.click(screen.getByRole('button', { name: 'Remove avatar' }))

    expect(screen.queryByAltText('Avatar')).not.toBeInTheDocument()
    expect(screen.getByText('Unsaved changes: avatar removed')).toBeInTheDocument()
    expect(callsTo('DELETE', '/characters/c-1/avatar')).toHaveLength(0)
  })

  it('opens the crop dialog for an image pasted while no text field is focused', async () => {
    await openEditor()
    const file = new File(['png'], 'paste.png', { type: 'image/png' })

    fireEvent.paste(document.body, {
      clipboardData: { items: [{ type: 'image/png', getAsFile: () => file }] },
    })

    expect(await screen.findByText('Confirm crop')).toBeInTheDocument()
  })

  it('leaves paste alone while a text field is focused', async () => {
    await openEditor()
    const file = new File(['png'], 'paste.png', { type: 'image/png' })

    fireEvent.paste(screen.getByLabelText('Notes'), {
      clipboardData: { items: [{ type: 'image/png', getAsFile: () => file }] },
    })

    expect(screen.queryByText('Confirm crop')).not.toBeInTheDocument()
  })
})

describe('CharacterEditForm — folders', () => {
  it('shows each folder with a Bookleaf link and a remove control', async () => {
    await openEditor()

    expect(await screen.findByRole('link', { name: 'Open Plant in Bookleaf' })).toHaveAttribute(
      'href',
      expect.stringMatching(/\/share\/tok-a$/),
    )
    expect(screen.getByRole('button', { name: 'Remove Plant' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Remove Cottagecore' })).toBeInTheDocument()
  })

  it('removing a chip excludes the folder from the save and the summary counts it', async () => {
    await openEditor()

    await userEvent.click(await screen.findByRole('button', { name: 'Remove Cottagecore' }))
    expect(screen.getByText('Unsaved changes: 1 folder removed')).toBeInTheDocument()
    await userEvent.click(screen.getByRole('button', { name: 'Save' }))

    await waitFor(() => expect(callsTo('PUT', '/characters/c-1')).toHaveLength(1))
    expect(callsTo('PUT', '/characters/c-1')[0].body).toMatchObject({ folder_ids: ['f-a'] })
  })

  it('warns about a folder that is no longer public and excludes it from the save', async () => {
    const stale = {
      ...baseCharacter,
      folders: [...baseCharacter.folders, { id: 'f-old', name: 'Old refs' }],
    }
    await openEditor(stale)

    expect(await screen.findByText('No longer public in Bookleaf.')).toBeInTheDocument()
    expect(screen.getByText('Old refs')).toBeInTheDocument()
    expect(screen.getByText('Unsaved changes: 1 folder removed')).toBeInTheDocument()

    await userEvent.click(screen.getByRole('button', { name: 'Save' }))

    await waitFor(() => expect(callsTo('PUT', '/characters/c-1')).toHaveLength(1))
    expect(callsTo('PUT', '/characters/c-1')[0].body).toMatchObject({ folder_ids: ['f-a', 'f-b'] })
  })

  it('keeps chips visible and removable, without links, when Bookleaf is unreachable', async () => {
    await openEditor(baseCharacter, { 'GET /folders': { ok: false, status: 502, body: {} } })

    expect(await screen.findByText('Plant')).toBeInTheDocument()
    await waitFor(() => expect(callsTo('GET', '/folders').length).toBeGreaterThan(0))
    expect(screen.queryByRole('link', { name: /open plant in bookleaf/i })).not.toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Remove Plant' })).toBeEnabled()
  })

  it('opens and closes the Choose folders drawer without losing edits', async () => {
    await openEditor()
    await renameTo('Tan')

    await userEvent.click(screen.getByRole('button', { name: 'Choose folders' }))
    expect(await screen.findByText('Public folders from Bookleaf')).toBeInTheDocument()

    await userEvent.click(screen.getByRole('button', { name: /close folder panel/i }))

    await waitFor(() => expect(screen.queryByText('Public folders from Bookleaf')).not.toBeInTheDocument())
    expect(screen.getByLabelText('Name')).toHaveValue('Tan')
  })

  it('shows folders picked in the drawer as chips and counts them as added', async () => {
    await openEditor()

    await userEvent.click(screen.getByRole('button', { name: 'Choose folders' }))
    await userEvent.click(await screen.findByRole('checkbox', { name: 'Tavern' }))
    await userEvent.click(screen.getByRole('button', { name: /close folder panel/i }))

    expect(await screen.findByRole('button', { name: 'Remove Tavern' })).toBeInTheDocument()
    expect(screen.getByText('Unsaved changes: 1 folder added')).toBeInTheDocument()
  })

  it('submits original folders plus added minus removed', async () => {
    await openEditor()

    await userEvent.click(await screen.findByRole('button', { name: 'Remove Cottagecore' }))
    await userEvent.click(screen.getByRole('button', { name: 'Choose folders' }))
    await userEvent.click(await screen.findByRole('checkbox', { name: 'Tavern' }))
    await userEvent.click(screen.getByRole('button', { name: /close folder panel/i }))
    await userEvent.click(screen.getByRole('button', { name: 'Save' }))

    await waitFor(() => expect(callsTo('PUT', '/characters/c-1')).toHaveLength(1))
    expect(callsTo('PUT', '/characters/c-1')[0].body).toMatchObject({ folder_ids: ['f-a', 'f-c'] })
  })
})

describe('CharacterEditForm — saving', () => {
  it('saves metadata, toasts, and returns to the view with the new values', async () => {
    await openEditor()
    await renameTo('Tanelle Wren')

    await userEvent.click(screen.getByRole('button', { name: 'Save' }))

    expect(await screen.findByRole('heading', { name: 'Tanelle Wren' })).toBeInTheDocument()
    expect(toast.success).toHaveBeenCalledWith('Character updated')
    expect(callsTo('PUT', '/characters/c-1')[0].body).toEqual({
      name: 'Tanelle Wren',
      notes: 'Apothecary from Dyffrynne',
      is_public: false,
      folder_ids: ['f-a', 'f-b'],
      verse_id: 'v-1',
    })
  })

  it('uploads a new avatar before updating the character', async () => {
    await openEditor(baseCharacter, {
      'POST /characters/c-1/avatar/init': {
        body: { id: 'up-1', upload_url: 'https://r2/put', expires_at: '' },
      },
    })
    await pickAvatar()

    await userEvent.click(screen.getByRole('button', { name: 'Save' }))

    await waitFor(() => expect(callsTo('PUT', '/characters/c-1')).toHaveLength(1))
    const order = apiCalls()
      .map((c) => `${c.method} ${c.path}`)
      .filter((p) => p.startsWith('POST /characters/c-1/avatar') || p === 'PUT /characters/c-1')
    expect(order).toEqual([
      'POST /characters/c-1/avatar/init',
      'POST /characters/c-1/avatar/up-1/complete',
      'PUT /characters/c-1',
    ])
    expect(fetchMock).toHaveBeenCalledWith('https://r2/put', expect.objectContaining({ method: 'PUT' }))
  })

  it('removes the avatar before updating the character', async () => {
    await openEditor()
    await userEvent.click(screen.getByRole('button', { name: 'Remove avatar' }))

    await userEvent.click(screen.getByRole('button', { name: 'Save' }))

    await waitFor(() => expect(callsTo('PUT', '/characters/c-1')).toHaveLength(1))
    const order = apiCalls()
      .map((c) => `${c.method} ${c.path}`)
      .filter((p) => p === 'DELETE /characters/c-1/avatar' || p === 'PUT /characters/c-1')
    expect(order).toEqual(['DELETE /characters/c-1/avatar', 'PUT /characters/c-1'])
  })

  it('stops, toasts and stays in edit mode when the avatar upload fails', async () => {
    await openEditor(baseCharacter, {
      'POST /characters/c-1/avatar/init': {
        body: { id: 'up-1', upload_url: 'https://r2/put', expires_at: '' },
      },
    })
    fetchMock.mockResolvedValue({ ok: false } as Response)
    await pickAvatar()

    await userEvent.click(screen.getByRole('button', { name: 'Save' }))

    await waitFor(() => expect(toast.error).toHaveBeenCalledWith('Failed to upload avatar'))
    expect(callsTo('PUT', '/characters/c-1')).toHaveLength(0)
    expect(screen.getByRole('button', { name: 'Save' })).toBeInTheDocument()
    expect(screen.getByText('Unsaved changes: avatar changed')).toBeInTheDocument()
  })

  it('toasts and keeps the edits when the update fails', async () => {
    await openEditor(baseCharacter, { 'PUT /characters/c-1': { ok: false, status: 500, body: {} } })
    await renameTo('Tanelle Wren')

    await userEvent.click(screen.getByRole('button', { name: 'Save' }))

    await waitFor(() => expect(toast.error).toHaveBeenCalledWith('Failed to update character'))
    expect(screen.getByLabelText('Name')).toHaveValue('Tanelle Wren')
    expect(screen.getByText('Unsaved changes: name changed')).toBeInTheDocument()
  })
})

describe('CharacterEditForm — leaving with unsaved changes', () => {
  it('asks before the back link discards changes, and Keep editing stays', async () => {
    await openEditor()
    await renameTo('Renamed')

    await userEvent.click(screen.getByRole('button', { name: 'Tanelle' }))
    expect(await screen.findByText('Discard changes?')).toBeInTheDocument()
    await userEvent.click(screen.getByRole('button', { name: 'Keep editing' }))

    await waitFor(() => expect(screen.queryByText('Discard changes?')).not.toBeInTheDocument())
    expect(screen.getByLabelText('Name')).toHaveValue('Renamed')
  })

  it('Discard from the back link returns to the view', async () => {
    await openEditor()
    await renameTo('Renamed')

    await userEvent.click(screen.getByRole('button', { name: 'Tanelle' }))
    await userEvent.click(await screen.findByRole('button', { name: 'Discard' }))

    expect(await screen.findByRole('heading', { name: 'Tanelle' })).toBeInTheDocument()
  })

  it('returns to the view immediately when nothing changed', async () => {
    await openEditor()

    await userEvent.click(screen.getByRole('button', { name: 'Tanelle' }))

    expect(await screen.findByRole('heading', { name: 'Tanelle' })).toBeInTheDocument()
    expect(screen.queryByText('Discard changes?')).not.toBeInTheDocument()
  })

  it('guards in-app navigation away from the page, and Discard lets it through', async () => {
    const { router } = await openEditor()
    await renameTo('Renamed')

    await act(async () => {
      void router.navigate('/app/other')
    })
    expect(await screen.findByText('Discard changes?')).toBeInTheDocument()
    expect(screen.queryByText('Other page')).not.toBeInTheDocument()

    await userEvent.click(screen.getByRole('button', { name: 'Discard' }))

    expect(await screen.findByText('Other page')).toBeInTheDocument()
  })

  it('does not guard navigation when nothing changed', async () => {
    const { router } = await openEditor()

    await act(async () => {
      void router.navigate('/app/other')
    })

    expect(await screen.findByText('Other page')).toBeInTheDocument()
  })
})

describe('CharacterEditForm — delete', () => {
  it('deletes after confirmation, toasts and returns to the list without a discard prompt', async () => {
    await openEditor()
    await renameTo('Unsaved rename')

    await userEvent.click(screen.getByRole('button', { name: /delete character/i }))
    expect(await screen.findByText('Delete character?')).toBeInTheDocument()
    await userEvent.click(screen.getByRole('button', { name: 'Delete' }))

    expect(await screen.findByText('Characters list')).toBeInTheDocument()
    expect(callsTo('DELETE', '/characters/c-1')).toHaveLength(1)
    expect(toast.success).toHaveBeenCalledWith('Character deleted')
    expect(screen.queryByText('Discard changes?')).not.toBeInTheDocument()
  })

  it('leaves the character alone when the confirmation is cancelled', async () => {
    await openEditor()

    await userEvent.click(screen.getByRole('button', { name: /delete character/i }))
    await userEvent.click(await screen.findByRole('button', { name: 'Cancel', hidden: false }))

    await waitFor(() => expect(screen.queryByText('Delete character?')).not.toBeInTheDocument())
    expect(callsTo('DELETE', '/characters/c-1')).toHaveLength(0)
    expect(screen.getByLabelText('Name')).toHaveValue('Tanelle')
  })

  it('toasts and stays on the character when the delete fails', async () => {
    await openEditor(baseCharacter, { 'DELETE /characters/c-1': { ok: false, status: 500, body: {} } })

    await userEvent.click(screen.getByRole('button', { name: /delete character/i }))
    await userEvent.click(await screen.findByRole('button', { name: 'Delete' }))

    await waitFor(() => expect(toast.error).toHaveBeenCalledWith('Failed to delete character'))
    expect(screen.queryByText('Characters list')).not.toBeInTheDocument()
    expect(screen.getByLabelText('Name')).toBeInTheDocument()
  })
})
