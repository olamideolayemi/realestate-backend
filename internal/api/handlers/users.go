package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/olamideolayemi/realestate-backend/internal/models"
)

type UsersHandler struct {
	DB *gorm.DB
}

type UserResponse struct {
	ID			uuid.UUID	`json:"id"`
	Email     	string    	`json:"email"`
	Name      	string    	`json:"name"`
	Role   	 	string      `json:"role"`
	IsVerified  bool      	`json:"is_verified"`
	CreatedAt 	time.Time 	`json:"created_at"`
	UpdatedAt 	time.Time 	`json:"updated_at"`
}

// List all Users
func (h *UsersHandler) ListUsers(c *gin.Context) {
	page := 1
	limit := 10
	if p := c.Query("page"); p != "" {
		fmt.Sscan(p, "%d", &page)
		if page < 1 {
			page = 1
		}
	}
	if l := c.Query("limit"); l != "" {
		fmt.Sscan(l, "%d", &limit)
		if limit < 1 {
			limit = 10
		}
	}

	offset := (page - 1) * limit

	//Fetch total count for pagination
	var total int64
	h.DB.Model(&models.User{}).Count(&total)

	//Fetch paginated results
	var users []models.User
	if err := h.DB.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Build safe response with timestamps
	var response []UserResponse
	for _, u := range users {
		response = append(response, UserResponse{
			ID:        u.ID,
			Email:     u.Email,
			Name:      u.Name,
			Role:	   u.Role,
			IsVerified:  u.IsVerified,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"page": page,
		"limit": limit,
		"total": total,
		"users": response,
	})
}

func (h *UsersHandler) GetUser(c *gin.Context) {
	id := c.Param("id")

	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if err := h.DB.First(&user, "id = ?", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	response := UserResponse{
		ID:         user.ID,
		Email:      user.Email,
		Name:       user.Name,
		Role:       user.Role,
		IsVerified: user.IsVerified,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
	c.JSON(http.StatusOK, response)
}

func (h *UsersHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.DB.Delete(&models.User{}, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		// "user":    response,
		"message": "User account deleted",
	})
	c.Status(http.StatusNoContent)
}
