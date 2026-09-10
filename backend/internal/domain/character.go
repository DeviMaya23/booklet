package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CharacterFolder struct {
	CharacterID uuid.UUID `gorm:"type:uuid;primaryKey;column:character_id"`
	FolderID    uuid.UUID `gorm:"type:uuid;primaryKey;column:folder_id"`
}

type Character struct {
	ID              uuid.UUID        `gorm:"type:uuid;primaryKey"`
	UserID          uuid.UUID        `gorm:"type:uuid;not null;column:user_id"`
	Name            string           `gorm:"type:text;not null;column:name"`
	AvatarR2Path    *string          `gorm:"type:text;column:avatar_r2_path"`
	Notes           *string          `gorm:"type:text;column:notes"`
	IsPublic        bool             `gorm:"column:is_public;default:false"`
	CreatedAt       time.Time        `gorm:"column:created_at"`
	UpdatedAt       time.Time        `gorm:"column:updated_at"`
	DeletedAt       gorm.DeletedAt   `gorm:"column:deleted_at;index"`
	Folders         []CharacterFolder `gorm:"foreignKey:CharacterID"`
}
