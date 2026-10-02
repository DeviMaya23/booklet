package domain

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey"`
	UserID          uuid.UUID      `gorm:"type:uuid;not null;column:user_id"`
	ArtpieceID      *uuid.UUID     `gorm:"type:uuid;column:artpiece_id"`
	FileR2Path      string         `gorm:"type:text;not null;column:file_r2_path"`
	MimeType        string         `gorm:"type:text;not null;column:mime_type"`
	ThumbnailR2Path *string        `gorm:"type:text;column:thumbnail_r2_path"`
	Name            *string        `gorm:"type:text;column:name"`
	Notes           *string        `gorm:"type:text;column:notes"`
	ImageMetadata   *ImageMetadata `gorm:"foreignKey:FileID"`
	CreatedAt       time.Time      `gorm:"column:created_at"`
	UpdatedAt       time.Time      `gorm:"column:updated_at"`
}

type ImageMetadata struct {
	FileID uuid.UUID `gorm:"type:uuid;primaryKey;column:file_id"`
	Width  int       `gorm:"column:width"`
	Height int       `gorm:"column:height"`
}
