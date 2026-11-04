import { parseFilterPayload } from './parseFilterPayload'
import { type SavedFilter } from '../api/useCreateSavedFilter'

export function filterToHref(filter: SavedFilter): string {
  const payload = parseFilterPayload(filter.filter_payload)
  if (!payload) return '/app/artpieces'
  const params = new URLSearchParams()
  if (payload.artist_id) params.set('artist', payload.artist_id)
  if (payload.character_ids.length > 0) params.set('characters', payload.character_ids.join(','))
  if (payload.match_mode === 'any') params.set('match', 'any')
  if (payload.is_reference) params.set('ref', 'true')
  if (payload.verse_id) params.set('verse', payload.verse_id)
  if (payload.tag_ids.length > 0) params.set('tags', payload.tag_ids.join(','))
  if (payload.tag_match_mode === 'any') params.set('tag_match', 'any')
  const qs = params.toString()
  return qs ? `/app/artpieces?${qs}` : '/app/artpieces'
}
