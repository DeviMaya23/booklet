import { useState } from 'react'
import ResourceCard from '@/components/ResourceCard'
import { type ArtpieceSummary } from '../api/useArtpieces'
import DeleteArtpieceDialog from './DeleteArtpieceDialog'

interface ArtpiecesGridProps {
  artpieces: ArtpieceSummary[]
  onArtpieceOpen: (id: string) => void
}

export default function ArtpiecesGrid({ artpieces, onArtpieceOpen }: ArtpiecesGridProps) {
  const [pendingDeleteId, setPendingDeleteId] = useState<string | null>(null)

  if (artpieces.length === 0) {
    return <p className="text-sm text-muted-foreground">No artpieces yet.</p>
  }

  return (
    <>
      <div className="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6">
        {artpieces.map((artpiece) => (
          <ResourceCard
            key={artpiece.id}
            imageUrl={artpiece.thumbnail_url}
            label={artpiece.title ?? 'Untitled'}
            sublabel={artpiece.artist_name ?? undefined}
            onDeleteClick={() => setPendingDeleteId(artpiece.id)}
            onDoubleClick={() => onArtpieceOpen(artpiece.id)}
          />
        ))}
      </div>

      <DeleteArtpieceDialog
        artpieceId={pendingDeleteId ?? ''}
        open={pendingDeleteId !== null}
        onOpenChange={(open) => { if (!open) setPendingDeleteId(null) }}
        onSuccess={() => setPendingDeleteId(null)}
      />
    </>
  )
}
