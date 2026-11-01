import { describe, it, expect } from 'vitest'
import { filterToHref } from './filterToHref'
import { type SavedFilter } from '../api/useCreateSavedFilter'

function makeSavedFilter(overrides: Partial<SavedFilter['filter_payload']> = {}): SavedFilter {
  return {
    id: 'sf-1',
    name: 'My Filter',
    thumbnail_url: null,
    filter_payload: {
      artist_id: null,
      character_ids: [],
      match_mode: 'all',
      is_reference: false,
      verse_id: null,
      ...overrides,
    },
    created_at: '',
    updated_at: '',
  }
}

describe('filterToHref — verse param', () => {
  it('includes a verse param when verse_id is set', () => {
    const href = filterToHref(makeSavedFilter({ verse_id: 'verse-1' }))
    expect(href).toContain('verse=verse-1')
  })

  it('omits the verse param when verse_id is null', () => {
    const href = filterToHref(makeSavedFilter({ verse_id: null }))
    expect(href).not.toContain('verse=')
  })

  it('omits the verse param when verse_id is absent from the payload', () => {
    const filter = makeSavedFilter()
    const legacyPayload: Record<string, unknown> = { ...filter.filter_payload }
    delete legacyPayload.verse_id
    const legacyFilter = { ...filter, filter_payload: legacyPayload } as unknown as SavedFilter
    const href = filterToHref(legacyFilter)
    expect(href).not.toContain('verse=')
  })

  it('combines verse with other active params', () => {
    const href = filterToHref(makeSavedFilter({ artist_id: 'artist-1', verse_id: 'verse-1' }))
    expect(href).toContain('artist=artist-1')
    expect(href).toContain('verse=verse-1')
  })
})
