import type { Tag } from '../api/useTags'

// Finds another tag whose name equals `name`, ignoring case and surrounding
// whitespace. `excludeId` skips the tag being renamed, so a case-only change to
// its own name is not a collision.
export function findTagByName(tags: Tag[], name: string, excludeId?: string): Tag | undefined {
  const wanted = name.trim().toLowerCase()
  return tags.find((t) => t.id !== excludeId && t.name.trim().toLowerCase() === wanted)
}
