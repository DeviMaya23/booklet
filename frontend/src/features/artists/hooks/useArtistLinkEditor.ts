import { useState, useRef } from 'react'
import { type ArtistLink } from '../api/useArtists'
import { normalizeUrl, isValidUrl } from '../lib/linkUtils'

export interface LinkDraft {
  _key: number
  url: string
  is_primary: boolean
}

export interface ArtistLinkEditorResult {
  links: LinkDraft[]
  linkErrors: (string | null)[]
  addLink: () => void
  removeLink: (index: number) => void
  setPrimary: (index: number) => void
  updateUrl: (index: number, url: string) => void
  handleBlur: (index: number) => void
  buildLinks: () => { url: string; is_primary: boolean }[]
  hasErrors: () => boolean
  reset: () => void
}

export function useArtistLinkEditor(
  initialLinks: Pick<ArtistLink, 'url' | 'is_primary'>[],
): ArtistLinkEditorResult {
  const nextKey = useRef(initialLinks.length)
  const initialLinksRef = useRef(initialLinks)

  const [links, setLinks] = useState<LinkDraft[]>(() =>
    initialLinks.map((l, i) => ({ ...l, _key: i })),
  )
  const [linkErrors, setLinkErrors] = useState<(string | null)[]>(
    () => initialLinks.map(() => null),
  )

  function setLinkError(index: number, error: string | null) {
    setLinkErrors((prev) => {
      const next = [...prev]
      next[index] = error
      return next
    })
  }

  function addLink() {
    setLinks((prev) => [
      ...prev,
      { _key: nextKey.current++, url: '', is_primary: prev.length === 0 },
    ])
    setLinkErrors((prev) => [...prev, null])
  }

  function removeLink(index: number) {
    setLinks((prev) => {
      const wasPrimary = prev[index].is_primary
      const next = prev.filter((_, i) => i !== index)
      if (wasPrimary && next.length > 0) {
        next[0] = { ...next[0], is_primary: true }
      }
      return next
    })
    setLinkErrors((prev) => prev.filter((_, i) => i !== index))
  }

  function setPrimary(index: number) {
    setLinks((prev) => prev.map((l, i) => ({ ...l, is_primary: i === index })))
  }

  function updateUrl(index: number, url: string) {
    setLinks((prev) => prev.map((l, i) => (i === index ? { ...l, url } : l)))
    if (!url) setLinkError(index, null)
  }

  function handleBlur(index: number) {
    const url = links[index].url
    const normalized = normalizeUrl(url)
    if (normalized !== url) {
      setLinks((prev) => {
        const next = [...prev]
        next[index] = { ...next[index], url: normalized }
        return next
      })
    }
    if (normalized && !isValidUrl(normalized)) {
      setLinkError(index, 'Enter a valid URL')
    } else {
      setLinkError(index, null)
    }
  }

  function buildLinks() {
    return links
      .filter((l) => l.url.trim())
      .map(({ url, is_primary }) => ({ url, is_primary }))
  }

  function hasErrors() {
    return buildLinks().some((l) => !isValidUrl(l.url))
  }

  function reset() {
    const initial = initialLinksRef.current
    setLinks(initial.map((l, i) => ({ ...l, _key: i })))
    setLinkErrors(initial.map(() => null))
    nextKey.current = initial.length
  }

  return { links, linkErrors, addLink, removeLink, setPrimary, updateUrl, handleBlur, buildLinks, hasErrors, reset }
}
