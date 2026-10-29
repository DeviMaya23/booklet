import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi } from 'vitest'
import VerseCombobox from './VerseCombobox'
import type { Verse } from '../api/useVerses'

const verse: Verse = {
  id: 'v-1',
  name: 'Solaris',
  sort_order: 0,
  last_used_at: null,
  created_at: '2026-01-01T00:00:00Z',
  updated_at: '2026-01-01T00:00:00Z',
}

vi.mock('../api/useVerses', () => ({
  useVerses: vi.fn(),
}))

import { useVerses } from '../api/useVerses'

function mockVerses(overrides: Partial<ReturnType<typeof useVerses>> = {}) {
  vi.mocked(useVerses).mockReturnValue({
    data: [verse],
    isError: false,
    isPending: false,
    isSuccess: true,
  } as ReturnType<typeof useVerses>, ...[] as never[])
  Object.assign(vi.mocked(useVerses).mock.results[0] ?? {}, overrides)

  // simple override via re-mock
  vi.mocked(useVerses).mockReturnValue({
    data: [verse],
    isError: false,
    isPending: false,
    isSuccess: true,
    ...overrides,
  } as ReturnType<typeof useVerses>)
}

describe('VerseCombobox', () => {
  it('shows error message when verse data fails to load', () => {
    mockVerses({ data: undefined, isError: true })

    render(<VerseCombobox value={null} onChange={() => {}} />)

    expect(screen.getByText('Failed to load verses.')).toBeInTheDocument()
  })

  it('shows the clear button when a verse is selected', () => {
    mockVerses()

    render(<VerseCombobox value={verse} onChange={() => {}} />)

    expect(screen.getByRole('button', { name: /clear verse/i })).toBeInTheDocument()
  })

  it('calls onChange(null) when the clear button is clicked', async () => {
    mockVerses()
    const onChange = vi.fn()

    render(<VerseCombobox value={verse} onChange={onChange} />)
    await userEvent.click(screen.getByRole('button', { name: /clear verse/i }))

    expect(onChange).toHaveBeenCalledWith(null)
  })

  it('does not show the clear button when no verse is selected', () => {
    mockVerses()

    render(<VerseCombobox value={null} onChange={() => {}} />)

    expect(screen.queryByRole('button', { name: /clear verse/i })).toBeNull()
  })
})
