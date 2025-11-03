package api

import (
	"gorm.io/gorm"

	"github.com/olamideolayemi/realestate-backend/internal/api/handlers"
)

type Dependencies struct {
	DB *gorm.DB

	AuthHandler     *handlers.AuthHandler
	PropertyHandler *handlers.PropertyHandler
	BookingHandler  *handlers.BookingHandler
	HealthHandler   *handlers.HealthHandler
}

func NewDependencies(db *gorm.DB) *Dependencies {
	deps := &Dependencies{DB: db}

	auth := &handlers.AuthHandler{DB: db}
	prop := &handlers.PropertyHandler{DB: db}
	book := &handlers.BookingHandler{DB: db}
	health := &handlers.HealthHandler{DB: db}

	deps.AuthHandler = auth
	deps.PropertyHandler = prop
	deps.BookingHandler = book
	deps.HealthHandler = health

	return deps
}
