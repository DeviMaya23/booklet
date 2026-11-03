import { describe, it, expect, beforeEach } from 'vitest'
import { act, renderHook } from '@testing-library/react'
import { ThemeProvider, useTheme } from './useTheme'

describe('useTheme', () => {
  beforeEach(() => {
    localStorage.clear()
    document.documentElement.removeAttribute('data-theme')
  })

  it('defaults to the warm theme when no preference is stored', () => {
    const { result } = renderHook(() => useTheme(), { wrapper: ThemeProvider })

    expect(result.current.theme).toBe('warm')
    expect(document.documentElement.dataset.theme).toBe('warm')
  })

  it.each(['warm', 'lumen', 'sunless'] as const)('restores a stored %s preference on load', (stored) => {
    localStorage.setItem('booklet-theme', stored)

    const { result } = renderHook(() => useTheme(), { wrapper: ThemeProvider })

    expect(result.current.theme).toBe(stored)
    expect(document.documentElement.dataset.theme).toBe(stored)
  })

  it('falls back to warm when the stored value is not a supported theme', () => {
    localStorage.setItem('booklet-theme', 'neon')

    const { result } = renderHook(() => useTheme(), { wrapper: ThemeProvider })

    expect(result.current.theme).toBe('warm')
    expect(document.documentElement.dataset.theme).toBe('warm')
  })

  it('setTheme updates the data-theme attribute and localStorage', () => {
    const { result } = renderHook(() => useTheme(), { wrapper: ThemeProvider })

    act(() => result.current.setTheme('sunless'))

    expect(result.current.theme).toBe('sunless')
    expect(document.documentElement.dataset.theme).toBe('sunless')
    expect(localStorage.getItem('booklet-theme')).toBe('sunless')
  })

  it('throws when used outside of ThemeProvider', () => {
    expect(() => renderHook(() => useTheme())).toThrow('useTheme must be used within a ThemeProvider')
  })
})
