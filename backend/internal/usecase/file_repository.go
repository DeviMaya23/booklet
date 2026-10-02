package usecase

import (
	"context"

	"github.com/devi/booklet/internal/domain"
	"github.com/google/uuid"
)

type FileRepository interface {
	GetByIDAndUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.File, error)
	GetByIDsAndUserID(ctx context.Context, ids []uuid.UUID, userID uuid.UUID) ([]*domain.File, error)
	List(ctx context.Context, userID uuid.UUID, unassigned bool) ([]*domain.File, error)
	Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, name *string, notes *string) error
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	BulkDelete(ctx context.Context, ids []uuid.UUID, userID uuid.UUID) error
}

type FileStorageService interface {
	DeleteObject(ctx context.Context, key string) error
}
