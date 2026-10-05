package domain

import (
	"time"

	"github.com/google/uuid"
)

type Artpiece struct {
	ID           uuid.UUID   `gorm:"type:uuid;primaryKey"`
	UserID       uuid.UUID   `gorm:"type:uuid;not null;column:user_id"`
	Title        *string     `gorm:"type:text;column:title"`
	ArtistID     *uuid.UUID  `gorm:"type:uuid;column:artist_id"`
	CoverFileID  *uuid.UUID  `gorm:"type:uuid;column:cover_file_id"`
	CommissionID *uuid.UUID  `gorm:"type:uuid;column:commission_id"`
	Notes        *string     `gorm:"type:text;column:notes"`
	Artist       *Artist     `gorm:"foreignKey:ArtistID"`
	CoverFile    *File       `gorm:"foreignKey:CoverFileID"`
	Files        []File      `gorm:"foreignKey:ArtpieceID"`
	Characters   []Character `gorm:"many2many:artpiece_characters;"`
	CreatedAt    time.Time   `gorm:"column:created_at"`
	UpdatedAt    time.Time   `gorm:"column:updated_at"`
}
