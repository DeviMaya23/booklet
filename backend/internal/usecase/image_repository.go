package usecase

import (
	"context"
	"errors"

	"github.com/devi/booklet/internal/domain"
)

var ErrCharacterNotOwned = errors.New("one or more character IDs do not belong to the user")

type UpdateImageParams struct {
	Title           *string
	ThumbnailR2Path *string
	ArtistName      *string
	ArtistLink      *string
	Notes           *string
	CharacterIDs    *[]string
}

type ImageRepository interface {
	GetByID(ctx context.Context, id, userID string) (*domain.Image, error)
	List(ctx context.Context, userID string) ([]*domain.Image, error)
	Update(ctx context.Context, id, userID string, params UpdateImageParams) (*domain.Image, error)
	Delete(ctx context.Context, id, userID string) error
}
