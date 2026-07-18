package repository

import (
	"context"
	"fmt"

	"github.com/devi/booklet/internal/domain"
	"github.com/devi/booklet/internal/usecase"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type imageRepository struct {
	db *gorm.DB
}

func NewImageRepository(db *gorm.DB) *imageRepository {
	return &imageRepository{db: db}
}

func (r *imageRepository) Create(ctx context.Context, image *domain.Image) (*domain.Image, error) {
	if err := dbFromContext(ctx, r.db).Create(image).Error; err != nil {
		return nil, fmt.Errorf("insert image: %w", err)
	}
	return image, nil
}

func (r *imageRepository) GetByID(ctx context.Context, id, userID string) (*domain.Image, error) {
	var image domain.Image
	err := dbFromContext(ctx, r.db).
		Preload("Characters").
		Where("id = ? AND user_id = ?", id, userID).
		First(&image).Error
	if err != nil {
		return nil, fmt.Errorf("get image: %w", err)
	}
	return &image, nil
}

func (r *imageRepository) List(ctx context.Context, userID string) ([]*domain.Image, error) {
	var images []*domain.Image
	err := dbFromContext(ctx, r.db).
		Preload("Characters").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&images).Error
	if err != nil {
		return nil, fmt.Errorf("list images: %w", err)
	}
	return images, nil
}

func (r *imageRepository) Update(ctx context.Context, id, userID string, params usecase.UpdateImageParams) (*domain.Image, error) {
	var image domain.Image
	err := dbFromContext(ctx, r.db).
		Where("id = ? AND user_id = ?", id, userID).
		First(&image).Error
	if err != nil {
		return nil, fmt.Errorf("get image: %w", err)
	}

	updates := map[string]interface{}{}
	if params.Title != nil {
		updates["title"] = *params.Title
	}
	if params.ThumbnailR2Path != nil {
		updates["thumbnail_r2_path"] = *params.ThumbnailR2Path
	}
	if params.ArtistName != nil {
		updates["artist_name"] = *params.ArtistName
	}
	if params.ArtistLink != nil {
		updates["artist_link"] = *params.ArtistLink
	}
	if params.Notes != nil {
		updates["notes"] = *params.Notes
	}

	if len(updates) > 0 {
		if err := dbFromContext(ctx, r.db).Model(&image).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("update image: %w", err)
		}
	}

	if params.CharacterIDs != nil {
		charIDs := *params.CharacterIDs
		if len(charIDs) > 0 {
			var count int64
			dbFromContext(ctx, r.db).Model(&domain.Character{}).
				Where("id IN ? AND user_id = ?", charIDs, userID).
				Count(&count)
			if count != int64(len(charIDs)) {
				return nil, usecase.ErrCharacterNotOwned
			}
		}

		characters := make([]domain.Character, len(charIDs))
		for i, cid := range charIDs {
			parsed, _ := uuid.Parse(cid)
			characters[i] = domain.Character{ID: parsed}
		}
		if err := dbFromContext(ctx, r.db).Model(&image).Association("Characters").Replace(characters); err != nil {
			return nil, fmt.Errorf("replace characters: %w", err)
		}
	}

	return r.GetByID(ctx, id, userID)
}

func (r *imageRepository) Delete(ctx context.Context, id, userID string) error {
	result := dbFromContext(ctx, r.db).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&domain.Image{})
	if result.Error != nil {
		return fmt.Errorf("delete image: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
