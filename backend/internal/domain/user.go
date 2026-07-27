package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                 string         `gorm:"type:text;primaryKey"`
	IsPendingDeletion  bool           `gorm:"column:is_pending_deletion;default:false"`
	CreatedAt          time.Time      `gorm:"column:created_at"`
	UpdatedAt          time.Time      `gorm:"column:updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"column:deleted_at;index"`
}
