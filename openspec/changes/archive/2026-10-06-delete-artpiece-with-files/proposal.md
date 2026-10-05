## Why

When a user deletes an artpiece, its files remain orphaned in their library with no artpiece context, which is rarely what they want. There is currently no way to delete an artpiece and its associated files in a single action.

## What Changes

- `DELETE /artpieces/:id` gains an optional `?delete_files=true` query parameter
- When `delete_files=true`, all files attached to the artpiece are permanently deleted (DB rows removed, R2 objects purged asynchronously) before the artpiece is deleted
- When the parameter is absent or false, behavior is unchanged — files remain with `artpiece_id` set to NULL

## Capabilities

### New Capabilities

_(none)_

### Modified Capabilities

- `artpiece-management`: The delete requirement changes — deletion can now optionally also delete attached files permanently

## Impact

- **Handler**: `artpiece_handler.go` — `DeleteArtpiece` reads the new query param and passes a flag to the usecase
- **Usecase**: `artpiece_usecase.go` — `Delete` receives a `deleteFiles bool`; when true, fetches file list, bulk-deletes DB rows, enqueues R2 purge job, then deletes the artpiece
- **Usecase interface** (`artpiece_handler.go`): `ArtpieceUsecase.Delete` signature changes from `Delete(ctx, id, userID)` to `Delete(ctx, id, userID, deleteFiles bool)`
- **New dependency on `ArtpieceUsecase`**: a `FileBulkDeleter` interface (wrapping `fileRepo.BulkDelete`) and `JobInserter` are injected
- **`ArtpieceFileRepository`**: gains a `BulkDelete(ctx, ids, userID)` method
- **`main.go`**: `ArtpieceUsecase` constructor call updated with new dependencies
- **Bruno**: new collection file for the updated delete endpoint
