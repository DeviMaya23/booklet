## Context

File records have a `thumbnail_r2_path` column that is null until the `generate_file_thumbnail` River worker writes the path on success. The API derives `thumbnail_url` from this path via presigning. `null` means either the thumbnail is still being generated, or it never will be (non-image file, or an image whose job failed after exhausting retries).

The frontend polls every 2s while any file has `thumbnail_url === null` and shows a spinner on those tiles. This loop runs indefinitely for non-image files and for image files whose worker fails (already observed with certain JPEGs that `imaging.Decode` cannot process).

The fix is to persist the generation pipeline outcome so the API can return an unambiguous signal.

## Goals / Non-Goals

**Goals:**
- Add `thumbnail_gen_state` to the `files` table to track auto-generation outcome.
- Expose `thumbnail_gen_state` in all file API responses.
- Worker writes the final state on both success and failure.
- Frontend uses `thumbnail_url` + `thumbnail_gen_state` together to decide spinner vs. fallback icon.

**Non-Goals:**
- Manual thumbnail upload (separate future change).
- Retrying failed thumbnail generation after the fact.
- Changes to the image thumbnail worker (`generate_thumbnail`) — that is for the old image module.

## Decisions

### Decision: Persist state in the DB, not derived in the handler

**Chosen:** Add a `thumbnail_gen_state` column to `files`.

**Alternative:** Derive state in `toFileResponse` — `IsImage && r2Path == nil → "pending"`, otherwise `"unavailable"`. No migration needed.

**Rationale:** Derived state cannot represent image files whose generation job failed — they'd be permanently `"pending"` since `IsImage` stays true and `r2Path` stays null. With real JPEGs already failing, this is not a theoretical edge case. Persisted state is the only way to transition a failed image to `"failed"`.

---

### Decision: `thumbnail_gen_state` values: `"pending"` | `"done"` | `"failed"`

**Chosen:** Three-value enum.

**Alternative A:** Boolean `pending_thumbnail`. Cheaper, but `false` collapses "done" and "failed" — losing the distinction when manual upload is added later. Adding a second boolean to recover it would be a worse API.

**`done`** is chosen over `"ready"` to avoid implying it describes the thumbnail itself (which `thumbnail_url` already does). `thumbnail_gen_state` describes the pipeline run, not the asset.

**`not_applicable`** is a distinct fourth value for non-image files — generation was never attempted. Separating it from `"failed"` keeps the semantics honest and makes the values self-documenting.

---

### Decision: Worker detects final attempt via `job.Attempt >= job.MaxAttempts`

**Chosen:** In `Work()`, when an error occurs, check `job.Attempt >= job.MaxAttempts`. If true, call `UpdateThumbnailGenState(ctx, fileID, "failed")` and return `nil` to stop River from re-enqueuing. Otherwise return the error normally for retry.

**Alternative:** River error hook / completion callback. More indirection, harder to test, and the repository call is the same either way.

**`return nil` on final attempt:** Signals to River that the job completed without error — it won't be retried or marked as errored. This is the standard River pattern for graceful terminal failure.

---

### Decision: Initial state set at complete-upload time, not at insert

**Chosen:** `thumbnail_gen_state` is written when `CompleteFileUpload` creates the file row:
- Image mime type → `"pending"`
- Non-image mime type → `"not_applicable"`

**Alternative:** DB column default of `"pending"` for all rows. Requires a secondary update for non-image files, or a CHECK constraint to enforce it. More complex migration.

---

### Decision: FE derives spinner/fallback, BE returns raw state

The rendering decision (`spinner` vs `fallback icon`) is:
```
showSpinner  = thumbnail_url === null && thumbnail_gen_state === "pending"
showFallback = thumbnail_url === null && thumbnail_gen_state !== "pending"
```

BE returns raw data. FE owns the rendering logic. Returning a derived `show_spinner` field from BE would encode a UI concern in the API and create a redundant third source of truth.

## Risks / Trade-offs

- **Existing rows** — after migration, all existing file rows will have `thumbnail_gen_state = NULL`. `toFileResponse` treats `NULL` as `"failed"`. No backfill needed.

## Migration Plan

1. Add `thumbnail_gen_state TEXT` column to `files` (nullable, no DB-level default).
2. Deploy BE with updated `toFileResponse` that handles `NULL` gracefully (treat as `"pending"` for images, `"failed"` for non-images) before any writes use the new column.
3. New uploads and worker completions start writing the column.
4. Optional follow-up: backfill existing rows based on `mime_type` and `thumbnail_r2_path`.

Rollback: column is additive — removing it reverts to old behavior with no data loss on other fields.
