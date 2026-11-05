import type { InitAvatarUploadInput, InitAvatarUploadResult } from '../api/useInitAvatarUpload'
import type { CompleteAvatarUploadInput } from '../api/useCompleteAvatarUpload'

type InitMutate = (input: InitAvatarUploadInput) => Promise<InitAvatarUploadResult>
type CompleteMutate = (input: CompleteAvatarUploadInput) => Promise<unknown>

/** init → PUT to the presigned URL → complete. Rejects if any step fails. */
export async function uploadCharacterAvatar(
  characterId: string,
  file: File,
  initMutate: InitMutate,
  completeMutate: CompleteMutate,
): Promise<void> {
  const init = await initMutate({ characterId, mimeType: file.type, fileSizeBytes: file.size })
  const res = await fetch(init.upload_url, {
    method: 'PUT',
    headers: { 'Content-Type': file.type },
    body: file,
  })
  if (!res.ok) throw new Error('R2 upload failed')
  await completeMutate({ characterId, uploadId: init.id })
}
