package usecase

import (
	"context"
	"time"

	"github.com/devi/booklet/internal/domain"
	"github.com/google/uuid"
)

type UploadRepository interface {
	Create(ctx context.Context, p *domain.PendingUpload) (*domain.PendingUpload, error)
	GetByID(ctx context.Context, id uuid.UUID, userID string) (*domain.PendingUpload, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ListStale(ctx context.Context, olderThan time.Time) ([]*domain.PendingUpload, error)
}

type UploadCharacterRepository interface {
	GetByIDsAndUserID(ctx context.Context, ids []uuid.UUID, userID string) ([]domain.Character, error)
}

type UploadImageRepository interface {
	Create(ctx context.Context, image *domain.Image) (*domain.Image, error)
}

type StorageService interface {
	GeneratePresignedPutURL(ctx context.Context, key, contentType string, ttl time.Duration) (string, error)
	DeleteObject(ctx context.Context, key string) error
}
