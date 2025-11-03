package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/olamideolayemi/realestate-backend/internal/models"
)

type PropertyHandler struct {
	DB *gorm.DB
}

type CreatePropertyRequest struct {
	Title        string  `json:"title" binding:"required"`
	Description  string  `json:"description"`
	Category     string  `json:"category" binding:"required"` // buy|rent|shortlet
	Price        float64 `json:"price" binding:"required"`
	Currency     string  `json:"currency"`
	Address      string  `json:"address"`
	Area         string  `json:"area"`
	Bedrooms     int     `json:"bedrooms"`
	Bathrooms    int     `json:"bathrooms"`
	Furnished    bool    `json:"furnished"`
	PartyAllowed bool    `json:"party_allowed"`
	InstantBook  bool    `json:"instant_book"`
}

func (h *PropertyHandler) CreateProperty(c *gin.Context) {
	var req CreatePropertyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p := models.Property{
		ID:           uuid.New(),
		Title:        req.Title,
		Description:  req.Description,
		Category:     req.Category,
		Price:        req.Price,
		Currency:     req.Currency,
		Address:      req.Address,
		Area:         req.Area,
		Bedrooms:     req.Bedrooms,
		Bathrooms:    req.Bathrooms,
		Furnished:    req.Furnished,
		PartyAllowed: req.PartyAllowed,
		InstantBook:  req.InstantBook,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := h.DB.Create(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create property"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"property": p})
}

func (h *PropertyHandler) ListProperties(c *gin.Context) {
	// Basic filter implementation
	category := c.Query("category") // buy|rent|shortlet
	area := c.Query("area")
	minBeds := c.Query("min_beds")
	checkin := c.Query("checkin")
	checkout := c.Query("checkout")

	// Build query
	q := h.DB.Model(&models.Property{}).Preload("Images")
	if category != "" {
		q = q.Where("category = ?", category)
	}
	if area != "" {
		q = q.Where("area = ?", area)
	}
	if minBeds != "" {
		q = q.Where("bedrooms >= ?", minBeds)
	}

	// If category is shortlet and dates present, filter out properties with conflicting bookings
	if category == "shortlet" && checkin != "" && checkout != "" {
		q = q.Where("NOT EXISTS (SELECT 1 FROM bookings b WHERE b.property_id = properties.id AND b.status IN ('confirmed','pending') AND NOT (b.checkout <= ? OR b.checkin >= ?))", checkin, checkout)
	}

	var props []models.Property
	if err := q.Limit(100).Find(&props).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch properties"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"properties": props})
}

func (h *PropertyHandler) GetProperty(c *gin.Context) {
	id := c.Param("id")
	var p models.Property
	if err := h.DB.Preload("Images").First(&p, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "property not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"property": p})
}

func (h *PropertyHandler) UpdateProperty(c *gin.Context) {
	id := c.Param("id")
	var req CreatePropertyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var p models.Property
	if err := h.DB.First(&p, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "property not found"})
		return
	}
	p.Title = req.Title
	p.Description = req.Description
	p.Category = req.Category
	p.Price = req.Price
	p.Currency = req.Currency
	p.Address = req.Address
	p.Area = req.Area
	p.Bedrooms = req.Bedrooms
	p.Bathrooms = req.Bathrooms
	p.Furnished = req.Furnished
	p.PartyAllowed = req.PartyAllowed
	p.InstantBook = req.InstantBook
	p.UpdatedAt = time.Now()

	if err := h.DB.Save(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update property"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"property": p})
}

func (h *PropertyHandler) DeleteProperty(c *gin.Context) {
	id := c.Param("id")
	if err := h.DB.Delete(&models.Property{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete property"})
		return
	}
	c.Status(http.StatusNoContent)
}
