package repository

import (
	"context"
	"fmt"

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

func (r *uploadRepository) CreatePendingUpload(ctx context.Context, p *domain.PendingUpload) error {
	if err := r.db.WithContext(ctx).Create(p).Error; err != nil {
		return fmt.Errorf("create pending upload: %w", err)
	}
	return nil
}

func (r *uploadRepository) GetPendingUpload(ctx context.Context, pendingID, userID string) (*domain.PendingUpload, error) {
	var p domain.PendingUpload
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", pendingID, userID).
		First(&p).Error
	if err != nil {
		return nil, fmt.Errorf("get pending upload: %w", err)
	}
	return &p, nil
}

func (r *uploadRepository) CompleteUpload(ctx context.Context, pendingID, userID string, validCharIDs []uuid.UUID) (*domain.Image, error) {
	var result *domain.Image

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var pending domain.PendingUpload
		if err := tx.Where("id = ? AND user_id = ?", pendingID, userID).First(&pending).Error; err != nil {
			return fmt.Errorf("get pending upload: %w", err)
		}

		if err := tx.Delete(&pending).Error; err != nil {
			return fmt.Errorf("delete pending upload: %w", err)
		}

		imageID := uuid.New()
		image := &domain.Image{
			ID:          imageID,
			UserID:      userID,
			ImageR2Path: pending.R2Key,
			MimeType:    pending.MimeType,
			Title:       pending.Title,
			ArtistName:  pending.ArtistName,
			ArtistLink:  pending.ArtistLink,
			Notes:       pending.Notes,
		}
		if err := tx.Create(image).Error; err != nil {
			return fmt.Errorf("create image: %w", err)
		}

		if len(validCharIDs) > 0 {
			characters := make([]domain.Character, len(validCharIDs))
			for i, cid := range validCharIDs {
				characters[i] = domain.Character{ID: cid}
			}
			if err := tx.Model(image).Association("Characters").Replace(characters); err != nil {
				return fmt.Errorf("associate characters: %w", err)
			}
		}

		if err := tx.Preload("Characters").First(image, imageID).Error; err != nil {
			return fmt.Errorf("reload image: %w", err)
		}

		result = image
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
