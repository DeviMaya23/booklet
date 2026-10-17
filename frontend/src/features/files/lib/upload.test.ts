import { describe, it, expect, vi, beforeEach } from 'vitest'
import { uploadFile } from './upload'
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

    expect(initMutate).toHaveBeenCalledWith({ mimeType: 'application/octet-stream', name: 'my_artwork', extension: 'clip' })
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
    expect(initMutate).toHaveBeenCalledWith({ mimeType: 'text/plain', name: 'README', extension: undefined })
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
