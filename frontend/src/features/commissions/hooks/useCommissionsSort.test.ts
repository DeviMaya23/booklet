import { renderHook, act } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import { useCommissionsSort } from './useCommissionsSort'
import { type Commission } from '../api/useCommissions'

function makeCommission(overrides: Partial<Commission>): Commission {
  return {
    id: 'c-default',
    title: null,
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
    ...overrides,
  }
}

const cAlpha = makeCommission({ id: 'c-alpha', title: 'Alpha' })
const cBeta = makeCommission({ id: 'c-beta', title: 'Beta' })
const cZeta = makeCommission({ id: 'c-zeta', title: 'Zeta' })

describe('useCommissionsSort — initial state', () => {
  it('defaults to title ascending', () => {
    const { result } = renderHook(() => useCommissionsSort([cZeta, cAlpha, cBeta]))
    expect(result.current.sortKey).toBe('title')
    expect(result.current.sortDir).toBe('asc')
    expect(result.current.sorted.map(c => c.id)).toEqual(['c-alpha', 'c-beta', 'c-zeta'])
  })
})

describe('useCommissionsSort — toggleSort', () => {
  it('clicking the active column flips direction to descending', () => {
    const { result } = renderHook(() => useCommissionsSort([cAlpha, cBeta]))
    act(() => result.current.toggleSort('title'))
    expect(result.current.sortDir).toBe('desc')
  })

  it('clicking the active column again flips back to ascending', () => {
    const { result } = renderHook(() => useCommissionsSort([cAlpha, cBeta]))
    act(() => result.current.toggleSort('title'))
    act(() => result.current.toggleSort('title'))
    expect(result.current.sortDir).toBe('asc')
  })

  it('clicking a different column sets the new key and resets to ascending', () => {
    const { result } = renderHook(() => useCommissionsSort([cAlpha, cBeta]))
    act(() => result.current.toggleSort('title'))
    act(() => result.current.toggleSort('artist_name'))
    expect(result.current.sortKey).toBe('artist_name')
    expect(result.current.sortDir).toBe('asc')
  })
})

describe('useCommissionsSort — sort output', () => {
  it('sorts by paid (checked last when asc)', () => {
    const unpaid = makeCommission({ id: 'c-unpaid', paid: false })
    const paid = makeCommission({ id: 'c-paid', paid: true })
    const { result } = renderHook(() => useCommissionsSort([paid, unpaid]))
    act(() => result.current.toggleSort('paid'))
    expect(result.current.sorted[0].id).toBe('c-unpaid')
    expect(result.current.sorted[1].id).toBe('c-paid')
  })
})
