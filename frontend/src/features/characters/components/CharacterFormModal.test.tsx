import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { toast } from 'sonner'
import CharacterFormModal from './CharacterFormModal'
import { createWrapper } from '@/test/queryClient'
import { mockApiRoutes, apiCalls, callsTo } from '@/test/mockApi'

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

let fetchMock: ReturnType<typeof vi.fn>

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

function renderModal() {
  const onOpenChange = vi.fn()
  render(<CharacterFormModal open onOpenChange={onOpenChange} />, { wrapper: createWrapper() })
  return { onOpenChange }
}

async function pickAvatar() {
  const input = document.querySelector('input[type="file"]') as HTMLInputElement
  await userEvent.upload(input, new File(['png'], 'src.png', { type: 'image/png' }))
  await userEvent.click(await screen.findByText('Confirm crop'))
}

describe('CharacterFormModal (create-only)', () => {
  it('titles the dialog "New character" and offers no Delete button', () => {
    mockApiRoutes({})
    renderModal()

    expect(screen.getByText('New character')).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: /delete/i })).not.toBeInTheDocument()
    expect(screen.queryByRole('button', { name: /remove avatar/i })).not.toBeInTheDocument()
  })

  it('disables Save while the name is blank', async () => {
    mockApiRoutes({})
    renderModal()

    expect(screen.getByRole('button', { name: 'Save' })).toBeDisabled()
    await userEvent.type(screen.getByLabelText(/name/i), 'Tanelle')
    expect(screen.getByRole('button', { name: 'Save' })).toBeEnabled()
    expect(callsTo('POST', '/characters')).toHaveLength(0)
  })

  it('creates a character with a name only and closes on success', async () => {
    mockApiRoutes({ 'POST /characters': { body: { id: 'c-new' } } })
    const { onOpenChange } = renderModal()

    await userEvent.type(screen.getByLabelText(/name/i), 'Tanelle')
    await userEvent.click(screen.getByRole('button', { name: 'Save' }))

    await waitFor(() => expect(onOpenChange).toHaveBeenCalledWith(false))
    expect(callsTo('POST', '/characters')[0].body).toEqual({ name: 'Tanelle', folder_ids: [], verse_id: null })
    expect(toast.success).toHaveBeenCalledWith('Character created')
  })

  it('uploads the avatar after creating the character', async () => {
    mockApiRoutes({
      'POST /characters': { body: { id: 'c-new' } },
      'POST /characters/c-new/avatar/init': {
        body: { id: 'up-1', upload_url: 'https://r2/put', expires_at: '' },
      },
    })
    const { onOpenChange } = renderModal()

    await userEvent.type(screen.getByLabelText(/name/i), 'Tanelle')
    await pickAvatar()
    await userEvent.click(screen.getByRole('button', { name: 'Save' }))

    await waitFor(() => expect(onOpenChange).toHaveBeenCalledWith(false))
    expect(apiCalls().map((c) => `${c.method} ${c.path}`).filter((p) => p.includes('/characters/') || p === 'POST /characters')).toEqual([
      'POST /characters',
      'POST /characters/c-new/avatar/init',
      'POST /characters/c-new/avatar/up-1/complete',
    ])
    expect(fetchMock).toHaveBeenCalledWith('https://r2/put', expect.objectContaining({ method: 'PUT' }))
  })

  it('deletes the new character and stays open when the avatar upload fails', async () => {
    mockApiRoutes({
      'POST /characters': { body: { id: 'c-new' } },
      'POST /characters/c-new/avatar/init': {
        body: { id: 'up-1', upload_url: 'https://r2/put', expires_at: '' },
      },
    })
    fetchMock.mockResolvedValue({ ok: false } as Response)
    const { onOpenChange } = renderModal()

    await userEvent.type(screen.getByLabelText(/name/i), 'Tanelle')
    await pickAvatar()
    await userEvent.click(screen.getByRole('button', { name: 'Save' }))

    await waitFor(() => expect(toast.error).toHaveBeenCalledWith('Failed to upload avatar'))
    expect(callsTo('DELETE', '/characters/c-new')).toHaveLength(1)
    expect(onOpenChange).not.toHaveBeenCalledWith(false)
  })

  it('shows the chosen avatar with a Change overlay', async () => {
    mockApiRoutes({})
    renderModal()

    await pickAvatar()

    expect(await screen.findByAltText('Avatar preview')).toBeInTheDocument()
    expect(screen.getByText('Change')).toBeInTheDocument()
  })
})
