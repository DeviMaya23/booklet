import type { Area } from 'react-easy-crop'

const MAX_OUTPUT_PX = 1024

export function getCroppedFile(sourceUrl: string, pixelCrop: Area, mimeType: 'image/jpeg'): Promise<File> {
  return new Promise((resolve, reject) => {
    const image = new Image()
    image.onload = () => {
      // react-easy-crop returns croppedAreaPixels in natural image coordinates,
      // which match image.naturalWidth/Height for an unattached Image element.
      const cropW = pixelCrop.width
      const cropH = pixelCrop.height
      const cropX = pixelCrop.x
      const cropY = pixelCrop.y

      const outputSize = Math.min(cropW, MAX_OUTPUT_PX)

      const canvas = document.createElement('canvas')
      canvas.width = outputSize
      canvas.height = outputSize

      const ctx = canvas.getContext('2d')
      if (!ctx) {
        reject(new Error('Failed to get canvas context'))
        return
      }

      ctx.drawImage(image, cropX, cropY, cropW, cropH, 0, 0, outputSize, outputSize)

      canvas.toBlob(
        (blob) => {
          if (!blob) {
            reject(new Error('Canvas toBlob returned null'))
            return
          }
          resolve(new File([blob], 'avatar.jpg', { type: mimeType }))
        },
        mimeType,
        0.9,
      )
    }
    image.onerror = () => reject(new Error('Failed to load image'))
    image.src = sourceUrl
  })
}
