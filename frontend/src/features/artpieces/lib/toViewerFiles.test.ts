import { describe, it, expect } from 'vitest'
import { toViewerFiles } from './toViewerFiles'
import type { ArtpieceSummary } from '../api/useArtpieces'

function artpiece(overrides: Partial<ArtpieceSummary>): ArtpieceSummary {
  return {
    id: 'a-1',
    title: null,
    artist_id: null,
    artist_name: null,
    artist_link: null,
    cover_file_id: null,
    commission_id: null,
    thumbnail_url: null,
    cover_file_url: null,
    cover_file_mime_type: null,
    cover_file_name: null,
    notes: null,
    is_reference: false,
    characters: [],
    tags: [],
    created_at: '',
    updated_at: '',
    ...overrides,
  }
}

describe('toViewerFiles', () => {
  it('maps the cover file fields of an artpiece', () => {
    const files = toViewerFiles([
      artpiece({
        id: 'a-1',
        title: 'Sketch',
        cover_file_id: 'f-1',
        cover_file_mime_type: 'image/png',
        thumbnail_url: 'https://t/1',
        cover_file_url: 'https://f/1',
      }),
    ])

    expect(files).toEqual([
      { id: 'f-1', name: 'Sketch', mimeType: 'image/png', thumbnailUrl: 'https://t/1', previewUrl: 'https://f/1' },
    ])
  })

  it('falls back to the artpiece id, cover file name, then "Untitled"', () => {
    const files = toViewerFiles([
      artpiece({ id: 'a-2', title: null, cover_file_name: 'cover.png' }),
      artpiece({ id: 'a-3', title: null, cover_file_name: null }),
    ])

    expect(files[0]).toMatchObject({ id: 'a-2', name: 'cover.png' })
    expect(files[1]).toMatchObject({ id: 'a-3', name: 'Untitled' })
  })

  it('defaults the mime type and preview url when there is no cover file', () => {
    const [file] = toViewerFiles([artpiece({})])

    expect(file.mimeType).toBe('application/octet-stream')
    expect(file.previewUrl).toBe('')
  })
})
