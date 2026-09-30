package worker

import (
	"context"
	"time"

	"github.com/riverqueue/river"
)

type fileCleanupUsecase interface {
	CleanupStaleUploads(ctx context.Context, threshold time.Duration) error
}

type PurgeExpiredFileUploadsArgs struct{}

func (PurgeExpiredFileUploadsArgs) Kind() string { return "purge_expired_file_uploads" }

type PurgeExpiredFileUploadsWorker struct {
	river.WorkerDefaults[PurgeExpiredFileUploadsArgs]
	usecase   fileCleanupUsecase
	threshold time.Duration
}

func NewPurgeExpiredFileUploadsWorker(uc fileCleanupUsecase, threshold time.Duration) *PurgeExpiredFileUploadsWorker {
	return &PurgeExpiredFileUploadsWorker{usecase: uc, threshold: threshold}
}

func (w *PurgeExpiredFileUploadsWorker) Work(ctx context.Context, _ *river.Job[PurgeExpiredFileUploadsArgs]) error {
	return w.usecase.CleanupStaleUploads(ctx, w.threshold)
}
