package usecase

import (
	"context"
	"errors"

	"github.com/devi/booklet/internal/domain"
	"github.com/google/uuid"
)

type ArtpieceRepository interface {
	Create(ctx context.Context, a *domain.Artpiece) (*domain.Artpiece, error)
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Artpiece, error)
	List(ctx context.Context, userID uuid.UUID, filters ListArtpieceFilters) ([]*domain.Artpiece, error)
	Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, params UpdateArtpieceParams) (*domain.Artpiece, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	UpdateCover(ctx context.Context, artpieceID uuid.UUID, coverFileID *uuid.UUID) error
	GetFilesForArtpiece(ctx context.Context, artpieceID uuid.UUID) ([]*domain.File, error)
}

type ArtpieceArtistRepository interface {
	GetByIDAndUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Artist, error)
}

type ArtpieceCharacterRepository interface {
	GetByIDsAndUserID(ctx context.Context, ids []uuid.UUID, userID uuid.UUID) ([]domain.Character, error)
}

type ArtpieceFileRepository interface {
	GetByIDAndUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.File, error)
	UpdateArtpieceID(ctx context.Context, fileID uuid.UUID, artpieceID *uuid.UUID) error
	GetFilesForArtpiece(ctx context.Context, artpieceID uuid.UUID) ([]*domain.File, error)
}

type ListArtpieceFilters struct {
	CharacterIDs []string `query:"character_ids" validate:"omitempty,dive,uuid4"`
	ArtistIDs    []string `query:"artist_ids"    validate:"omitempty,dive,uuid4"`
}

type CreateArtpieceParams struct {
	Title        *string
	ArtistID     *uuid.UUID
	Notes        *string
	CharacterIDs []uuid.UUID
}

type UpdateArtpieceParams struct {
	Title        *string
	ArtistID     *uuid.UUID
	Notes        *string
	CharacterIDs []uuid.UUID
}

var (
	ErrFileNotInArtpiece  = errors.New("file is not attached to this artpiece")
	ErrFileNotOwned       = errors.New("file does not exist or does not belong to the user")
	ErrFileAlreadyAttached = errors.New("file is already attached to another artpiece")
)
