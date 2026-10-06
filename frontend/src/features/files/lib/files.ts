const FORBIDDEN_CHARS = /[/\\:*?"<>|]/

export function stripFileExtension(name: string): string {
  const idx = name.lastIndexOf('.')
  if (idx <= 0) return name
  const stripped = name.slice(0, idx)
  return stripped === '' ? name : stripped
}

export function validateFileName(name: string): string | null {
  if (name.trim() === '') return 'Name is required'
  if (FORBIDDEN_CHARS.test(name)) return 'Name contains invalid characters (/ \\ : * ? " < > |)'
  if (name.length > 255) return 'Name must be 255 characters or fewer'
  return null
}
