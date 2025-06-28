// @title Caregiver Shift Tracker API
// @version 1.0
// @description API for Electronic Visit Verification and caregiver scheduling
// @host localhost:6000
// @BasePath /
// @contact.name Kiprono Bera
// @contact.email kiprono@example.com

package main

import (
	"caregiver-shift-tracker/config"
	"caregiver-shift-tracker/controller"
	"caregiver-shift-tracker/database"
	_ "caregiver-shift-tracker/docs"
	"caregiver-shift-tracker/logger"
	"caregiver-shift-tracker/routes"
	"caregiver-shift-tracker/utils"

	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files" //
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	fmt.Println("🚀 Starting Caregiver Shift Tracker service...")

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(gin.Recovery(), gin.Logger(), logger.Logger())

	cfg := config.LoadConfig()
	fmt.Printf("✅ Loaded config: DB_HOST=%s | REDIS_HOST=%s\n", cfg.DBHost, cfg.RedisHost)

	utils.InitJWTConfig(cfg)

	// DB Init
	db, err := database.InitializeDB(cfg)
	if err != nil {
		fmt.Printf("❌ DB init error: %v\n", err)
		logger.ErrorLogger.Fatalf("Failed to initialize database: %v", err)
	}
	fmt.Println("✅ Database initialized.")

	database.RedisConn()
	rdb := database.RedisInstance()
	fmt.Println("✅ Redis connected.")

	r.Use(database.DBMiddleware(db))

	authService := &controller.Controller{DB: db, RDB: rdb}
	routes.SetUpRoutes(r, authService, db)

	// Swagger docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	fmt.Println("✅ Swagger route set.")

	// Start server
	port := os.Getenv("SYSTEM_PORT")
	if port == "" {
		port = "6000"
	}
	fmt.Println("✅ Server starting at port: " + port)
	if err := r.Run(":" + port); err != nil {
		logger.ErrorLogger.Fatalf("Failed to start server: %v", err)
	}
}
