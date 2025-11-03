package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/olamideolayemi/realestate-backend/internal/models"
)

type BookingHandler struct {
	DB *gorm.DB
}

type CreateBookingRequest struct {
	PropertyID string `json:"property_id" binding:"required,uuid"`
	Checkin    string `json:"checkin" binding:"required"`  // "YYYY-MM-DD"
	Checkout   string `json:"checkout" binding:"required"` // "YYYY-MM-DD"
	Guests     int    `json:"guests" binding:"required"`
}

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var req CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pid, err := uuid.Parse(req.PropertyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid property id"})
		return
	}

	// parse dates
	layout := "2006-01-02"
	checkin, err := time.Parse(layout, req.Checkin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid checkin date; use YYYY-MM-DD"})
		return
	}
	checkout, err := time.Parse(layout, req.Checkout)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid checkout date; use YYYY-MM-DD"})
		return
	}
	if !checkout.After(checkin) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "checkout must be after checkin"})
		return
	}

	// find property
	var prop models.Property
	if err := h.DB.First(&prop, "id = ?", pid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "property not found"})
		return
	}
	if prop.Category != "shortlet" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "property not bookable as shortlet"})
		return
	}

	// compute nights
	nights := int(checkout.Sub(checkin).Hours() / 24)
	if nights <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date range"})
		return
	}
	total := float64(nights) * prop.Price

	// Get user id from context (set by middleware)
	userIDval, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}
	userID, ok := userIDval.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user in context"})
		return
	}

	// DB transaction: availability check + create booking
	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
	}()

	var conflictCount int64
	if err := tx.Model(&models.Booking{}).
		Where("property_id = ? AND status IN ? AND NOT (checkout <= ? OR checkin >= ?)",
			prop.ID, []string{"confirmed", "pending"}, checkin, checkout).
		Count(&conflictCount).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed availability check"})
		return
	}
	if conflictCount > 0 {
		tx.Rollback()
		c.JSON(http.StatusConflict, gin.H{"error": "property not available for selected dates"})
		return
	}

	booking := models.Booking{
		ID:          uuid.New(),
		PropertyID:  prop.ID,
		UserID:      &userID,
		Checkin:     checkin,
		Checkout:    checkout,
		Nights:      nights,
		Guests:      req.Guests,
		TotalAmount: total,
		Status:      "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := tx.Create(&booking).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create booking"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit booking"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"booking": booking})
}

func (h *BookingHandler) ListUserBookings(c *gin.Context) {
	userIDval, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}
	userID, ok := userIDval.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user in context"})
		return
	}
	var bookings []models.Booking
	if err := h.DB.Where("user_id = ?", userID).Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch bookings"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"bookings": bookings})
}
