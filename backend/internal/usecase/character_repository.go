package usecase

import (
	"context"

	"github.com/devi/booklet/internal/domain"
)

type UpdateCharacterParams struct {
	Name            *string
	HeroImageR2Path *string
	Biography       *string
	IsPublic        *bool
}

type CharacterRepository interface {
	Create(ctx context.Context, character *domain.Character) error
	GetByID(ctx context.Context, id, userID string) (*domain.Character, error)
	List(ctx context.Context, userID string) ([]*domain.Character, error)
	Update(ctx context.Context, id, userID string, params UpdateCharacterParams) (*domain.Character, error)
	Delete(ctx context.Context, id, userID string) error
}
