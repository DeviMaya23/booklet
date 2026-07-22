## 1. Domain

- [x] 1.1 In `internal/domain/pending_upload.go`: change `CharacterIDs` from `[]string` to `[]uuid.UUID` (keep gorm tag `type:jsonb;serializer:json;column:character_ids`)

## 2. characterRepository

- [x] 2.1 In `internal/repository/character_repository.go`: change `GetByIDsAndUserID` signature from `ids []string` to `ids []uuid.UUID`; update the GORM `IN ?` clause accordingly
- [x] 2.2 In the usecase character repository interface (currently in `usecase/upload_repository.go`): update `GetByIDsAndUserID` to take `[]uuid.UUID`

## 3. Transactor Infrastructure

- [x] 3.1 Create `internal/usecase/transactor.go` — define `Transactor` interface: `InTransaction(ctx context.Context, fn func(ctx context.Context) error) error`
- [x] 3.2 Create `internal/repository/transactor.go` — implement `GormTransactor` struct with `db *gorm.DB`; implement `InTransaction` by wrapping `db.WithContext(ctx).Transaction`, storing the `*gorm.DB` tx in context under a package-private `txKey{}`; add package-private `dbFromContext(ctx context.Context, fallback *gorm.DB) *gorm.DB` helper used by all repos

## 4. Repository Updates

- [x] 4.1 Rewrite `internal/repository/upload_repository.go`:
  - Rename `CreatePendingUpload` → `Create`; return `(*domain.PendingUpload, error)`; use `dbFromContext`
  - Rename `GetPendingUpload` → `GetByID`; change `pendingID string` to `id uuid.UUID`; use `dbFromContext`
  - Rename existing `Delete`; change param to `id uuid.UUID`; use `dbFromContext`
  - Remove `CompleteUpload`
- [x] 4.2 In `internal/repository/image_repository.go`: add `Create(ctx context.Context, image *domain.Image) (*domain.Image, error)` method using `dbFromContext(ctx, r.db)` so it participates in transactions

## 5. Upload Usecase Interfaces

- [x] 5.1 Rewrite `internal/usecase/upload_repository.go`:
  - `UploadRepository` interface: `Create(ctx, *domain.PendingUpload) (*domain.PendingUpload, error)`, `GetByID(ctx, uuid.UUID, userID string) (*domain.PendingUpload, error)`, `Delete(ctx, uuid.UUID) error`
  - `UploadCharacterRepository` interface: `GetByIDsAndUserID(ctx, []uuid.UUID, userID string) ([]domain.Character, error)`
  - `UploadImageRepository` interface: `Create(ctx, *domain.Image) (*domain.Image, error)`
  - Keep `StorageService` interface here

## 6. Upload Usecase

- [x] 6.1 Rewrite `internal/usecase/upload_usecase.go`:
  - Replace `charLookup UploadCharacterRepository` with `characterRepo UploadCharacterRepository`
  - Add `imageRepo UploadImageRepository` and `transactor Transactor` fields
  - Add `uploadCount metric.Int64Counter` field; initialize via `tel.Meter` in constructor
  - Update `InitialUpload` to call `uploadRepo.Create` (new signature)
  - Rewrite `CompleteUpload`: call `uploadRepo.GetByID` → filter chars via `characterRepo.GetByIDsAndUserID` → emit counter + log → call `transactor.InTransaction` which calls `imageRepo.Create` then `uploadRepo.Delete`; assert characters via `image.Characters` association after create
  - Update `NewUploadUsecase` constructor to accept `characterRepo`, `imageRepo`, `transactor`
- [x] 6.2 Rewrite `internal/usecase/upload_usecase_test.go` — update spies for new interfaces (`[]uuid.UUID` in character spy; `Transactor` spy that executes `fn` directly; `UploadRepository` spy with `Create`/`GetByID`/`Delete`); update existing scenarios; add scenario: `CompleteUpload` happy path asserts returned image is non-nil

## 7. Upload Repository Integration Tests

- [x] 7.1 Rewrite `internal/repository/upload_repository_integration_test.go`:
  - Rename `TestUploadRepository_CreatePendingUpload` → `TestUploadRepository_Create`: assert row fields
  - Replace `CompleteUpload_*` tests with: `TestUploadRepository_GetByID_Success` (returns correct row), `TestUploadRepository_GetByID_WrongUser` (returns `gorm.ErrRecordNotFound`), `TestUploadRepository_GetByID_NotFound` (returns `gorm.ErrRecordNotFound`), `TestUploadRepository_Delete` (row no longer present after call)

## 8. Wire Up

- [x] 8.1 In `cmd/server/main.go` `initApp`: create `repository.NewGormTransactor(db)`; update `usecase.NewUploadUsecase` call to pass `characterRepository`, `imageRepository`, `GormTransactor`

## 9. Lint

- [x] 9.1 Run `golangci-lint run ./...` from `backend/` and fix any reported issues
