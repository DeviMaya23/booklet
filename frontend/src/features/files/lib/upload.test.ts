import { describe, it, expect, vi, beforeEach } from 'vitest'
import { uploadFile, filterBySize } from './upload'
import { generateThumbnail } from './thumbnail'

vi.mock('./thumbnail', () => ({ generateThumbnail: vi.fn() }))

const mockFetch = vi.fn()
vi.stubGlobal('fetch', mockFetch)

const thumbnailBlob = new Blob(['thumb'], { type: 'image/jpeg' })

const makeInitResult = (thumbnailUploadUrl: string | null = null) => ({
  id: 'file-123',
  upload_url: 'https://r2.example.com/upload',
  thumbnail_upload_url: thumbnailUploadUrl,
  expires_at: '2026-10-18T00:00:00Z',
})

function makeFile(name: string, type: string): File {
  return new File(['data'], name, { type })
}

beforeEach(() => {
  vi.clearAllMocks()
  mockFetch.mockResolvedValue({ ok: true })
  vi.mocked(generateThumbnail).mockResolvedValue({ blob: thumbnailBlob, width: 1200, height: 800 })
})

describe('uploadFile — non-image', () => {
  it('PUTs the file only and calls completeMutate with thumbnailUploaded false', async () => {
    const initMutate = vi.fn().mockResolvedValue(makeInitResult(null))
    const completeMutate = vi.fn().mockResolvedValue(undefined)

    await uploadFile(makeFile('artwork.psd', 'image/vnd.adobe.photoshop'), initMutate, completeMutate)

    expect(mockFetch).toHaveBeenCalledOnce()
    expect(mockFetch).toHaveBeenCalledWith('https://r2.example.com/upload', { method: 'PUT', body: expect.any(File) })
    expect(completeMutate).toHaveBeenCalledWith({ pendingId: 'file-123', width: undefined, height: undefined, thumbnailUploaded: false })
  })

  it('returns null thumbnailBlob for a non-image', async () => {
    const initMutate = vi.fn().mockResolvedValue(makeInitResult(null))
    const completeMutate = vi.fn().mockResolvedValue(undefined)

    const result = await uploadFile(makeFile('archive.zip', 'application/zip'), initMutate, completeMutate)

    expect(result.thumbnailBlob).toBeNull()
    expect(result.width).toBeUndefined()
    expect(result.height).toBeUndefined()
  })
})

describe('uploadFile — image with thumbnail URL', () => {
  it('generates thumbnail, parallel PUTs both, and calls completeMutate with thumbnailUploaded true', async () => {
    const initMutate = vi.fn().mockResolvedValue(makeInitResult('https://r2.example.com/thumb'))
    const completeMutate = vi.fn().mockResolvedValue(undefined)

    await uploadFile(makeFile('photo.jpg', 'image/jpeg'), initMutate, completeMutate)

    expect(generateThumbnail).toHaveBeenCalledOnce()
    expect(mockFetch).toHaveBeenCalledTimes(2)
    expect(mockFetch).toHaveBeenCalledWith('https://r2.example.com/upload', { method: 'PUT', body: expect.any(File) })
    expect(mockFetch).toHaveBeenCalledWith('https://r2.example.com/thumb', {
      method: 'PUT',
      body: thumbnailBlob,
      headers: { 'Content-Type': thumbnailBlob.type },
    })
    expect(completeMutate).toHaveBeenCalledWith({ pendingId: 'file-123', width: 1200, height: 800, thumbnailUploaded: true })
  })

  it('returns the thumbnail blob and original dimensions', async () => {
    const initMutate = vi.fn().mockResolvedValue(makeInitResult('https://r2.example.com/thumb'))
    const completeMutate = vi.fn().mockResolvedValue(undefined)

    const result = await uploadFile(makeFile('photo.jpg', 'image/jpeg'), initMutate, completeMutate)

    expect(result.thumbnailBlob).toBe(thumbnailBlob)
    expect(result.width).toBe(1200)
    expect(result.height).toBe(800)
  })
})

describe('uploadFile — thumbnail generation failure', () => {
  it('falls back to single PUT and calls completeMutate with thumbnailUploaded false', async () => {
    vi.mocked(generateThumbnail).mockRejectedValueOnce(new Error('Unsupported format'))
    const initMutate = vi.fn().mockResolvedValue(makeInitResult('https://r2.example.com/thumb'))
    const completeMutate = vi.fn().mockResolvedValue(undefined)

    await uploadFile(makeFile('photo.heic', 'image/heic'), initMutate, completeMutate)

    expect(mockFetch).toHaveBeenCalledOnce()
    expect(mockFetch).toHaveBeenCalledWith('https://r2.example.com/upload', { method: 'PUT', body: expect.any(File) })
    expect(completeMutate).toHaveBeenCalledWith({ pendingId: 'file-123', width: undefined, height: undefined, thumbnailUploaded: false })
  })

  it('returns null thumbnailBlob on generation failure', async () => {
    vi.mocked(generateThumbnail).mockRejectedValueOnce(new Error('Unsupported format'))
    const initMutate = vi.fn().mockResolvedValue(makeInitResult('https://r2.example.com/thumb'))
    const completeMutate = vi.fn().mockResolvedValue(undefined)

    const result = await uploadFile(makeFile('photo.heic', 'image/heic'), initMutate, completeMutate)

    expect(result.thumbnailBlob).toBeNull()
  })
})

