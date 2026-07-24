## Context

When `InitialUpload` runs it creates a `PendingUpload` row and a presigned R2 URL valid for 15 minutes. If the client never calls `CompleteUpload`, both the DB record and any partially-uploaded R2 object are left behind with no cleanup path. The app needs a background sweep to hard-delete these after the presign window expires.

`internal/worker` is already defined as a first-class directory in CONVENTIONS.md. No new architectural pattern is being introduced.

## Goals / Non-Goals

**Goals:**
- Periodically find and delete `PendingUpload` records older than a configurable threshold
- For each, attempt R2 object deletion before removing the DB record, so failures are always retryable
- Log the count of deletions per run in the usecase

**Non-Goals:**
- Alerting or metrics on cleanup volume (observability can be added later)
- Cleaning up R2 objects not tracked in `pending_uploads` (orphan scan is a separate concern)
- Rate-limiting or batching large sweeps (volume does not warrant it now)

## Decisions

### 1. River for job scheduling

River is chosen over a hand-rolled goroutine ticker because it gives retry semantics, job persistence through restarts, and built-in periodic scheduling without additional infrastructure. A raw `time.Ticker` would lose in-flight work on process restart and has no retry on failure.

**Alternative considered:** `gocron` — simpler, but no persistence or retry. A failed cleanup run is silently dropped.

### 2. Periodic scan, not per-upload deferred job

A single periodic River job scans for all stale records on each tick, rather than enqueueing one deferred job per upload at initiation time.

**Why:** Per-upload jobs require `InitialUpload` to successfully enqueue a job as a second operation — if that enqueue fails or River is temporarily down, the record is never scheduled for cleanup. A periodic scan is self-healing: it picks up any record that matches the age predicate regardless of how it got there.

**Trade-off:** Cleanup fires up to one scan-interval late (e.g. a record eligible at T+15min may not be purged until T+15min+5min). Acceptable for this use case.

### 3. Dual pgxpool alongside GORM

River requires a `pgxpool.Pool` — it cannot use `database/sql`. `pgx/v5` is already an indirect dependency (pulled by GORM's postgres driver), so no new transitive dep is added. A second pool is opened in `main.go` pointing at the same DB URL, used exclusively by River. GORM retains its own pool.

**Alternative considered:** Extracting GORM's underlying `*sql.DB` and wrapping it. River does not support `database/sql`; this path does not exist.

### 4. R2 deletion before DB deletion

For each stale record the usecase:
1. Calls `storage.DeleteObject` (R2)
2. On success (or `NoSuchKey` / 404, treated as success): calls `uploadRepo.Delete` (DB)

If R2 deletion fails → job retries → R2 deletion attempted again → proceeds to DB on success.
If R2 deletion succeeds but DB deletion fails → job retries → R2 deletion returns 404 (treated as success) → DB deletion retried → succeeds.

Reversing the order creates a permanent leak: if DB deletion commits but R2 deletion then fails, the record is gone so no retry can find it, and the R2 object remains indefinitely.

**No GORM transaction needed.** R2 is not a transactional resource, and the DB operation is a single-row delete — there is nothing to roll back atomically.

### 5. Threshold passed at construction, not hardcoded in worker

The worker receives the cleanup threshold (duration) when constructed in `main.go`. This is where `presignTTL` is already known. The worker itself is generic — it does not import or reference the presign constant.

### 6. Logging in the usecase

`CleanupStaleUploads` logs the count of deleted records via the existing `tel.Logger` pattern, consistent with `InitialUpload` and `CompleteUpload`. The worker layer only handles the error return.

## Risks / Trade-offs

**Partial batch failure** → If the job fails mid-batch (e.g. after deleting 3 of 10 records), River retries the whole job. The 3 already-deleted records will have their R2 deletion attempted again (NoSuchKey → treated as success) and their DB deletion attempted again (record not found → GORM returns no error for a zero-row delete). The retry is idempotent.

**R2 object never uploaded** → `InitialUpload` creates the DB record before the client uploads to R2. If the client never uploads, `DeleteObject` on a non-existent key returns 404 — treated as success. No special handling needed.

**River migration** → River requires its own schema tables in the database. These are managed by River's own migrator, called once at startup before the worker starts. This is a one-time addition to the DB; no application migration file is needed.

## Migration Plan

1. Add `riverqueue/river` and `riverqueue/river/riverdriver/riverpgxv5` to `go.mod`
2. On startup, run River's migrator against the existing DB before starting the Echo server
3. Register the periodic worker and start River's client in `main.go` alongside Echo
4. River schema tables are additive — no existing tables are modified; rollback is dropping the `river_*` tables
