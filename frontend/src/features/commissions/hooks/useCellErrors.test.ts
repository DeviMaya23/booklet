import { renderHook, act } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import { useCellErrors } from './useCellErrors'

describe('useCellErrors', () => {
  it('hasCellError returns false initially', () => {
    const { result } = renderHook(() => useCellErrors())
    expect(result.current.hasCellError('c-1', 'paid')).toBe(false)
  })

  it('setCellError(true) makes hasCellError return true', () => {
    const { result } = renderHook(() => useCellErrors())
    act(() => result.current.setCellError('c-1', 'paid', true))
    expect(result.current.hasCellError('c-1', 'paid')).toBe(true)
  })

  it('setCellError(false) clears the error', () => {
    const { result } = renderHook(() => useCellErrors())
    act(() => result.current.setCellError('c-1', 'paid', true))
    act(() => result.current.setCellError('c-1', 'paid', false))
    expect(result.current.hasCellError('c-1', 'paid')).toBe(false)
  })

  it('errors are isolated per id and field', () => {
    const { result } = renderHook(() => useCellErrors())
    act(() => result.current.setCellError('c-1', 'paid', true))
    expect(result.current.hasCellError('c-2', 'paid')).toBe(false)
    expect(result.current.hasCellError('c-1', 'status')).toBe(false)
  })
})
