package usecase

import (
	"context"

	"github.com/devi/booklet/internal/domain"
	"github.com/google/uuid"
)

type DashboardRepository interface {
	GetRecentArtpieces(ctx context.Context, userID uuid.UUID) ([]*domain.Artpiece, error)
	GetInProgressCommissions(ctx context.Context, userID uuid.UUID) ([]*domain.Commission, error)
	GetCommissionsNoArtist(ctx context.Context, userID uuid.UUID) ([]DashboardHousekeepingItem, error)
	GetArtpiecesNoArtist(ctx context.Context, userID uuid.UUID) ([]DashboardHousekeepingItem, error)
	GetDoneCommissionsNoArtpieces(ctx context.Context, userID uuid.UUID) ([]DashboardHousekeepingItem, error)
	GetArtpiecesNoFiles(ctx context.Context, userID uuid.UUID) ([]DashboardHousekeepingItem, error)
}

type DashboardHousekeepingItem struct {
	ID    uuid.UUID
	Title *string
}
