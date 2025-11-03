package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthHandler struct {
	DB *gorm.DB
}

func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(200, gin.H{"ok": true})
}
