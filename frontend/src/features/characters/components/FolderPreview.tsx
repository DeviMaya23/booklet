import { useState } from 'react'
import { ExternalLink, RefreshCw } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { useGetFolderImages } from '../api/useGetFolderImages'

const PAGE_SIZE = 24

interface FolderPreviewProps {
  folderID: string
  /** Bookleaf share URL for the folder; enables the "View all in Bookleaf" link. */
  shareUrl?: string
}

export function FolderPreview({ folderID, shareUrl }: FolderPreviewProps) {
  const [visibleCount, setVisibleCount] = useState(PAGE_SIZE)
  const { data, isLoading, isError, refetch } = useGetFolderImages(folderID)

  if (isLoading) {
    return (
      <div className="flex flex-col gap-2">
        <div className="grid grid-cols-8 gap-1">
          {Array.from({ length: PAGE_SIZE }).map((_, i) => (
            <div key={i} className="aspect-square rounded bg-muted animate-pulse" />
          ))}
        </div>
        <p className="text-xs text-muted-foreground">Loading images…</p>
      </div>
    )
  }

  if (isError) {
    return (
      <div className="flex items-center justify-between rounded-md bg-destructive/5 px-3 py-2 text-sm">
        <span className="text-destructive">Couldn't load the images for this folder.</span>
        <Button
          type="button"
          variant="outline"
          size="sm"
          onClick={() => refetch()}
          className="ml-3 shrink-0 gap-1.5"
        >
          <RefreshCw className="size-3.5" />
          Retry
        </Button>
      </div>
    )
  }

  if (!data || data.length === 0) {
    return <p className="text-xs text-muted-foreground">No images in this folder.</p>
  }

  const visible = data.slice(0, visibleCount)
  const hasMore = visibleCount < data.length

  return (
    <div className="flex flex-col gap-2">
      <div className="grid grid-cols-8 gap-1">
        {visible.map((img) => (
          <img
            key={img.image_id}
            src={img.thumbnail_url}
            alt="thumbnail"
            className="aspect-square rounded object-cover"
          />
        ))}
      </div>
      <div className="flex items-center gap-2">
        <p className="mr-auto text-xs text-muted-foreground">
          {hasMore ? `Showing the first ${visible.length}` : `Showing all ${data.length}`}
        </p>
        {hasMore && (
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() => setVisibleCount((prev) => prev + PAGE_SIZE)}
          >
            Show more
          </Button>
        )}
        {shareUrl && (
          <a
            href={shareUrl}
            target="_blank"
            rel="noreferrer"
            className="inline-flex items-center gap-1 text-xs font-medium hover:underline"
          >
            View all in Bookleaf
            <ExternalLink className="size-3" />
          </a>
        )}
      </div>
    </div>
  )
}
