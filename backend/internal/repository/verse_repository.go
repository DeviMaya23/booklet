package repository

import (
	"context"
	"fmt"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type verseRepository struct {
	db *gorm.DB
}

func NewVerseRepository(db *gorm.DB) *verseRepository {
	return &verseRepository{db: db}
}

func (r *verseRepository) Create(ctx context.Context, verse *domain.Verse) (*domain.Verse, error) {
	if err := dbFromContext(ctx, r.db).Create(verse).Error; err != nil {
		if isUniqueConstraintViolation(err) {
			return nil, usecase.ErrVerseNameConflict
		}
		return nil, fmt.Errorf("insert verse: %w", err)
	}
	return verse, nil
}

func (r *verseRepository) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Verse, error) {
	var verse domain.Verse
	err := dbFromContext(ctx, r.db).
		Where("id = ? AND user_id = ?", id, userID).
		First(&verse).Error
	if err != nil {
		return nil, fmt.Errorf("get verse: %w", err)
	}
	return &verse, nil
}

func (r *verseRepository) List(ctx context.Context, userID uuid.UUID, _ usecase.ListVerseFilters) ([]*domain.Verse, error) {
	var verses []*domain.Verse
	err := dbFromContext(ctx, r.db).
		Where("user_id = ?", userID).
		Order("last_used_at DESC NULLS LAST, name ASC").
		Find(&verses).Error
	if err != nil {
		return nil, fmt.Errorf("list verses: %w", err)
	}
	return verses, nil
}

func (r *verseRepository) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, params usecase.UpdateVerseParams) (*domain.Verse, error) {
	result := dbFromContext(ctx, r.db).
		Model(&domain.Verse{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("name", params.Name)
	if result.Error != nil {
		if isUniqueConstraintViolation(result.Error) {
			return nil, usecase.ErrVerseNameConflict
		}
		return nil, fmt.Errorf("update verse: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return r.GetByID(ctx, id, userID)
}

func (r *verseRepository) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	result := dbFromContext(ctx, r.db).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&domain.Verse{})
	if result.Error != nil {
		return fmt.Errorf("delete verse: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *verseRepository) UpdateLastUsedAt(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	result := dbFromContext(ctx, r.db).
		Model(&domain.Verse{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("last_used_at", "NOW()")
	if result.Error != nil {
		return fmt.Errorf("update verse last_used_at: %w", result.Error)
	}
	return nil
}
