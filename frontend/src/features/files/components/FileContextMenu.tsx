import { useEffect, useRef } from 'react'

interface FileContextMenuProps {
  x: number
  y: number
  onNewArtpiece: () => void
  onAddToExisting: () => void
  onDelete: () => void
  onClose: () => void
}

export default function FileContextMenu({ x, y, onNewArtpiece, onAddToExisting, onDelete, onClose }: FileContextMenuProps) {
  const ref = useRef<HTMLDivElement>(null)

  useEffect(() => {
    function handleClickOutside(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) {
        onClose()
      }
    }
    function handleKeyDown(e: KeyboardEvent) {
      if (e.key === 'Escape') onClose()
    }
    document.addEventListener('mousedown', handleClickOutside)
    document.addEventListener('keydown', handleKeyDown)
    return () => {
      document.removeEventListener('mousedown', handleClickOutside)
      document.removeEventListener('keydown', handleKeyDown)
    }
  }, [onClose])

  return (
    <div
      ref={ref}
      style={{ position: 'fixed', top: y, left: x, zIndex: 50 }}
      className="min-w-40 rounded-md border bg-popover text-popover-foreground shadow-md py-1"
    >
      <button
        className="w-full px-3 py-1.5 text-sm text-left hover:bg-accent"
        onClick={() => { onNewArtpiece(); onClose() }}
      >
        New Artpiece
      </button>
      <button
        className="w-full px-3 py-1.5 text-sm text-left hover:bg-accent"
        onClick={() => { onAddToExisting(); onClose() }}
      >
        Add to Existing Artpiece
      </button>
      <div className="my-1 h-px bg-border" />
      <button
        className="w-full px-3 py-1.5 text-sm text-left text-destructive hover:bg-accent"
        onClick={() => { onDelete(); onClose() }}
      >
        Delete
      </button>
    </div>
  )
}
