## 1. Database Migration

- [x] 1.1 Create migration `000022_add_thumbnail_gen_state_to_files` adding `thumbnail_gen_state TEXT` column (nullable, no default) to the `files` table

## 2. Domain & Repository

- [x] 2.1 Add `ThumbnailGenState *string` field to `domain.File`
- [x] 2.2 Add `UpdateThumbnailGenState(ctx, id, state)` method to `fileRepository` and its interface
- [x] 2.3 Add integration test for `UpdateThumbnailGenState`

## 3. Backend — File Upload Usecase

- [x] 3.1 Update `CompleteFileUpload` to set `ThumbnailGenState = "pending"` on image files and `"not_applicable"` on non-image files at file row creation

## 4. Backend — Worker

- [x] 4.1 Extend `fileThumbnailFileRepository` interface with `UpdateThumbnailGenState`
- [x] 4.2 Update `GenerateFileThumbnailWorker.Work` to set `thumbnail_gen_state = "done"` alongside `thumbnail_r2_path` on success
- [x] 4.3 Update `GenerateFileThumbnailWorker.Work` to set `thumbnail_gen_state = "failed"` and return `nil` when `job.Attempt >= job.MaxAttempts` on error
- [x] 4.4 Write unit tests for the worker: success path sets state to `"done"`; final attempt failure sets state to `"failed"` and returns nil; non-final failure returns error and leaves state unchanged

## 5. Backend — API Response

- [x] 5.1 Add `ThumbnailGenState string` to `fileResponse` struct in `file_upload_handler.go`
- [x] 5.2 Update `toFileResponse` to populate `thumbnail_gen_state`: use the file's value if non-null, otherwise `"failed"`
- [x] 5.3 Update bruno collection files (`get_file.bru`, `list_files.bru`, `list_files_unassigned.bru`, `complete_upload.bru`, `update_file.bru`) to include `thumbnail_gen_state` in response examples

## 6. Frontend — API Type & Polling

- [x] 6.1 Add `thumbnail_gen_state: "pending" | "done" | "failed" | "not_applicable"` to the `File` interface in `useFiles.ts`
- [x] 6.2 Update `refetchInterval` condition in `useFiles` to poll only while any file has `thumbnail_gen_state === "pending"`

## 7. Frontend — FileTile

- [x] 7.1 Update `FileTile` to show a fallback icon when `thumbnail_url === null && thumbnail_gen_state !== "pending"` instead of the shimmer skeleton

## 8. Frontend — ArtpieceFilesInput

- [x] 8.1 Update `DraftFile` uploaded type: replace `thumbnailUrl: string | null` with `thumbnailUrl: string | null` and `thumbnailGenState: "pending" | "done" | "failed" | "not_applicable"`
- [x] 8.2 Update the `useFiles` sync effect to also sync `thumbnail_gen_state` into draft state
- [x] 8.3 Update `thumbnailPending` derivation to use `thumbnailGenState === "pending"` instead of `thumbnailUrl === null`
- [x] 8.4 Update the thumbnail slot: show fallback icon when `thumbnailUrl === null && thumbnailGenState !== "pending"` (currently only renders spinner or image)
- [x] 8.5 Update name/notes fields disabled condition to use `thumbnailGenState === "pending"`

## 9. Lint & Build

- [x] 9.1 Run `golangci-lint run` and fix any issues
- [x] 9.2 Run `npm run build` and `npm run lint` in `frontend/` and fix any issues
