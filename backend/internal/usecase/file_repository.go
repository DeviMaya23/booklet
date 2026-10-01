package usecase

import (
	"context"

	"github.com/devi/booklet/internal/domain"
	"github.com/google/uuid"
)

type FileRepository interface {
	GetByIDAndUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.File, error)
	List(ctx context.Context, userID uuid.UUID, unassigned bool) ([]*domain.File, error)
	UpdateNotes(ctx context.Context, id uuid.UUID, userID uuid.UUID, notes *string) error
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type FileStorageService interface {
	DeleteObject(ctx context.Context, key string) error
}
