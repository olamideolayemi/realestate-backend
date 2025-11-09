package api

import (
	"github.com/gin-gonic/gin"
	"github.com/olamideolayemi/realestate-backend/internal/api/middleware"
)

func RegisterRoutes(r *gin.Engine, deps *Dependencies) {
	api := r.Group("/api/v1")

	// API Health Check
	api.GET("/health", deps.HealthHandler.Health)

	// Auth
	api.POST("/auth/register", deps.AuthHandler.Register)
	api.POST("/auth/login", deps.AuthHandler.Login)
	api.POST("/auth/verify", deps.AuthHandler.VerifyEmail)
	api.POST("/auth/resend-ver", deps.AuthHandler.ResendVerification)

	// Properties
	props := api.Group("/properties")
	props.GET("", deps.PropertyHandler.ListProperties)
	props.GET("/:id", deps.PropertyHandler.GetProperty)

	// Admin routes (protect with auth + admin check)
	admin := api.Group("/admin", middleware.AuthMiddleware(deps.DB), middleware.AdminOnly())
	admin.POST("/properties", deps.PropertyHandler.CreateProperty)
	admin.PATCH("/properties/:id", deps.PropertyHandler.UpdateProperty)
	admin.DELETE("/properties/:id", deps.PropertyHandler.DeleteProperty)
	// Users (Admin)
	admin.GET("/users", deps.UsersHandler.ListUsers)
	admin.GET("/users/:id", deps.UsersHandler.GetUser)
	admin.DELETE("/users/:id", deps.UsersHandler.DeleteUser)

	// Bookings
	api.POST("/bookings", middleware.AuthMiddleware(deps.DB), deps.BookingHandler.CreateBooking)
	api.GET("/bookings", middleware.AuthMiddleware(deps.DB), deps.BookingHandler.ListUserBookings)
}
