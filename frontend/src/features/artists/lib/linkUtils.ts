export function normalizeUrl(url: string): string {
  if (url && !url.startsWith('http://') && !url.startsWith('https://')) {
    return 'https://' + url
  }
  return url
}

export function isValidUrl(value: string): boolean {
  try {
    new URL(value)
    return true
  } catch {
    return false
  }
}
