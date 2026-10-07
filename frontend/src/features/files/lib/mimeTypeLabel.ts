const MIME_LABELS: Record<string, string> = {
  'image/jpeg': 'JPG',
  'image/png': 'PNG',
  'image/gif': 'GIF',
  'image/webp': 'WEBP',
  'image/tiff': 'TIFF',
  'image/svg+xml': 'SVG',
  'image/vnd.adobe.photoshop': 'PSD',
  'application/x-photoshop': 'PSD',
  'application/pdf': 'PDF',
}

export function mimeTypeLabel(mimeType: string): string {
  return MIME_LABELS[mimeType] ?? ''
}

const MIME_DESCRIPTIONS: Record<string, string> = {
  'image/vnd.adobe.photoshop': 'Photoshop document',
  'application/x-photoshop': 'Photoshop document',
  'application/pdf': 'PDF document',
  'image/tiff': 'TIFF image',
  'image/svg+xml': 'SVG image',
}

export function mimeTypeDescription(mimeType: string): string {
  return MIME_DESCRIPTIONS[mimeType] ?? (mimeTypeLabel(mimeType) || mimeType)
}
