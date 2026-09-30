package domain

import (
	"time"

	"github.com/google/uuid"
)

type PendingFileUpload struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey"`
	UserID     uuid.UUID  `gorm:"type:uuid;not null;column:user_id"`
	R2Key      string     `gorm:"type:text;not null;column:r2_key"`
	MimeType   string     `gorm:"type:text;not null;column:mime_type"`
	ArtpieceID *uuid.UUID `gorm:"type:uuid;column:artpiece_id"`
	Notes      *string    `gorm:"type:text;column:notes"`
	CreatedAt  time.Time  `gorm:"column:created_at"`
}
