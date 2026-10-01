## 1. Repository: file_repository additions

- [x] 1.1 Add `GetByIDsAndUserID(ctx, ids []uuid.UUID, userID uuid.UUID) ([]*domain.File, error)` to `fileRepository` — bulk ownership check for create-with-files and bulk replace
- [x] 1.2 Add `List(ctx, userID uuid.UUID, unassigned bool) ([]*domain.File, error)` to `fileRepository` — returns all user files; when `unassigned=true` filters to `artpiece_id IS NULL`
- [x] 1.3 Add `UpdateNotes(ctx, id uuid.UUID, userID uuid.UUID, notes *string) error` to `fileRepository`
- [x] 1.4 Add `Delete(ctx, id uuid.UUID, userID uuid.UUID) error` to `fileRepository` — hard delete, returns `gorm.ErrRecordNotFound` if not found/not owned
- [x] 1.5 Add `BulkUpdateArtpieceID(ctx, fileIDs []uuid.UUID, artpieceID *uuid.UUID) error` to `fileRepository` — batch-sets `artpiece_id` for the given IDs; used by bulk replace and create-with-files
- [x] 1.6 Add integration tests for all new `fileRepository` methods

## 2. File usecase

- [x] 2.1 Define `FileRepository` interface in `usecase/file_repository.go` covering the methods needed by `FileUsecase`: `GetByIDAndUserID`, `List`, `UpdateNotes`, `Delete`
- [x] 2.2 Define `FileStorageService` interface in the same file with `DeleteObject(ctx, key string) error`
- [x] 2.3 Implement `FileUsecase` in `usecase/file_usecase.go` with `GetByID`, `List`, `Delete`, `UpdateNotes`; `Delete` calls `storage.DeleteObject` before removing the DB row, logs R2 errors but does not fail the operation
- [x] 2.4 Add unit tests for `FileUsecase`: `GetByID` (file not owned → error), `List` (no filter, unassigned filter), `UpdateNotes` (file not owned → error), `Delete` (success with R2 cleanup, R2 error still deletes DB row)

## 3. File handler

- [x] 3.1 Implement `FileHandler` in `handler/file_handler.go` with `GetFile`, `ListFiles`, `UpdateFile`, `DeleteFile`; use `c.Bind` + `c.Validate` for query/body params; use `toFileResponse` (already package-level in handler package)
- [x] 3.2 Define the `FileUsecase` interface on `FileHandler` (get, list, update, delete)

## 4. Artpiece usecase: transactor + create-with-files + bulk replace

- [x] 4.1 Add `Transactor` dependency to `ArtpieceUsecase` and update `NewArtpieceUsecase` constructor signature
- [x] 4.2 Add `GetByIDsAndUserID` and `BulkUpdateArtpieceID` to the `ArtpieceFileRepository` interface in `usecase/artpiece_repository.go`
- [x] 4.3 Update `CreateArtpieceParams` to include `FileIDs []uuid.UUID`; update `ArtpieceUsecase.Create` to validate file ownership, attach all files in a transaction, and apply cover logic (prefer first image in set) when `FileIDs` is non-empty
- [x] 4.4 Implement `ArtpieceUsecase.ReplaceFiles(ctx, artpieceID, userID uuid.UUID, fileIDs []uuid.UUID) (*domain.Artpiece, error)`: validate artpiece ownership, validate all file IDs belong to the user, check none are attached to a different artpiece, compute attach/detach diff, run batch updates in a transaction, apply cover logic on the resulting set
- [x] 4.5 Add unit tests for `ArtpieceUsecase.Create` with file IDs: empty file_ids (no-op files path), valid file IDs (files attached + cover set), file not owned (returns error), file attached to another artpiece (returns error)
- [x] 4.6 Add unit tests for `ArtpieceUsecase.ReplaceFiles`: full replace with a new set, empty set (all detached + cover cleared), current cover stays when still in set, cover reassigned when current cover removed, file not owned (returns error), file attached to another artpiece (returns error), artpiece not found (returns error)

## 5. Artpiece handler: create request + ReplaceFiles handler

- [x] 5.1 Add `FileIDs []uuid.UUID` field to `createArtpieceRequest` and pass it through to `CreateArtpieceParams`
- [x] 5.2 Add `ReplaceFiles` handler method to `ArtpieceHandler`; accept `{ "file_ids": [...] }` body; path param is artpiece ID
- [x] 5.3 Add `ReplaceFiles` to the `ArtpieceUsecase` interface in `handler/artpiece_handler.go`

## 6. Route wiring

- [x] 6.1 Instantiate `FileUsecase` and `FileHandler` in `cmd/server/main.go`; inject `fileRepository` and `r2Storage`
- [x] 6.2 Register new file routes: `GET /files/:id`, `GET /files`, `PUT /files/:id`, `DELETE /files/:id`
- [x] 6.3 Register `PUT /artpieces/:id/files` → `artpieceHandler.ReplaceFiles`
- [x] 6.4 Wire `Transactor` into `NewArtpieceUsecase` call in `main.go`

## 7. Bruno files

- [x] 7.1 Create bruno file for `GET /files/:id`
- [x] 7.2 Create bruno file for `GET /files` (with and without `unassigned=true`)
- [x] 7.3 Create bruno file for `PUT /files/:id`
- [x] 7.4 Create bruno file for `DELETE /files/:id`
- [x] 7.5 Create bruno file for `PUT /artpieces/:id/files`

## 8. Lint

- [x] 8.1 Run `golangci-lint run ./...` and fix all reported issues
