package repository

import (
	"context"
	"fmt"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type savedFilterRepository struct {
	db *gorm.DB
}

func NewSavedFilterRepository(db *gorm.DB) *savedFilterRepository {
	return &savedFilterRepository{db: db}
}

func (r *savedFilterRepository) Create(ctx context.Context, sf *domain.SavedFilter) (*domain.SavedFilter, error) {
	if err := dbFromContext(ctx, r.db).Create(sf).Error; err != nil {
		return nil, fmt.Errorf("insert saved filter: %w", err)
	}
	return sf, nil
}

func (r *savedFilterRepository) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.SavedFilter, error) {
	var sf domain.SavedFilter
	err := dbFromContext(ctx, r.db).
		Where("id = ? AND user_id = ?", id, userID).
		First(&sf).Error
	if err != nil {
		return nil, fmt.Errorf("get saved filter: %w", err)
	}
	return &sf, nil
}

func (r *savedFilterRepository) List(ctx context.Context, userID uuid.UUID) ([]*domain.SavedFilter, error) {
	var filters []*domain.SavedFilter
	err := dbFromContext(ctx, r.db).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&filters).Error
	if err != nil {
		return nil, fmt.Errorf("list saved filters: %w", err)
	}
	return filters, nil
}

func (r *savedFilterRepository) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, params usecase.UpdateSavedFilterParams) (*domain.SavedFilter, error) {
	updates := make(map[string]interface{})

	if params.Name.Set && params.Name.Value != nil {
		updates["name"] = *params.Name.Value
	}
	if params.ThumbnailR2Path.Set {
		if params.ThumbnailR2Path.Value == nil {
			updates["thumbnail_r2_path"] = nil
		} else {
			updates["thumbnail_r2_path"] = *params.ThumbnailR2Path.Value
		}
	}
	if params.FilterPayload != nil {
		updates["filter_payload"] = *params.FilterPayload
	}

	if len(updates) > 0 {
		result := dbFromContext(ctx, r.db).
			Model(&domain.SavedFilter{}).
			Where("id = ? AND user_id = ?", id, userID).
			Updates(updates)
		if result.Error != nil {
			return nil, fmt.Errorf("update saved filter: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return nil, gorm.ErrRecordNotFound
		}
	} else {
		var count int64
		if err := dbFromContext(ctx, r.db).Model(&domain.SavedFilter{}).Where("id = ? AND user_id = ?", id, userID).Count(&count).Error; err != nil {
			return nil, fmt.Errorf("update saved filter: %w", err)
		}
		if count == 0 {
			return nil, gorm.ErrRecordNotFound
		}
	}

	return r.GetByID(ctx, id, userID)
}

func (r *savedFilterRepository) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	result := dbFromContext(ctx, r.db).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&domain.SavedFilter{})
	if result.Error != nil {
		return fmt.Errorf("delete saved filter: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
