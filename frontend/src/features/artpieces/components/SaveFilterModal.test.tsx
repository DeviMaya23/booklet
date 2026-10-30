import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { vi, describe, it, expect, beforeEach } from 'vitest'
import { toast } from 'sonner'
import SaveFilterModal from './SaveFilterModal'
import { type SavedFilter, type DuplicateNameError } from '../api/useCreateSavedFilter'
import { type ArtpieceSummary } from '../api/useArtpieces'

vi.mock('sonner', () => ({
  toast: { success: vi.fn(), error: vi.fn() },
}))

const mockMutateAsync = vi.fn()
let mockIsPending = false

vi.mock('../api/useCreateSavedFilter', () => ({
  useCreateSavedFilter: () => ({ mutateAsync: mockMutateAsync, isPending: mockIsPending }),
}))

function makeArtpiece(id: string, title: string, thumbnailUrl: string | null = null): ArtpieceSummary {
  return {
    id,
    title,
    artist_id: null,
    artist_name: null,
    artist_link: null,
    cover_file_id: null,
    commission_id: null,
    thumbnail_url: thumbnailUrl,
    cover_file_url: null,
    cover_file_mime_type: null,
    cover_file_name: null,
    notes: null,
    is_reference: false,
    characters: [],
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
  }
}

const baseFilter = {
  artist: null,
  characters: [] as { id: string; name: string }[],
  matchMode: 'all' as const,
  isReference: false,
}

const onOpenChange = vi.fn()
const onSuccess = vi.fn()

beforeEach(() => {
  vi.clearAllMocks()
  mockIsPending = false
})

describe('SaveFilterModal — no-results state', () => {
  it('shows empty state and Save is disabled when no artpieces match', () => {
    render(
      <SaveFilterModal
        open
        onOpenChange={onOpenChange}
        filteredArtpieces={[]}
        filterState={baseFilter}
        onSuccess={onSuccess}
      />,
    )
    expect(screen.getByText('No artpieces match this filter')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Save' })).toBeDisabled()
  })
})

describe('SaveFilterModal — duplicate name', () => {
  it('shows inline name error and keeps modal open when server returns 422 duplicate', async () => {
    const duplicateErr: Error & DuplicateNameError = Object.assign(
      new Error('A saved filter named "My Filter" already exists.'),
      { type: 'duplicate_name' as const },
    )
    mockMutateAsync.mockRejectedValueOnce(duplicateErr)

    render(
      <SaveFilterModal
        open
        onOpenChange={onOpenChange}
        filteredArtpieces={[makeArtpiece('ap-1', 'Sunrise')]}
        filterState={baseFilter}
        onSuccess={onSuccess}
      />,
    )

    await userEvent.type(screen.getByRole('textbox', { name: /name/i }), 'My Filter')
    await userEvent.click(screen.getByRole('button', { name: 'Save' }))

    expect(await screen.findByRole('alert')).toHaveTextContent(/already exists/)
    expect(onOpenChange).not.toHaveBeenCalled()
    expect(onSuccess).not.toHaveBeenCalled()
  })
})

describe('SaveFilterModal — cancel', () => {
  it('closes modal without calling mutation', async () => {
    render(
      <SaveFilterModal
        open
        onOpenChange={onOpenChange}
        filteredArtpieces={[makeArtpiece('ap-1', 'Sunrise')]}
        filterState={baseFilter}
        onSuccess={onSuccess}
      />,
    )

    await userEvent.click(screen.getByRole('button', { name: 'Cancel' }))

    expect(mockMutateAsync).not.toHaveBeenCalled()
    expect(onOpenChange).toHaveBeenCalledWith(false)
  })
})

describe('SaveFilterModal — server error', () => {
  it('shows error toast and keeps modal open on non-422 server error', async () => {
    mockMutateAsync.mockRejectedValueOnce(new Error('Internal server error'))

    render(
      <SaveFilterModal
        open
        onOpenChange={onOpenChange}
        filteredArtpieces={[makeArtpiece('ap-1', 'Sunrise')]}
        filterState={baseFilter}
        onSuccess={onSuccess}
      />,
    )

    await userEvent.type(screen.getByRole('textbox', { name: /name/i }), 'My Filter')
    await userEvent.click(screen.getByRole('button', { name: 'Save' }))

    expect(await screen.findByRole('button', { name: 'Save' })).toBeInTheDocument()
    expect(toast.error).toHaveBeenCalledWith('Failed to save filter')
    expect(onOpenChange).not.toHaveBeenCalled()
    expect(onSuccess).not.toHaveBeenCalled()
  })
})

describe('SaveFilterModal — successful save', () => {
  it('calls onSuccess and shows a success toast on 201', async () => {
    const created: SavedFilter = {
      id: 'sf-1',
      name: 'My Filter',
      thumbnail_url: null,
      filter_payload: { artist_id: null, character_ids: [], match_mode: 'all', is_reference: false },
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
    }
    mockMutateAsync.mockResolvedValueOnce(created)

    render(
      <SaveFilterModal
        open
        onOpenChange={onOpenChange}
        filteredArtpieces={[makeArtpiece('ap-1', 'Sunrise')]}
        filterState={baseFilter}
        onSuccess={onSuccess}
      />,
    )

    await userEvent.type(screen.getByRole('textbox', { name: /name/i }), 'My Filter')
    await userEvent.click(screen.getByRole('button', { name: 'Save' }))

    expect(await screen.findByRole('button', { name: 'Save' })).toBeInTheDocument()
    expect(onSuccess).toHaveBeenCalledWith(created)
    expect(toast.success).toHaveBeenCalledWith('Filter "My Filter" saved')
  })
})
