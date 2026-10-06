const unsafeChars = /[/\\:*?"<>|]+/g

export function sanitizeFilename(s: string): string {
  return s.replace(unsafeChars, '').trim() || 'artpiece'
}
