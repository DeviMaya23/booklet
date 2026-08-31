package usecase

import (
	"context"
	"errors"

	"github.com/devi/booklet/internal/domain"
	"github.com/google/uuid"
)

var (
	ErrArtistNotOwned    = errors.New("artist does not belong to the user")
	ErrArtistNameConflict = errors.New("an artist with that name already exists")
)

type UpdateArtistParams struct {
	Name       *string
	Notes      *string
	ArtistLink *string
}

type ArtistRepository interface {
	Create(ctx context.Context, artist *domain.Artist) (*domain.Artist, error)
	GetByID(ctx context.Context, id string, userID uuid.UUID) (*domain.Artist, error)
	GetByIDAndUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Artist, error)
	List(ctx context.Context, userID uuid.UUID) ([]*domain.Artist, error)
	Update(ctx context.Context, id string, userID uuid.UUID, params UpdateArtistParams) (*domain.Artist, error)
	Delete(ctx context.Context, id string, userID uuid.UUID) error
}
