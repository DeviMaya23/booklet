## Why

Files uploaded without an artpiece assignment land in a kind of limbo — they exist in the DB but there is no UI surface to see, manage, or sort them. This change adds the "Inbox" view: a grid of all unassigned files where the user can upload, search, select, and delete files before eventually wiring them to artpieces.

## What Changes

- **New Inbox UI** at `/app/files` — grid of unassigned files with drag-to-upload (multi-file), click/ctrl/shift/right-click selection, per-tile and bulk delete with confirm step, notes-based search, and shimmer placeholder tiles while uploads are in flight
- **New `DELETE /files` endpoint** — accepts `{ "ids": ["uuid", ...] }`, verifies ownership, deletes DB rows in bulk, enqueues a single async job to purge all associated R2 objects (file path + thumbnail path per file)
- **`FileUsecase`** gains a `jobInserter` dependency to support the async R2 cleanup
- **`FileRepository` interface** gains `GetByIDsAndUserID` (already in the repo layer) and `BulkDelete`
- **`PurgeUserStorageWorker` renamed to `PurgeR2ObjectsWorker`** (`PurgeUserStorageArgs` → `PurgeR2ObjectsArgs`, `Kind()` → `"purge_r2_objects"`) — same behavior, name no longer implies it is exclusive to user-deletion flows; call sites in `user_usecase.go` and `main.go` updated

## Capabilities

### New Capabilities

- `web-file-inbox`: Inbox UI — grid of unassigned files, drag-to-upload (multiple files in parallel), selection model (click/ctrl/shift/right-click), single and bulk delete with confirm, notes search, placeholder tiles during upload, toolbar with no-op "Add to artpiece" actions
- `file-bulk-delete`: `DELETE /files` endpoint — accepts a list of file IDs, verifies ownership of all, deletes DB rows in bulk, and enqueues a single async job to purge all associated R2 objects (file path + thumbnail path per file)

### Modified Capabilities

<!-- none -->

## Impact

### Backend
- Modified: `internal/usecase/file_usecase.go`, `internal/usecase/file_repository.go`, `internal/repository/file_repository.go`, `internal/worker/purge_user_storage.go` (rename → `purge_r2_objects.go`), `internal/usecase/user_usecase.go`, `cmd/server/main.go`
- New: `internal/handler/file_handler.go` (bulk delete handler method), Bruno file for `DELETE /files`

### Frontend
- New: `src/pages/FilesPage.tsx`
- New: `src/features/files/api/useFiles.ts`, `useInitFileUpload.ts`, `useCompleteFileUpload.ts`, `useDeleteFile.ts`, `useBulkDeleteFiles.ts`
- New: `src/features/files/components/FileInboxGrid.tsx`, `FileTile.tsx`, `FileContextMenu.tsx`, `FileDeleteDialog.tsx`
- Modified: `src/App.tsx`, `src/components/AppSidebar.tsx`
