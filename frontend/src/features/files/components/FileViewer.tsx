import { useCallback, useEffect, useState } from 'react'
import { Dialog as DialogPrimitive } from '@base-ui/react/dialog'
import { ChevronLeft, ChevronRight, Download, FileIcon, X } from 'lucide-react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { mimeTypeDescription, mimeTypeLabel } from '../lib/mimeTypeLabel'
import { useDownloadFile } from '../api/useDownloadFile'

export interface ViewerFile {
  id: string
  name: string | null
  mimeType: string
  thumbnailUrl: string | null
  previewUrl: string
  width?: number | null
  height?: number | null
}

interface FileViewerProps {
  files: ViewerFile[]
  initialIndex: number
  onClose: () => void
}

function isPreviewable(file: ViewerFile): boolean {
  return file.mimeType.startsWith('image/') && file.thumbnailUrl !== null
}

export default function FileViewer({ files, initialIndex, onClose }: FileViewerProps) {
  const [index, setIndex] = useState(initialIndex)
  const download = useDownloadFile()

  const prev = useCallback(() => setIndex((i) => Math.max(0, i - 1)), [])
  const next = useCallback(() => setIndex((i) => Math.min(files.length - 1, i + 1)), [files.length])

  useEffect(() => {
    function onKey(e: KeyboardEvent) {
      if (e.key === 'ArrowLeft') prev()
      if (e.key === 'ArrowRight') next()
    }
    document.addEventListener('keydown', onKey)
    return () => document.removeEventListener('keydown', onKey)
  }, [prev, next])

  const file = files[index]

  if (!file) {
    onClose()
    return null
  }

  const hasPrev = index > 0
  const hasNext = index < files.length - 1

  function handleDownload() {
    download.mutate(file.id, {
      onError: () => toast.error('Failed to download file'),
    })
  }

  const label = mimeTypeLabel(file.mimeType)
  const preview = isPreviewable(file)

  return (
    <DialogPrimitive.Root open onOpenChange={(open) => { if (!open) onClose() }}>
      <DialogPrimitive.Portal>
        <DialogPrimitive.Popup
          aria-label="File viewer"
          className="fixed inset-0 z-50 flex flex-col bg-black/90 outline-none"
        >
      {/* Top bar */}
      <div className="flex shrink-0 items-center gap-3 border-b border-white/10 px-4 py-3">
        <span className="truncate text-sm font-medium text-white">
          {file.name ?? 'Untitled'}
        </span>
        {label && (
          <span className="shrink-0 rounded bg-white/20 px-1.5 py-0.5 text-[11px] font-semibold text-white">
            {label}
          </span>
        )}
        {!!file.width && !!file.height && (
          <span className="shrink-0 text-xs text-white/50">
            {file.width} × {file.height} px
          </span>
        )}
        <span className="ml-auto shrink-0 text-xs text-white/50">
          {index + 1} of {files.length}
        </span>
        <Button
          size="sm"
          variant="outline"
          className="shrink-0 border-white/20 bg-transparent text-white hover:bg-white/10 hover:text-white"
          onClick={handleDownload}
          disabled={download.isPending}
        >
          <Download className="mr-1.5 size-3.5" />
          Download
        </Button>
        <button
          autoFocus
          type="button"
          className="shrink-0 rounded p-1 text-white/70 hover:bg-white/10 hover:text-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/30"
          onClick={onClose}
          aria-label="Close viewer"
        >
          <X className="size-5" />
        </button>
      </div>

      {/* Main area */}
      <div className="relative flex flex-1 items-center justify-center overflow-hidden">
        {/* Prev arrow */}
        {hasPrev && (
          <button
            type="button"
            onClick={prev}
            className="absolute left-4 z-10 flex size-10 items-center justify-center rounded-full bg-black/50 text-white hover:bg-black/70 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/30"
            aria-label="Previous file"
          >
            <ChevronLeft className="size-6" />
          </button>
        )}

        {preview ? (
          /* Image preview with checkerboard tile */
          <div
            className="relative max-h-full max-w-full overflow-hidden rounded-md"
            style={{
              backgroundImage:
                'repeating-conic-gradient(#555 0% 25%, #333 0% 50%)',
              backgroundSize: '16px 16px',
            }}
          >
            <img
              key={file.id}
              src={file.previewUrl}
              alt={file.name ?? ''}
              className="block max-h-[calc(100vh-10rem)] max-w-[calc(100vw-8rem)] object-contain"
            />
          </div>
        ) : (
          /* No-preview fallback card */
          <div className="flex w-72 flex-col items-center gap-4 rounded-xl bg-white/10 p-8 text-center">
            <FileIcon className="size-12 text-white/50" />
            <div className="flex flex-col gap-1">
              <p className="text-sm font-medium text-white">{file.name ?? 'Untitled'}</p>
              <p className="text-xs text-white/50">
                {mimeTypeDescription(file.mimeType)} · no preview available
              </p>
            </div>
            <Button
              variant="outline"
              className="border-white/20 bg-transparent text-white hover:bg-white/10 hover:text-white"
              onClick={handleDownload}
              disabled={download.isPending}
            >
              <Download className="mr-1.5 size-4" />
              Download file
            </Button>
          </div>
        )}

        {/* Next arrow */}
        {hasNext && (
          <button
            type="button"
            onClick={next}
            className="absolute right-4 z-10 flex size-10 items-center justify-center rounded-full bg-black/50 text-white hover:bg-black/70 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/30"
            aria-label="Next file"
          >
            <ChevronRight className="size-6" />
          </button>
        )}
      </div>

      {/* Filmstrip */}
      {files.length > 1 && (
        <div className="flex shrink-0 justify-center gap-2 overflow-x-auto px-4 py-3">
          {files.map((f, i) => (
            <button
              key={f.id}
              type="button"
              onClick={() => setIndex(i)}
              className={`size-12 shrink-0 overflow-hidden rounded focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/50 ${
                i === index ? 'ring-2 ring-white' : 'ring-1 ring-white/20'
              }`}
              aria-label={`Go to file ${i + 1}`}
            >
              {f.thumbnailUrl ? (
                <img src={f.thumbnailUrl} alt="" className="h-full w-full object-cover" />
              ) : (
                <div className="flex h-full w-full items-center justify-center bg-white/10">
                  <FileIcon className="size-4 text-white/50" />
                </div>
              )}
            </button>
          ))}
        </div>
      )}
        </DialogPrimitive.Popup>
      </DialogPrimitive.Portal>
    </DialogPrimitive.Root>
  )
}
