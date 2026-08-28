package usecase

import (
	"context"

	"github.com/devi/booklet/internal/domain"
	"github.com/google/uuid"
)

type UpdateCharacterParams struct {
	Name            *string
	HeroImageR2Path *string
	Biography       *string
	IsPublic        *bool
	FolderIDs       *[]uuid.UUID
}

type CharacterRepository interface {
	Create(ctx context.Context, character *domain.Character) error
	GetByID(ctx context.Context, id string, userID uuid.UUID) (*domain.Character, error)
	List(ctx context.Context, userID uuid.UUID) ([]*domain.Character, error)
	Update(ctx context.Context, id string, userID uuid.UUID, params UpdateCharacterParams) (*domain.Character, error)
	Delete(ctx context.Context, id string, userID uuid.UUID) error
}
