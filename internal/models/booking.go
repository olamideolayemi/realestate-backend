package models

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	PropertyID  uuid.UUID `gorm:"type:uuid;index" json:"property_id"`
	UserID      *uuid.UUID `gorm:"type:uuid" json:"user_id"`
	Checkin     time.Time `gorm:"type:date" json:"checkin"`
	Checkout    time.Time `gorm:"type:date" json:"checkout"`
	Nights      int       `json:"nights"`
	Guests      int       `json:"guests"`
	TotalAmount float64   `json:"total_amount"`
	Status      string    `gorm:"default:pending" json:"status"` // pending|confirmed|cancelled
	PaymentRef  string    `json:"payment_ref"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
