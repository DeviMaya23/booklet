package domain

import (
	"time"

	"github.com/google/uuid"
)

type Tag struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey"`
	UserID     uuid.UUID  `gorm:"type:uuid;not null;column:user_id"`
	Name       string     `gorm:"type:text;not null;column:name"`
	LastUsedAt *time.Time `gorm:"column:last_used_at"`
	CreatedAt  time.Time  `gorm:"column:created_at"`
	UpdatedAt  time.Time  `gorm:"column:updated_at"`
}

// TagWithUsage is a Tag plus the number of artpieces it is assigned to.
// It is a read model scanned from a grouped query, not a table of its own.
type TagWithUsage struct {
	Tag
	ArtpieceCount int `gorm:"column:artpiece_count"`
}
