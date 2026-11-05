import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { uploadCharacterAvatar } from './uploadCharacterAvatar'

const file = new File(['abc'], 'avatar.jpg', { type: 'image/jpeg' })
const initResult = { id: 'up-1', upload_url: 'https://r2/put', expires_at: '' }

let fetchMock: ReturnType<typeof vi.fn>

beforeEach(() => {
  fetchMock = vi.fn().mockResolvedValue({ ok: true } as Response)
  vi.stubGlobal('fetch', fetchMock)
})

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('uploadCharacterAvatar', () => {
  it('runs init, PUT and complete in order with the file details', async () => {
    const calls: string[] = []
    const init = vi.fn().mockImplementation(async () => {
      calls.push('init')
      return initResult
    })
    fetchMock.mockImplementation(async () => {
      calls.push('put')
      return { ok: true } as Response
    })
    const complete = vi.fn().mockImplementation(async () => {
      calls.push('complete')
    })

    await uploadCharacterAvatar('c-1', file, init, complete)

    expect(calls).toEqual(['init', 'put', 'complete'])
    expect(init).toHaveBeenCalledWith({ characterId: 'c-1', mimeType: 'image/jpeg', fileSizeBytes: file.size })
    expect(fetchMock).toHaveBeenCalledWith('https://r2/put', {
      method: 'PUT',
      headers: { 'Content-Type': 'image/jpeg' },
      body: file,
    })
    expect(complete).toHaveBeenCalledWith({ characterId: 'c-1', uploadId: 'up-1' })
  })

  it('rejects before complete when the PUT fails', async () => {
    const init = vi.fn().mockResolvedValue(initResult)
    const complete = vi.fn()
    fetchMock.mockResolvedValue({ ok: false } as Response)

    await expect(uploadCharacterAvatar('c-1', file, init, complete)).rejects.toThrow('R2 upload failed')

    expect(complete).not.toHaveBeenCalled()
  })
})
