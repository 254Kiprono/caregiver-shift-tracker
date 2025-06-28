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
	gin.SetMode(gin.ReleaseMode)

	// Set up Gin router
	r := gin.Default()
	r.Use(gin.Recovery(), gin.Logger(), logger.Logger())

	// Load DB config
	cfg := config.LoadConfig()

	utils.InitJWTConfig(cfg)

	// Initialize DB
	db, err := database.InitializeDB(cfg)
	if err != nil {
		logger.ErrorLogger.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize Redis connection
	database.RedisConn()
	rdb := database.RedisInstance()

	// Attach DB to context
	r.Use(database.DBMiddleware(db))

	// Set up the Service and route
	authService := &controller.Controller{DB: db, RDB: rdb}
	routes.SetUpRoutes(r, authService, db)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	port := os.Getenv("SYSTEM_PORT")
	if port == "" {
		port = "6000"
	}
	fmt.Println("Server starting at port: " + port)
	if err := r.Run(":" + port); err != nil {
		logger.ErrorLogger.Fatalf("Failed to start server: %v", err)
	}
}
