package domain

import (
	"time"

	"github.com/google/uuid"
)

type Commission struct {
	ID         uuid.UUID   `gorm:"type:uuid;primaryKey"`
	UserID     uuid.UUID   `gorm:"type:uuid;not null;column:user_id"`
	ArtistID   *uuid.UUID  `gorm:"type:uuid;column:artist_id"`
	Status     string      `gorm:"type:text;not null;column:status;default:waitlist"`
	Price      *float64    `gorm:"type:numeric;column:price"`
	Paid       bool        `gorm:"column:paid;default:false"`
	PaidDate   *time.Time  `gorm:"type:date;column:paid_date"`
	FinishDate *time.Time  `gorm:"type:date;column:finish_date"`
	Title      *string     `gorm:"type:text;column:title"`
	Notes      *string     `gorm:"type:text;column:notes"`
	Artist     *Artist     `gorm:"foreignKey:ArtistID"`
	Characters []Character `gorm:"many2many:commission_characters;"`
	Artpieces  []Artpiece  `gorm:"foreignKey:CommissionID"`
	CreatedAt  time.Time   `gorm:"column:created_at"`
	UpdatedAt  time.Time   `gorm:"column:updated_at"`
}
