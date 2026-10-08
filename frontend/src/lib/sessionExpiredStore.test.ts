import { describe, it, expect, beforeEach } from 'vitest'
import { renderHook, act } from '@testing-library/react'
import { setSessionExpired, useSessionExpired, subscribe } from './sessionExpiredStore'

describe('sessionExpiredStore', () => {
  beforeEach(() => {
    setSessionExpired(false)
  })

  it('notifies listeners when set to true', () => {
    const { result } = renderHook(() => useSessionExpired())
    expect(result.current).toBe(false)

    act(() => setSessionExpired(true))

    expect(result.current).toBe(true)
  })

  it('is a no-op when already true', () => {
    act(() => setSessionExpired(true))

    let notifyCount = 0
    const unsubscribe = subscribe(() => { notifyCount++ })

    act(() => setSessionExpired(true))

    expect(notifyCount).toBe(0)
    unsubscribe()
  })

  it('resets to false', () => {
    act(() => setSessionExpired(true))
    act(() => setSessionExpired(false))

    const { result } = renderHook(() => useSessionExpired())
    expect(result.current).toBe(false)
  })
})
