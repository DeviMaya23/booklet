## Why

When a client initiates an upload but never calls the complete endpoint (crash, network drop, intentional abort), a `PendingUpload` record and a potentially-uploaded R2 object are left behind indefinitely. These accumulate silently and waste storage with no cleanup path.

## What Changes

- Introduce a River-based background worker that runs on a periodic schedule
- Add a new `internal/worker` package to house the River client setup and job definitions
- Add `CleanupStaleUploads(ctx, threshold) error` to the upload usecase — fetches all pending uploads older than the threshold, deletes each from R2 then from the DB, and logs the count
- Add `ListStale(ctx, olderThan time.Time) ([]*domain.PendingUpload, error)` to the upload repository
- Wire up River's `pgxpool` connection alongside the existing GORM pool in `main.go`

## Capabilities

### New Capabilities

- `stale-upload-purge`: Periodic background job that finds and hard-deletes pending uploads (and their R2 objects) that have exceeded the presign TTL without being completed

### Modified Capabilities

- `image-upload`: The usecase gains a cleanup method; the repository gains a stale-listing query. No changes to the upload initiation or completion flows, and no API surface changes.

## Impact

- **New dependency**: `riverqueue/river` (requires `pgxpool` — a second DB connection pool is added alongside GORM)
- **New package**: `internal/worker`
- **Modified**: `usecase.uploadUsecase` — new `CleanupStaleUploads` method
- **Modified**: `usecase.UploadRepository` interface — new `ListStale` method
- **Modified**: `repository.uploadRepository` — implements `ListStale`
- **Modified**: `cmd/server/main.go` — wires River pool and starts worker alongside Echo
- **No API changes**, no migration required (uses existing `pending_uploads` table)
