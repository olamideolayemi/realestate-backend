package models

import (
	"time"

	"github.com/google/uuid"
)

type PropertyImage struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	PropertyID uuid.UUID `gorm:"type:uuid;index" json:"property_id"`
	URL        string    `json:"url"`
	AltText    string    `json:"alt_text"`
	CreatedAt  time.Time `json:"created_at"`
}
