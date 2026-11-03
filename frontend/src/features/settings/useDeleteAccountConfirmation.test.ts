import { act, renderHook } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import { useDeleteAccountConfirmation } from './useDeleteAccountConfirmation'

describe('useDeleteAccountConfirmation', () => {
  it('does not match while the typed text differs from the email', () => {
    const { result } = renderHook(() => useDeleteAccountConfirmation('me@example.com'))

    act(() => result.current.updateTypedEmail('someone@example.com'))

    expect(result.current.emailMatches).toBe(false)
  })

  it('matches ignoring case and surrounding whitespace', () => {
    const { result } = renderHook(() => useDeleteAccountConfirmation('Me@Example.com'))

    act(() => result.current.updateTypedEmail('  me@example.COM '))

    expect(result.current.emailMatches).toBe(true)
  })

  it('never matches when the account email is unknown', () => {
    const { result } = renderHook(() => useDeleteAccountConfirmation(undefined))

    act(() => result.current.updateTypedEmail(''))

    expect(result.current.emailMatches).toBe(false)
  })

  it('cancel clears the typed text and exits confirming', () => {
    const { result } = renderHook(() => useDeleteAccountConfirmation('me@example.com'))

    act(() => result.current.start())
    act(() => result.current.updateTypedEmail('me@example.com'))
    expect(result.current.confirming).toBe(true)

    act(() => result.current.cancel())

    expect(result.current.confirming).toBe(false)
    expect(result.current.typedEmail).toBe('')
    expect(result.current.emailMatches).toBe(false)
  })
})
