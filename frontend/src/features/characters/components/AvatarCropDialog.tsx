import { useCallback, useState } from 'react'
import Cropper from 'react-easy-crop'
import type { Area } from 'react-easy-crop'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog'
import { getCroppedFile } from '../lib/cropCanvas'

interface AvatarCropDialogProps {
  open: boolean
  sourceUrl: string
  onCrop: (file: File) => void
  onCancel: () => void
}

export default function AvatarCropDialog({
  open,
  sourceUrl,
  onCrop,
  onCancel,
}: AvatarCropDialogProps) {
  const [crop, setCrop] = useState({ x: 0, y: 0 })
  const [zoom, setZoom] = useState(1)
  const [croppedAreaPixels, setCroppedAreaPixels] = useState<Area | null>(null)
  const [isExporting, setIsExporting] = useState(false)

  const onCropComplete = useCallback((_: Area, pixels: Area) => {
    setCroppedAreaPixels(pixels)
  }, [])

  async function handleCrop() {
    if (!croppedAreaPixels) return
    setIsExporting(true)
    try {
      const file = await getCroppedFile(sourceUrl, croppedAreaPixels, 'image/jpeg')
      onCrop(file)
    } finally {
      setIsExporting(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={(next) => { if (!next) onCancel() }}>
      <DialogContent className="max-w-sm">
        <DialogHeader>
          <DialogTitle>Crop avatar</DialogTitle>
        </DialogHeader>

        <div className="relative h-72 w-full overflow-hidden rounded-md bg-black">
          <Cropper
            image={sourceUrl}
            crop={crop}
            zoom={zoom}
            aspect={1}
            cropShape="rect"
            onCropChange={setCrop}
            onZoomChange={setZoom}
            onCropComplete={onCropComplete}
          />
        </div>

        <div className="flex items-center gap-3">
          <span className="text-xs text-muted-foreground">Zoom</span>
          <input
            type="range"
            min={0.5}
            max={3}
            step={0.05}
            value={zoom}
            onChange={(e) => setZoom(Number(e.target.value))}
            className="h-1.5 w-full cursor-pointer appearance-none rounded-full bg-muted accent-primary"
          />
        </div>

        <DialogFooter>
          <Button variant="outline" type="button" onClick={onCancel} disabled={isExporting}>
            Cancel
          </Button>
          <Button type="button" onClick={handleCrop} disabled={isExporting || !croppedAreaPixels}>
            Crop
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
