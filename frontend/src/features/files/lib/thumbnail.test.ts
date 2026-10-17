import { describe, it, expect, vi, beforeEach } from 'vitest'
import { generateThumbnail } from './thumbnail'

const mockClose = vi.fn()
const mockDrawImage = vi.fn()
const mockConvertToBlob = vi.fn()
const mockGetContext = vi.fn()
const mockCreateImageBitmap = vi.fn()

beforeEach(() => {
  vi.clearAllMocks()

  mockCreateImageBitmap.mockResolvedValue({ width: 1200, height: 800, close: mockClose })
  mockGetContext.mockReturnValue({ drawImage: mockDrawImage })
  mockConvertToBlob.mockResolvedValue(new Blob(['x'], { type: 'image/jpeg' }))

  ;(globalThis as unknown as Record<string, unknown>).createImageBitmap = mockCreateImageBitmap
  ;(globalThis as unknown as Record<string, unknown>).OffscreenCanvas = class {
    getContext = mockGetContext
    convertToBlob = mockConvertToBlob
  }
})

function makeFile(name: string, type: string): File {
  return new File(['data'], name, { type })
}

describe('generateThumbnail — output format', () => {
  it('encodes as JPEG for a JPEG source', async () => {
    await generateThumbnail(makeFile('photo.jpg', 'image/jpeg'))
    expect(mockConvertToBlob).toHaveBeenCalledWith({ type: 'image/jpeg', quality: 0.9 })
  })

  it('encodes as PNG for a PNG source (alpha-capable)', async () => {
    mockConvertToBlob.mockResolvedValueOnce(new Blob(['x'], { type: 'image/png' }))
    await generateThumbnail(makeFile('art.png', 'image/png'))
    expect(mockConvertToBlob).toHaveBeenCalledWith({ type: 'image/png', quality: undefined })
  })

  it('encodes as PNG for a GIF source', async () => {
    mockConvertToBlob.mockResolvedValueOnce(new Blob(['x'], { type: 'image/png' }))
    await generateThumbnail(makeFile('anim.gif', 'image/gif'))
    expect(mockConvertToBlob).toHaveBeenCalledWith({ type: 'image/png', quality: undefined })
  })

  it('encodes as PNG for a WebP source', async () => {
    mockConvertToBlob.mockResolvedValueOnce(new Blob(['x'], { type: 'image/png' }))
    await generateThumbnail(makeFile('art.webp', 'image/webp'))
    expect(mockConvertToBlob).toHaveBeenCalledWith({ type: 'image/png', quality: undefined })
  })
})

describe('generateThumbnail — dimensions', () => {
  it('returns original source dimensions, not scaled ones', async () => {
    const result = await generateThumbnail(makeFile('photo.jpg', 'image/jpeg'))
    expect(result.width).toBe(1200)
    expect(result.height).toBe(800)
  })

  it('scales down a large image to fit within 600px on the long side', async () => {
    // 1200×800 → scale = 600/1200 = 0.5 → 600×400
    let ctorWidth: number | undefined
    let ctorHeight: number | undefined
    ;(globalThis as unknown as Record<string, unknown>).OffscreenCanvas = class {
      constructor(w: number, h: number) { ctorWidth = w; ctorHeight = h }
      getContext = mockGetContext
      convertToBlob = mockConvertToBlob
    }
    await generateThumbnail(makeFile('photo.jpg', 'image/jpeg'))
    expect(ctorWidth).toBe(600)
    expect(ctorHeight).toBe(400)
  })

  it('does not upscale an image smaller than 600px', async () => {
    mockCreateImageBitmap.mockResolvedValueOnce({ width: 300, height: 200, close: mockClose })
    let ctorWidth: number | undefined
    let ctorHeight: number | undefined
    ;(globalThis as unknown as Record<string, unknown>).OffscreenCanvas = class {
      constructor(w: number, h: number) { ctorWidth = w; ctorHeight = h }
      getContext = mockGetContext
      convertToBlob = mockConvertToBlob
    }
    await generateThumbnail(makeFile('small.jpg', 'image/jpeg'))
    expect(ctorWidth).toBe(300)
    expect(ctorHeight).toBe(200)
  })
})

describe('generateThumbnail — resource cleanup', () => {
  it('closes the ImageBitmap after drawing', async () => {
    await generateThumbnail(makeFile('photo.jpg', 'image/jpeg'))
    expect(mockClose).toHaveBeenCalledOnce()
  })
})

describe('generateThumbnail — error cases', () => {
  it('throws when createImageBitmap rejects (unsupported format)', async () => {
    mockCreateImageBitmap.mockRejectedValueOnce(new Error('Unsupported image type'))
    await expect(generateThumbnail(makeFile('photo.heic', 'image/heic'))).rejects.toThrow()
  })

  it('throws a descriptive error when getContext returns null', async () => {
    mockGetContext.mockReturnValueOnce(null)
    await expect(generateThumbnail(makeFile('photo.jpg', 'image/jpeg'))).rejects.toThrow(
      'Failed to acquire 2D rendering context',
    )
  })
})
