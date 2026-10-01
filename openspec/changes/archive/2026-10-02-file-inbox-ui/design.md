## Context

Files uploaded via `POST /files` (no `artpiece_id`) land in the DB with `artpiece_id IS NULL` but are currently invisible in the UI. The backend already has all the primitives needed: `GET /files?unassigned=true`, single-file delete, and the two-phase upload flow (`initiate → PUT → complete`). What is missing is the inbox UI surface, a notes search param on the list endpoint, a bulk delete endpoint, and a rename of `PurgeUserStorageWorker` to make it reusable without implying it is specific to account deletion.

## Goals / Non-Goals

**Goals:**
- Add `/app/files` Inbox page with grid, search, multi-select, single/bulk delete, and drag-to-upload
- Add `DELETE /files` bulk delete with async R2 cleanup
- Rename `PurgeUserStorageWorker` → `PurgeR2ObjectsWorker`

**Non-Goals:**
- File detail view (no-op everywhere it is triggered)
- "Add to artpiece" modal flows (no-op; button and menu items are present but inert)
- The `+` button upload modal (button rendered, no-op)
- Searching by artist name (inbox files have no artist association since `artpiece_id IS NULL`)
- Server-side search — no pagination exists yet, so the full file list is always fetched and search is filtered client-side on `notes`

## Decisions

### D1: Bulk delete R2 cleanup is async via existing job worker

Bulk delete deletes DB rows synchronously then enqueues one `PurgeR2ObjectsArgs` job carrying all collected R2 keys (file path + thumbnail path per deleted file). The worker iterates and deletes each key, logging failures without failing the job.

**Why not delete R2 synchronously in the handler or usecase?**  
Holding a DB transaction open while making multiple R2 calls (one per file × two objects) is fragile and slow. The existing `PurgeR2ObjectsWorker` pattern (used by account deletion) already handles this gracefully: best-effort, logged, non-blocking to the caller.

**Trade-off:** If the job enqueue fails, R2 objects are orphaned. This is acceptable — same trade-off the single-delete path already makes for thumbnail cleanup failures.

### D2: Ownership check before bulk delete — reject entire batch on any violation

`BulkDelete` usecase fetches all requested IDs scoped to `user_id` first. If the returned count doesn't match the requested count, the entire operation is rejected with 422.

**Why not silently delete only the owned subset?**  
Deleting a partial set without telling the caller is surprising behavior. A mismatch almost certainly means a bug or a crafted request — rejecting loudly is safer.

### D3: `FileUsecase` gains `jobInserter` dependency

Rather than creating a new usecase, `jobInserter JobInserter` is added to the existing `FileUsecase` struct. This follows the same pattern as `FileUploadUsecase`.

### D4: Multi-file drag-and-drop runs N independent parallel upload flows

Each dropped file triggers its own `initiate → PUT → complete` sequence. There is no batch upload API. Placeholder shimmer tiles are inserted immediately into local React state (keyed by a transient client-side ID); they are replaced by real tiles once the `useFiles` poll detects `thumbnail_url` is non-null for that file.

**Why not wait for all uploads to finish before showing anything?**  
Large files could take a while. Per-file placeholders give immediate feedback and let fast files resolve independently.

### D5: Selection state lives in `FileInboxGrid` component state

A `Set<string>` of selected file IDs is sufficient. Shift-click range is tracked with a `lastClickedId` ref. No global state manager needed — selection is ephemeral UI state scoped to the page.

### D6: `PurgeUserStorageWorker` renamed to `PurgeR2ObjectsWorker`

The worker's behavior (delete a list of R2 keys) is not specific to user storage deletion. Renaming makes it safe to reuse from any usecase without implying an account-deletion context. `Kind()` changes from `"purge_user_storage"` to `"purge_r2_objects"` — this is a breaking change for any jobs already enqueued under the old kind string. Since River job kinds are matched at dequeue time, any in-flight `purge_user_storage` jobs at deploy time will become unroutable. Acceptable because: (a) these are best-effort cleanup jobs, (b) deploy windows are short, and (c) the account deletion flow that enqueues them is rare.

### D7: Notes search is client-side

The full unassigned file list is already fetched (no pagination exists), so search is a `.filter()` on the in-memory `notes` field — case-insensitive, no extra round-trip per keystroke. When pagination is introduced later, search can be moved server-side at that point.

## Risks / Trade-offs

- **Orphaned R2 objects on job enqueue failure** → Acceptable; same trade-off as existing single-delete thumbnail cleanup. A future sweep job could reconcile orphans if this becomes a problem.
- **`purge_user_storage` jobs in-flight at rename deploy time** → Acceptable; these are best-effort cleanup jobs and the deploy window is short. Document in migration plan.
- **Placeholder tiles keyed by transient client ID can't be matched back to the real file after complete** → After `complete` the `useFiles` query is invalidated; React reconciles by `thumbnail_url` poll. The shimmer stays until the poll resolves, which is by design.

## Migration Plan

1. Deploy backend with `PurgeR2ObjectsWorker` registered alongside (or in place of) `PurgeUserStorageWorker`. No DB migration needed.
2. The rename changes `Kind()` — any `purge_user_storage` jobs in the River queue at cutover become unroutable. These are best-effort; let them expire.
3. Frontend changes are purely additive (new route, new sidebar item) — no existing routes or components are modified.
