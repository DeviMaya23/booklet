## Why

The v1 upload flow has four structural issues: the upload repository inserts into the `images` table (cross-table coupling at the repo layer), character IDs are typed as `[]string` instead of `[]uuid.UUID` (inconsistent with every other entity ID in the codebase), repository methods use long names (`CreatePendingUpload`, `GetPendingUpload`) instead of the short naming convention already established by `imageRepository`, and a separate `charLookup` interface duplicates the existing `characterRepository`. These are all fixable without changing the external API contract.

## What Changes

- Replace the monolithic `uploadRepository.CompleteUpload` (which directly inserted images) with a `Transactor` that coordinates `uploadRepository` and `imageRepository` independently inside a single transaction
- Add `usecase.Transactor` interface and `repository.GormTransactor` implementation; repos detect a transaction in context via a package-private key
- Change `PendingUpload.CharacterIDs` from `[]string` to `[]uuid.UUID`; GORM JSON serializer handles the JSONB round-trip transparently
- Change `uploadRepository` method names to `Create`, `GetByID`, `Delete` (matching `imageRepository`)
- Rename `charLookup` field to `characterRepo` in the upload usecase struct; `UploadCharacterRepository` interface stays in the usecase package per convention, updated to take `[]uuid.UUID`
- Remove the now-redundant `GetPendingUpload` call inside `uploadRepository.CompleteUpload` (the double-fetch)
- Add telemetry to `CompleteUpload` usecase: upload counter metric and structured log on completion

## Capabilities

### New Capabilities

- `transactor`: Context-based transaction coordinator used by the upload usecase

### Modified Capabilities

- `image-upload`: `PendingUpload.CharacterIDs` type changes from `[]string` to `[]uuid.UUID`; `GetByIDsAndUserID` signature changes to accept `[]uuid.UUID`; no external API changes

## Impact

- `internal/usecase/upload_repository.go` — `UploadRepository` interface rewritten; `UploadCharacterRepository` updated to `[]uuid.UUID`
- `internal/usecase/upload_usecase.go` — `charLookup` replaced with `characterRepository`; `CompleteUpload` rewritten to use transactor
- `internal/repository/upload_repository.go` — `CompleteUpload` removed; `Create`/`GetByID`/`Delete` added; `dbFromContext` helper added
- `internal/repository/transactor.go` — new file; `GormTransactor` implementation and context key
- `internal/domain/pending_upload.go` — `CharacterIDs []string` → `[]uuid.UUID`
- `internal/repository/character_repository.go` — `GetByIDsAndUserID` signature: `[]string` → `[]uuid.UUID`
- `internal/repository/image_repository.go` — add `Create` method (with `dbFromContext` support)
- `cmd/server/main.go` — inject `imageRepository`, `characterRepository`, and `GormTransactor` into `NewUploadUsecase`
- No changes to HTTP endpoints, request/response shapes, or database schema
