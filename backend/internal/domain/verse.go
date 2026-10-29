package domain

import (
	"time"

	"github.com/google/uuid"
)

type Verse struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey"`
	UserID     uuid.UUID  `gorm:"type:uuid;not null;column:user_id"`
	Name       string     `gorm:"type:text;not null;column:name"`
	SortOrder  int        `gorm:"column:sort_order;not null;default:0"`
	LastUsedAt *time.Time `gorm:"column:last_used_at"`
	CreatedAt  time.Time  `gorm:"column:created_at"`
	UpdatedAt  time.Time  `gorm:"column:updated_at"`
}
