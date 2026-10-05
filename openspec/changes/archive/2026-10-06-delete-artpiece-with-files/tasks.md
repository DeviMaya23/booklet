## 1. Repository Layer

- [x] 1.1 Add `BulkDelete(ctx context.Context, ids []uuid.UUID, userID uuid.UUID) error` to `ArtpieceFileRepository` interface in `artpiece_repository.go`
- [x] 1.2 Implement `BulkDelete` on the SQL file repository in `repository/file_repository.go`

## 2. Usecase Layer

- [x] 2.1 Add `JobInserter` as a dependency to `ArtpieceUsecase` (struct field + constructor param in `artpiece_usecase.go`)
- [x] 2.2 Update `ArtpieceUsecase.Delete` signature to accept `deleteFiles bool`
- [x] 2.3 Implement `deleteFiles=true` path: fetch files, run transaction (`fileRepo.BulkDelete` + `artpieceRepo.Delete`), enqueue `PurgeR2ObjectsArgs` job after commit
- [x] 2.4 Handle empty file set (skip bulk delete + enqueue, go straight to artpiece delete)

## 3. Handler Layer

- [x] 3.1 Update `ArtpieceUsecase` interface in `artpiece_handler.go` — `Delete` signature gains `deleteFiles bool`
- [x] 3.2 Update `DeleteArtpiece` handler to bind `delete_files` query param and pass it to the usecase

## 4. Wiring

- [x] 4.1 Update `NewArtpieceUsecase` call in `main.go` to pass `enqueuer` as the new `JobInserter` argument

## 5. Unit Tests

- [x] 5.1 `ArtpieceUsecase.Delete` — `delete_files=false` (or absent): artpiece deleted, no file deletion, no job enqueue
- [x] 5.2 `ArtpieceUsecase.Delete` — `delete_files=true`, artpiece has files: transaction runs `BulkDelete` + artpiece delete, `PurgeR2ObjectsArgs` job enqueued with correct R2 keys
- [x] 5.3 `ArtpieceUsecase.Delete` — `delete_files=true`, artpiece has no files: artpiece deleted, no `BulkDelete` called, no job enqueued
- [x] 5.4 `ArtpieceUsecase.Delete` — `delete_files=true`, job enqueue fails: returns nil (204), error logged

## 6. Bruno

- [x] 6.1 Add Bruno request file for `DELETE /artpieces/:id?delete_files=true`

## 7. Lint

- [x] 7.1 Run `golangci-lint run ./...` and fix any issues