describe('uploadFile — file metadata', () => {
  it('passes stripped name and extension to initMutate', async () => {
    const initMutate = vi.fn().mockResolvedValue(makeInitResult(null))
    const completeMutate = vi.fn().mockResolvedValue(undefined)

    await uploadFile(makeFile('my_artwork.clip', 'application/octet-stream'), initMutate, completeMutate)

    expect(initMutate).toHaveBeenCalledWith({ mimeType: 'application/octet-stream', name: 'my_artwork', extension: 'clip', fileSizeBytes: 4 })
  })

  it('returns the stripped name and extension in the result', async () => {
    const initMutate = vi.fn().mockResolvedValue(makeInitResult(null))
    const completeMutate = vi.fn().mockResolvedValue(undefined)

    const result = await uploadFile(makeFile('portrait.jpg', 'image/jpeg'), initMutate, completeMutate)

    expect(result.name).toBe('portrait')
    expect(result.extension).toBe('jpg')
  })

  it('returns null extension for a file with no dot in the name', async () => {
    const initMutate = vi.fn().mockResolvedValue(makeInitResult(null))
    const completeMutate = vi.fn().mockResolvedValue(undefined)

    const result = await uploadFile(makeFile('README', 'text/plain'), initMutate, completeMutate)

    expect(result.extension).toBeNull()
    expect(initMutate).toHaveBeenCalledWith({ mimeType: 'text/plain', name: 'README', extension: undefined, fileSizeBytes: 4 })
  })

  it('falls back to application/octet-stream when file.type is empty', async () => {
    const initMutate = vi.fn().mockResolvedValue(makeInitResult(null))
    const completeMutate = vi.fn().mockResolvedValue(undefined)

    await uploadFile(makeFile('data.bin', ''), initMutate, completeMutate)

    expect(initMutate).toHaveBeenCalledWith(expect.objectContaining({ mimeType: 'application/octet-stream' }))
  })

  it('returns the fileId from the init result', async () => {
    const initMutate = vi.fn().mockResolvedValue(makeInitResult(null))
    const completeMutate = vi.fn().mockResolvedValue(undefined)

    const result = await uploadFile(makeFile('photo.png', 'image/png'), initMutate, completeMutate)

    expect(result.fileId).toBe('file-123')
  })
})

describe('filterBySize', () => {
  function makeFileWithSize(name: string, sizeBytes: number): File {
    const blob = new Blob([new ArrayBuffer(sizeBytes)])
    return new File([blob], name)
  }

  it('accepts a file within its category limit', () => {
    // image limit is 50 MB; 1 MB is well within
    const file = makeFileWithSize('photo.png', 1 * 1024 * 1024)
    const { accepted, rejected } = filterBySize([file])
    expect(accepted).toHaveLength(1)
    expect(rejected).toHaveLength(0)
  })

  it('rejects a file exceeding its category limit with the correct limitBytes and categoryName', () => {
    // image limit is 50 MB; 51 MB exceeds it
    const limit = 50 * 1024 * 1024
    const file = makeFileWithSize('photo.png', limit + 1)
    const { accepted, rejected } = filterBySize([file])
    expect(accepted).toHaveLength(0)
    expect(rejected).toHaveLength(1)
    expect(rejected[0].file).toBe(file)
    expect(rejected[0].limitBytes).toBe(limit)
    expect(rejected[0].categoryName).toBe('image')
  })

  it('handles a mixed batch, accepting and rejecting independently', () => {
    const imageLimit = 50 * 1024 * 1024
    const small = makeFileWithSize('icon.png', 1024)
    const large = makeFileWithSize('render.png', imageLimit + 1)
    const { accepted, rejected } = filterBySize([small, large])
    expect(accepted).toHaveLength(1)
    expect(accepted[0]).toBe(small)
    expect(rejected).toHaveLength(1)
    expect(rejected[0].file).toBe(large)
  })

  it('returns empty arrays for an empty input', () => {
    const { accepted, rejected } = filterBySize([])
    expect(accepted).toHaveLength(0)
    expect(rejected).toHaveLength(0)
  })

  it('uses the fallback limit and categoryName "other" for a file with no extension', () => {
    // fallback limit is 100 MB
    const fallbackLimit = 100 * 1024 * 1024
    const file = makeFileWithSize('README', fallbackLimit + 1)
    const { rejected } = filterBySize([file])
    expect(rejected).toHaveLength(1)
    expect(rejected[0].limitBytes).toBe(fallbackLimit)
    expect(rejected[0].categoryName).toBe('other')
  })

  it('accepts a file at the exact category limit (boundary is inclusive)', () => {
    const limit = 50 * 1024 * 1024
    const file = makeFileWithSize('photo.png', limit)
    const { accepted, rejected } = filterBySize([file])
    expect(accepted).toHaveLength(1)
    expect(rejected).toHaveLength(0)
  })
})
