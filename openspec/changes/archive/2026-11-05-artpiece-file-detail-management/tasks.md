# Tasks

## 1. Schema

- [x] 1.1 Add migration `000039_add_is_visible_to_files` with `up.sql` (`ALTER TABLE files ADD COLUMN is_visible boolean NOT NULL DEFAULT true;`) and `down.sql` (`ALTER TABLE files DROP COLUMN is_visible;`); verify by running the migration up and down locally against a dev database
- [x] 1.2 Add `IsVisible bool` (`gorm:"column:is_visible;not null;default:true"`) to `domain.File`; verify the column tag matches the migration exactly

## 2. Backend — cover selection and visibility logic

- [x] 2.1 In `artpiece_usecase.go`, filter candidates to `IsVisible` files in `maybeAutoSetCover`, `reassignCover`, and `setCoverFromFileSet`, per the `artpiece-files` delta scenarios (hidden files skipped on attach, detach, create-with-`file_ids`, and full-replace); verify with usecase unit tests covering: attaching a hidden image to a coverless artpiece leaves the cover unset; attaching a non-image when the only other image is hidden sets the non-image as cover; detaching the cover when the only remaining file is hidden clears the cover; full-replace with an all-hidden new set leaves the cover unset
- [x] 2.2 Add `FileUsecase.Update` support for an optional `is_visible` field (extend its signature to accept `*bool` alongside `name`/`notes`); verify with a usecase unit test asserting the returned file's `IsVisible` matches the request when provided, and is unchanged when the field is omitted
- [x] 2.3 In `FileUsecase.Update`, when `is_visible` transitions to `false` on a file that is the artpiece's current cover (via `ArtpieceFileRepository` lookup), call the artpiece-side cover-reassignment path used by detach; verify with a usecase unit test asserting the artpiece's `cover_file_id` is reassigned to another visible file (or cleared if none remain) when the current cover is hidden, and is left unchanged when a non-cover file is hidden or when visibility is set to `true`

## 3. Backend — handler and response changes

- [x] 3.1 Add `IsVisible *bool` to `updateFileRequest` in `file_handler.go` and pass it through to `FileUsecase.Update`; verify with a handler test asserting `PUT /files/:id` with only `is_visible` in the body updates visibility and leaves `name`/`notes` untouched
- [x] 3.2 Add `is_visible` to `fileResponse` and `toFileResponse` in `file_upload_handler.go`; verify with a handler test asserting `GET /files/:id`, `GET /files`, `PUT /files/:id`, and the complete-upload response all include `is_visible` matching the stored value
- [x] 3.3 Add `notes` and `is_visible` to the artpiece-detail per-file response struct (alongside the existing `name`/`extension`/`width`/`height`); verify with a handler test asserting `GET /artpieces/:id` includes `notes` and `is_visible` for each attached file, matching the `artpiece-files` delta's response scenarios
- [x] 3.4 Update the Bruno collection's `PUT /files/:id` request example to include `is_visible` in its sample body; verify the request runs against a local server and returns the updated field

## 4. Backend — repository integration tests

- [x] 4.1 Add an integration test for `ArtpieceRepository.GetFilesForArtpiece` asserting the returned files include `IsVisible` (it's consumed directly by the cover-selection filtering in task 2.1, so a stale/unselected column here would silently break that logic)
- [x] 4.2 Add an integration test for `FileRepository`'s update path asserting `IsVisible` persists through a save/reload round trip

## 5. Backend — lint

- [x] 5.1 Run `golangci-lint run` from `backend/` and fix any findings introduced by this change

## 6. Frontend — types and data

- [x] 6.1 Add `notes: string | null` and `is_visible: boolean` to the `ArtpieceFile` type in `useArtpiece.ts`; verify the type compiles against the updated backend response
- [x] 6.2 Add an optional `isVisible?: boolean` parameter to the `useUpdateFile` mutation's input, threaded to the `PUT /files/:id` body; verify with a test asserting the request body includes `is_visible` only when provided

## 7. Frontend — grid visibility filtering

- [x] 7.1 In `ArtpieceDetailFileGrid.tsx`, filter the rendered file list to `is_visible: true` files in both `view` and `edit` modes, per the `web-artpiece-detail` delta; verify with a component test asserting a hidden file never renders as a tile in either mode, and that an artpiece with only hidden files shows the empty-grid state
- [x] 7.2 In `ArtpieceDetailView.tsx`, add the non-interactive hidden-count hint to the files-section header (e.g. "6 files · 2 hidden"), shown only when at least one attached file has `is_visible: false`; verify with a component test asserting the hint is absent when no files are hidden and present with the correct count otherwise

## 8. Frontend — Edit details modal

- [x] 8.1 Create `ArtpieceFileDetailsModal.tsx` modeled on `AddFilesModal.tsx`/`FileTileEditOverlay.tsx`: one row per attached file (including hidden ones) with thumbnail, name input, notes input, and a visibility toggle, each pre-filled from the artpiece's already-loaded file list; verify with a component test asserting all attached files render as rows, including hidden ones, with correct pre-filled values
- [x] 8.2 Wire name/notes auto-save on blur via `useUpdateFile`, with per-field inline error display and retry-clears-error behavior, matching the `FileTileEditOverlay` pattern; verify with a component test covering a failed save showing the inline error and a subsequent successful save clearing it, for both name and notes independently
- [x] 8.3 Wire the visibility toggle to call `useUpdateFile` with `is_visible` immediately on toggle (no staging), reverting the toggle and showing an inline error on failure; verify with a component test asserting a successful toggle updates the row's state and a failed toggle reverts it with an error shown
- [x] 8.4 After a successful visibility toggle, invalidate/refetch the artpiece query so a cover reassignment triggered server-side (per the `artpiece-files` delta) is reflected in the modal and the underlying detail view without a separate user action; verify with a component test asserting the artpiece query is invalidated after a successful `is_visible` mutation
- [x] 8.5 Add the "Edit details" entry point to the files-section header in `ArtpieceDetailView.tsx`, available in both view and edit mode and independent of any pending staged edit-mode changes; verify with a component test asserting the button opens the modal from both modes

## 9. Frontend — build and lint

- [x] 9.1 Run `npm run build` from `frontend/` and fix any type errors introduced by this change
- [x] 9.2 Run `npm run lint` from `frontend/` and fix any findings introduced by this change

## Workflow follow-up

- Archive the change after review, per the project's OpenSpec workflow.
- Verify the archived result with `/opsx:verify`.
