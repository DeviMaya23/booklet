package repository

import (
	"context"
	"fmt"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *tagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) Create(ctx context.Context, tag *domain.Tag) (*domain.Tag, error) {
	if err := dbFromContext(ctx, r.db).Create(tag).Error; err != nil {
		if isUniqueConstraintViolation(err) {
			return nil, usecase.ErrTagNameConflict
		}
		return nil, fmt.Errorf("insert tag: %w", err)
	}
	return tag, nil
}

func (r *tagRepository) List(ctx context.Context, userID uuid.UUID, filters usecase.ListTagFilters) ([]*domain.Tag, error) {
	order := "name ASC"
	if filters.Sort != nil && *filters.Sort == "recent" {
		order = "last_used_at DESC NULLS LAST, name ASC"
	}
	var tags []*domain.Tag
	err := dbFromContext(ctx, r.db).
		Where("user_id = ?", userID).
		Order(order).
		Find(&tags).Error
	if err != nil {
		return nil, fmt.Errorf("list tags: %w", err)
	}
	return tags, nil
}

func (r *tagRepository) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, params usecase.UpdateTagParams) (*domain.Tag, error) {
	result := dbFromContext(ctx, r.db).
		Model(&domain.Tag{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("name", params.Name)
	if result.Error != nil {
		if isUniqueConstraintViolation(result.Error) {
			return nil, usecase.ErrTagNameConflict
		}
		return nil, fmt.Errorf("update tag: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	var tag domain.Tag
	if err := dbFromContext(ctx, r.db).
		Where("id = ? AND user_id = ?", id, userID).
		First(&tag).Error; err != nil {
		return nil, fmt.Errorf("get tag: %w", err)
	}
	return &tag, nil
}

func (r *tagRepository) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	result := dbFromContext(ctx, r.db).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&domain.Tag{})
	if result.Error != nil {
		return fmt.Errorf("delete tag: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *tagRepository) GetByIDsAndUserID(ctx context.Context, ids []uuid.UUID, userID uuid.UUID) ([]domain.Tag, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var tags []domain.Tag
	err := dbFromContext(ctx, r.db).
		Where("id IN ? AND user_id = ?", ids, userID).
		Find(&tags).Error
	if err != nil {
		return nil, fmt.Errorf("get tags by ids: %w", err)
	}
	return tags, nil
}

func (r *tagRepository) UpdateLastUsedAt(ctx context.Context, ids []uuid.UUID, userID uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}
	result := dbFromContext(ctx, r.db).
		Model(&domain.Tag{}).
		Where("id IN ? AND user_id = ?", ids, userID).
		Update("last_used_at", gorm.Expr("NOW()"))
	if result.Error != nil {
		return fmt.Errorf("update tags last_used_at: %w", result.Error)
	}
	return nil
}
