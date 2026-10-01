package worker

import (
	"context"

	"github.com/riverqueue/river"
	"go.uber.org/zap"
)

type storageDeleter interface {
	DeleteObject(ctx context.Context, key string) error
}

type PurgeR2ObjectsArgs struct {
	R2Keys []string `json:"r2_keys"`
}

func (PurgeR2ObjectsArgs) Kind() string { return "purge_r2_objects" }

type PurgeR2ObjectsWorker struct {
	river.WorkerDefaults[PurgeR2ObjectsArgs]
	storage storageDeleter
	logger  *zap.Logger
}

func NewPurgeR2ObjectsWorker(storage storageDeleter, logger *zap.Logger) *PurgeR2ObjectsWorker {
	return &PurgeR2ObjectsWorker{storage: storage, logger: logger}
}

func (w *PurgeR2ObjectsWorker) Work(ctx context.Context, job *river.Job[PurgeR2ObjectsArgs]) error {
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
