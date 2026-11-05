import type { Character } from '../api/useCharacters'

export interface CharacterEditSummaryInput {
  character: Character
  name: string
  notes: string
  verseId: string | null
  avatar: { hasNewFile: boolean; cleared: boolean }
  selectedFolderIds: string[]
}

function plural(n: number, noun: string): string {
  return `${n} ${noun}${n === 1 ? '' : 's'}`
}

/** "Unsaved changes: …" for the edit form, or null when nothing differs from the saved character. */
export function characterEditSummary({
  character,
  name,
  notes,
  verseId,
  avatar,
  selectedFolderIds,
}: CharacterEditSummaryInput): string | null {
  const parts: string[] = []

  if (name.trim() !== character.name) parts.push('name changed')
  if (verseId !== (character.verse_id ?? null)) parts.push('verse changed')
  if (notes.trim() !== (character.notes ?? '').trim()) parts.push('notes changed')
  if (avatar.hasNewFile) parts.push('avatar changed')
  else if (avatar.cleared) parts.push('avatar removed')

  const original = new Set(character.folders.map((f) => f.id))
  const selected = new Set(selectedFolderIds)
  const added = selectedFolderIds.filter((id) => !original.has(id)).length
  const removed = [...original].filter((id) => !selected.has(id)).length
  if (added > 0) parts.push(`${plural(added, 'folder')} added`)
  if (removed > 0) parts.push(`${plural(removed, 'folder')} removed`)

  return parts.length > 0 ? `Unsaved changes: ${parts.join(', ')}` : null
}
