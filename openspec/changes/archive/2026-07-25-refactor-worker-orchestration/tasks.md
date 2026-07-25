## 1. Worker Package

- [x] 1.1 Add `NewPurgeExpiredUploadsWorker(uc cleanupUsecase, threshold time.Duration) *PurgeExpiredUploadsWorker` constructor to `internal/worker/purge_upload.go`
- [x] 1.2 Remove `worker.New` and `internal/worker/worker.go` entirely; move migration and River client creation into `initRiverClient` local function in `cmd/server/main.go` (mirrors `initDB` pattern)

## 2. Main Wiring

- [x] 2.1 In `cmd/server/main.go`, build `*river.Workers` and `[]*river.PeriodicJob` before calling `worker.New`: call `river.AddWorker(workers, worker.NewPurgeExpiredUploadsWorker(uploadUsecase, usecase.PresignTTL))` and construct the `PeriodicJob` slice with the 5-minute interval and `RunOnStart: true`
- [x] 2.2 Call `initRiverClient(ctx, riverPool, workers, periodicJobs, logger)` in `initApp`; update `server` struct and `initApp` return type to use `*river.Client[pgx.Tx]` directly

## 3. Lint

- [x] 3.1 Run `golangci-lint run ./...` from `backend/` and fix any issues
