package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/devi/booklet/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type characterAvatarRepository struct {
	db *gorm.DB
}

func NewCharacterAvatarRepository(db *gorm.DB) *characterAvatarRepository {
	return &characterAvatarRepository{db: db}
}

func (r *characterAvatarRepository) Create(ctx context.Context, p *domain.PendingCharacterAvatarUpload) (*domain.PendingCharacterAvatarUpload, error) {
	if err := dbFromContext(ctx, r.db).Create(p).Error; err != nil {
		return nil, fmt.Errorf("insert pending character avatar upload: %w", err)
	}
	return p, nil
}

func (r *characterAvatarRepository) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.PendingCharacterAvatarUpload, error) {
	var p domain.PendingCharacterAvatarUpload
	err := dbFromContext(ctx, r.db).
		Where("id = ? AND user_id = ?", id, userID).
		First(&p).Error
	if err != nil {
		return nil, fmt.Errorf("select pending character avatar upload: %w", err)
	}
	return &p, nil
}

func (r *characterAvatarRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := dbFromContext(ctx, r.db).
		Where("id = ?", id).
		Delete(&domain.PendingCharacterAvatarUpload{}).Error; err != nil {
		return fmt.Errorf("delete pending character avatar upload: %w", err)
	}
	return nil
}

func (r *characterAvatarRepository) ListStale(ctx context.Context, olderThan time.Time) ([]*domain.PendingCharacterAvatarUpload, error) {
	var records []*domain.PendingCharacterAvatarUpload
	if err := dbFromContext(ctx, r.db).
		Where("created_at < ?", olderThan).
		Find(&records).Error; err != nil {
		return nil, fmt.Errorf("list stale pending character avatar uploads: %w", err)
	}
	return records, nil
}
