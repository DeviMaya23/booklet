package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/devi/booklet/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type uploadRepository struct {
	db *gorm.DB
}

func NewUploadRepository(db *gorm.DB) *uploadRepository {
	return &uploadRepository{db: db}
}

func (r *uploadRepository) Create(ctx context.Context, p *domain.PendingUpload) (*domain.PendingUpload, error) {
	if err := dbFromContext(ctx, r.db).Create(p).Error; err != nil {
		return nil, fmt.Errorf("insert pending upload: %w", err)
	}
	return p, nil
}

func (r *uploadRepository) GetByID(ctx context.Context, id uuid.UUID, userID string) (*domain.PendingUpload, error) {
	var p domain.PendingUpload
	err := dbFromContext(ctx, r.db).
		Where("id = ? AND user_id = ?", id, userID).
		First(&p).Error
	if err != nil {
		return nil, fmt.Errorf("select pending upload: %w", err)
	}
	return &p, nil
}

func (r *uploadRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := dbFromContext(ctx, r.db).
		Where("id = ?", id).
		Delete(&domain.PendingUpload{}).Error; err != nil {
		return fmt.Errorf("delete pending upload: %w", err)
	}
	return nil
}

func (r *uploadRepository) ListStale(ctx context.Context, olderThan time.Time) ([]*domain.PendingUpload, error) {
	var records []*domain.PendingUpload
	if err := dbFromContext(ctx, r.db).
		Where("created_at < ?", olderThan).
		Find(&records).Error; err != nil {
		return nil, fmt.Errorf("list stale pending uploads: %w", err)
	}
	return records, nil
}
