package repository

import (
	"context"
	"fmt"
	"strings"

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

func (r *imageRepository) GetByID(ctx context.Context, id string, userID uuid.UUID) (*domain.Image, error) {
	var image domain.Image
	err := dbFromContext(ctx, r.db).
		Preload("Characters").
		Preload("Artist").
		Where("id = ? AND user_id = ?", id, userID).
		First(&image).Error
	if err != nil {
		return nil, fmt.Errorf("get image: %w", err)
	}
	return &image, nil
}

func (r *imageRepository) List(ctx context.Context, userID uuid.UUID, filters usecase.ListImageFilters) ([]*domain.Image, error) {
	var images []*domain.Image
	q := dbFromContext(ctx, r.db).
		Preload("Characters").
		Preload("Artist").
		Where("images.user_id = ?", userID)
	if filters.Q != nil {
		q = q.Where("LOWER(images.title) LIKE ?", "%"+strings.ToLower(*filters.Q)+"%")
	}
	if len(filters.CharacterIDs) > 0 {
		q = q.Distinct().
			Joins("JOIN image_characters ON image_characters.image_id = images.id").
			Where("image_characters.character_id::text IN ?", filters.CharacterIDs)
	}
	if len(filters.ArtistIDs) > 0 {
		q = q.Where("images.artist_id::text IN ?", filters.ArtistIDs)
	}
	err := q.Order("images.created_at DESC").Find(&images).Error
	if err != nil {
		return nil, fmt.Errorf("list images: %w", err)
	}
	return images, nil
}

func (r *imageRepository) Update(ctx context.Context, id string, userID uuid.UUID, params usecase.UpdateImageParams) (*domain.Image, error) {
	var image domain.Image
	err := dbFromContext(ctx, r.db).
		Where("id = ? AND user_id = ?", id, userID).
		First(&image).Error
	if err != nil {
		return nil, fmt.Errorf("get image: %w", err)
	}

	var title, notes interface{}
	if params.Title != nil {
		title = *params.Title
	}
	if params.Notes != nil {
		notes = *params.Notes
	}
	updates := map[string]interface{}{
		"title": title,
		"notes": notes,
	}

	if params.ArtistID == nil {
		updates["artist_id"] = nil
	} else {
		var count int64
		dbFromContext(ctx, r.db).Model(&domain.Artist{}).
			Where("id = ? AND user_id = ?", *params.ArtistID, userID).
			Count(&count)
		if count == 0 {
			return nil, usecase.ErrArtistNotOwned
		}
		updates["artist_id"] = *params.ArtistID
	}

	if err := dbFromContext(ctx, r.db).Model(&image).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("update image: %w", err)
	}

	charIDs := params.CharacterIDs
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

	return r.GetByID(ctx, id, userID)
}

func (r *imageRepository) GetByIDForWorker(ctx context.Context, id uuid.UUID) (*domain.Image, error) {
	var image domain.Image
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&image).Error
	if err != nil {
		return nil, fmt.Errorf("get image for worker: %w", err)
	}
	return &image, nil
}

func (r *imageRepository) UpdateThumbnailPath(ctx context.Context, id uuid.UUID, r2Path string) error {
	result := r.db.WithContext(ctx).
		Model(&domain.Image{}).
		Where("id = ?", id).
		Update("thumbnail_r2_path", r2Path)
	if result.Error != nil {
		return fmt.Errorf("update thumbnail_r2_path: %w", result.Error)
	}
	return nil
}

func (r *imageRepository) ListByCharacterID(ctx context.Context, characterID uuid.UUID, userID uuid.UUID) ([]*domain.Image, error) {
	var images []*domain.Image
	err := dbFromContext(ctx, r.db).
		Joins("JOIN image_characters ON image_characters.image_id = images.id").
		Where("image_characters.character_id = ? AND images.user_id = ?", characterID, userID).
		Find(&images).Error
	if err != nil {
		return nil, fmt.Errorf("list images by character: %w", err)
	}
	return images, nil
}

func (r *imageRepository) Delete(ctx context.Context, id string, userID uuid.UUID) error {
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
