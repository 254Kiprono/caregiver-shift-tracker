package controller

import (
	"caregiver-shift-tracker/models"
	"caregiver-shift-tracker/service"
	"caregiver-shift-tracker/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// RegisterUser handles user registration
func (c *Controller) RegisterUser(ctx *gin.Context) {
	var req models.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Password encryption failed"})
		return
	}

	user := &models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		FullName: req.FullName,
		Mobile:   req.Mobile,
		RoleID:   models.ROLE_CAREGIVER,
	}

	if _, err := service.RegisterUser(c.DB, user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			ctx.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user", "details": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// LoginUser handles user login

func (c *Controller) LoginUser(ctx *gin.Context) {
	var req models.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login payload"})
		return
	}

	// Fetch user
	var user models.User
	if err := c.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate access and refresh tokens
	accessToken, refreshToken, err := utils.GenerateJWT(int(user.ID), user.RoleID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens", "details": err.Error()})
		return
	}

	// Store refresh token in DB (optional for future logout/refresh logic)
	user.RefreshToken = &refreshToken
	if err := c.DB.Save(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store refresh token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":       "Login successful",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": gin.H{
			"id":        user.ID,
			"email":     user.Email,
			"full_name": user.FullName,
			"mobile":    user.Mobile,
			"role_id":   user.RoleID,
		},
	})
}

func (c *Controller) RegAdmin(ctx *gin.Context) {
	var req models.RegisterUserRequest

	// Bind and validate request body
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Password encryption failed"})
		return
	}

	// Prepare admin user data
	user := &models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		FullName: req.FullName,
		Mobile:   req.Mobile,
		RoleID:   models.ROLE_ADMIN,
	}

	// Save to database
	if _, err := service.RegisterUser(c.DB, user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			ctx.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register admin", "details": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Admin registered successfully"})
}
