package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/devi/booklet/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetOrCreate(ctx context.Context, idpSubject string) (*domain.User, error) {
	err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "idp_subject"}},
			DoNothing: true,
		}).
		Create(&domain.User{
			ID:         uuid.New(),
			IDPSubject: idpSubject,
		}).
		Error
	if err != nil {
		return nil, fmt.Errorf("insert user: %w", err)
	}

	return r.GetByIDPSubject(ctx, idpSubject)
}

func (r *userRepository) GetByIDPSubject(ctx context.Context, idpSubject string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("idp_subject = ?", idpSubject).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("select user by idp_subject: %w", err)
	}
	return &user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("select user: %w", err)
	}
	return &user, nil
}

func (r *userRepository) SetPendingDeletion(ctx context.Context, id uuid.UUID) error {
	result := dbFromContext(ctx, r.db).WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ?", id).
		Update("account_state", domain.AccountStatePendingDeletion)
	if result.Error != nil {
		return fmt.Errorf("set pending deletion: %w", result.Error)
	}
	return nil
}

func (r *userRepository) DeleteAllUserData(ctx context.Context, userID uuid.UUID) ([]string, error) {
	db := dbFromContext(ctx, r.db).WithContext(ctx)

	var keys []string

	var chars []domain.Character
	if err := db.Select("avatar_r2_path").Where("user_id = ?", userID).Find(&chars).Error; err != nil {
		return nil, fmt.Errorf("collect character keys: %w", err)
	}
	for _, c := range chars {
		if c.AvatarR2Path != nil {
			keys = append(keys, *c.AvatarR2Path)
		}
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
