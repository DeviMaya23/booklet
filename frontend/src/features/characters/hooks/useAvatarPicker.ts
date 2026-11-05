import { useCallback, useEffect, useState } from 'react'
import { toast } from 'sonner'
import { FILE_SIZE_LIMITS } from '@/features/files/lib/fileCategory'

interface PasteLike {
  target: EventTarget | null
  clipboardData: DataTransfer | null
}

interface UseAvatarPickerOptions {
  /** The avatar already saved on the server, shown until a local change is made. */
  currentUrl?: string | null
  /** Also accept pastes anywhere on the page (a page has no focused container, unlike a dialog). */
  listenForPaste?: boolean
}

export function useAvatarPicker({ currentUrl = null, listenForPaste = false }: UseAvatarPickerOptions = {}) {
  const [localFile, setLocalFile] = useState<File | null>(null)
  const [localPreviewUrl, setLocalPreviewUrl] = useState<string | null>(null)
  const [avatarCleared, setAvatarCleared] = useState(false)
  const [cropSourceUrl, setCropSourceUrl] = useState<string | null>(null)
  const [dragActive, setDragActive] = useState(false)

  useEffect(() => {
    return () => {
      if (localPreviewUrl) URL.revokeObjectURL(localPreviewUrl)
    }
  }, [localPreviewUrl])

  useEffect(() => {
    return () => {
      if (cropSourceUrl) URL.revokeObjectURL(cropSourceUrl)
    }
  }, [cropSourceUrl])

  const displayUrl = localPreviewUrl ?? (avatarCleared ? null : currentUrl)
  const canClear = !avatarCleared && (localPreviewUrl !== null || currentUrl !== null)
  const isChanged = localFile !== null || avatarCleared

  const stageFile = useCallback((file: File) => {
    const imageLimit = FILE_SIZE_LIMITS['image']
    if (file.size > imageLimit) {
      toast.error(`Image exceeds the ${Math.round(imageLimit / (1024 * 1024))} MB limit`)
      return
    }
    setCropSourceUrl(URL.createObjectURL(file))
  }, [])

  function handleFileChange(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0]
    if (!file) return
    stageFile(file)
    e.target.value = ''
  }

  function handleDrop(e: React.DragEvent) {
    e.preventDefault()
    setDragActive(false)
    const file = e.dataTransfer.files?.[0]
    if (file) stageFile(file)
  }

  const dragHandlers = {
    onDragOver: (e: React.DragEvent) => {
      e.preventDefault()
      setDragActive(true)
    },
    onDragLeave: () => setDragActive(false),
    onDrop: handleDrop,
  }

  // Stable identity so callers can bind it to `document` without re-subscribing every render.
  const handlePaste = useCallback(
    (e: PasteLike) => {
      if (e.target instanceof HTMLInputElement || e.target instanceof HTMLTextAreaElement) return
      const item = Array.from(e.clipboardData?.items ?? []).find((i) => i.type.startsWith('image/'))
      const file = item?.getAsFile()
      if (file) stageFile(file)
    },
    [stageFile],
  )

  useEffect(() => {
    if (!listenForPaste) return
    document.addEventListener('paste', handlePaste)
    return () => document.removeEventListener('paste', handlePaste)
  }, [listenForPaste, handlePaste])

  function confirmCrop(cropped: File) {
    setLocalFile(cropped)
    setLocalPreviewUrl(URL.createObjectURL(cropped))
    setAvatarCleared(false)
    setCropSourceUrl(null)
  }

  function cancelCrop() {
    setCropSourceUrl(null)
  }

  function clear() {
    setLocalFile(null)
    setLocalPreviewUrl(null)
    setAvatarCleared(true)
  }

  function reset() {
    setLocalFile(null)
    setLocalPreviewUrl(null)
    setAvatarCleared(false)
    setCropSourceUrl(null)
    setDragActive(false)
  }

  return {
    localFile,
    displayUrl,
    avatarCleared,
    canClear,
    isChanged,
    cropSourceUrl,
    dragActive,
    dragHandlers,
    handleFileChange,
    handlePaste,
    confirmCrop,
    cancelCrop,
    clear,
    reset,
  }
}

export type AvatarPicker = ReturnType<typeof useAvatarPicker>
