package worker

import (
	"context"
	"time"

	"github.com/riverqueue/river"
)

type cleanupUsecase interface {
	CleanupStaleUploads(ctx context.Context, threshold time.Duration) error
}

type PurgeExpiredUploadsArgs struct{}

func (PurgeExpiredUploadsArgs) Kind() string { return "purge_expired_uploads" }

type PurgeExpiredUploadsWorker struct {
	river.WorkerDefaults[PurgeExpiredUploadsArgs]
	usecase   cleanupUsecase
	threshold time.Duration
}

func (w *PurgeExpiredUploadsWorker) Work(ctx context.Context, _ *river.Job[PurgeExpiredUploadsArgs]) error {
	return w.usecase.CleanupStaleUploads(ctx, w.threshold)
}
