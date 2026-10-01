## 1. Rename PurgeUserStorage → PurgeR2Objects

- [x] 1.1 Rename `internal/worker/purge_user_storage.go` to `purge_r2_objects.go`; rename `PurgeUserStorageArgs` → `PurgeR2ObjectsArgs`, `PurgeUserStorageWorker` → `PurgeR2ObjectsWorker`, `NewPurgeUserStorageWorker` → `NewPurgeR2ObjectsWorker`, and `Kind()` return value from `"purge_user_storage"` to `"purge_r2_objects"`
- [x] 1.2 Update `internal/usecase/user_usecase.go` — replace `worker.PurgeUserStorageArgs` with `worker.PurgeR2ObjectsArgs`
- [x] 1.3 Update `cmd/server/main.go` — replace `worker.NewPurgeUserStorageWorker` with `worker.NewPurgeR2ObjectsWorker`

## 2. FileRepository — Bulk Delete

- [x] 2.1 Add `GetByIDsAndUserID(ctx context.Context, ids []uuid.UUID, userID uuid.UUID) ([]*domain.File, error)` and `BulkDelete(ctx context.Context, ids []uuid.UUID, userID uuid.UUID) error` to the `FileRepository` interface in `internal/usecase/file_repository.go`
- [x] 2.2 `GetByIDsAndUserID` already exists on the concrete `fileRepository` in `internal/repository/file_repository.go` — verify the signature matches the interface; add `BulkDelete` implementation: `DELETE FROM files WHERE id IN ? AND user_id = ?`

## 3. FileUsecase — BulkDelete

- [x] 3.1 Add `jobInserter JobInserter` field to `FileUsecase` struct and update `NewFileUsecase` constructor in `internal/usecase/file_usecase.go`
- [x] 3.2 Implement `BulkDelete(ctx context.Context, ids []uuid.UUID, userID uuid.UUID) error` on `FileUsecase`: fetch files via `GetByIDsAndUserID`, reject with error if returned count doesn't match requested count (ownership violation), call `BulkDelete` on repo, collect all R2 keys (file path + thumbnail path per file where non-nil), enqueue one `worker.PurgeR2ObjectsArgs{R2Keys: keys}` job (log enqueue failure, don't return error)
- [x] 3.3 Write unit tests for `BulkDelete` usecase: (a) all IDs owned — repo delete called, job enqueued with correct keys; (b) one ID not owned — returns error, repo delete not called; (c) job enqueue fails — still returns nil (success)

## 4. FileHandler — Bulk Delete Endpoint

- [x] 4.1 Add `BulkDeleteFiles` handler method to `internal/handler/file_handler.go`: bind `{ "ids": [] }` body, validate non-empty, call usecase `BulkDelete`, return 204 on success; return 422 on ownership violation or empty IDs, 500 on other errors
- [x] 4.2 Register `protected.DELETE("/files", fileHandler.BulkDeleteFiles)` in `cmd/server/main.go`
- [x] 4.3 Write unit tests for `BulkDeleteFiles` handler: (a) happy path — 204; (b) empty ids — 422; (c) ownership violation — 422; (d) missing body / bind error — 400
- [x] 4.4 Create `collection/files/bulk_delete_files.bru` — `DELETE /files` with JSON body `{ "ids": ["{{file_id}}"] }`

## 5. FileUsecase — Wire jobInserter in main.go

- [x] 5.1 Update `cmd/server/main.go` to pass `jobInserter` when constructing `FileUsecase` (the `riverEnqueuer` instance already exists; just add it as an argument)

## 6. Frontend — API Hooks

- [x] 6.1 Create `frontend/src/features/files/api/useFiles.ts` — `GET /files?unassigned=true`; poll every 2s when any file has `thumbnail_url: null` (same `refetchInterval` pattern as `useImages`)
- [x] 6.2 Create `frontend/src/features/files/api/useInitFileUpload.ts` — `POST /files` with `{ mime_type }` only
- [x] 6.3 Create `frontend/src/features/files/api/useCompleteFileUpload.ts` — `POST /files/:id/complete`
- [x] 6.4 Create `frontend/src/features/files/api/useDeleteFile.ts` — `DELETE /files/:id`; invalidates `useFiles` query on success
- [x] 6.5 Create `frontend/src/features/files/api/useBulkDeleteFiles.ts` — `DELETE /files` with `{ ids: string[] }`; invalidates `useFiles` query on success

## 7. Frontend — File Inbox Components

- [x] 7.1 Create `frontend/src/features/files/components/FileTile.tsx` — renders thumbnail (or shimmer if null), "..." menu with "View Detail" (no-op) and "Delete" items; accepts `selected` boolean for highlight styling
- [x] 7.2 Create `frontend/src/features/files/components/FileDeleteDialog.tsx` — alert dialog confirming deletion of N files; used for both single and bulk delete
- [x] 7.3 Create `frontend/src/features/files/components/FileContextMenu.tsx` — right-click context menu with "New Artpiece" (no-op), "Add to Existing Artpiece" (no-op), divider, "Delete" (triggers delete confirm)
- [x] 7.4 Create `frontend/src/features/files/components/FileInboxGrid.tsx`:
  - Renders grid of `FileTile` components wrapping the fetched file list
  - Manages selection state (`Set<string>` of IDs, `lastClickedId` ref)
  - Handles click / ctrl+cmd-click / shift-click selection logic
  - Wraps right-click on tile with `FileContextMenu`, applying selection replacement logic (in-selection → act on full selection; out-of-selection → replace then act)
  - Handles drag-over and drop events: for each dropped file, insert a shimmer placeholder tile immediately, then run `initiate → PUT → complete` in parallel; invalidate `useFiles` query on each complete; no auto-select on completion
  - Renders persistent "drag files here to upload" hint text below the grid
  - Exposes current selection to parent via callback or renders toolbar inline

## 8. Frontend — FilesPage

- [x] 8.1 Create `frontend/src/pages/FilesPage.tsx`:
  - Fetches files via `useFiles`
  - Maintains search string state; filters file list client-side on `notes` (case-insensitive)
  - Renders search bar, `FileInboxGrid`, floating "+" button (no-op), and toolbar
  - Toolbar: shows "Add to artpiece ▾" split-button with "New Artpiece" and "Add to Existing Artpiece" dropdown items (both no-op) when selection is non-empty; hides when empty
- [x] 8.2 Add route `<Route path="/app/files" element={<FilesPage />} />` in `frontend/src/App.tsx`
- [x] 8.3 Add `{ label: 'Inbox', to: '/app/files' }` to `navItems` in `frontend/src/components/AppSidebar.tsx`

## 9. Lint & Build

- [x] 9.1 Run `golangci-lint run ./...` from `backend/` and fix any reported issues
- [x] 9.2 Run `npm run build` and `npm run lint` from `frontend/` and fix any reported issues
