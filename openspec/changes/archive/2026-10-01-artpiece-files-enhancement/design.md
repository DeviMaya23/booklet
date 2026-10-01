## Context

Files uploaded to the system are created via a two-step upload flow (`FileUploadUsecase`) and associated with an artpiece through single attach/detach endpoints on `ArtpieceUsecase`. Beyond those entry points, no management endpoints exist for files. The artpiece create endpoint also does not accept initial file IDs.

The codebase follows a consistent layered pattern: handler → usecase → repository, with interface-segregated repository contracts per usecase. A `GormTransactor` is already established and used by several usecases for multi-step atomic operations.

## Goals / Non-Goals

**Goals:**
- Expose CRUD endpoints for the files module (get, list with unassigned filter, update notes, delete)
- Add a bulk full-replace endpoint for an artpiece's file set, running atomically
- Allow artpiece create to accept an optional initial file array
- Reuse existing cover-assignment logic without duplicating its rules

**Non-Goals:**
- Custom/user-supplied thumbnails for non-image files
- Touching or depending on the images module or the old pending-upload infrastructure
- Pagination on the file list (not needed at this scale)
- Bulk file operations outside of the artpiece context

## Decisions

### New `FileUsecase` and `FileHandler` (not extending `FileUploadUsecase`)

File upload is a distinct flow: initiate → client uploads to R2 → complete. CRUD operations (get, list, update, delete) are independent lifecycle concerns. Merging them into `FileUploadUsecase` would blur the boundary and make both harder to test.

A new `FileUsecase` with its own `FileHandler` follows the exact same pattern as every other module and keeps each usecase focused on a single responsibility.

### `ArtpieceUsecase` gains a `Transactor` dependency

The bulk replace touches N file rows plus the artpiece's cover in one logical operation. Without a transaction, a partial failure (e.g., cover update succeeds but some detaches fail) leaves the artpiece inconsistent. The `GormTransactor` pattern is already used by user, character, and upload usecases — this is wiring an existing pattern, not introducing a new one.

The artpiece create-with-files path also benefits from a transaction: files are batch-updated only after the artpiece row is successfully created.

### Cover logic for bulk replace: compute final state once

The existing `maybeAutoSetCover` and `reassignCover` helpers implement per-file cover logic. For bulk replace, we compute the final file set after all attaches/detaches and apply a single cover decision:

1. If the current cover is still in the new set → no change.
2. If the current cover was removed → `reassignCover` over the remaining files (prefer image, else first, else nil).
3. If the artpiece had no cover → `maybeAutoSetCover` using the first image in the new set (or first file if no images).

This avoids re-running the logic N times and produces a deterministic result.

### `unassigned=true` query param for file list filter

Alternatives considered:
- `artpiece_id=null` (string "null") — awkward for a nullable UUID; query params don't have a native null type
- `artpiece_id=none` — non-standard

`unassigned=true` is a boolean flag that expresses clear semantic intent (show me the inbox). It is consistent with how the list endpoint will be used: the FE needs a filtered inbox view, not a generic artpiece_id filter.

### File delete includes R2 cleanup

Deleting a file row without removing the R2 object leaks storage permanently. The `StorageService` interface (`DeleteObject`) is already available in the codebase. The delete usecase method calls `DeleteObject` before (or after) removing the DB row; failure to delete from R2 is logged but does not block the DB deletion, consistent with how other storage cleanup is handled.

### `toFileResponse` stays package-level in the `handler` package

`toFileResponse` is already defined (unexported) in `file_upload_handler.go`. Since both `FileUploadHandler` and the new `FileHandler` live in the same `handler` package, the function is accessible without any change. No need to move or export it.

## Risks / Trade-offs

- **Bulk replace is not idempotent on cover if called concurrently** → Acceptable: cover rules are deterministic given a file set; concurrent calls on the same artpiece are not a designed-for scenario.
- **R2 delete failure on file delete** → Mitigation: log the error, still delete the DB row. The file is unreachable at the API level; orphaned R2 objects can be cleaned up separately if needed.
- **Artpiece create-with-files: files are validated but not locked** → A race between two concurrent creates referencing the same file ID is structurally impossible to prevent without row-locking, but is not a realistic scenario for a single-user application.
