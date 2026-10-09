import { describe, it, expect } from 'vitest'
import { renderHook, act } from '@testing-library/react'
import { useArtistLinkEditor } from './useArtistLinkEditor'

describe('useArtistLinkEditor', () => {
  it('first link added is primary', () => {
    const { result } = renderHook(() => useArtistLinkEditor([]))
    act(() => result.current.addLink())
    expect(result.current.links[0].is_primary).toBe(true)
  })

  it('subsequent links added are not primary', () => {
    const { result } = renderHook(() => useArtistLinkEditor([{ url: 'https://a.com', is_primary: true }]))
    act(() => result.current.addLink())
    expect(result.current.links[1].is_primary).toBe(false)
  })

  it('removing the primary link promotes the first remaining', () => {
    const { result } = renderHook(() =>
      useArtistLinkEditor([
        { url: 'https://a.com', is_primary: true },
        { url: 'https://b.com', is_primary: false },
      ]),
    )
    act(() => result.current.removeLink(0))
    expect(result.current.links).toHaveLength(1)
    expect(result.current.links[0].url).toBe('https://b.com')
    expect(result.current.links[0].is_primary).toBe(true)
  })

  it('removing a non-primary link leaves the primary unchanged', () => {
    const { result } = renderHook(() =>
      useArtistLinkEditor([
        { url: 'https://a.com', is_primary: true },
        { url: 'https://b.com', is_primary: false },
      ]),
    )
    act(() => result.current.removeLink(1))
    expect(result.current.links[0].is_primary).toBe(true)
  })

  it('clearing a url clears its validation error', () => {
    const { result } = renderHook(() => useArtistLinkEditor([{ url: 'https://a.com', is_primary: true }]))
    act(() => result.current.updateUrl(0, 'hello world'))
    act(() => result.current.handleBlur(0))
    expect(result.current.linkErrors[0]).toBe('Enter a valid URL')
    act(() => result.current.updateUrl(0, ''))
    expect(result.current.linkErrors[0]).toBeNull()
  })

  it('handleBlur sets an error for an invalid url', () => {
    const { result } = renderHook(() => useArtistLinkEditor([{ url: '', is_primary: true }]))
    act(() => result.current.updateUrl(0, 'hello world'))
    act(() => result.current.handleBlur(0))
    expect(result.current.linkErrors[0]).toBe('Enter a valid URL')
  })

  it('handleBlur prepends https:// when scheme is absent', () => {
    const { result } = renderHook(() => useArtistLinkEditor([{ url: '', is_primary: true }]))
    act(() => result.current.updateUrl(0, 'example.com'))
    act(() => result.current.handleBlur(0))
    expect(result.current.links[0].url).toBe('https://example.com')
    expect(result.current.linkErrors[0]).toBeNull()
  })

  it('buildLinks excludes blank entries', () => {
    const { result } = renderHook(() => useArtistLinkEditor([{ url: 'https://a.com', is_primary: true }]))
    act(() => result.current.addLink())
    const built = result.current.buildLinks()
    expect(built).toHaveLength(1)
    expect(built[0].url).toBe('https://a.com')
  })

  it('each link has a stable key unaffected by deletion', () => {
    const { result } = renderHook(() =>
      useArtistLinkEditor([
        { url: 'https://a.com', is_primary: true },
        { url: 'https://b.com', is_primary: false },
      ]),
    )
    const keyOfB = result.current.links[1]._key
    act(() => result.current.removeLink(0))
    expect(result.current.links[0]._key).toBe(keyOfB)
  })
})
