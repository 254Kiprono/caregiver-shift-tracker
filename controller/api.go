package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Controller struct {
	DB  *gorm.DB
	GIN *gin.Engine
	RDB *redis.Client
}
