package usecase

import (
	"context"
	"errors"

	"github.com/devi/booklet/internal/domain"
	"github.com/google/uuid"
)

const (
	CommissionStatusWaitlist = "waitlist"
	CommissionStatusWIP      = "wip"
	CommissionStatusDone     = "done"
)

func isValidCommissionStatus(s string) bool {
	return s == CommissionStatusWaitlist || s == CommissionStatusWIP || s == CommissionStatusDone
}

type CommissionRepository interface {
	Create(ctx context.Context, c *domain.Commission) (*domain.Commission, error)
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Commission, error)
	List(ctx context.Context, userID uuid.UUID) ([]*domain.Commission, error)
	Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, params UpdateCommissionParams) (*domain.Commission, error)
	Patch(ctx context.Context, id uuid.UUID, userID uuid.UUID, params PatchCommissionParams) (*domain.Commission, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type CommissionArtistRepository interface {
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Artist, error)
	UpdateLastUsedAt(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type CommissionCharacterRepository interface {
	GetByIDsAndUserID(ctx context.Context, ids []uuid.UUID, userID uuid.UUID) ([]domain.Character, error)
}

type CommissionArtpieceRepository interface {
	GetByIDsAndUserID(ctx context.Context, ids []uuid.UUID, userID uuid.UUID) ([]*domain.Artpiece, error)
	GetByIDAndUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Artpiece, error)
	BulkUpdateCommissionID(ctx context.Context, artpieceIDs []uuid.UUID, commissionID *uuid.UUID) error
	GetArtpiecesForCommission(ctx context.Context, commissionID uuid.UUID) ([]*domain.Artpiece, error)
}

type CreateCommissionParams struct {
	Title        *string
	ArtistID     *uuid.UUID
	Status       string
	Price        *float64
	Paid         bool
	PaidDate     *string
	FinishDate   *string
	Notes        *string
	CharacterIDs []uuid.UUID
	ArtpieceIDs  []uuid.UUID
}

type UpdateCommissionParams struct {
	Title        *string
	ArtistID     *uuid.UUID
	Status       string
	Price        *float64
	Paid         bool
	PaidDate     *string
	FinishDate   *string
	Notes        *string
	CharacterIDs []uuid.UUID
}

type PatchCommissionParams struct {
	Status          *string
	Paid            *bool
	PaidDate        *string
	LastContactedAt *string
}

var (
	ErrArtpieceAlreadyAttached = errors.New("artpiece is already attached to another commission")
	ErrInvalidCommissionStatus = errors.New("invalid commission status: must be one of waitlist, wip, done")
	ErrCommissionNotFound      = errors.New("commission not found")
)
