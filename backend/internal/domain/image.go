package domain

import (
	"time"

	"github.com/google/uuid"
)

type Image struct {
	ID               uuid.UUID  `gorm:"type:uuid;primaryKey"`
	UserID           uuid.UUID  `gorm:"type:uuid;not null;column:user_id"`
	ImageR2Path      string     `gorm:"type:text;not null;column:image_r2_path"`
	MimeType         string     `gorm:"type:text;not null;column:mime_type"`
	Title            *string    `gorm:"type:text;column:title"`
	ThumbnailR2Path  *string    `gorm:"type:text;column:thumbnail_r2_path"`
	ArtistName       *string    `gorm:"type:text;column:artist_name"`
	ArtistLink       *string    `gorm:"type:text;column:artist_link"`
	Notes            *string    `gorm:"type:text;column:notes"`
	CreatedAt        time.Time  `gorm:"column:created_at"`
	UpdatedAt        time.Time  `gorm:"column:updated_at"`
	Characters       []Character `gorm:"many2many:image_characters;"`
}
