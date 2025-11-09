package handlers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/olamideolayemi/realestate-backend/internal/models"
	"github.com/olamideolayemi/realestate-backend/internal/utils"
)

type AuthHandler struct {
	DB *gorm.DB
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists
	var existing models.User
	if err := h.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "A user with this email already exists"})
		return
	}

	// Hash password
	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
		return
	}

	// Create user
	expiry := time.Now().Add(24 * time.Hour)
	u := models.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: hashed,
		Name:         req.Name,
		Role:         "user",
		IsVerified:   false,
		ExpiresAt:    &expiry,
	}

	if err := h.DB.Create(&u).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	// Generate OTP
	code := utils.GenerateOTP()
	verification := models.EmailVerification{
		Email:     req.Email,
		Code:      code,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}

	if err := h.DB.Create(&verification).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create verification record"})
		return
	}

	// Send OTP via email
	body := fmt.Sprintf("<p>Your verification code is <b>%s</b></p>", code)
	if err := utils.SendMail(req.Email, "Verify your account", body); err != nil {
		fmt.Println("Mail error:", err)
	}

	// create token
	// secret := os.Getenv("JWT_SECRET")
	// if secret == "" {
	// 	secret = "dev_secret"
	// }

	// token, err := utils.GenerateJWT(u.ID.String(), u.Role, secret, time.Hour*24*7)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
	// 	return
	// }

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration successful. Please verify your email.",
		// "user":  gin.H{"id": u.ID, "email": u.Email, "name": u.Name, "role": u.Role},
		// "token": token,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var u models.User
	if err := h.DB.Where("email = ?", req.Email).First(&u).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if !utils.CheckPasswordHash(req.Password, u.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if !u.IsVerified {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Please verify your email before logging in"})
		return
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev_secret"
	}
	token, err := utils.GenerateJWT(u.ID.String(), u.Role, secret, time.Hour*24*7)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user":  gin.H{"id": u.ID, "email": u.Email, "name": u.Name, "role": u.Role},
		"token": token,
	})
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
		Code  string `json:"code" binding:"required,len=6"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var record models.EmailVerification
	if err := h.DB.Where("email = ? AND code = ?", input.Email, input.Code).First(&record).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification code"})
		return
	}

	if time.Now().After(record.ExpiresAt) {
		h.DB.Delete(&record)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verification code expired"})
		return
	}

	var user models.User
	if err := h.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	// Check if user expired (not verified within 24h)
	if user.ExpiresAt != nil && time.Now().After(*user.ExpiresAt) {
		h.DB.Delete(&user)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Registration expired. Please register again."})
		return
	}

	// Mark verified
	user.IsVerified = true
	user.ExpiresAt = nil
	if err := h.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify user"})
		return
	}

	// Delete the verification record
	h.DB.Delete(&record)

	// Send success email
	successBody := fmt.Sprintf("<p>Hi %s,<br>Your email has been successfully verified. You can now log in to your account.</p>", user.Name)
	if err := utils.SendMail(user.Email, "Email Verification Successful", successBody); err != nil {
		fmt.Println("Mail error:", err)
	}

	// Generate token now
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev_secret"
	}
	token, err := utils.GenerateJWT(user.ID.String(), user.Role, secret, time.Hour*24*7)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Email verified successfully!",
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
			"role":  user.Role,
		},
		"token": token,
	})
}

func (h *AuthHandler) ResendVerification(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	if user.IsVerified {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is already verified"})
		return
	}

	// Delete existing verification record if any
	h.DB.Unscoped().Where("email = ?", input.Email).Delete(&models.EmailVerification{})

	// Generate new OTP
	code := utils.GenerateOTP()
	verification := models.EmailVerification{
		Email:     input.Email,
		Code:      code,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}

	if err := h.DB.Create(&verification).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create verification record"})
		return
	}

	// Send OTP via email
	body := fmt.Sprintf("<p>Your new verification code is <b>%s</b></p>", code)
	if err := utils.SendMail(input.Email, "Resend: Verify your account", body); err != nil {
		fmt.Println("Mail error:", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification code resent. Please check your email."})
}