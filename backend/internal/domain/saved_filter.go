package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type SavedFilter struct {
	ID               uuid.UUID       `gorm:"type:uuid;primaryKey"`
	UserID           uuid.UUID       `gorm:"type:uuid;not null;column:user_id"`
	Name             string          `gorm:"type:text;not null;column:name"`
	ThumbnailR2Path  *string         `gorm:"type:text;column:thumbnail_r2_path"`
	FilterPayload    json.RawMessage `gorm:"type:jsonb;column:filter_payload"`
	CreatedAt        time.Time       `gorm:"column:created_at"`
	UpdatedAt        time.Time       `gorm:"column:updated_at"`
}
