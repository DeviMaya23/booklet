## Why

Images cannot currently be uploaded through the API — there is no creation endpoint. This change introduces a two-phase presigned URL upload flow so clients can upload images directly to R2 without proxying through the backend, while keeping metadata and character associations fully server-controlled.

## What Changes

- **New endpoint** `POST /images` — accepts image metadata, persists a `pending_upload` row, and returns a presigned R2 PUT URL for the client to upload to
- **New endpoint** `POST /images/:id/complete` — called after the client finishes the upload; atomically promotes the pending upload into a permanent `images` row
- **New table** `pending_upload` — transient scratch-pad holding metadata and character IDs while the client uploads to R2
- **New column** `mime_type TEXT NOT NULL` on `images` — stored at upload time, derived from client-provided value
- **R2 key generation** — BE generates the storage key (`users/{userID}/images/{uuid}.{ext}`) at `InitialUpload` time; stored in `pending_upload` and carried forward to `images.image_r2_path`
- **R2 health probe** — `/health` response gains an `r2` field; the existing commented-out scaffold in `HealthHandler` is wired up
- Stale/failed `pending_upload` cleanup is deferred
- Thumbnail generation job is deferred (TODO stub left at `CompleteUpload` site)

## Capabilities

### New Capabilities

- `image-upload`: Two-phase presigned URL upload flow — `InitialUpload` creates a pending record and returns a presigned PUT URL; `CompleteUpload` atomically promotes it to a permanent image
- `health-check`: `/health` endpoint behaviour including DB and R2 probes

### Modified Capabilities

- `image-management`: Image data model gains a required `mime_type` column; image response shape gains a `mime_type` field

## Impact

- New Go files: `internal/domain/pending_upload.go`, `internal/usecase/upload_usecase.go`, `internal/usecase/upload_repository.go`, `internal/handler/upload_handler.go`, `internal/repository/upload_repository.go`
- Modified files: `internal/domain/image.go`, `internal/handler/health.go`, `cmd/server/main.go`
- New migrations: `pending_upload` table, `mime_type` column on `images`
- New dependency already added: `github.com/aws/aws-sdk-go-v2/service/s3` (used via existing `internal/storage/r2.go`)
- New Bruno collection files for both endpoints
