import { describe, it, expect } from 'vitest'
import { parseFilterPayload } from './parseFilterPayload'

describe('parseFilterPayload', () => {
  it('parses a valid full payload', () => {
    const result = parseFilterPayload({
      artist_id: 'artist-uuid',
      character_ids: ['char-1', 'char-2'],
      match_mode: 'any',
    })
    expect(result).toEqual({
      artist_id: 'artist-uuid',
      character_ids: ['char-1', 'char-2'],
      match_mode: 'any',
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
})
