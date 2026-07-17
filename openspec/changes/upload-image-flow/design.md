## Context

Images are already stored on Cloudflare R2 — the `images` table holds `image_r2_path`, and `internal/storage/r2.go` is a fully instrumented R2 adapter with tracing, metrics, and logging. However there is no creation path: no `POST /images` endpoint exists, no upload mechanism, and no way for a client to get an object into R2 through the API.

The upload flow must avoid proxying binary data through the backend. Presigned PUT URLs let the client upload directly to R2, keeping the backend out of the data path.

The existing `HealthHandler` already has a commented-out R2 probe scaffold — wiring it up is a small change that is natural to bundle here.

## Goals / Non-Goals

**Goals:**
- Two-phase presigned URL upload: client receives a PUT URL, uploads directly to R2, then signals completion
- `mime_type` stored on every image
- R2 health probe in `/health`

**Non-Goals:**
- Stale/failed `pending_upload` cleanup (deferred)
- Thumbnail generation job (deferred — a TODO stub is left at CompleteUpload)
- Server-side content-type verification after upload

## Decisions

### Separate upload handler and usecase

The upload flow depends on a `StorageService` (for presign) that the existing image CRUD handlers never need. Merging them would pollute the image handler's dependency set. `upload_handler.go` and `upload_usecase.go` are new files parallel to the existing image files.

### `pending_upload` stores character IDs as JSONB

`pending_upload` is a short-lived scratch pad that exists only while the client is uploading. A dedicated join table (`pending_upload_characters`) would be a permanent relational structure for a transient entity. JSONB `character_ids []string` is simpler and sufficient — the IDs are validated at CompleteUpload time anyway.

### R2 key generated at InitialUpload time

Pattern: `users/{userID}/images/{uuid}.{ext}` (ext derived from mime_type, e.g. `image/jpeg` → `.jpg`). The key is stored in `pending_upload.r2_key` and copied verbatim to `images.image_r2_path` at CompleteUpload. This means the presigned URL and the final image path are always consistent.

### CompleteUpload transaction encapsulated in the repository

The CompleteUpload operation — delete `pending_upload`, insert `images` row, associate characters — must be atomic. Exposing a `*gorm.DB` transaction handle to the usecase would leak ORM concerns across layers. Instead, `UploadRepository` exposes a single `CompleteUpload(ctx, pendingID, userID string, validCharIDs []uuid.UUID) (*domain.Image, error)` method that owns the transaction internally.

Character filtering (drop invalid IDs) happens in the usecase before calling into the repo — the usecase queries characters by ID scoped to the user, takes the intersection, and passes the validated slice down.

### Narrow `StorageService` interface in upload usecase

The upload usecase only needs one method from R2:

```go
type StorageService interface {
    GeneratePresignedPutURL(ctx context.Context, key, contentType string, ttl time.Duration) (string, error)
}
```

`*r2Storage` satisfies this without modification. Future usecases define their own narrow slices of the same struct.

### Presigned URL TTL: 15 minutes

Short enough to limit the exposure window if a URL leaks, long enough for a typical client upload flow.

### HealthHandler wires R2 via a `Pinger` interface

`HealthHandler` gains a second dependency: a narrow `R2Pinger` interface with a single `Ping(ctx) error` method. `*r2Storage` already implements `Ping`. This keeps the health handler testable without pulling in the full R2 client.

## Risks / Trade-offs

- **Orphaned R2 objects** — if CompleteUpload is never called, the object lands in R2 but no `images` row exists. Cleanup is deferred; accepted for now as the blast radius is storage cost only.
- **Pending rows accumulate** — same root cause; deferred stale-row cleanup. The table is small and reads/writes are cheap.
- **Presigned URL replay** — a leaked PUT URL allows uploading a different object to the same key within 15 minutes. R2's server-side content-type enforcement (set via the presign input) limits this to the declared mime_type.

## Migration Plan

1. Run migration `000004`: create `pending_upload` table
2. Run migration `000005`: add `mime_type TEXT NOT NULL` to `images`
   - Existing rows: backfill with `'image/jpeg'` (safe default — all current images are JPEGs from the previous manual ingestion)
3. Deploy new binary
4. Rollback: revert migrations (both are reversible), deploy previous binary
