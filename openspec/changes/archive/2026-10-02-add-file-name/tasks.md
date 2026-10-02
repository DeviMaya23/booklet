## 1. Database Migration

- [x] 1.1 Create `000021_add_name_to_files.up.sql` — `ALTER TABLE files ADD COLUMN name TEXT` and `ALTER TABLE pending_file_uploads ADD COLUMN name TEXT`
- [x] 1.2 Create `000021_add_name_to_files.down.sql` — `ALTER TABLE files DROP COLUMN name` and `ALTER TABLE pending_file_uploads DROP COLUMN name`

## 2. Domain

- [x] 2.1 Add `Name *string` to `domain.File` with gorm tag `column:name`
- [x] 2.2 Add `Name *string` to `domain.PendingFileUpload` with gorm tag `column:name`

## 3. Repository

- [x] 3.1 Rename `UpdateNotes` → `Update` in `repository/file_repository.go`; update the SQL to set both `name` and `notes` columns and preserve the `RowsAffected == 0 → ErrRecordNotFound` guard
- [x] 3.2 Update `usecase/file_repository.go` `FileRepository` interface: rename `UpdateNotes` → `Update(ctx, id, userID, name, notes *string) error`
- [x] 3.3 Update `repository/file_repository_integration_test.go`: rename `TestFileRepository_UpdateNotes` tests to cover `Update`, add a scenario asserting `name` is persisted

## 4. File Upload Usecase

- [x] 4.1 Add `Name *string` to `usecase.InitiateFileUploadParams`
- [x] 4.2 In `FileUploadUsecase.InitiateUpload`: pass `params.Name` when constructing `domain.PendingFileUpload`
- [x] 4.3 In `FileUploadUsecase.CompleteUpload`: copy `pending.Name` to the constructed `domain.File`

## 5. File Usecase

- [x] 5.1 Rename `FileUsecase.UpdateNotes` → `Update`; update signature to `(ctx, id, userID uuid.UUID, name, notes *string) (*domain.File, error)`; call `fileRepo.Update` with both fields
- [x] 5.2 Update `usecase/file_usecase_test.go`: update `fakeFileRepository.UpdateNotes` → `Update`; rename `TestFileUsecase_UpdateNotes_NotOwned` test to `TestFileUsecase_Update_NotOwned`

## 6. Handlers

- [x] 6.1 In `handler/file_upload_handler.go`: add `Name *string json:"name"` to `initiateFileUploadRequest`; pass `Name: req.Name` in `InitiateFileUploadParams`; add `Name *string json:"name"` to `fileResponse`; set `Name: f.Name` in `toFileResponse`
- [x] 6.2 In `handler/file_handler.go`: update `FileUsecase` interface — rename `UpdateNotes` → `Update(ctx, id, userID, name, notes *string) (*domain.File, error)`; add `Name *string json:"name"` to `updateFileRequest`; call `h.fileUsecase.Update(ctx, id, userID, req.Name, req.Notes)` in `UpdateFile`
- [x] 6.3 Update `handler/file_handler_test.go`: rename `spyFileUsecase.UpdateNotes` → `Update` with the new signature
- [x] 6.4 Update `handler/file_upload_handler_test.go` if the spy needs changes for the new `Name` field path

## 7. Bruno Collection

- [x] 7.1 Update `collection/files/initiate_upload.bru` — add `name` to the request body example
- [x] 7.2 Update `collection/files/update_file.bru` — add `name` to the request body example

## 8. Frontend

- [x] 8.1 Add `name: string | null` to the `File` interface in `useFiles.ts`
- [x] 8.2 Change `useInitFileUpload.ts` mutation param from `mimeType: string` to `{ mimeType: string; name?: string }`; include `name` in the POST body
- [x] 8.3 Update `FileInboxGrid.tsx` call site: replace `initUpload.mutateAsync(file.type || 'application/octet-stream')` with `initUpload.mutateAsync({ mimeType: file.type || 'application/octet-stream', name: file.name })`
- [x] 8.4 Update `FilesPage.tsx` search filter to match on `f.name` in addition to `f.notes`

## 9. Quality

- [x] 9.1 Run `golangci-lint run ./...` from `backend/` and fix any issues
- [x] 9.2 Run `npm run build && npm run lint` from `frontend/` and fix any issues
