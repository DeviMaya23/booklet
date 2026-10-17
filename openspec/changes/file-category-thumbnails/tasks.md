# Tasks

## 1. Backend — Migration and Domain

- [x] 1.1 Create migration `000030_add_extension_to_files` — ALTER TABLE `files` and `pending_file_uploads` to add a nullable `extension text` column; write matching `.down.sql` to drop both columns. Verify by running `make migrate` and checking both tables have the column.
- [x] 1.2 Add `Extension *string` to `domain.File` and `domain.PendingFileUpload` structs with the correct GORM `column:extension` tag. Verify the project compiles without errors.

## 2. Backend — Upload Pipeline

- [x] 2.1 Add `Extension *string` to `usecase.InitiateFileUploadParams`; update `InitiateUpload` to store `Extension` on the `PendingFileUpload` record; update `CompleteUpload` to copy `Extension` from pending to the new file row. Verify by updating the existing `file_upload_usecase_test.go` — add a scenario to `TestCompleteUpload` that asserts `extension` is copied from pending to the created file.
- [x] 2.2 Add `Extension *string` to `initiateFileUploadRequest` in `file_upload_handler.go` and pass it through to `InitiateFileUploadParams`. Add `Extension` to `fileResponse` and populate it in `toFileResponse`. Update `file_upload_handler_test.go` — add a scenario asserting `extension` is forwarded to the usecase from a request that includes it; add a scenario asserting `extension` appears in the `CompleteUpload` response.

## 3. Backend — Artpiece File Response

- [x] 3.1 Add `Extension *string` to the `fileRef` type in `artpiece_handler.go`; populate it from `f.Extension` in the mapping. Update `artpiece_handler_test.go` — add a scenario asserting the artpiece detail response includes `extension` for each attached file.

## 4. Backend — Repository Integration Tests

- [x] 4.1 Update `file_upload_repository_integration_test.go`: add a test asserting that `Create` and `GetByID` on `PendingFileUpload` round-trip the `Extension` field. Update `file_repository_integration_test.go`: add a test asserting that `Create` and `GetByIDAndUserID` on `File` round-trip the `Extension` field.

## 5. Backend — Bruno and Lint

- [x] 5.1 Update `collection/files/initiate_upload.bru` to include `extension` in the request body example. Verify the file parses correctly in Bruno.
- [x] 5.2 Run `golangci-lint run ./...` from the backend directory and fix any issues introduced by this change.

## 6. Frontend — Types and Upload Call Sites

- [x] 6.1 Add `extension: string | null` to the `File` interface in `features/files/api/useFiles.ts` and to `ArtpieceFile` in `features/artpieces/api/useArtpiece.ts`. Verify `npm run build` passes with no new type errors.
- [x] 6.2 Update `useInitFileUpload` mutation to accept and forward an optional `extension?: string` field in its input and request body.
- [x] 6.3 Update `FileInboxGrid.tsx` drag-and-drop upload: derive `extension` from the original `file.name` (lowercase of the last `.`-separated segment) and pass it to `initUpload`. Verify via `npm run build`.
- [x] 6.4 Update `ArtpieceFilesInput.tsx` upload path: derive `extension` from the original file name and pass it to `initUpload`. Verify via `npm run build`.
- [x] 6.5 Update `ArtpieceDetailFileGrid.tsx` edit-mode upload (the drop zone in the file grid): derive `extension` and pass it to `initUpload`. Verify via `npm run build`.

## 7. Frontend — fileCategory Utility

- [x] 7.1 Create `features/files/lib/fileCategory.ts` with a `getFileCategory(mimeType: string, extension: string | null)` function. The function returns `{ icon: LucideIcon, label: string | null, tinted: boolean, cssVar: string | null }`. Categories and mappings per the design (extension-first, mime prefix fallback). Verify with unit tests covering: a known proprietary extension (`"clip"` → Layered), a mime-matched type (`"video/mp4"` → Video), null extension with image mime (`"image/png"` → Image), and unknown extension with generic mime (→ Fallback, no label).

## 8. Frontend — CSS Tint Tokens

- [x] 8.1 Add `--file-cat-video-bg`, `--file-cat-video-fg`, `--file-cat-audio-bg`, `--file-cat-audio-fg`, `--file-cat-layered-bg`, `--file-cat-layered-fg`, `--file-cat-3d-bg`, `--file-cat-3d-fg` CSS custom properties to the `:root` / `[data-theme="warm"]` and `[data-theme="sunless"]` blocks in `index.css`. Add the same tokens to `[data-theme="lumen"]`. Verify by inspecting the computed styles in the browser for each theme.

## 9. Frontend — FileThumb Component

- [x] 9.1 Create `features/files/components/FileThumb.tsx`. Props: `{ mimeType: string, thumbnailUrl: string | null, extension: string | null, thumbnailGenState?: 'pending' | 'done' | 'failed' | 'not_applicable' }`. Renders: (1) `<img>` when `thumbnailUrl` is non-null; (2) `<Skeleton>` when `thumbnailGenState === 'pending'`; (3) category tile otherwise. The category tile uses `w-[28%] aspect-square` for the icon container, CSS var inline styles for tinted categories, `bg-muted text-muted-foreground` for neutral ones, and an uppercase extension label pill centered at the bottom. Verify `npm run build` passes.
- [x] 9.2 Replace the inline fallback in `FileTile.tsx` with `<FileThumb>`, passing `file.thumbnail_url`, `file.mime_type`, `file.extension`, and `file.thumbnail_gen_state`. Verify the inbox grid renders category tiles for non-image files and the skeleton for pending ones.
- [x] 9.3 Replace the inline fallback in `ArtpieceDetailFileGrid.tsx` with `<FileThumb>`, passing `file.thumbnail_url`, `file.extension`, and `file.thumbnail_gen_state`. `thumbnail_gen_state` was added to the backend `fileRef` struct and `ArtpieceFile` type post-spec so the artpiece detail grid also shows the loading spinner for pending files.
- [x] 9.4 Replace the inline fallback in `ArtpieceFilesInput.tsx` with `<FileThumb>`, passing `thumbnailUrl`, `mimeType`, `extension` (derived at upload time from the original filename and stored in component state alongside the other draft fields), and `thumbnailGenState`. Verify the uploader form shows category tiles once upload completes.
- [x] 9.5 Replace the inline fallback in `FileTileEditOverlay.tsx` with `<FileThumb>`, passing `file.thumbnail_url`, `file.mime_type`, `file.extension`, and no `thumbnailGenState`. Verify the edit overlay displays a category tile for non-image files.

## 10. Frontend — Build and Lint

- [x] 10.1 Run `npm run build` from the frontend directory and verify it succeeds with no errors.
- [x] 10.2 Run `npm run lint` from the frontend directory and fix any issues introduced by this change.

## Workflow follow-up

- Archive the change after review.
