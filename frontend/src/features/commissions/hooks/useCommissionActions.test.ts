import { renderHook, act } from '@testing-library/react'
import { describe, it, expect, vi } from 'vitest'
import { useCommissionActions } from './useCommissionActions'

function makeDeps(mutateResolved = true) {
  const mutateAsync = mutateResolved
    ? vi.fn().mockResolvedValue({})
    : vi.fn().mockRejectedValue(new Error('network error'))
  const setCellError = vi.fn()
  return { mutateAsync, setCellError }
}

describe('useCommissionActions — patchPaid', () => {
  it('paid=true includes lastContactedAt in the mutation call', async () => {
    const deps = makeDeps()
    const { result } = renderHook(() => useCommissionActions(deps))
    await act(() => result.current.patchPaid('c-1', true))
    expect(deps.mutateAsync).toHaveBeenCalledWith(expect.objectContaining({
      id: 'c-1',
      paid: true,
      lastContactedAt: expect.any(String),
    }))
  })

  it('paid=false does not include lastContactedAt', async () => {
    const deps = makeDeps()
    const { result } = renderHook(() => useCommissionActions(deps))
    await act(() => result.current.patchPaid('c-1', false))
    const call = deps.mutateAsync.mock.calls[0][0]
    expect(call.paid).toBe(false)
    expect(call.lastContactedAt).toBeUndefined()
  })

  it('success clears the cell error', async () => {
    const deps = makeDeps()
    const { result } = renderHook(() => useCommissionActions(deps))
    await act(() => result.current.patchPaid('c-1', true))
    expect(deps.setCellError).toHaveBeenCalledWith('c-1', 'paid', false)
  })

  it('failure sets the cell error', async () => {
    const deps = makeDeps(false)
    const { result } = renderHook(() => useCommissionActions(deps))
    await act(() => result.current.patchPaid('c-1', true))
    expect(deps.setCellError).toHaveBeenCalledWith('c-1', 'paid', true)
  })
})

describe('useCommissionActions — patchPaidDate', () => {
  it('builds the date string from local date components', async () => {
    const deps = makeDeps()
    const { result } = renderHook(() => useCommissionActions(deps))
    const date = new Date(2025, 5, 15) // June 15, 2025 local time
    await act(() => result.current.patchPaidDate('c-1', date))
    expect(deps.mutateAsync).toHaveBeenCalledWith(expect.objectContaining({
      paidDate: '2025-06-15',
    }))
  })
})
