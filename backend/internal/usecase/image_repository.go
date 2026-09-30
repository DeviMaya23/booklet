package usecase

import (
	"context"
	"errors"

	"github.com/devi/booklet/internal/domain"
	"github.com/google/uuid"
)

var ErrCharacterNotOwned = errors.New("one or more character IDs do not belong to the user")

type UpdateImageParams struct {
	Title        *string
	ArtistID     *uuid.UUID
	Notes        *string
	CharacterIDs []uuid.UUID
}

type ListImageFilters struct {
	Q            *string     `query:"q"`
	CharacterIDs []uuid.UUID `query:"character_ids"`
	ArtistIDs    []uuid.UUID `query:"artist_ids"`
}

type ImageRepository interface {
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Image, error)
	List(ctx context.Context, userID uuid.UUID, filters ListImageFilters) ([]*domain.Image, error)
	Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, params UpdateImageParams) (*domain.Image, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	ListByCharacterID(ctx context.Context, characterID uuid.UUID, userID uuid.UUID) ([]*domain.Image, error)
}
