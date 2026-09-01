## Context

`thumbnail_r2_path` exists on the `images` table and in the image response shape but has always been null — no code generates or stores thumbnails. River is already wired for background jobs (purge workers). The upload usecase owns CompleteUpload, which is where thumbnail generation should be triggered. The transactor is GORM-based (`GormTransactor`), not pgx-native.

## Goals / Non-Goals

**Goals:**
- Generate a 600px max long edge JPEG thumbnail for every image on upload completion
- Store the thumbnail in R2 and persist its path to `thumbnail_r2_path`
- Keep CompleteUpload fast (async generation)
- Null `thumbnail_r2_path` is a valid, handled state for callers

**Non-Goals:**
- Retroactive thumbnail generation for images that already exist
- Polling endpoint or push mechanism to notify FE when thumbnail is ready
- Thumbnail generation for image updates (only on initial upload)

## Decisions

### 1. Async River job over synchronous generation

Generating a thumbnail synchronously in CompleteUpload adds R2 fetch + resize + R2 write latency to every upload completion call. With bulk upload coming (FE will orchestrate concurrent single uploads), this would hold each concurrent slot longer. Async keeps CompleteUpload fast regardless of scale and isolates thumbnail failures from upload completion.

**Alternative considered**: Sync generation in CompleteUpload — rejected because it couples upload success to R2 read/write latency and blocks concurrent upload slots longer.

### 2. Post-commit enqueue, not transactional

River supports transactional job insertion via `InsertTx(pgx.Tx)`, which would guarantee the job only enqueues if the image record commits. However, the existing `Transactor` abstraction is GORM-based and does not surface a pgx.Tx — extracting one would require breaking the abstraction or adding a parallel pgx transaction path.

Instead, the job is enqueued immediately after `InTransaction` returns successfully. The failure window is narrow (image committed, enqueue failed) and the consequence is acceptable: `thumbnail_r2_path` stays null permanently, which callers already handle. Enqueue failures are logged.

**Alternative considered**: Extend Transactor to support pgx.Tx — rejected as unnecessary complexity for a low-risk failure case.

**Init-cycle note**: `uploadUsecase` needs a `JobInserter` (to enqueue), but `riverClient` (the concrete `JobInserter`) is built from workers that include `PurgeExpiredUploadsWorker`, which needs `uploadUsecase`. This creates a `uploadUsecase ↔ riverClient` init cycle. It is broken via a `riverEnqueuer` deferred wrapper in `main.go`: the struct is created before both usecases and workers, passed to `uploadUsecase` as the `JobInserter`, then its `client` field is set to `riverClient` after `river.NewClient` returns.

### 3. Enqueue from the usecase layer

The `PurgeUserStorage` job is enqueued from the handler layer because it depends on handler-level output (which R2 keys to delete). Thumbnail enqueue is unconditional on CompleteUpload success — it belongs in the usecase, which already owns the image ID and the full upload completion logic. `uploadUsecase` gains a `JobInserter` dependency (same interface already defined in `user_handler.go`).

### 4. Thumbnail always encoded as JPEG

Input images are JPEG or PNG. Output thumbnails are always JPEG regardless of input. This normalises the thumbnail format, simplifies the R2 path (always `.jpg`), and is the conventional choice for thumbnails.

### 5. R2 path: `users/{user_id}/thumbnails/{image_id}.jpg`

Consistent with the existing `users/{user_id}/images/{image_id}.ext` pattern. Keeps per-user content namespaced. The thumbnail key is deterministic from the image ID — no need to pass it through the job args beyond the image ID.

### 6. `disintegration/imaging` for decode, resize, and encode

Pure Go, no CGO, handles JPEG and PNG decode/encode with correct aspect-ratio scaling. `imaging.Decode` handles format detection from the reader, `imaging.Fit(img, 600, 600, imaging.Lanczos)` produces a max 600px long edge output preserving aspect ratio, and `imaging.Encode` handles JPEG output. Using the library for the full pipeline (not mixing with stdlib `image/jpeg`) keeps the encode path consistent with the spec constraint. No alternative libraries considered — this is the standard choice for this use case in Go.

## Risks / Trade-offs

- **Enqueue failure after commit**: image exists with `thumbnail_r2_path = null` permanently. Acceptable — null thumbnail is a handled state. Mitigated by logging the enqueue error so it's observable.
- **River job failure / retry exhaustion**: same outcome as above. River's default retry behaviour (25 attempts with exponential backoff) makes permanent failure unlikely for transient R2 errors. Permanent failures (corrupt image, unsupported encoding) will exhaust retries and leave null — acceptable.
- **R2 read in the worker**: worker fetches the original from R2 by `image_r2_path`. If the original is deleted before the job runs (e.g. stale upload purge race), the job fails and retries. This edge case is unlikely given the purge threshold is 15 minutes and the job runs immediately after CompleteUpload.
- **No backfill**: existing images with `thumbnail_r2_path = null` remain null. Out of scope — FE must handle null thumbnails regardless.

## Migration Plan

- No DB migration required — `thumbnail_r2_path` column already exists.
- Deploy is additive: new worker registered, new job kind added. No existing behaviour changes.
- Rollback: remove worker registration and the enqueue call from CompleteUpload. In-flight jobs will fail (worker gone) — River will retry then exhaust; `thumbnail_r2_path` stays null. Safe.
