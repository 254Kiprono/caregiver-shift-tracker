package main

import (
	"caregiver-shift-tracker/config"
	"caregiver-shift-tracker/controller"
	"caregiver-shift-tracker/database"
	"caregiver-shift-tracker/logger"
	"caregiver-shift-tracker/routes"
	"caregiver-shift-tracker/utils"

	_ "caregiver-shift-tracker/docs"

	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Caregiver Shift Tracker API
// @version 1.0
// @description API for Electronic Visit Verification and caregiver scheduling
// @contact.name Devs In Kenya
// @contact.url http://devsinkenya.com
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(gin.Recovery(), gin.Logger(), logger.Logger())

	cfg := config.LoadConfig()
	fmt.Printf("Loaded config successfully")

	utils.InitJWTConfig(cfg)

	// DB Init
	db, err := database.InitializeDB(cfg)
	if err != nil {
		fmt.Printf("DB init error: %v\n", err)
		logger.ErrorLogger.Fatalf("Failed to initialize database: %v", err)
	}
	fmt.Println("Database initialized.")

	database.RedisConn()
	rdb := database.RedisInstance()
	fmt.Println("Redis connected.")

	r.Use(database.DBMiddleware(db))

	authService := &controller.Controller{DB: db, RDB: rdb}
	routes.SetUpRoutes(r, authService, db)

	//Swagger endpoint
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
