package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/usecase"
	"github.com/devi/booklet/pkg/timeutil"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type commissionRepository struct {
	db *gorm.DB
}

func NewCommissionRepository(db *gorm.DB) *commissionRepository {
	return &commissionRepository{db: db}
}

func (r *commissionRepository) Create(ctx context.Context, c *domain.Commission) (*domain.Commission, error) {
	if err := dbFromContext(ctx, r.db).Create(c).Error; err != nil {
		return nil, fmt.Errorf("insert commission: %w", err)
	}
	return r.GetByID(ctx, c.ID, c.UserID)
}

func (r *commissionRepository) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Commission, error) {
	var c domain.Commission
	err := dbFromContext(ctx, r.db).
		Preload("Artist.Links").
		Preload("Characters").
		Preload("Artpieces.CoverFile").
		Where("id = ? AND user_id = ?", id, userID).
		First(&c).Error
	if err != nil {
		return nil, fmt.Errorf("get commission: %w", err)
	}
	return &c, nil
}

func (r *commissionRepository) List(ctx context.Context, userID uuid.UUID) ([]*domain.Commission, error) {
	var commissions []*domain.Commission
	err := dbFromContext(ctx, r.db).
		Preload("Artist.Links").
		Preload("Characters").
		Preload("Artpieces.CoverFile").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&commissions).Error
	if err != nil {
		return nil, fmt.Errorf("list commissions: %w", err)
	}
	return commissions, nil
}

func (r *commissionRepository) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, params usecase.UpdateCommissionParams) (*domain.Commission, error) {
	var c domain.Commission
	err := dbFromContext(ctx, r.db).
		Where("id = ? AND user_id = ?", id, userID).
		First(&c).Error
	if err != nil {
		return nil, fmt.Errorf("get commission: %w", err)
	}

	updates := map[string]interface{}{
		"title":       params.Title,
		"artist_id":   params.ArtistID,
		"status":      params.Status,
		"price":       params.Price,
		"paid":        params.Paid,
		"paid_date":   timeutil.ParseDate(params.PaidDate),
		"finish_date": timeutil.ParseDate(params.FinishDate),
		"notes":       params.Notes,
	}

	if err := dbFromContext(ctx, r.db).Model(&c).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("update commission: %w", err)
	}

	characters := make([]domain.Character, len(params.CharacterIDs))
	for i, cid := range params.CharacterIDs {
		characters[i] = domain.Character{ID: cid}
	}
	if err := dbFromContext(ctx, r.db).Model(&c).Association("Characters").Replace(characters); err != nil {
		return nil, fmt.Errorf("replace characters: %w", err)
	}

	return r.GetByID(ctx, id, userID)
}

func (r *commissionRepository) Patch(ctx context.Context, id uuid.UUID, userID uuid.UUID, params usecase.PatchCommissionParams) (*domain.Commission, error) {
	var c domain.Commission
	err := dbFromContext(ctx, r.db).
		Where("id = ? AND user_id = ?", id, userID).
		First(&c).Error
	if err != nil {
		return nil, fmt.Errorf("get commission: %w", err)
	}

	updates := make(map[string]interface{})
	if params.Status != nil {
		updates["status"] = *params.Status
	}
	if params.Paid != nil {
		updates["paid"] = *params.Paid
	}
	if params.PaidDate != nil {
		updates["paid_date"] = timeutil.ParseDate(params.PaidDate)
	}
	if params.LastContactedAt != nil {
		t, err := time.Parse(time.RFC3339, *params.LastContactedAt)
		if err != nil {
			return nil, fmt.Errorf("invalid last_contacted_at: %w", err)
		}
		updates["last_contacted_at"] = t
	}

	if len(updates) > 0 {
		if err := dbFromContext(ctx, r.db).Model(&c).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("patch commission: %w", err)
		}
	}

	return r.GetByID(ctx, id, userID)
}

func (r *commissionRepository) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	result := dbFromContext(ctx, r.db).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&domain.Commission{})
	if result.Error != nil {
		return fmt.Errorf("delete commission: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

