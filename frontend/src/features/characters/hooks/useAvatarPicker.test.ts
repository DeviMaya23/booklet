import { renderHook, act } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { toast } from 'sonner'
import { FILE_SIZE_LIMITS } from '@/features/files/lib/fileCategory'
import { useAvatarPicker } from './useAvatarPicker'

vi.mock('sonner', () => ({ toast: { error: vi.fn() } }))

let urlCounter = 0

beforeEach(() => {
  urlCounter = 0
  vi.clearAllMocks()
  URL.createObjectURL = vi.fn(() => `blob:mock-${++urlCounter}`)
  URL.revokeObjectURL = vi.fn()
})

afterEach(() => {
  vi.restoreAllMocks()
})

function imageFile(size = 10) {
  const file = new File(['x'], 'a.png', { type: 'image/png' })
  Object.defineProperty(file, 'size', { value: size })
  return file
}

function fileInputEvent(file: File) {
  return { target: { files: [file], value: 'a.png' } } as unknown as React.ChangeEvent<HTMLInputElement>
}

function pasteEvent(target: EventTarget | null, file: File | null) {
  return {
    target,
    clipboardData: {
      items: file ? [{ type: file.type, getAsFile: () => file }] : [{ type: 'text/plain', getAsFile: () => null }],
    } as unknown as DataTransfer,
  }
}

describe('useAvatarPicker', () => {
  it('shows the current avatar until something changes', () => {
    const { result } = renderHook(() => useAvatarPicker({ currentUrl: 'https://img/current' }))

    expect(result.current.displayUrl).toBe('https://img/current')
    expect(result.current.isChanged).toBe(false)
    expect(result.current.canClear).toBe(true)
  })

  it('stages a picked image for cropping', () => {
    const { result } = renderHook(() => useAvatarPicker())

    act(() => result.current.handleFileChange(fileInputEvent(imageFile())))

    expect(result.current.cropSourceUrl).toBe('blob:mock-1')
  })

  it('rejects an oversized image with a toast and no crop', () => {
    const { result } = renderHook(() => useAvatarPicker())

    act(() => result.current.handleFileChange(fileInputEvent(imageFile(FILE_SIZE_LIMITS['image'] + 1))))

    expect(result.current.cropSourceUrl).toBeNull()
    expect(toast.error).toHaveBeenCalledWith(expect.stringContaining('MB limit'))
  })

  it('stages a pasted image when no text field is the target', () => {
    const { result } = renderHook(() => useAvatarPicker())

    act(() => result.current.handlePaste(pasteEvent(document.body, imageFile())))

    expect(result.current.cropSourceUrl).toBe('blob:mock-1')
  })

  it('ignores paste from inputs and textareas', () => {
    const { result } = renderHook(() => useAvatarPicker())

    act(() => result.current.handlePaste(pasteEvent(document.createElement('input'), imageFile())))
    act(() => result.current.handlePaste(pasteEvent(document.createElement('textarea'), imageFile())))

    expect(result.current.cropSourceUrl).toBeNull()
  })

  it('ignores paste with no image on the clipboard', () => {
    const { result } = renderHook(() => useAvatarPicker())

    act(() => result.current.handlePaste(pasteEvent(document.body, null)))

    expect(result.current.cropSourceUrl).toBeNull()
  })

  it('confirming a crop sets the file and preview and clears the cleared flag', () => {
    const { result } = renderHook(() => useAvatarPicker({ currentUrl: 'https://img/current' }))
    act(() => result.current.clear())
    expect(result.current.avatarCleared).toBe(true)
    const cropped = imageFile()

    act(() => result.current.handleFileChange(fileInputEvent(imageFile())))
    act(() => result.current.confirmCrop(cropped))

    expect(result.current.localFile).toBe(cropped)
    expect(result.current.displayUrl).toBe('blob:mock-2')
    expect(result.current.avatarCleared).toBe(false)
    expect(result.current.cropSourceUrl).toBeNull()
    expect(result.current.isChanged).toBe(true)
  })

  it('cancelling a crop changes nothing but the crop source', () => {
    const { result } = renderHook(() => useAvatarPicker({ currentUrl: 'https://img/current' }))

    act(() => result.current.handleFileChange(fileInputEvent(imageFile())))
    act(() => result.current.cancelCrop())

    expect(result.current.cropSourceUrl).toBeNull()
    expect(result.current.localFile).toBeNull()
    expect(result.current.displayUrl).toBe('https://img/current')
  })

  it('clear drops the local file and hides the current avatar', () => {
    const { result } = renderHook(() => useAvatarPicker({ currentUrl: 'https://img/current' }))

    act(() => result.current.confirmCrop(imageFile()))
    act(() => result.current.clear())

    expect(result.current.localFile).toBeNull()
    expect(result.current.displayUrl).toBeNull()
    expect(result.current.avatarCleared).toBe(true)
    expect(result.current.canClear).toBe(false)
  })

  it('reset restores the untouched state', () => {
    const { result } = renderHook(() => useAvatarPicker({ currentUrl: 'https://img/current' }))

    act(() => result.current.confirmCrop(imageFile()))
    act(() => result.current.reset())

    expect(result.current.isChanged).toBe(false)
    expect(result.current.displayUrl).toBe('https://img/current')
  })

  it('revokes the preview URL when it is replaced', () => {
    const { result } = renderHook(() => useAvatarPicker())

    act(() => result.current.confirmCrop(imageFile()))
    act(() => result.current.confirmCrop(imageFile()))

    expect(URL.revokeObjectURL).toHaveBeenCalledWith('blob:mock-1')
  })
})
