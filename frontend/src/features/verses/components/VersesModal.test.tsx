import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import VersesModal from './VersesModal'
import type { Verse } from '../api/useVerses'
import type { Character } from '@/features/characters/api/useCharacters'

// --- mutation mocks -----------------------------------------------------------

type MutateWithCallbacks = (input: unknown, cbs?: { onSuccess?: () => void; onError?: (err: Error) => void }) => void

let createMutate: MutateWithCallbacks
let updateMutate: MutateWithCallbacks
let deleteMutate: MutateWithCallbacks

vi.mock('../api/useCreateVerse', () => ({
  useCreateVerse: () => ({ mutate: (...args: Parameters<MutateWithCallbacks>) => createMutate(...args), isPending: false }),
}))
vi.mock('../api/useUpdateVerse', () => ({
  useUpdateVerse: () => ({ mutate: (...args: Parameters<MutateWithCallbacks>) => updateMutate(...args), isPending: false }),
}))
vi.mock('../api/useDeleteVerse', () => ({
  useDeleteVerse: () => ({ mutate: (...args: Parameters<MutateWithCallbacks>) => deleteMutate(...args), isPending: false }),
}))
vi.mock('../api/useReorderVerses', () => ({
  useReorderVerses: () => ({ mutate: vi.fn(), isPending: false }),
}))
vi.mock('sonner', () => ({ toast: { success: vi.fn(), error: vi.fn() } }))

// --- fixtures -----------------------------------------------------------------

const verse1: Verse = { id: 'v-1', name: 'Solaris', sort_order: 0, last_used_at: null, created_at: '', updated_at: '' }
const verse2: Verse = { id: 'v-2', name: 'Terran', sort_order: 1, last_used_at: null, created_at: '', updated_at: '' }

function makeCharacter(id: string, verseId: string | null): Character {
  return { id, name: `Char ${id}`, verse_id: verseId, avatar_url: null, notes: null, is_public: false, folders: [], created_at: '', updated_at: '' }
}

const noop = () => {}

// --- helpers ------------------------------------------------------------------

beforeEach(() => {
  createMutate = vi.fn()
  updateMutate = vi.fn()
  deleteMutate = vi.fn()
})

function open(verses: Verse[] = [], characters: Character[] = []) {
  render(<VersesModal open onOpenChange={noop} verses={verses} characters={characters} />)
}

// --- tests --------------------------------------------------------------------

describe('VersesModal', () => {
  it('shows empty state when there are no verses', () => {
    open()
    expect(screen.getByText('No verses yet')).toBeInTheDocument()
  })

  it('renders each verse name', () => {
    open([verse1, verse2])
    expect(screen.getByText('Solaris')).toBeInTheDocument()
    expect(screen.getByText('Terran')).toBeInTheDocument()
  })

  it('shows the correct character count next to each verse', () => {
    const characters = [
      makeCharacter('c1', 'v-1'),
      makeCharacter('c2', 'v-1'),
      makeCharacter('c3', 'v-2'),
    ]
    open([verse1, verse2], characters)
    expect(screen.getByText('Solaris')).toBeInTheDocument()
    // "2 characters" next to Solaris row
    const rows = screen.getAllByText(/characters?/)
    const counts = rows.map((el) => el.textContent)
    expect(counts).toContain('2 characters')
    expect(counts).toContain('1 character')
  })

  it('shows the new verse input when "+ New verse" is clicked', async () => {
    open()
    await userEvent.click(screen.getByRole('button', { name: /new verse/i }))
    expect(screen.getByPlaceholderText('Verse name')).toBeInTheDocument()
  })

  it('shows the rename input pre-populated with the verse name', async () => {
    open([verse1])
    await userEvent.click(screen.getByRole('button', { name: /rename solaris/i }))
    expect(screen.getByDisplayValue('Solaris')).toBeInTheDocument()
  })

  it('shows delete confirmation with verse name and character count', async () => {
    const characters = [makeCharacter('c1', 'v-1'), makeCharacter('c2', 'v-1')]
    open([verse1], characters)
    await userEvent.click(screen.getByRole('button', { name: /delete solaris/i }))
    expect(screen.getByText(/"Solaris"/)).toBeInTheDocument()
    expect(screen.getByText(/2 characters will move to Unsorted/)).toBeInTheDocument()
  })

  it('shows a conflict error when create returns 409', async () => {
    createMutate = vi.fn((_input, cbs) => {
      const err = Object.assign(new Error(), { status: 409 })
      cbs?.onError?.(err)
    })
    open()
    await userEvent.click(screen.getByRole('button', { name: /new verse/i }))
    await userEvent.type(screen.getByPlaceholderText('Verse name'), 'Solaris')
    await userEvent.click(screen.getByRole('button', { name: /^save$/i }))
    expect(screen.getByText('A verse with that name already exists.')).toBeInTheDocument()
  })

  it('shows a conflict error when rename returns 409', async () => {
    updateMutate = vi.fn((_input, cbs) => {
      const err = Object.assign(new Error(), { status: 409 })
      cbs?.onError?.(err)
    })
    open([verse1])
    await userEvent.click(screen.getByRole('button', { name: /rename solaris/i }))
    await userEvent.click(screen.getByRole('button', { name: /^save$/i }))
    expect(screen.getByText('A verse with that name already exists.')).toBeInTheDocument()
  })
})
