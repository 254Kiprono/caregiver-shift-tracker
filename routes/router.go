package routes

import (
	"caregiver-shift-tracker/controller"
	"caregiver-shift-tracker/logger"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetUpRoutes(r *gin.Engine, ctrl *controller.Controller, DB *gorm.DB) {
	allowedMethods := []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete}
	AllowOrigins := []string{"*"}

	// CORS
	corsConfig := cors.Config{
		AllowOrigins: AllowOrigins,
		AllowHeaders: AllowOrigins,
		AllowMethods: allowedMethods,
	}
	r.Use(cors.New(corsConfig))

	r.GET("/status", func(c *gin.Context) {
		logger.InfoLogger.Println("System Health Status Check Successful")
		c.JSON(http.StatusOK, gin.H{"message": "System Health Status Check Successful"})
	})

	userRoutes := r.Group("/api")
	{
		userRoutes.POST("/user/register", ctrl.RegisterUser)
		userRoutes.POST("/login", ctrl.LoginUser)
		userRoutes.GET("/user/schedules/today", ctrl.GetTodaySchedules)
		userRoutes.GET("/user/schedules", ctrl.GetAllSchedules)
		userRoutes.GET("/user/schedules/:id", ctrl.GetScheduleDetails)
		userRoutes.POST("/user/schedules/:id/start", ctrl.StartVisit)
		userRoutes.POST("/user/schedules/:id/end", ctrl.EndVisit)
	}

}
