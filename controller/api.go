package controller

import (
	"caregiver-shift-tracker/utils"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Controller struct {
	DB  *gorm.DB
	GIN *gin.Engine
	RDB *redis.Client
}

// getUserIDFromJWT extracts user_id from JWT token
func GetUserIDFromJWT(ctx *gin.Context) (int, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return 0, fmt.Errorf("authorization header missing or invalid")
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	userID, _, err := utils.ExtractJWT(token, false)
	if err != nil {
		return 0, fmt.Errorf("invalid or expired token: %v", err)
	}
	return userID, nil
}

func GetUserIDAndRoleFromJWT(ctx *gin.Context) (int, int, error) {
	authHeader := ctx.GetHeader("authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return 0, 0, fmt.Errorf("authorization header missing or invalid")
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	userID, roleID, err := utils.ExtractJWT(token, false)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid or expired token: %v", err)
	}
	return userID, roleID, nil
}
