import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { vi, describe, it, expect, beforeEach } from 'vitest'
import CommissionFormModal from './CommissionFormModal'
import { type Commission } from '../api/useCommissions'

const mockCreateMutateAsync = vi.fn().mockResolvedValue({})
const mockUpdateMutateAsync = vi.fn().mockResolvedValue({})
const mockReplaceMutateAsync = vi.fn().mockResolvedValue({})

vi.mock('../api/useCreateCommission', () => ({
  useCreateCommission: () => ({ mutateAsync: mockCreateMutateAsync, isPending: false }),
}))
vi.mock('../api/useUpdateCommission', () => ({
  useUpdateCommission: () => ({ mutateAsync: mockUpdateMutateAsync, isPending: false }),
}))
vi.mock('../api/useReplaceArtpieces', () => ({
  useReplaceArtpieces: () => ({ mutateAsync: mockReplaceMutateAsync, isPending: false }),
}))
vi.mock('../api/useCommission', () => {
  const detail = { id: 'c1', artpieces: [] as { id: string; thumbnail_url: string | null }[] }
  return {
    useCommission: () => ({ data: detail, isLoading: false }),
  }
})
vi.mock('@/features/artpieces/api/useArtpieces', () => ({
  useArtpieces: () => ({ data: [], isPending: false, isError: false }),
}))
vi.mock('@/features/characters/api/useCharacters', () => ({
  useCharacters: () => ({ data: [], isPending: false, isError: false }),
}))
vi.mock('@/features/artists/api/useArtists', () => ({
  useArtists: () => ({ data: [], isError: false }),
}))
vi.mock('@/features/artists/components/ArtistFormModal', () => ({
  default: () => null,
}))

const baseCommission: Commission = {
  id: 'c1',
  title: 'Test Commission',
  artist_id: null,
  artist_name: null,
  artist_link: null,
  status: 'waitlist',
  price: null,
  paid: false,
  paid_date: null,
  finish_date: null,
  last_contacted_at: null,
  notes: null,
  created_at: '2024-01-01T00:00:00Z',
  updated_at: '2024-01-01T00:00:00Z',
}

beforeEach(() => {
  vi.clearAllMocks()
})

describe('CommissionFormModal — create mode', () => {
  it('calls createCommission on save and not updateCommission', async () => {
    render(
      <CommissionFormModal
        open
        onOpenChange={() => {}}
        mode="create"
      />,
    )

    await userEvent.click(screen.getByRole('button', { name: /save/i }))

    expect(mockCreateMutateAsync).toHaveBeenCalledOnce()
    expect(mockUpdateMutateAsync).not.toHaveBeenCalled()
    expect(mockReplaceMutateAsync).not.toHaveBeenCalled()
  })
})

describe('CommissionFormModal — edit mode', () => {
  it('calls updateCommission on save and not createCommission', async () => {
    render(
      <CommissionFormModal
        open
        onOpenChange={() => {}}
        mode="edit"
        commission={baseCommission}
      />,
    )

    await userEvent.click(screen.getByRole('button', { name: /save/i }))

    expect(mockUpdateMutateAsync).toHaveBeenCalledOnce()
    expect(mockCreateMutateAsync).not.toHaveBeenCalled()
  })

  it('does not call replaceArtpieces when artpiece set is unchanged', async () => {
    render(
      <CommissionFormModal
        open
        onOpenChange={() => {}}
        mode="edit"
        commission={baseCommission}
      />,
    )

    await userEvent.click(screen.getByRole('button', { name: /save/i }))

    expect(mockReplaceMutateAsync).not.toHaveBeenCalled()
  })
})
