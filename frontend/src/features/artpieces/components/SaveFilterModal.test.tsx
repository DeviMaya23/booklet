import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { vi, describe, it, expect, beforeEach } from 'vitest'
import { toast } from 'sonner'
import SaveFilterModal from './SaveFilterModal'
import { type SavedFilter, type DuplicateNameError } from '../api/useCreateSavedFilter'
import { type ArtpieceSummary } from '../api/useArtpieces'
import { type Verse } from '@/features/verses/api/useVerses'

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
    tags: [],
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
  }
}

const baseFilter = {
  artist: null,
  characters: [] as { id: string; name: string }[],
  matchMode: 'all' as const,
  isReference: false,
  verse: null as Verse | null,
  tags: [] as { id: string; name: string }[],
  tagMatch: 'all' as 'all' | 'any',
}

const verse: Verse = {
  id: 'verse-1',
  name: 'Elowen',
  sort_order: 0,
  last_used_at: null,
  created_at: '',
  updated_at: '',
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
      filter_payload: { artist_id: null, character_ids: [], match_mode: 'all', is_reference: false, verse_id: null, tag_ids: [], tag_match_mode: 'all' },
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

describe('SaveFilterModal — verse filter', () => {
  it('includes verse_id in the create request payload when a verse is active', async () => {
    const created: SavedFilter = {
      id: 'sf-1',
      name: 'My Filter',
      thumbnail_url: null,
      filter_payload: { artist_id: null, character_ids: [], match_mode: 'all', is_reference: false, verse_id: verse.id, tag_ids: [], tag_match_mode: 'all' },
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
    }
    mockMutateAsync.mockResolvedValueOnce(created)

    render(
      <SaveFilterModal
        open
        onOpenChange={onOpenChange}
        filteredArtpieces={[makeArtpiece('ap-1', 'Sunrise')]}
        filterState={{ ...baseFilter, verse }}
        onSuccess={onSuccess}
      />,
    )

    await userEvent.type(screen.getByRole('textbox', { name: /name/i }), 'My Filter')
    await userEvent.click(screen.getByRole('button', { name: 'Save' }))

    expect(await screen.findByRole('button', { name: 'Save' })).toBeInTheDocument()
    expect(mockMutateAsync).toHaveBeenCalledWith(
      expect.objectContaining({
        filterPayload: expect.objectContaining({ verse_id: verse.id }),
      }),
    )
  })

  it('includes the verse name in the filter summary chip when active', () => {
    render(
      <SaveFilterModal
        open
        onOpenChange={onOpenChange}
        filteredArtpieces={[makeArtpiece('ap-1', 'Sunrise')]}
        filterState={{ ...baseFilter, verse }}
        onSuccess={onSuccess}
      />,
    )
    expect(screen.getByText(/Verse: Elowen/)).toBeInTheDocument()
  })
})

describe('SaveFilterModal — tags', () => {
  const tags = [
    { id: 'tag-1', name: 'outfit' },
    { id: 'tag-2', name: 'sketch' },
  ]

  function renderModal(filterState: typeof baseFilter) {
    render(
      <SaveFilterModal
        open
        onOpenChange={onOpenChange}
        filteredArtpieces={[makeArtpiece('ap-1', 'Sunrise')]}
        filterState={filterState}
        onSuccess={onSuccess}
      />,
    )
  }

  async function save() {
    await userEvent.type(screen.getByRole('textbox', { name: /name/i }), 'My Filter')
    await userEvent.click(screen.getByRole('button', { name: 'Save' }))
  }

  it('includes the tag names and match mode in the summary chip', () => {
    renderModal({ ...baseFilter, tags, tagMatch: 'all' })
    expect(screen.getByText(/Tags: outfit, sketch \(all\)/)).toBeInTheDocument()
  })

  it('shows tags alongside other active filters in the chip', () => {
    renderModal({ ...baseFilter, tags: [tags[0]], tagMatch: 'any', verse })
    expect(screen.getByText(/Tags: outfit \(any\) · Verse: Elowen/)).toBeInTheDocument()
  })

  it('omits tags from the chip when none are selected', () => {
    renderModal({ ...baseFilter, verse })
    expect(screen.queryByText(/Tags:/)).not.toBeInTheDocument()
  })

  it('sends tag_ids and tag_match_mode in the payload', async () => {
    renderModal({ ...baseFilter, tags, tagMatch: 'any' })
    await save()
    expect(mockMutateAsync).toHaveBeenCalledWith(
      expect.objectContaining({
        filterPayload: expect.objectContaining({ tag_ids: ['tag-1', 'tag-2'], tag_match_mode: 'any' }),
      }),
    )
  })

  it('sends an empty tag_ids array when no tags are selected', async () => {
    renderModal({ ...baseFilter, verse })
    await save()
    expect(mockMutateAsync).toHaveBeenCalledWith(
      expect.objectContaining({
        filterPayload: expect.objectContaining({ tag_ids: [], tag_match_mode: 'all' }),
      }),
    )
  })
})
