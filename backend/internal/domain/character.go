package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Character struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey"`
	UserID          string         `gorm:"type:text;not null;column:user_id"`
	Name            string         `gorm:"type:text;not null;column:name"`
	HeroImageR2Path *string        `gorm:"type:text;column:hero_image_r2_path"`
	Biography       *string        `gorm:"type:text;column:biography"`
	IsPublic        bool           `gorm:"column:is_public;default:false"`
	CreatedAt       time.Time      `gorm:"column:created_at"`
	UpdatedAt       time.Time      `gorm:"column:updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;index"`
}
