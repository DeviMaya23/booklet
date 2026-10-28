# Tasks

## 1. BE: `pkg/filesize` package

- [x] 1.1 Create `backend/pkg/filesize/filesize.go` with `MaxBytesForExtension(ext string) int64` mapping extensions to category limits (video 2 GB, layered 800 MB, 3d 500 MB, audio 100 MB, archive 500 MB, document 25 MB, image 50 MB, fallback 100 MB) with a cross-reference comment pointing to `frontend/src/features/files/lib/fileCategory.ts FILE_SIZE_LIMITS`. Verify by running `go build ./...` with no errors.
- [x] 1.2 Write unit tests in `backend/pkg/filesize/filesize_test.go` covering: a known video extension returns 2 GB, a known layered extension returns 800 MB, an unknown extension returns the 100 MB fallback. Verify with `go test ./pkg/filesize/...`.

## 2. BE: Storage layer

- [x] 2.1 Add `contentLength int64` parameter to `GeneratePresignedPutURL` in `backend/internal/storage/r2.go` and the `StorageService` interface; pass it as `ContentLength` on `PutObjectInput`. Verify the existing `r2_test.go` tests still pass with `go test ./internal/storage/...`.
- [x] 2.2 Update the `StorageService` fake in `backend/internal/usecase/fakes_test.go` to match the updated `GeneratePresignedPutURL` signature. Verify with `go build ./...`.

## 3. BE: File upload init endpoint

- [x] 3.1 Add `FileSizeBytes int64` to `InitiateFileUploadParams` in `backend/internal/usecase/file_upload_usecase.go`. In `InitiateUpload`, call `filesize.MaxBytesForExtension` (deriving the extension from mime type via the existing `bookmime.MimeTypeToExt` or from the `Extension` param) and return `ErrFileTooLarge` (new sentinel) if `FileSizeBytes` exceeds the limit. Pass `FileSizeBytes` as `contentLength` to `GeneratePresignedPutURL`. Verify by checking the usecase test file compiles.
- [x] 3.2 Add `file_size` (`validate:"required,min=1"`) to `initiateFileUploadRequest` in `backend/internal/handler/file_upload_handler.go`; pass it as `FileSizeBytes` to the usecase; map `ErrFileTooLarge` to HTTP 422 with message `"file size exceeds the limit for this file type"`. Verify by running `go build ./...`.
- [x] 3.3 In `backend/internal/handler/file_upload_handler_test.go`, add a handler test: `POST /files` with `file_size` exceeding the limit returns 422 with the expected message. Assert the response body, not just the status code. Verify with `go test ./internal/handler/...`.
- [x] 3.4 In `backend/internal/usecase/file_upload_usecase_test.go`, add a usecase test: `InitiateUpload` with `FileSizeBytes` exceeding `MaxBytesForExtension` for the given extension returns `ErrFileTooLarge`. Assert the specific error. Verify with `go test ./internal/usecase/...`.

## 4. BE: Character avatar upload init endpoint

- [x] 4.1 Add `FileSizeBytes int64` to `InitAvatarUpload` in `backend/internal/usecase/character_usecase.go`; validate against `filesize.MaxBytesForExtension` for the image category limit; return `ErrFileTooLarge` if exceeded; pass `FileSizeBytes` as `contentLength` to `GeneratePresignedPutURL`. Verify with `go build ./...`.
- [x] 4.2 Add `file_size` (`validate:"required,min=1"`) to `initAvatarUploadRequest` in `backend/internal/handler/character_handler.go`; pass it to the usecase; map `ErrFileTooLarge` to HTTP 422 with message `"file size exceeds the limit for this file type"`. Verify with `go build ./...`.
- [x] 4.3 In `backend/internal/handler/character_handler_test.go`, add a handler test: avatar init with `file_size` exceeding the image limit returns 422. Assert the response body. Verify with `go test ./internal/handler/...`.
- [x] 4.4 In `backend/internal/usecase/character_usecase_test.go`, add a usecase test: `InitAvatarUpload` with `FileSizeBytes` exceeding the image limit returns `ErrFileTooLarge`. Verify with `go test ./internal/usecase/...`.

## 5. BE: Bruno + lint

- [x] 5.1 Update the Bruno request file for `POST /files` init to include a `file_size` field. Verify the file is valid JSON/Bruno format.
- [x] 5.2 Run `golangci-lint run ./...` from `backend/` and fix any reported issues. Verify with a clean lint run.

## 6. FE: Size constants and utility

- [x] 6.1 Add `FILE_SIZE_LIMITS` (`Record<FileCategory, number>`) to `frontend/src/features/files/lib/fileCategory.ts` with values matching the BE constants (video 2 GB, layered 800 MB, 3d 500 MB, audio 100 MB, archive 500 MB, document 25 MB, image 50 MB, fallback 100 MB) and a cross-reference comment pointing to `backend/pkg/filesize/filesize.go`. Export a helper `getFileSizeLimit(extension: string | null): number` that resolves category then limit. Verify with `npm run build` with no TS errors.
- [x] 6.2 Add `filterBySize(files: File[]): { accepted: File[]; rejected: { file: File; limitBytes: number }[] }` to `frontend/src/features/files/lib/upload.ts`. Verify with `npm run build`.
- [x] 6.3 Add tests for `filterBySize` in `frontend/src/features/files/lib/upload.test.ts`: a file within its category limit is accepted; a file exceeding its category limit is rejected with the correct `limitBytes`. Verify with `npm test -- upload.test`.

## 7. FE: `ArtpieceFilesInput` and `AddFilesModal`

- [x] 7.1 In `ArtpieceFilesInput.uploadFiles()`, call `filterBySize(list)` before queuing; for each rejected file show `toast.error` naming the file and the limit (format: `"<filename> exceeds the <X MB> limit for <category> files"`). Only accepted files proceed. Verify manually: drop an oversized file → toast appears, file not queued; drop an in-limit file → uploads normally.
- [x] 7.2 Apply the same change to `AddFilesModal.uploadFiles()`. Verify same behavior in the modal.

## 8. FE: `useInitFileUpload` and `useInitAvatarUpload`

- [x] 8.1 Add `fileSizeBytes: number` to `useInitFileUpload` mutation params and include it as `file_size` in the request body. Verify with `npm run build`.
- [x] 8.2 In `upload.ts uploadFile()`, pass `file.size` as `fileSizeBytes` to `initMutate`. Verify with `npm run build`.
- [x] 8.3 Add `fileSizeBytes: number` to `useInitAvatarUpload` mutation params and include it as `file_size` in the request body. Verify with `npm run build`.

## 9. FE: `CharacterFormModal` avatar size check

- [x] 9.1 In `CharacterFormModal.stageFileForCrop()`, check the source file's size against the image category limit before calling `URL.createObjectURL`. If exceeded, show a toast (`"Image exceeds the <X MB> limit"`) and return without opening the crop dialog. Verify manually: paste/drop an oversized image → toast shown, crop dialog does not open; paste/drop a normal image → crop dialog opens.
- [x] 9.2 After the crop dialog produces the cropped `File` in `onCrop`, pass `croppedFile.size` (not the source size) to `initAvatarUpload`. Verify with `npm run build`.

## 10. FE: Build and lint

- [x] 10.1 Run `npm run build` from `frontend/` and fix any TypeScript errors. Verify with a clean build.
- [x] 10.2 Run `npm run lint` from `frontend/` and fix any reported issues. Verify with a clean lint run.

## Workflow follow-up

- Archive the change once all tasks are complete and reviewed.
