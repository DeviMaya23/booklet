package domain

import (
	"time"

	"github.com/google/uuid"
)

type ArtistLink struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	ArtistID  uuid.UUID `gorm:"type:uuid;not null;column:artist_id"`
	URL       string    `gorm:"type:text;not null;column:url"`
	IsPrimary bool      `gorm:"column:is_primary;default:false"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

type Artist struct {
	ID         uuid.UUID    `gorm:"type:uuid;primaryKey"`
	UserID     uuid.UUID    `gorm:"type:uuid;not null;column:user_id"`
	Name       string       `gorm:"type:text;not null;column:name"`
	Notes      *string      `gorm:"type:text;column:notes"`
	LastUsedAt *time.Time   `gorm:"column:last_used_at"`
	Links      []ArtistLink `gorm:"foreignKey:ArtistID"`
	CreatedAt  time.Time    `gorm:"column:created_at"`
	UpdatedAt  time.Time    `gorm:"column:updated_at"`
}
