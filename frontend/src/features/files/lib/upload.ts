import { generateThumbnail } from './thumbnail'
import { stripFileExtension } from './files'
import type { InitFileUploadResult } from '../api/useInitFileUpload'

export interface UploadFileResult {
  fileId: string
  name: string
  extension: string | null
  thumbnailBlob: Blob | null
  width: number | undefined
  height: number | undefined
}

type InitMutate = (params: { mimeType: string; name?: string; extension?: string }) => Promise<InitFileUploadResult>
type CompleteMutate = (params: { pendingId: string; width?: number; height?: number; thumbnailUploaded: boolean }) => Promise<unknown>

export async function uploadFile(
  file: File,
  initMutate: InitMutate,
  completeMutate: CompleteMutate,
): Promise<UploadFileResult> {
  const mimeType = file.type || 'application/octet-stream'
  const dotIdx = file.name.lastIndexOf('.')
  const extension = dotIdx > 0 ? file.name.slice(dotIdx + 1).toLowerCase() : null
  const name = stripFileExtension(file.name)

  const init = await initMutate({ mimeType, name, extension: extension ?? undefined })

  const isImage = mimeType.startsWith('image/')
  let width: number | undefined
  let height: number | undefined
  let thumbnailBlob: Blob | null = null
  let thumbnailUploaded = false

  if (isImage && init.thumbnail_upload_url) {
    try {
      const thumb = await generateThumbnail(file)
      width = thumb.width
      height = thumb.height
      await Promise.all([
        fetch(init.upload_url, { method: 'PUT', body: file }),
        fetch(init.thumbnail_upload_url, { method: 'PUT', body: thumb.blob, headers: { 'Content-Type': thumb.blob.type } }),
      ])
      thumbnailBlob = thumb.blob
      thumbnailUploaded = true
    } catch {
      await fetch(init.upload_url, { method: 'PUT', body: file })
    }
  } else {
    await fetch(init.upload_url, { method: 'PUT', body: file })
  }

  await completeMutate({ pendingId: init.id, width, height, thumbnailUploaded })

  return { fileId: init.id, name, extension, thumbnailBlob, width, height }
}
