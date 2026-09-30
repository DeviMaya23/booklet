package repository

import (
	"context"
	"fmt"

	"github.com/devi/booklet/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type fileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) *fileRepository {
	return &fileRepository{db: db}
}

func (r *fileRepository) Create(ctx context.Context, f *domain.File) (*domain.File, error) {
	if err := dbFromContext(ctx, r.db).Create(f).Error; err != nil {
		return nil, fmt.Errorf("insert file: %w", err)
	}
	return f, nil
}

func (r *fileRepository) GetByIDAndUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.File, error) {
	var f domain.File
	err := dbFromContext(ctx, r.db).
		Where("id = ? AND user_id = ?", id, userID).
		First(&f).Error
	if err != nil {
		return nil, fmt.Errorf("get file: %w", err)
	}
	return &f, nil
}

func (r *fileRepository) GetByIDForWorker(ctx context.Context, id uuid.UUID) (*domain.File, error) {
	var f domain.File
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&f).Error
	if err != nil {
		return nil, fmt.Errorf("get file for worker: %w", err)
	}
	return &f, nil
}

func (r *fileRepository) UpdateThumbnailPath(ctx context.Context, id uuid.UUID, r2Path string) error {
	result := r.db.WithContext(ctx).
		Model(&domain.File{}).
		Where("id = ?", id).
		Update("thumbnail_r2_path", r2Path)
	if result.Error != nil {
		return fmt.Errorf("update thumbnail_r2_path: %w", result.Error)
	}
	return nil
}

func (r *fileRepository) GetFilesForArtpiece(ctx context.Context, artpieceID uuid.UUID) ([]*domain.File, error) {
	var files []*domain.File
	err := dbFromContext(ctx, r.db).
		Where("artpiece_id = ?", artpieceID).
		Order("created_at ASC").
		Find(&files).Error
	if err != nil {
		return nil, fmt.Errorf("get files for artpiece: %w", err)
	}
	return files, nil
}

func (r *fileRepository) UpdateArtpieceID(ctx context.Context, fileID uuid.UUID, artpieceID *uuid.UUID) error {
	result := dbFromContext(ctx, r.db).
		Model(&domain.File{}).
		Where("id = ?", fileID).
		Update("artpiece_id", artpieceID)
	if result.Error != nil {
		return fmt.Errorf("update artpiece_id: %w", result.Error)
	}
	return nil
}

func (r *fileRepository) CreateImageMetadata(ctx context.Context, m *domain.ImageMetadata) error {
	if err := dbFromContext(ctx, r.db).Create(m).Error; err != nil {
		return fmt.Errorf("insert image_metadata: %w", err)
	}
	return nil
}
