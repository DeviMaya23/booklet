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
	GetByID(ctx context.Context, id, userID string) (*domain.Character, error)
	List(ctx context.Context, userID string) ([]*domain.Character, error)
	Update(ctx context.Context, id, userID string, params UpdateCharacterParams) (*domain.Character, error)
	Delete(ctx context.Context, id, userID string) error
}
