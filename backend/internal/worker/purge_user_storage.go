package worker

import (
	"context"

	"github.com/riverqueue/river"
	"go.uber.org/zap"
)

type storageDeleter interface {
	DeleteObject(ctx context.Context, key string) error
}

type PurgeUserStorageArgs struct {
	R2Keys []string `json:"r2_keys"`
}

func (PurgeUserStorageArgs) Kind() string { return "purge_user_storage" }

type PurgeUserStorageWorker struct {
	river.WorkerDefaults[PurgeUserStorageArgs]
	storage storageDeleter
	logger  *zap.Logger
}

func NewPurgeUserStorageWorker(storage storageDeleter, logger *zap.Logger) *PurgeUserStorageWorker {
	return &PurgeUserStorageWorker{storage: storage, logger: logger}
}

func (w *PurgeUserStorageWorker) Work(ctx context.Context, job *river.Job[PurgeUserStorageArgs]) error {
	for _, key := range job.Args.R2Keys {
		if err := w.storage.DeleteObject(ctx, key); err != nil {
			w.logger.Error("failed to delete R2 object",
				zap.String("key", key),
				zap.Error(err),
			)
		}
	}
	return nil
}
