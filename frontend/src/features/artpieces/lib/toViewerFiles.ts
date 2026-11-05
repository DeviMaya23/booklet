import type { ViewerFile } from '@/features/files/components/FileViewer'
import type { ArtpieceSummary } from '../api/useArtpieces'

export function toViewerFiles(artpieces: ArtpieceSummary[]): ViewerFile[] {
  return artpieces.map((a) => ({
    id: a.cover_file_id ?? a.id,
    name: a.title ?? a.cover_file_name ?? 'Untitled',
    mimeType: a.cover_file_mime_type ?? 'application/octet-stream',
    thumbnailUrl: a.thumbnail_url,
    previewUrl: a.cover_file_url ?? '',
  }))
}
