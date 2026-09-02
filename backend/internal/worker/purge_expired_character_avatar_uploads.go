package worker

import (
	"context"
	"time"

	"github.com/riverqueue/river"
)

type cleanupAvatarUsecase interface {
	CleanupStaleAvatarUploads(ctx context.Context, threshold time.Duration) error
}

type PurgeExpiredCharacterAvatarUploadsArgs struct{}

func (PurgeExpiredCharacterAvatarUploadsArgs) Kind() string {
	return "purge_expired_character_avatar_uploads"
}

type PurgeExpiredCharacterAvatarUploadsWorker struct {
	river.WorkerDefaults[PurgeExpiredCharacterAvatarUploadsArgs]
	usecase   cleanupAvatarUsecase
	threshold time.Duration
}

func NewPurgeExpiredCharacterAvatarUploadsWorker(uc cleanupAvatarUsecase, threshold time.Duration) *PurgeExpiredCharacterAvatarUploadsWorker {
	return &PurgeExpiredCharacterAvatarUploadsWorker{usecase: uc, threshold: threshold}
}

func (w *PurgeExpiredCharacterAvatarUploadsWorker) Work(ctx context.Context, _ *river.Job[PurgeExpiredCharacterAvatarUploadsArgs]) error {
	return w.usecase.CleanupStaleAvatarUploads(ctx, w.threshold)
}
