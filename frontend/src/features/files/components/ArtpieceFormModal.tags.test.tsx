import { render, screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { vi, describe, it, expect, beforeEach } from 'vitest'
import ArtpieceFormModal from './ArtpieceFormModal'
import { createTestQueryClient, createWrapper } from '@/test/queryClient'
import { apiFetch } from '@/lib/api'

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

const tagList = [
  { id: 't-a', name: 'outfit' },
  { id: 't-b', name: 'sketch' },
]
const createdTag = { id: 't-new', name: 'wip' }

let createStatus = 201

function mockApi() {
  vi.mocked(apiFetch).mockImplementation(async (path: string, _t: unknown, init?: RequestInit) => {
    if (path === '/tags' && init?.method === 'POST') {
      return { ok: createStatus < 400, status: createStatus, json: async () => createdTag } as Response
    }
    if (path.startsWith('/tags')) return { ok: true, status: 200, json: async () => tagList } as Response
    if (path === '/artpieces' && init?.method === 'POST') {
      return { ok: true, status: 201, json: async () => ({ id: 'ap-1' }) } as Response
    }
    return { ok: true, status: 200, json: async () => [] } as Response
  })
}

function tagsField() {
  return within(screen.getByText('Tags').parentElement!)
}

function artpiecePost() {
  const call = vi.mocked(apiFetch).mock.calls.find((c) => c[0] === '/artpieces' && c[2]?.method === 'POST')
  return call ? JSON.parse(call[2]!.body as string) : null
}

async function renderModal(onOpenChange = vi.fn()) {
  const client = createTestQueryClient()
  const utils = render(<ArtpieceFormModal open onOpenChange={onOpenChange} />, { wrapper: createWrapper(client) })
  await screen.findByText('Tags')
  await waitFor(() => expect(apiFetch).toHaveBeenCalledWith('/tags?sort=recent', expect.any(Function)))
  return { ...utils, client, onOpenChange }
}

async function pickTag(name: string) {
  await userEvent.click(tagsField().getByRole('textbox'))
  await userEvent.click(await screen.findByText(name))
}

beforeEach(() => {
  vi.clearAllMocks()
  createStatus = 201
  mockApi()
})

describe('ArtpieceFormModal — tags', () => {
  it('requests tags with sort=recent when the modal opens', async () => {
    await renderModal()

    expect(vi.mocked(apiFetch).mock.calls.map((c) => c[0])).toContain('/tags?sort=recent')
  })

  it('selecting a suggestion adds a chip and its remove button deselects it', async () => {
    await renderModal()

    await pickTag('outfit')
    expect(tagsField().getByText('outfit')).toBeInTheDocument()

    await userEvent.click(tagsField().getByRole('button', { name: /remove outfit/i }))
    expect(tagsField().queryByText('outfit')).not.toBeInTheDocument()
  })

  it('saves with the selected tag ids', async () => {
    await renderModal()
    await userEvent.type(screen.getByLabelText(/title/i), 'Piece')
    await pickTag('outfit')
    await pickTag('sketch')

    await userEvent.click(screen.getByRole('button', { name: 'Save' }))

    await waitFor(() => expect(artpiecePost()).not.toBeNull())
    expect(artpiecePost().tag_ids).toEqual(['t-a', 't-b'])
  })

  it('saves an empty tag_ids array when no tags are selected', async () => {
    await renderModal()
    await userEvent.type(screen.getByLabelText(/title/i), 'Piece')

    await userEvent.click(screen.getByRole('button', { name: 'Save' }))

    await waitFor(() => expect(artpiecePost()).not.toBeNull())
    expect(artpiecePost().tag_ids).toEqual([])
  })

  it('submits an inline-created tag with the artpiece', async () => {
    await renderModal()
    await userEvent.type(screen.getByLabelText(/title/i), 'Piece')
    await userEvent.type(tagsField().getByRole('textbox'), 'wip{Enter}')
    await waitFor(() => expect(tagsField().getByText('wip')).toBeInTheDocument())

    await userEvent.click(screen.getByRole('button', { name: 'Save' }))

    await waitFor(() => expect(artpiecePost()).not.toBeNull())
    expect(artpiecePost().tag_ids).toEqual(['t-new'])
  })

  it('invalidates the tags query after a successful save', async () => {
    const { client } = await renderModal()
    const invalidate = vi.spyOn(client, 'invalidateQueries')
    await userEvent.type(screen.getByLabelText(/title/i), 'Piece')
    await pickTag('outfit')

    await userEvent.click(screen.getByRole('button', { name: 'Save' }))

    await waitFor(() => expect(invalidate).toHaveBeenCalledWith({ queryKey: ['tags'] }))
  })

  it('keeps the selected tags when the save fails', async () => {
    vi.mocked(apiFetch).mockImplementation(async (path: string, _t: unknown, init?: RequestInit) => {
      if (path === '/artpieces' && init?.method === 'POST') return { ok: false, status: 500 } as Response
      if (path.startsWith('/tags')) return { ok: true, status: 200, json: async () => tagList } as Response
      return { ok: true, status: 200, json: async () => [] } as Response
    })
    await renderModal()
    await userEvent.type(screen.getByLabelText(/title/i), 'Piece')
    await pickTag('outfit')

    await userEvent.click(screen.getByRole('button', { name: 'Save' }))

    await waitFor(() => expect(artpiecePost()).not.toBeNull())
    expect(tagsField().getByText('outfit')).toBeInTheDocument()
  })

  it('treats a tags-only form as dirty when dismissed with Escape', async () => {
    const { onOpenChange } = await renderModal()
    await pickTag('outfit')

    await userEvent.keyboard('{Escape}')

    expect(await screen.findByText('Discard changes?')).toBeInTheDocument()
    expect(onOpenChange).not.toHaveBeenCalledWith(false)
  })

  it('clears the tags when the modal is closed and reopened', async () => {
    const client = createTestQueryClient()
    const Wrapper = createWrapper(client)
    const { rerender } = render(<ArtpieceFormModal open onOpenChange={vi.fn()} />, { wrapper: Wrapper })
    await screen.findByText('Tags')
    await pickTag('outfit')
    expect(tagsField().getByText('outfit')).toBeInTheDocument()

    rerender(<ArtpieceFormModal open={false} onOpenChange={vi.fn()} />)
    rerender(<ArtpieceFormModal open onOpenChange={vi.fn()} />)

    await screen.findByText('Tags')
    expect(tagsField().queryByText('outfit')).not.toBeInTheDocument()
  })
})
