import { toast } from 'sonner'
import { Link2, Pencil } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { type Artist } from '../api/useArtists'

interface ArtistsListProps {
  artists: Artist[]
  onEditClick: (artist: Artist) => void
}

export default function ArtistsList({ artists, onEditClick }: ArtistsListProps) {
  function handleCopyLink(artist: Artist) {
    if (!artist.artist_link) return
    navigator.clipboard.writeText(artist.artist_link).then(() => {
      toast.success('Link copied to clipboard')
    })
  }

  return (
    <div className="flex flex-col">
      {artists.map((artist, index) => (
        <div key={artist.id}>
          {index > 0 && <hr className="border-border" />}
          <div className="flex items-center gap-3 py-3">
            <span className="text-muted-foreground" aria-hidden="true">•</span>
            <span className="flex-1 text-sm font-medium">{artist.name}</span>
            <Button
              variant="ghost"
              size="icon-sm"
              onClick={() => handleCopyLink(artist)}
              disabled={!artist.artist_link}
              title={artist.artist_link ? 'Copy link' : 'No link available'}
            >
              <Link2 className="size-4" />
              <span className="sr-only">Copy link</span>
            </Button>
            <Button
              variant="ghost"
              size="icon-sm"
              onClick={() => onEditClick(artist)}
              title="Edit artist"
            >
              <Pencil className="size-4" />
              <span className="sr-only">Edit artist</span>
            </Button>
          </div>
        </div>
      ))}
    </div>
  )
}
