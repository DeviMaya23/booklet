# Tasks

## 1. Database migration

- [x] 1.1 Create `backend/migration/000031_drop_thumbnail_gen_state_from_files.up.sql` (DROP COLUMN thumbnail_gen_state from files) and `000031_drop_thumbnail_gen_state_from_files.down.sql` (ADD COLUMN thumbnail_gen_state TEXT DEFAULT 'done'); verify migration runs clean with `migrate up` and rolls back cleanly with `migrate down`

## 2. Backend domain and repository cleanup

- [x] 2.1 Remove `ThumbnailGenState *string` from `domain.File` and verify the project builds with `go build ./...`
- [x] 2.2 Remove `UpdateThumbnailGenState` from `file_repository.go` and its interface declaration in `worker/generate_file_thumbnail.go` (the `fileThumbnailFileRepository` interface); verify `go build ./...` passes
- [x] 2.3 Remove `ThumbnailGenState` field and the `thumbnailGenState` / null-guard logic from `fileResponse` and `toFileResponse` in `file_upload_handler.go`; verify `go build ./...` passes and unit tests for `toFileResponse` pass

## 3. Backend upload usecase — initiate

- [x] 3.1 Add `ThumbnailUploadURL *string` and `ThumbnailR2Key *string` to `InitiateFileUploadResult` in `file_upload_usecase.go`
- [x] 3.2 In `InitiateUpload`, for image mime types compute `thumbnailKey = users/{userID}/thumbnails/{fileID}.<ext>` using the existing `formatForMimeType` logic (move it to a shared package or inline it), generate a presigned PUT URL for that key (same TTL as the original), and populate `ThumbnailUploadURL` and `ThumbnailR2Key` on the result; non-image types leave both as nil; verify the unit tests for `InitiateUpload` cover both branches
- [x] 3.3 Update `initiateFileUploadResponse` in `file_upload_handler.go` to include `ThumbnailUploadURL *string \`json:"thumbnail_upload_url"\`` and return it from `InitiateUpload`; verify the existing handler unit test (or add one) asserts the field is non-null for image types and null for non-image types

## 4. Backend upload usecase — complete

- [x] 4.1 Add optional `Width *int` and `Height *int` to a new `CompleteFileUploadParams` struct and update `CompleteUpload` signature to accept it (or add them as optional fields on an existing params type); update the handler to bind an optional JSON body `{ "width": int, "height": int }` and pass them through
- [x] 4.2 In `CompleteUpload`, for image mime types set `ThumbnailR2Path` on the new `domain.File` using the deterministic key (same derivation as initiate: `users/{userID}/thumbnails/{fileID}.<ext>`); remove the `thumbnail_gen_state` assignment; update `image_metadata` creation to use the provided width/height (defaulting to 0 when not provided)
- [x] 4.3 Remove the `u.jobInserter.Insert(...)` call and `JobInserter` dependency from `FileUploadUsecase`; verify `go build ./...` passes and usecase unit tests cover the image-type completion scenario (thumbnail path set, no job enqueued) and the non-image-type scenario (no thumbnail path, no image_metadata)

## 5. Remove the worker

- [x] 5.1 Delete `backend/internal/worker/generate_file_thumbnail.go` and `backend/internal/worker/generate_file_thumbnail_test.go`
- [x] 5.2 Remove `river.AddWorker(workers, worker.NewGenerateFileThumbnailWorker(...))` from `backend/cmd/server/main.go`; verify `go build ./...` passes

## 6. Update bruno collection

- [x] 6.1 Update the Bruno file for `POST /files` (initiate) to document the `thumbnail_upload_url` field in the response; update `POST /files/:id/complete` to include an optional body with `width` and `height` fields; verify the requests run successfully against a local server

## 7. Backend lint

- [x] 7.1 Run `golangci-lint run ./...` from `backend/` and fix all reported issues; verify it exits 0

## 8. Frontend — thumbnail generation lib

- [x] 8.1 Create `frontend/src/features/files/lib/thumbnail.ts` exporting `generateThumbnail(file: File): Promise<{ blob: Blob; width: number; height: number }>` using OffscreenCanvas: decode with `createImageBitmap`, scale to max 600px preserving aspect ratio, encode as PNG for `image/png`, `image/gif`, `image/webp` sources and as JPEG (quality 0.9) for all others; verify with a manual test that a PNG file produces a PNG blob and a JPEG file produces a JPEG blob

## 9. Frontend — API layer

- [x] 9.1 Update `useInitFileUpload.ts`: add `thumbnail_upload_url: string | null` to `InitFileUploadResult`
- [x] 9.2 Update `useFiles.ts`: remove `thumbnail_gen_state` from the `File` type (and exported `ThumbnailGenState` type if it exists); verify TypeScript compilation passes with `npm run build`
- [x] 9.3 Update `useCompleteFileUpload.ts`: accept optional `{ width?: number; height?: number }` payload and send it as a JSON body in the POST

## 10. Frontend — upload components

- [x] 10.1 Update `AddFilesModal.tsx`: in `uploadFiles`, after `initUpload` call (1) call `generateThumbnail(file)` for image mime types to get `{ blob, width, height }`, (2) `Promise.all` the original PUT and the thumbnail PUT (to `thumbnail_upload_url`), (3) call `completeUpload` with `{ width, height }`, (4) set `localThumbnailUrl: URL.createObjectURL(thumbnailBlob)` on the draft; update `DraftFile` 'uploaded' shape to replace `thumbnailGenState` with `localThumbnailUrl: string | null`; remove the `useEffect` that syncs `thumbnail_gen_state` from `useFiles`; revoke object URLs on component unmount; verify TypeScript compilation passes
- [x] 10.2 Update `ArtpieceFilesInput.tsx` with the same changes as 10.1; verify TypeScript compilation passes
- [x] 10.3 Update `FileThumb.tsx`: remove the `thumbnailGenState` prop and the `thumbnailGenState === 'pending'` spinner branch; component now renders `<img>` when `thumbnailUrl` is present and the file-category fallback icon otherwise; verify TypeScript compilation passes
- [x] 10.4 Update `FileTileEditOverlay.tsx`: remove the `thumbnailGenState` prop pass-through to `FileThumb`; verify TypeScript compilation passes

## 11. Frontend — build and lint

- [x] 11.1 Run `npm run build` from `frontend/` and fix all type errors and build warnings; verify it exits 0
- [x] 11.2 Run `npm run lint` from `frontend/` and fix all reported issues; verify it exits 0

## Workflow follow-up

- Archive the change once implementation is reviewed and merged.
