package usecase

import (
	"context"
	"time"

	"github.com/devi/booklet/internal/domain"
	"github.com/google/uuid"
	"github.com/riverqueue/river"
	rivertype "github.com/riverqueue/river/rivertype"
)

type UploadRepository interface {
	Create(ctx context.Context, p *domain.PendingUpload) (*domain.PendingUpload, error)
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.PendingUpload, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ListStale(ctx context.Context, olderThan time.Time) ([]*domain.PendingUpload, error)
}

type UploadCharacterRepository interface {
	GetByIDsAndUserID(ctx context.Context, ids []uuid.UUID, userID uuid.UUID) ([]domain.Character, error)
}

type UploadArtistRepository interface {
	GetByIDAndUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Artist, error)
}

type UploadImageRepository interface {
	Create(ctx context.Context, image *domain.Image) (*domain.Image, error)
}

type StorageService interface {
	GeneratePresignedPutURL(ctx context.Context, key, contentType string, ttl time.Duration) (string, error)
	DeleteObject(ctx context.Context, key string) error
}

type JobInserter interface {
	Insert(ctx context.Context, args river.JobArgs, opts *river.InsertOpts) (*rivertype.JobInsertResult, error)
}
