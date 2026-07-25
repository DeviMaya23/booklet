## Why

`worker.New` currently bundles River infrastructure setup (migrations, client creation) with domain-specific worker registration (adding workers, wiring periodic jobs). Every new worker added in the future would require growing `worker.New`'s parameter list with that worker's dependencies, making the function a coupling point rather than a stable infrastructure boundary.

## What Changes

- **BREAKING (internal)**: `worker.New` signature changes — drops all domain-specific params, accepts only `ctx`, `pool`, `workers`, `periodicJobs`, and `logger`; returns a `*river.Client` ready to start
- `worker/purge_upload.go` gains a constructor `NewPurgeExpiredUploadsWorker(uc cleanupUsecase, threshold time.Duration) *PurgeExpiredUploadsWorker`
- `cmd/server/main.go` takes ownership of worker registration: calls `river.AddWorker` and builds the `PeriodicJobs` slice before passing to `worker.New` — mirroring how route registration works in `initApp`
- No changes to `PurgeExpiredUploadsArgs`, `cleanupUsecase` interface, or any usecase/repository code

## Capabilities

### New Capabilities
None.

### Modified Capabilities
None. This is a pure internal structural refactor. The stale-upload-purge behaviour is unchanged.

## Impact

- **Modified**: `internal/worker/worker.go` — signature of `New`
- **Modified**: `internal/worker/purge_upload.go` — adds constructor
- **Modified**: `cmd/server/main.go` — moves `river.AddWorker` and `PeriodicJob` wiring here
- **No API changes**, no schema changes, no new dependencies
