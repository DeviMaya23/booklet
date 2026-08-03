package domain

import (
	"time"

	"gorm.io/gorm"
)

type AccountState string

const (
	AccountStateActive          AccountState = "active"
	AccountStatePendingDeletion AccountState = "pending_deletion"
	AccountStatePurged          AccountState = "purged"
)

type User struct {
	ID           string         `gorm:"type:text;primaryKey"`
	AccountState AccountState   `gorm:"column:account_state;default:active"`
	PurgedAt     *time.Time     `gorm:"column:purged_at"`
	CreatedAt    time.Time      `gorm:"column:created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index"`
}
