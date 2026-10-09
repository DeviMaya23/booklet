package repository

import (
	"context"
	"fmt"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) *dashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) GetRecentArtpieces(ctx context.Context, userID uuid.UUID) ([]*domain.Artpiece, error) {
	var artpieces []*domain.Artpiece
	err := dbFromContext(ctx, r.db).
		Preload("Artist").
		Preload("CoverFile").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(5).
		Find(&artpieces).Error
	if err != nil {
		return nil, fmt.Errorf("get recent artpieces: %w", err)
	}
	return artpieces, nil
}

func (r *dashboardRepository) GetInProgressCommissions(ctx context.Context, userID uuid.UUID) ([]*domain.Commission, error) {
	var commissions []*domain.Commission
	err := dbFromContext(ctx, r.db).
		Preload("Artist").
		Where("user_id = ? AND status IN ?", userID, []string{"waitlist", "wip"}).
		Order("created_at ASC").
		Find(&commissions).Error
	if err != nil {
		return nil, fmt.Errorf("get in-progress commissions: %w", err)
	}
	return commissions, nil
}

func (r *dashboardRepository) GetCommissionsNoArtist(ctx context.Context, userID uuid.UUID) ([]usecase.DashboardHousekeepingItem, error) {
	return r.queryHousekeepingItems(ctx,
		"SELECT id, title FROM commissions WHERE user_id = ? AND artist_id IS NULL ORDER BY created_at DESC",
		userID,
	)
}

func (r *dashboardRepository) GetArtpiecesNoArtist(ctx context.Context, userID uuid.UUID) ([]usecase.DashboardHousekeepingItem, error) {
	return r.queryHousekeepingItems(ctx,
		"SELECT id, title FROM artpieces WHERE user_id = ? AND artist_id IS NULL ORDER BY created_at DESC",
		userID,
	)
}

func (r *dashboardRepository) GetDoneCommissionsNoArtpieces(ctx context.Context, userID uuid.UUID) ([]usecase.DashboardHousekeepingItem, error) {
	return r.queryHousekeepingItems(ctx,
		"SELECT id, title FROM commissions WHERE user_id = ? AND status = 'done' AND NOT EXISTS (SELECT 1 FROM artpieces WHERE commission_id = commissions.id) ORDER BY created_at DESC",
		userID,
	)
}

func (r *dashboardRepository) GetArtpiecesNoFiles(ctx context.Context, userID uuid.UUID) ([]usecase.DashboardHousekeepingItem, error) {
	return r.queryHousekeepingItems(ctx,
		"SELECT id, title FROM artpieces WHERE user_id = ? AND NOT EXISTS (SELECT 1 FROM files WHERE artpiece_id = artpieces.id) ORDER BY created_at DESC",
		userID,
	)
}

type housekeepingRow struct {
	ID    uuid.UUID `gorm:"column:id"`
	Title *string   `gorm:"column:title"`
}

func (r *dashboardRepository) queryHousekeepingItems(ctx context.Context, query string, userID uuid.UUID) ([]usecase.DashboardHousekeepingItem, error) {
	var rows []housekeepingRow
	if err := dbFromContext(ctx, r.db).Raw(query, userID).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("housekeeping query: %w", err)
	}
	items := make([]usecase.DashboardHousekeepingItem, len(rows))
	for i, row := range rows {
		items[i] = usecase.DashboardHousekeepingItem{ID: row.ID, Title: row.Title}
	}
	return items, nil
}
