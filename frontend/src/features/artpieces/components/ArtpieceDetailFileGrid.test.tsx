import { render, screen } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import ArtpieceDetailFileGrid from './ArtpieceDetailFileGrid'
import { type ArtpieceFile } from '../api/useArtpiece'

function makeFile(id: string, overrides: Partial<ArtpieceFile> = {}): ArtpieceFile {
  return {
    id,
    file_url: `https://cdn/${id}`,
    thumbnail_url: null,
    mime_type: 'image/png',
    name: id,
    notes: null,
    is_visible: true,
    extension: 'png',
    width: null,
    height: null,
    ...overrides,
  }
}

const editProps = {
  mode: 'edit' as const,
  locallyRemovedIds: new Set<string>(),
  pendingCoverFileId: null,
  onSetCover: () => {},
  onRemove: () => {},
  onUndoRemove: () => {},
  onAddFiles: () => {},
  dragHandlers: {
    dragOver: false,
    onDragEnter: () => {},
    onDragOver: () => {},
    onDragLeave: () => {},
    onDrop: () => {},
  },
  uploadPlaceholders: [],
}

describe('ArtpieceDetailFileGrid — hidden files', () => {
  const files = [
    makeFile('shown', { extension: 'png' }),
    makeFile('hidden', { is_visible: false, extension: 'psd' }),
  ]

  it('does not render a hidden file as a tile in view mode', () => {
    render(<ArtpieceDetailFileGrid files={files} coverFileId="shown" mode="view" />)
    expect(screen.getByText('PNG')).toBeInTheDocument()
    expect(screen.queryByText('PSD')).not.toBeInTheDocument()
  })

  it('does not render a hidden file as a tile in edit mode', () => {
    render(<ArtpieceDetailFileGrid files={files} coverFileId="shown" {...editProps} />)
    expect(screen.getByText('PNG')).toBeInTheDocument()
    expect(screen.queryByText('PSD')).not.toBeInTheDocument()
    expect(screen.getAllByTitle('Remove')).toHaveLength(1)
  })

  it('shows no file tiles in edit mode when every file is hidden', () => {
    render(
      <ArtpieceDetailFileGrid
        files={[makeFile('h1', { is_visible: false }), makeFile('h2', { is_visible: false })]}
        coverFileId={null}
        {...editProps}
      />,
    )
    expect(screen.queryByTitle('Remove')).not.toBeInTheDocument()
    expect(screen.queryByTitle('Set as cover')).not.toBeInTheDocument()
    expect(screen.queryByText('PNG')).not.toBeInTheDocument()
    expect(screen.getByRole('button', { name: /add files/i })).toBeInTheDocument()
  })

  it('shows the empty state in view mode when every file is hidden', () => {
    render(
      <ArtpieceDetailFileGrid
        files={[makeFile('h1', { is_visible: false })]}
        coverFileId={null}
        mode="view"
      />,
    )
    expect(screen.getByText('No files attached.')).toBeInTheDocument()
  })
})
