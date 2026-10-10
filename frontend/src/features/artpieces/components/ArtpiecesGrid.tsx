import ArtpieceTile from './ArtpieceTile'
import { type ArtpieceSummary } from '../api/useArtpieces'

export type TileSize = 'small' | 'medium' | 'large'

function columnClass(tileSize: TileSize): string {
  switch (tileSize) {
    case 'small':
      return 'grid-cols-4 sm:grid-cols-6 md:grid-cols-8'
    case 'medium':
      return 'grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5'
    case 'large':
      return 'grid-cols-1 sm:grid-cols-2 md:grid-cols-3'
  }
}

interface ArtpiecesGridProps {
  artpieces: ArtpieceSummary[]
  tileSize: TileSize
  showDetails: boolean
  onTileClick: (id: string, index: number) => void
  onInfoClick: (id: string) => void
}

export default function ArtpiecesGrid({
  artpieces,
  tileSize,
  showDetails,
  onTileClick,
  onInfoClick,
}: ArtpiecesGridProps) {
  if (artpieces.length === 0) {
    return <p className="text-sm text-muted-foreground">No artpieces yet.</p>
  }

  return (
    <div className={`grid gap-3 ${columnClass(tileSize)}`}>
      {artpieces.map((artpiece, index) => (
        <ArtpieceTile
          key={artpiece.id}
          title={artpiece.title}
          artistName={artpiece.artist_name}
          imageUrl={artpiece.thumbnail_url}
          showDetails={showDetails}
          onClick={() => onTileClick(artpiece.id, index)}
          onInfoClick={() => onInfoClick(artpiece.id)}
        />
      ))}
    </div>
  )
}
