# Proposal

## Why

There is currently no size cap on file uploads. A user can upload arbitrarily large files without any feedback or enforcement, wasting R2 storage and degrading the upload experience. FE-only validation is bypassable by any authenticated caller, so the design pairs a FE pre-upload check (instant UX feedback) with BE enforcement via a presigned PUT URL signed with a `ContentLength`, so the limit holds even if the FE check is skipped.

## What Changes

- New per-category size limit constants added to the FE (`fileCategory.ts`), keyed by `FileCategory`.
- FE `uploadFiles()` in `ArtpieceFilesInput` and `AddFilesModal` rejects oversized files before calling init, showing a toast per rejected file. Both components share the check via a utility in `upload.ts`.
- FE `stageFileForCrop()` in `CharacterFormModal` rejects oversized source images before opening the crop dialog.
- `POST /files` init request body gains a required `file_size` field (bytes as integer).
- BE init handler validates `file_size` against a per-extension limit map and returns 422 if exceeded.
- BE `GeneratePresignedPutURL` accepts a `contentLength` parameter and includes it in the signed `PutObjectInput`, so R2 rejects any PUT whose `Content-Length` header does not match.
- Same `file_size` enforcement applies to the character avatar init endpoint for consistency (FE check guards against loading a huge source image into the crop canvas; BE enforcement is on the cropped JPEG upload).
- Both the FE category→maxBytes constants and the BE extension→maxBytes map carry a cross-reference comment pointing to the other side, since the values must stay in sync manually.

## Capabilities

### New Capabilities

_(none)_

### Modified Capabilities

- `file-upload`: init endpoint now requires `file_size`, validates it against the limit for the file's extension, and signs `ContentLength` into the presigned PUT URL.
- `character-avatar-upload`: init endpoint gains the same `file_size` + `ContentLength` enforcement; FE source-image size check is added before the crop dialog opens.
- `web-file-inbox-upload-modal`: oversized files are silently dropped from the upload batch before init is called, with a per-file toast naming the limit.
- `web-artpiece-create-modal`: same FE rejection behavior for oversized files in the `ArtpieceFilesInput` component.
- `web-character-avatar-crop`: source image size is validated before the crop dialog opens; oversized source images are rejected with a toast.

## Impact

- `backend/internal/handler/file_upload_handler.go` — `initiateFileUploadRequest` gains `FileSizeBytes int64`; handler passes it to usecase.
- `backend/internal/usecase/file_upload_usecase.go` — `InitiateFileUploadParams` gains `FileSizeBytes int64`; usecase validates against the extension→maxBytes map before generating the presigned URL.
- `backend/internal/storage/r2.go` — `GeneratePresignedPutURL` signature gains `contentLength int64`; passed as `ContentLength` on `PutObjectInput`.
- New BE extension→maxBytes lookup (location decided in design) with a cross-reference comment pointing to FE constants.
- `backend/internal/handler/character_handler.go` — same `file_size` field added to avatar init request; same usecase/storage changes apply via `GeneratePresignedPutURL`.
- `frontend/src/features/files/lib/fileCategory.ts` — adds `FILE_SIZE_LIMITS` map (`FileCategory → number`) with cross-reference comment.
- `frontend/src/features/files/lib/upload.ts` — pre-upload size check utility; `uploadFile()` passes `file.size` to init.
- `frontend/src/features/files/api/useInitFileUpload.ts` — mutation params gain `fileSizeBytes`.
- `frontend/src/features/files/components/ArtpieceFilesInput.tsx` and `AddFilesModal.tsx` — call the size check before queuing uploads.
- `frontend/src/features/characters/components/CharacterFormModal.tsx` — source image size check in `stageFileForCrop()`.
- Bruno file for the updated `POST /files` init endpoint.
