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
      tag_ids: [],
      tag_match_mode: 'all',
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

describe('filterToHref — tag params', () => {
  it('writes tags as a comma-separated list', () => {
    const href = filterToHref(makeSavedFilter({ tag_ids: ['t-a', 't-b'] }))
    expect(new URL(href, 'http://x').searchParams.get('tags')).toBe('t-a,t-b')
  })

  it('writes tag_match=any for any mode', () => {
    const params = new URL(filterToHref(makeSavedFilter({ tag_ids: ['t-a', 't-b'], tag_match_mode: 'any' })), 'http://x').searchParams
    expect(params.get('tag_match')).toBe('any')
  })

  it('omits tag_match for all mode', () => {
    const params = new URL(filterToHref(makeSavedFilter({ tag_ids: ['t-a', 't-b'], tag_match_mode: 'all' })), 'http://x').searchParams
    expect(params.get('tags')).toBe('t-a,t-b')
    expect(params.has('tag_match')).toBe(false)
  })

  it('adds no tag params when there are no tags', () => {
    const href = filterToHref(makeSavedFilter({ tag_ids: [], tag_match_mode: 'all' }))
    expect(href).toBe('/app/artpieces')
  })

  it('treats a payload without tag fields as having no tags', () => {
    const filter = makeSavedFilter()
    delete (filter.filter_payload as Partial<typeof filter.filter_payload>).tag_ids
    delete (filter.filter_payload as Partial<typeof filter.filter_payload>).tag_match_mode
    expect(filterToHref(filter)).toBe('/app/artpieces')
  })

  it('keeps tag params alongside other params', () => {
    const params = new URL(
      filterToHref(makeSavedFilter({ artist_id: 'a1', verse_id: 'v1', tag_ids: ['t-a'] })),
      'http://x',
    ).searchParams
    expect(params.get('artist')).toBe('a1')
    expect(params.get('verse')).toBe('v1')
    expect(params.get('tags')).toBe('t-a')
  })

  it('still returns the bare route for an unparseable payload', () => {
    const filter = { ...makeSavedFilter(), filter_payload: null } as unknown as SavedFilter
    expect(filterToHref(filter)).toBe('/app/artpieces')
  })
})
