package models

import (
	"time"

	"github.com/google/uuid"
)

type Property struct {
	ID           uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	Category     string          `json:"category"` // buy|rent|shortlet
	Price        float64         `json:"price"`
	Currency     string          `gorm:"default:NGN" json:"currency"`
	Address      string          `json:"address"`
	Area         string          `json:"area"`
	Bedrooms     int             `json:"bedrooms"`
	Bathrooms    int             `json:"bathrooms"`
	Furnished    bool            `json:"furnished"`
	PartyAllowed bool            `json:"party_allowed"`
	InstantBook  bool            `json:"instant_book"`
	OwnerID      *uuid.UUID      `gorm:"type:uuid" json:"owner_id"`
	Images       []PropertyImage `gorm:"foreignKey:PropertyID" json:"images,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}
