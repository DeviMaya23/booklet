package usecase

import (
	"context"

	"github.com/devi/booklet/internal/domain"
	"github.com/google/uuid"
)

type UploadRepository interface {
	CreatePendingUpload(ctx context.Context, p *domain.PendingUpload) error
	GetPendingUpload(ctx context.Context, pendingID, userID string) (*domain.PendingUpload, error)
	CompleteUpload(ctx context.Context, pendingID, userID string, validCharIDs []uuid.UUID) (*domain.Image, error)
}
