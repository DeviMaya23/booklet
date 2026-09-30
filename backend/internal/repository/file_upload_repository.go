package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/devi/booklet/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type pendingFileUploadRepository struct {
	db *gorm.DB
}

func NewPendingFileUploadRepository(db *gorm.DB) *pendingFileUploadRepository {
	return &pendingFileUploadRepository{db: db}
}

func (r *pendingFileUploadRepository) Create(ctx context.Context, p *domain.PendingFileUpload) (*domain.PendingFileUpload, error) {
	if err := dbFromContext(ctx, r.db).Create(p).Error; err != nil {
		return nil, fmt.Errorf("insert pending file upload: %w", err)
	}
	return p, nil
}

func (r *pendingFileUploadRepository) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.PendingFileUpload, error) {
	var p domain.PendingFileUpload
	err := dbFromContext(ctx, r.db).
		Where("id = ? AND user_id = ?", id, userID).
		First(&p).Error
	if err != nil {
		return nil, fmt.Errorf("select pending file upload: %w", err)
	}
	return &p, nil
}

func (r *pendingFileUploadRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := dbFromContext(ctx, r.db).
		Where("id = ?", id).
		Delete(&domain.PendingFileUpload{}).Error; err != nil {
		return fmt.Errorf("delete pending file upload: %w", err)
	}
	return nil
}

func (r *pendingFileUploadRepository) ListStale(ctx context.Context, olderThan time.Time) ([]*domain.PendingFileUpload, error) {
	var records []*domain.PendingFileUpload
	if err := dbFromContext(ctx, r.db).
		Where("created_at < ?", olderThan).
		Find(&records).Error; err != nil {
		return nil, fmt.Errorf("list stale pending file uploads: %w", err)
	}
	return records, nil
}
