## 1. Dependencies

- [x] 1.1 Add `riverqueue/river` and `riverqueue/river/riverdriver/riverpgxv5` to `go.mod` via `go get`

## 2. Repository Layer

- [x] 2.1 Add `ListStale(ctx context.Context, olderThan time.Time) ([]*domain.PendingUpload, error)` to the `UploadRepository` interface in `usecase/upload_repository.go`
- [x] 2.2 Implement `ListStale` in `repository/upload_repository.go` — query `WHERE created_at < olderThan`
- [x] 2.3 Write integration tests for `ListStale`: records older than cutoff are returned; records newer than cutoff are not returned

## 3. Usecase Layer

- [x] 3.1 Add `DeleteObject(ctx context.Context, key string) error` to the `StorageService` interface in `usecase/upload_repository.go`
- [x] 3.2 Implement `CleanupStaleUploads(ctx context.Context, threshold time.Duration) error` on `uploadUsecase` — compute cutoff from threshold, call `ListStale`, iterate: delete R2 then delete DB for each record, log count via `tel.Logger`
- [x] 3.3 Write unit test: no stale records found → zero deletions, R2 and DB delete not called
- [x] 3.4 Write unit test: stale records exist → R2 deleted with correct key and DB deleted with correct ID for each record

## 4. Worker Package

- [x] 4.1 Create `internal/worker/purge_upload.go`: define `PurgeExpiredUploadsArgs` (implements `river.JobArgs`) and `PurgeExpiredUploadsWorker` (implements `river.Worker[PurgeExpiredUploadsArgs]`) — worker holds a narrow `cleanupUsecase` interface and threshold, calls `CleanupStaleUploads` in `Work`
- [x] 4.2 Create `internal/worker/worker.go`: implement `New(pool, uc cleanupUsecase, threshold, logger)` — runs River migrations, registers the periodic job (every 5 minutes), returns a `*river.Client` ready to start

## 5. Main Wiring

- [x] 5.1 Open a `pgxpool.Pool` in `main.go` using the existing `cfg.DB.URL`, sized independently from the GORM pool
- [x] 5.2 Call `worker.New(pool, uploadUsecase, presignTTL, logger)` in `initApp` and start the River client alongside Echo; add River client stop to the shutdown sequence

## 6. Lint

- [x] 6.1 Run `golangci-lint run ./...` from `backend/` and fix any issues
