package domain

import (
	"time"

	"github.com/google/uuid"
)

type PendingCharacterAvatarUpload struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;column:user_id"`
	CharacterID uuid.UUID `gorm:"type:uuid;not null;column:character_id"`
	R2Key       string    `gorm:"type:text;not null;column:r2_key"`
	MimeType    string    `gorm:"type:text;not null;column:mime_type"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}
