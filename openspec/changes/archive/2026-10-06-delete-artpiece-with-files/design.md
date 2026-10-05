## Context

`DELETE /artpieces/:id` currently removes the artpiece record and lets the FK `ON DELETE SET NULL` orphan attached files. There is no way to also delete those files in the same request. The `file-bulk-delete` capability already handles permanent file deletion (DB + async R2 purge via `purge_r2_objects` worker), so this change wires that mechanism into the artpiece delete flow behind an opt-in query param.

## Goals / Non-Goals

**Goals:**
- Add `?delete_files=true` to `DELETE /artpieces/:id` that permanently deletes all attached files alongside the artpiece
- Keep the existing no-param path completely unchanged

**Non-Goals:**
- Deleting files selectively (some but not all) — that's `DetachFile` / `ReplaceFiles`
- A new standalone endpoint for this operation
- Synchronous R2 deletion — async purge via job queue is accepted (consistent with `file-bulk-delete`)

## Decisions

### D1: Query param, not a new endpoint
`DELETE /artpieces/:id?delete_files=true` rather than `DELETE /artpieces/:id/files` (a separate endpoint).

**Rationale**: The operation is "delete artpiece, and optionally its files" — it's a modifier on an existing action, not a separate resource action. A query param keeps it as one atomic request from the caller's perspective.

**Alternative considered**: Separate `DELETE /artpieces/:id/files` endpoint called before delete. Rejected because it requires two round-trips, leaves a window where files are deleted but the artpiece still exists, and the UI would always call both together anyway.

### D2: Files deleted before artpiece, inside a single transaction
Sequence inside `ArtpieceUsecase.Delete` when `deleteFiles=true`:
1. Fetch artpiece (ownership check) — outside transaction
2. Fetch file list for artpiece (collect IDs + R2 keys) — outside transaction
3. **Transaction**: `fileRepo.BulkDelete(fileIDs)` → `artpieceRepo.Delete(artpieceID)`
4. **After commit**: `jobInserter.Insert(PurgeR2ObjectsArgs{R2Keys: keys})`

Files are deleted before the artpiece record inside the same transaction, ensuring the FK `ON DELETE SET NULL` never fires on already-deleted rows. R2 purge is enqueued after commit (same pattern as `file-bulk-delete`).

**Alternative considered**: Delete artpiece first, then bulk-delete files. Rejected because the FK SET NULL fires immediately, making it harder to collect the file list reliably.

### D3: `ArtpieceFileRepository` gains `BulkDelete`; no new usecase-level deleter interface
Rather than injecting a `FileBulkDeleter` interface that wraps `FileUsecase.BulkDelete`, we add `BulkDelete(ctx, ids, userID)` directly to `ArtpieceFileRepository`. The artpiece usecase already holds this repo; no new dependency is needed.

`FileUsecase.BulkDelete` owns its own ownership check, which we bypass here — ownership is already verified when we fetch the artpiece. The repo method deletes rows directly.

**Alternative considered**: Call `FileUsecase.BulkDelete`. Rejected because it re-fetches and re-validates ownership that the usecase already confirmed, and creates a cross-usecase dependency.

### D4: `JobInserter` added as a new dependency to `ArtpieceUsecase`
Already used by `FileUsecase` and others; injected the same way via `main.go`.

## Risks / Trade-offs

- **Orphaned R2 objects on job enqueue failure** → Accepted tradeoff, same as `file-bulk-delete`. Log the error, return 204. A future sweep job can recover orphans.
- **No files on artpiece** → Fetch returns empty list; skip DB delete and R2 enqueue; proceed to artpiece delete. Returns 204 silently — no error.
- **`Delete` signature change** → The `ArtpieceUsecase` interface in `artpiece_handler.go` changes. The handler is the only caller; update both together.
