package domain

import (
	"time"

	"github.com/google/uuid"
)

type PendingUpload struct {
	ID           uuid.UUID   `gorm:"type:uuid;primaryKey"`
	UserID       string      `gorm:"type:text;not null;column:user_id"`
	R2Key        string      `gorm:"type:text;not null;column:r2_key"`
	MimeType     string      `gorm:"type:text;not null;column:mime_type"`
	Title        *string     `gorm:"type:text;column:title"`
	ArtistName   *string     `gorm:"type:text;column:artist_name"`
	ArtistLink   *string     `gorm:"type:text;column:artist_link"`
	Notes        *string     `gorm:"type:text;column:notes"`
	CharacterIDs []uuid.UUID `gorm:"type:jsonb;serializer:json;column:character_ids"`
	CreatedAt    time.Time   `gorm:"column:created_at"`
}
