## Why

Files have no human-readable label beyond their notes field, making the inbox harder to navigate. Adding a dedicated `name` field — auto-populated from the original filename on drag-and-drop upload — gives each file an identity without requiring manual annotation.

## What Changes

- Add `name TEXT` (nullable) column to both `files` and `pending_file_uploads` tables
- Upload flow passes the original browser filename as `name` at initiate time, carrying it through `pending_file_uploads` into the created `file` record
- Rename `UpdateNotes` to a combined `Update(name, notes)` at the repository, usecase, and handler interface layers — one atomic PATCH instead of two sequential calls
- File API response includes `name`
- Inbox search extends to match on `name` in addition to `notes`

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `file-crud`: new `name` field on the file resource; `UpdateNotes` becomes `Update` accepting both `name` and `notes`
- `file-upload`: initiate request accepts `name`; name is stored in `pending_file_uploads` and copied to `file` on completion
- `web-file-inbox`: drag-and-drop sets `name` from the browser `File.name`; search bar matches against both `name` and `notes`

## Impact

**Backend**
- New migration: `ALTER TABLE files ADD COLUMN name TEXT` and `ALTER TABLE pending_file_uploads ADD COLUMN name TEXT`
- `domain.File` and `domain.PendingFileUpload`: add `Name *string`
- `usecase/file_repository.go` (`FileRepository` interface): rename `UpdateNotes` → `Update(name, notes *string)`
- `repository/file_repository.go`: rename + expand `UpdateNotes` impl
- `usecase/file_upload_usecase.go`: `InitiateFileUploadParams` and pending record creation get `Name *string`
- `usecase/file_usecase.go`: rename `UpdateNotes` → `Update`
- `handler/file_upload_handler.go`: `initiateFileUploadRequest` gets `Name *string`; `fileResponse` and `toFileResponse` gain `Name`
- `handler/file_handler.go`: `FileUsecase` interface and `updateFileRequest` updated; handler calls `Update` instead of `UpdateNotes`
- Tests: all `UpdateNotes` call sites and fakes updated

**Frontend**
- `useFiles.ts` `File` interface: add `name: string | null`
- `useInitFileUpload.ts`: param changes from `mimeType: string` to `{ mimeType: string; name?: string }`
- `FileInboxGrid.tsx`: pass `file.name` when calling `initUpload.mutateAsync`
- `FilesPage.tsx`: search filter extends to `f.name`
