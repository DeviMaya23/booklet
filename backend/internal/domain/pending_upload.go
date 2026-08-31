package domain

import (
	"time"

	"github.com/google/uuid"
)

type PendingUpload struct {
	ID           uuid.UUID   `gorm:"type:uuid;primaryKey"`
	UserID       uuid.UUID   `gorm:"type:uuid;not null;column:user_id"`
	R2Key        string      `gorm:"type:text;not null;column:r2_key"`
	MimeType     string      `gorm:"type:text;not null;column:mime_type"`
	Title        *string     `gorm:"type:text;column:title"`
	ArtistID     *uuid.UUID  `gorm:"type:uuid;column:artist_id"`
	Notes        *string     `gorm:"type:text;column:notes"`
	CharacterIDs []uuid.UUID `gorm:"type:jsonb;serializer:json;column:character_ids"`
	CreatedAt    time.Time   `gorm:"column:created_at"`
}
