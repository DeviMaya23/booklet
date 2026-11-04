export interface FilterPayload {
  artist_id: string | null
  character_ids: string[]
  match_mode: 'all' | 'any'
  is_reference: boolean
  verse_id: string | null
  tag_ids: string[]
  tag_match_mode: 'all' | 'any'
}

export function parseFilterPayload(raw: unknown): FilterPayload | null {
  if (!raw || typeof raw !== 'object' || Array.isArray(raw)) return null
  const p = raw as Record<string, unknown>

  const artist_id = typeof p.artist_id === 'string' ? p.artist_id : null
  const character_ids = Array.isArray(p.character_ids)
    ? p.character_ids.filter((id): id is string => typeof id === 'string')
    : []
  const match_mode: 'all' | 'any' = p.match_mode === 'any' ? 'any' : 'all'
  const is_reference = p.is_reference === true
  const verse_id = typeof p.verse_id === 'string' ? p.verse_id : null
  const tag_ids = Array.isArray(p.tag_ids)
    ? p.tag_ids.filter((id): id is string => typeof id === 'string')
    : []
  const tag_match_mode: 'all' | 'any' = p.tag_match_mode === 'any' ? 'any' : 'all'

  return { artist_id, character_ids, match_mode, is_reference, verse_id, tag_ids, tag_match_mode }
}
