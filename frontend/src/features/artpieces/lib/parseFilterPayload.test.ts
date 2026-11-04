import { describe, it, expect } from 'vitest'
import { parseFilterPayload } from './parseFilterPayload'

describe('parseFilterPayload', () => {
  it('parses a valid full payload', () => {
    const result = parseFilterPayload({
      artist_id: 'artist-uuid',
      character_ids: ['char-1', 'char-2'],
      match_mode: 'any',
      is_reference: true,
      verse_id: 'verse-uuid',
      tag_ids: ['tag-1', 'tag-2'],
      tag_match_mode: 'any',
    })
    expect(result).toEqual({
      artist_id: 'artist-uuid',
      character_ids: ['char-1', 'char-2'],
      match_mode: 'any',
      is_reference: true,
      verse_id: 'verse-uuid',
      tag_ids: ['tag-1', 'tag-2'],
      tag_match_mode: 'any',
    })
  })

  it('accepts null artist_id', () => {
    const result = parseFilterPayload({
      artist_id: null,
      character_ids: [],
      match_mode: 'all',
    })
    expect(result?.artist_id).toBeNull()
  })

  it('accepts empty character_ids', () => {
    const result = parseFilterPayload({
      artist_id: null,
      character_ids: [],
      match_mode: 'all',
    })
    expect(result?.character_ids).toEqual([])
  })

  it('returns null for non-object input', () => {
    expect(parseFilterPayload(null)).toBeNull()
    expect(parseFilterPayload('string')).toBeNull()
    expect(parseFilterPayload(42)).toBeNull()
    expect(parseFilterPayload([])).toBeNull()
  })

  it('coerces missing artist_id to null', () => {
    const result = parseFilterPayload({ character_ids: [], match_mode: 'all' })
    expect(result?.artist_id).toBeNull()
  })

  it('coerces missing character_ids to empty array', () => {
    const result = parseFilterPayload({ artist_id: null, match_mode: 'all' })
    expect(result?.character_ids).toEqual([])
  })

  it('defaults match_mode to all when absent', () => {
    const result = parseFilterPayload({ artist_id: null, character_ids: [] })
    expect(result?.match_mode).toBe('all')
  })

  it('defaults match_mode to all for unrecognised value', () => {
    const result = parseFilterPayload({ artist_id: null, character_ids: [], match_mode: 'unknown' })
    expect(result?.match_mode).toBe('all')
  })

  it('filters non-string values out of character_ids', () => {
    const result = parseFilterPayload({ artist_id: null, character_ids: ['valid', 123, null], match_mode: 'all' })
    expect(result?.character_ids).toEqual(['valid'])
  })

  it('is_reference true round-trips', () => {
    const result = parseFilterPayload({ artist_id: null, character_ids: [], match_mode: 'all', is_reference: true })
    expect(result?.is_reference).toBe(true)
  })

  it('is_reference false round-trips', () => {
    const result = parseFilterPayload({ artist_id: null, character_ids: [], match_mode: 'all', is_reference: false })
    expect(result?.is_reference).toBe(false)
  })

  it('absent is_reference coerces to false', () => {
    const result = parseFilterPayload({ artist_id: null, character_ids: [], match_mode: 'all' })
    expect(result?.is_reference).toBe(false)
  })

  it('non-boolean is_reference coerces to false', () => {
    const result = parseFilterPayload({ artist_id: null, character_ids: [], match_mode: 'all', is_reference: 'yes' })
    expect(result?.is_reference).toBe(false)
  })

  it('verse_id round-trips when present', () => {
    const result = parseFilterPayload({ artist_id: null, character_ids: [], match_mode: 'all', verse_id: 'verse-1' })
    expect(result?.verse_id).toBe('verse-1')
  })

  it('absent verse_id defaults to null', () => {
    const result = parseFilterPayload({ artist_id: null, character_ids: [], match_mode: 'all' })
    expect(result?.verse_id).toBeNull()
  })

  it('null verse_id round-trips', () => {
    const result = parseFilterPayload({ artist_id: null, character_ids: [], match_mode: 'all', verse_id: null })
    expect(result?.verse_id).toBeNull()
  })

  it('non-string non-null verse_id coerces to null', () => {
    const result = parseFilterPayload({ artist_id: null, character_ids: [], match_mode: 'all', verse_id: 123 })
    expect(result?.verse_id).toBeNull()
  })

  it('absent tag_ids defaults to an empty array', () => {
    const result = parseFilterPayload({ artist_id: null, character_ids: [], match_mode: 'all' })
    expect(result?.tag_ids).toEqual([])
  })

  it('absent tag_match_mode defaults to all', () => {
    const result = parseFilterPayload({ artist_id: null, character_ids: [], match_mode: 'all' })
    expect(result?.tag_match_mode).toBe('all')
  })

  it('tag_match_mode any round-trips', () => {
    const result = parseFilterPayload({ artist_id: null, character_ids: [], match_mode: 'all', tag_ids: ['t'], tag_match_mode: 'any' })
    expect(result?.tag_match_mode).toBe('any')
  })

  it('unrecognised tag_match_mode defaults to all', () => {
    const result = parseFilterPayload({ artist_id: null, character_ids: [], match_mode: 'all', tag_match_mode: 'some' })
    expect(result?.tag_match_mode).toBe('all')
  })

  it('filters non-string values out of tag_ids', () => {
    const result = parseFilterPayload({ artist_id: null, character_ids: [], match_mode: 'all', tag_ids: ['ok', 7, null] })
    expect(result?.tag_ids).toEqual(['ok'])
  })

  it('non-array tag_ids coerces to an empty array', () => {
    const result = parseFilterPayload({ artist_id: null, character_ids: [], match_mode: 'all', tag_ids: 'nope' })
    expect(result?.tag_ids).toEqual([])
  })

  it('tag match mode is independent of character match_mode', () => {
    const result = parseFilterPayload({ artist_id: null, character_ids: [], match_mode: 'any', tag_ids: ['t'] })
    expect(result?.match_mode).toBe('any')
    expect(result?.tag_match_mode).toBe('all')
  })
})
