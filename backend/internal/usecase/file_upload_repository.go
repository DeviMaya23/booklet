package usecase

import (
	"context"
	"time"

	"github.com/devi/booklet/internal/domain"
	"github.com/google/uuid"
)

type FileUploadPendingRepository interface {
	Create(ctx context.Context, p *domain.PendingFileUpload) (*domain.PendingFileUpload, error)
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.PendingFileUpload, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ListStale(ctx context.Context, olderThan time.Time) ([]*domain.PendingFileUpload, error)
}

type FileUploadArtpieceRepository interface {
	GetByIDAndUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Artpiece, error)
}

type FileUploadFileRepository interface {
	Create(ctx context.Context, f *domain.File) (*domain.File, error)
	GetFilesForArtpiece(ctx context.Context, artpieceID uuid.UUID) ([]*domain.File, error)
	UpdateArtpieceID(ctx context.Context, fileID uuid.UUID, artpieceID *uuid.UUID) error
}

type FileUploadImageMetadataRepository interface {
	CreateImageMetadata(ctx context.Context, m *domain.ImageMetadata) error
}
