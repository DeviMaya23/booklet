package domain

import (
	"time"

	"github.com/google/uuid"
)

type Artist struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;column:user_id"`
	Name       string    `gorm:"type:text;not null;column:name"`
	Notes      *string   `gorm:"type:text;column:notes"`
	ArtistLink *string   `gorm:"type:text;column:artist_link"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}
