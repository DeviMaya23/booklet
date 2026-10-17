# Proposal

## Why

Thumbnail generation currently runs as a River background worker on the backend, which downloads the original image from R2, decodes it with Go's imaging library, and re-uploads the result. On the production machine this causes context deadline errors and bottlenecks under load — the 30-second decode timeout is regularly exceeded and retries compound the problem. Moving thumbnail generation to the frontend eliminates the worker entirely: the browser generates the thumbnail before the PUT, uploads it in parallel with the original, and has the result in hand before the upload is even acknowledged.

## What Changes

- **BREAKING** `POST /files` (initiate) now returns a `thumbnail_upload_url` presigned PUT URL for image-type files; clients must upload the generated thumbnail to this URL before calling complete
- `POST /files/:id/complete` no longer enqueues a River job; it sets `thumbnail_r2_path` deterministically at completion time (trusting the client uploaded it)
- `thumbnail_gen_state` column is dropped from the `files` table; the field is removed from all API responses
- `GenerateFileThumbnailWorker` and its River job registration are removed
- Frontend upload flow (both `AddFilesModal` and `ArtpieceFilesInput`) gains a thumbnail generation step using OffscreenCanvas before the PUT: alpha-capable source formats (PNG, GIF, WebP) produce a PNG thumbnail; all others produce JPEG
- The FE can display the thumbnail from a local object URL immediately after upload, with no polling wait
- The thumbnail-pending spinner is removed from `FileThumb`; when no `thumbnail_url` is present the component shows the file type fallback icon

## Capabilities

### New Capabilities

_(none)_

### Modified Capabilities

- `file-upload`: initiate returns `thumbnail_upload_url` for image types; complete sets thumbnail path deterministically and no longer enqueues a worker job; thumbnail generation requirement replaced by client-side generation contract
- `file-thumbnail-gen-state`: all requirements retired; `thumbnail_gen_state` column and field are removed
- `file-crud`: `thumbnail_gen_state` field requirement removed from all file response contracts

## Impact

**Backend**
- `backend/internal/usecase/file_upload_usecase.go` — `InitiateUpload` generates presigned thumbnail PUT URL for image types; `CompleteUpload` sets `ThumbnailR2Path` deterministically, removes job enqueue and gen_state updates
- `backend/internal/handler/file_upload_handler.go` — initiate response gains `thumbnail_upload_url`
- `backend/internal/handler/file_handler.go` — `toFileResponse` removes `thumbnail_gen_state`
- `backend/internal/domain/file.go` — `ThumbnailGenState` field removed
- `backend/internal/repository/file_repository.go` — `UpdateThumbnailGenState` method removed
- `backend/internal/worker/generate_file_thumbnail.go` — deleted
- `backend/internal/worker/generate_file_thumbnail_test.go` — deleted
- DB migration: `ALTER TABLE files DROP COLUMN thumbnail_gen_state`

**Frontend**
- New `frontend/src/features/files/lib/thumbnail.ts` — `generateThumbnail(file: File)` using OffscreenCanvas, max 600px, alpha-preserving format selection
- `frontend/src/features/files/api/useInitFileUpload.ts` — `InitFileUploadResult` gains `thumbnail_upload_url: string | null`
- `frontend/src/features/files/components/AddFilesModal.tsx` — upload orchestration gains thumbnail generation step; `DraftFile` gains `localThumbnailUrl`, drops `thumbnailGenState`; `useEffect` thumbnail sync removed
- `frontend/src/features/files/components/ArtpieceFilesInput.tsx` — same changes as `AddFilesModal`
- `frontend/src/features/files/components/FileThumb.tsx` — spinner/pending branch removed; `thumbnailGenState` prop removed
- `frontend/src/features/files/api/useFiles.ts` — `thumbnail_gen_state` removed from `File` type
- `frontend/src/features/files/components/FileTileEditOverlay.tsx` — `thumbnailGenState` prop pass-through removed
