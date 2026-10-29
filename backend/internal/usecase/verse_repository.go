package usecase

import (
	"context"
	"errors"

	"github.com/devi/booklet/internal/domain"
	"github.com/google/uuid"
)

var (
	ErrVerseNameConflict = errors.New("a verse with that name already exists")
	ErrVerseNotOwned     = errors.New("verse does not exist or does not belong to the user")
)

type CreateVerseParams struct {
	Name string
}

type UpdateVerseParams struct {
	Name string
}

type ListVerseFilters struct{}

type VerseRepository interface {
	Create(ctx context.Context, verse *domain.Verse) (*domain.Verse, error)
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Verse, error)
	List(ctx context.Context, userID uuid.UUID, filters ListVerseFilters) ([]*domain.Verse, error)
	Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, params UpdateVerseParams) (*domain.Verse, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	UpdateLastUsedAt(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}
