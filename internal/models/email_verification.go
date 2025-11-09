package models

import (
	"time"

	"gorm.io/gorm"
)

type EmailVerification struct {
	ID        uint           `gorm:"primaryKey"`
	Email     string         `gorm:"uniqueIndex"`
	Code      string         `gorm:"size:6"`
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
