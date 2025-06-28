package database

import (
	"caregiver-shift-tracker/config"
	"caregiver-shift-tracker/models"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DBMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}

// InitializeDB initializes the database connection
func InitializeDB(cfg *config.Config) (*gorm.DB, error) {
	// Build the Data Source Name (DSN) using the provided config values
	dsn := cfg.DBUsername + ":" + cfg.DBPassword + "@tcp(" + cfg.DBHost + ":" + cfg.DBPort + ")/" +
		cfg.DBName + "?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&loc=UTC&timeout=30s"

	// Open the database connection with the constructed DSN
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Perform automatic migration for the User model
	err = db.AutoMigrate(&models.User{}, models.Schedule{}, models.Task{})
	if err != nil {
		log.Fatalf("failed to auto-migrate User model: %v", err)
	}

	// Return the database connection object
	return db, nil
}
