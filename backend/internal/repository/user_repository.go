package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/devi/booklet/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *userRepository) SetPendingDeletion(ctx context.Context, id string) error {
	result := dbFromContext(ctx, r.db).WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ?", id).
		Update("account_state", domain.AccountStatePendingDeletion)
	if result.Error != nil {
		return fmt.Errorf("set pending deletion: %w", result.Error)
	}
	return nil
}

func (r *userRepository) DeleteAllUserData(ctx context.Context, userID string) ([]string, error) {
	db := dbFromContext(ctx, r.db).WithContext(ctx)

	var keys []string

	var images []domain.Image
	if err := db.Select("image_r2_path", "thumbnail_r2_path").Where("user_id = ?", userID).Find(&images).Error; err != nil {
		return nil, fmt.Errorf("collect image keys: %w", err)
	}
	for _, img := range images {
		keys = append(keys, img.ImageR2Path)
		if img.ThumbnailR2Path != nil {
			keys = append(keys, *img.ThumbnailR2Path)
		}
	}

	var chars []domain.Character
	if err := db.Select("hero_image_r2_path").Where("user_id = ?", userID).Find(&chars).Error; err != nil {
		return nil, fmt.Errorf("collect character keys: %w", err)
	}
	for _, c := range chars {
		if c.HeroImageR2Path != nil {
			keys = append(keys, *c.HeroImageR2Path)
		}
	}

	var uploads []domain.PendingUpload
	if err := db.Select("r2_key").Where("user_id = ?", userID).Find(&uploads).Error; err != nil {
		return nil, fmt.Errorf("collect upload keys: %w", err)
	}
	for _, u := range uploads {
		keys = append(keys, u.R2Key)
	}

	if err := db.Where("user_id = ?", userID).Delete(&domain.PendingUpload{}).Error; err != nil {
		return nil, fmt.Errorf("delete pending uploads: %w", err)
	}

	if err := db.Exec("DELETE FROM image_characters WHERE image_id IN (SELECT id FROM images WHERE user_id = ?)", userID).Error; err != nil {
		return nil, fmt.Errorf("delete image characters: %w", err)
	}

	if err := db.Where("user_id = ?", userID).Delete(&domain.Image{}).Error; err != nil {
		return nil, fmt.Errorf("delete images: %w", err)
	}

	if err := db.Exec("DELETE FROM character_folders WHERE character_id IN (SELECT id FROM characters WHERE user_id = ?)", userID).Error; err != nil {
		return nil, fmt.Errorf("delete character folders: %w", err)
	}

	if err := db.Unscoped().Where("user_id = ?", userID).Delete(&domain.Character{}).Error; err != nil {
		return nil, fmt.Errorf("delete characters: %w", err)
	}

	now := time.Now()
	result := db.Model(&domain.User{}).
		Where("id = ?", userID).
		Updates(map[string]any{
			"account_state": domain.AccountStatePurged,
			"purged_at":     now,
		})
	if result.Error != nil {
		return nil, fmt.Errorf("tombstone user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return keys, nil
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetOrCreate(ctx context.Context, id string) (*domain.User, error) {
	err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&domain.User{ID: id}).
		Error
	if err != nil {
		return nil, fmt.Errorf("insert user: %w", err)
	}

	return r.GetByID(ctx, id)
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User

	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("select user: %w", err)
	}

	return &user, nil
}

func (r *userRepository) DeleteExpiredTombstones(ctx context.Context) error {
	cutoff := time.Now().Add(-24 * time.Hour)
	result := r.db.WithContext(ctx).Unscoped().
		Where("account_state = ? AND purged_at < ?", domain.AccountStatePurged, cutoff).
		Delete(&domain.User{})
	if result.Error != nil {
		return fmt.Errorf("delete expired tombstones: %w", result.Error)
	}
	return nil
}
