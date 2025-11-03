package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/olamideolayemi/realestate-backend/configs"
	"github.com/olamideolayemi/realestate-backend/internal/api"
	"github.com/olamideolayemi/realestate-backend/internal/models"
)

func main() {
	// Load .env for local development (optional)
	_ = godotenv.Load()

	db, err := configs.ConnectDB()
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	// Auto migrate models (dev convenience)
	if err := db.AutoMigrate(&models.User{}, &models.Property{}, &models.PropertyImage{}, &models.Booking{}); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}

	// Create one router instance
	r := gin.New()
	r.Use(gin.Recovery(), cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router := gin.Default()

	deps := api.NewDependencies(db)
	api.RegisterRoutes(router, deps)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
