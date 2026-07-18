## Context

The v1 upload flow has a `uploadRepository.CompleteUpload` method that inserts into the `images` table, associates characters via `image_characters`, and deletes the `pending_upload` row — all in one repo method. This works but breaks the principle that each repository owns only its own table. Additionally, a redundant `UploadCharacterRepository` interface duplicates `characterRepository.GetByIDsAndUserID`, character IDs flow through the system as `[]string` despite being UUID-typed everywhere else, and repository method names deviate from the short naming convention used by `imageRepository`.

## Goals / Non-Goals

**Goals:**
- Each repository touches only its own table; cross-table coordination lives in the usecase
- A `Transactor` pattern handles the atomic delete-pending + insert-image operation without leaking `*gorm.DB` across layers
- `PendingUpload.CharacterIDs` is `[]uuid.UUID`; `characterRepository.GetByIDsAndUserID` takes `[]uuid.UUID`
- Upload repository methods named `Create`, `GetByID`, `Delete`
- `characterRepository` injected directly into the upload usecase; no separate interface wrapper
- Telemetry (counter + structured log) on `CompleteUpload`

**Non-Goals:**
- Any change to the HTTP API contract (endpoints, request/response shapes)
- Any database migration
- Stale pending_upload cleanup
- Thumbnail generation

## Decisions

### Context-based Transactor

The `CompleteUpload` operation must atomically delete a `pending_upload` row and insert an `images` row. Giving that responsibility to either repository creates cross-table coupling. Instead:

- `usecase.Transactor` interface: `InTransaction(ctx, fn func(context.Context) error) error`
- `repository.GormTransactor` implements it: wraps `*gorm.DB.Transaction`, stores the `*gorm.DB` tx in context under a package-private key (`txKey{}`)
- Every repo method calls `dbFromContext(ctx, r.db)` (package-private to `repository`) instead of `r.db.WithContext(ctx)`; if a tx is in context they use it, otherwise the live DB

The usecase calls `transactor.InTransaction(ctx, func(txCtx) { imageRepo.Create(txCtx, img); uploadRepo.Delete(txCtx, id) })`. Neither repo knows a transaction is in progress — they just use whichever DB handle the context provides.

Alternative considered: callback with typed repo arguments (`Transaction(fn func(pendingRepo, imageRepo))`). Rejected because the transaction ownership is arbitrary (why does uploadRepo own it?), and the callback signature becomes a maintenance surface.

### characterRepository injected directly

The v1 `UploadCharacterRepository` interface with `GetByIDsAndUserID` is a one-method wrapper around `characterRepository`. The usecase defines a narrow interface for everything it depends on, and `characterRepository` already satisfies the narrow slice needed. Removing the wrapper reduces indirection without breaking testability — the spy in tests implements the same method directly.

### `[]uuid.UUID` for CharacterIDs

`PendingUpload.CharacterIDs` changes from `[]string` to `[]uuid.UUID`. GORM's JSON serializer serializes `[]uuid.UUID` to a JSON string array in JSONB — identical bytes on disk, no migration needed. The parse happens once at the handler boundary (already validated as `uuid4`). Downstream code receives typed values and never re-parses.

`characterRepository.GetByIDsAndUserID` signature changes from `[]string` to `[]uuid.UUID`. GORM passes UUID values to the `IN ?` clause correctly.

### Telemetry placement

The upload counter and log line fire after `uploadRepo.GetByID` succeeds (confirming R2 upload actually happened) and before the DB transaction. This is the correct semantic: we're recording that the object is in R2, not that the DB commit succeeded. The DB commit failure path is still recorded via span error.

## Risks / Trade-offs

- **`dbFromContext` is implicit** — A repo method called outside a transaction silently uses the live DB, which is correct. But a repo method called inside a transaction that forgets to propagate `txCtx` would silently use the live DB too, breaking atomicity. Convention: always pass `ctx` through, never capture it.
- **Character ID JSONB migration** — Existing `pending_upload` rows store character IDs as `["uuid-string"]`. GORM deserializes these correctly into `[]uuid.UUID` since `uuid.UUID` unmarshals from a JSON string. No data migration needed, but this relies on GORM's UUID JSON unmarshaling behavior.
