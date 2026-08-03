package worker

import (
	"context"

	"github.com/riverqueue/river"
)

type tombstoneUsecase interface {
	CleanupExpiredTombstones(ctx context.Context) error
}

type PurgeTombstonesArgs struct{}

func (PurgeTombstonesArgs) Kind() string { return "purge_tombstones" }

type PurgeTombstonesWorker struct {
	river.WorkerDefaults[PurgeTombstonesArgs]
	usecase tombstoneUsecase
}

func NewPurgeTombstonesWorker(uc tombstoneUsecase) *PurgeTombstonesWorker {
	return &PurgeTombstonesWorker{usecase: uc}
}

func (w *PurgeTombstonesWorker) Work(ctx context.Context, _ *river.Job[PurgeTombstonesArgs]) error {
	return w.usecase.CleanupExpiredTombstones(ctx)
}
