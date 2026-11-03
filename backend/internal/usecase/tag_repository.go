package usecase

import (
	"context"
	"errors"

	"github.com/devi/booklet/internal/domain"
	"github.com/google/uuid"
)

var (
	ErrTagNameConflict = errors.New("a tag with that name already exists")
	ErrTagNameInvalid  = errors.New("tag name must not be blank")
	ErrTagNotOwned     = errors.New("tag does not exist or does not belong to the user")
)

type CreateTagParams struct {
	Name string
}

type UpdateTagParams struct {
	Name string
}

type ListTagFilters struct {
	Sort *string `query:"sort" json:"sort" validate:"omitempty,oneof=recent"`
}

type TagRepository interface {
	Create(ctx context.Context, tag *domain.Tag) (*domain.Tag, error)
	List(ctx context.Context, userID uuid.UUID, filters ListTagFilters) ([]*domain.Tag, error)
	Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, params UpdateTagParams) (*domain.Tag, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	GetByIDsAndUserID(ctx context.Context, ids []uuid.UUID, userID uuid.UUID) ([]domain.Tag, error)
	UpdateLastUsedAt(ctx context.Context, ids []uuid.UUID, userID uuid.UUID) error
}
