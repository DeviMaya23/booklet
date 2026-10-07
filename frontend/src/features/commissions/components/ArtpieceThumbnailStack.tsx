import { Link } from 'react-router-dom'

interface ArtpieceThumbnail {
  id: string
  thumbnail_url: string | null
}

interface ArtpieceThumbnailStackProps {
  artpieces: ArtpieceThumbnail[]
}

export default function ArtpieceThumbnailStack({ artpieces }: ArtpieceThumbnailStackProps) {
  if (artpieces.length === 0) {
    return <span className="text-muted-foreground">—</span>
  }

  const visible = artpieces.slice(0, 2)
  const overflow = artpieces.length - visible.length

  return (
    <span className="inline-flex items-center gap-1">
      {visible.map((a) => (
        <Link
          key={a.id}
          to={`/app/artpieces/${a.id}`}
          className="block size-7 shrink-0 overflow-hidden rounded border border-border bg-muted hover:opacity-80 transition-opacity"
          aria-label="Open artpiece"
        >
          {a.thumbnail_url ? (
            <img src={a.thumbnail_url} alt="" className="size-full object-cover" />
          ) : (
            <div className="size-full bg-muted" />
          )}
        </Link>
      ))}
      {overflow > 0 && (
        <span className="rounded-full bg-muted px-1.5 py-0.5 text-xs text-muted-foreground">
          +{overflow}
        </span>
      )}
    </span>
  )
}
