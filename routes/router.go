package routes

import (
	"caregiver-shift-tracker/controller"
	"caregiver-shift-tracker/logger"

	// "caregiver-shift-tracker/utils"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetUpRoutes(r *gin.Engine, ctrl *controller.Controller, DB *gorm.DB) {
	allowedMethods := []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete}
	allowHeaders := []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Timezone"}

	// CORS
	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     allowHeaders,
		AllowMethods:     allowedMethods,
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(corsConfig))

	// Health check
	r.GET("/status", func(c *gin.Context) {
		logger.InfoLogger.Println("System Health Status Check Successful")
		c.JSON(http.StatusOK, gin.H{"message": "System Health Status Check Successful"})
	})

	// Open Task routes
	admin := r.Group("/tasks")
	{
		admin.POST("/", ctrl.CreateTask)
		admin.POST("/assign/:id", ctrl.AssignTasksToSchedule)
		admin.DELETE("/:id", ctrl.DeleteTask)
		admin.PUT("/:id", ctrl.UpdateTask)
		admin.POST("/create/schedule", ctrl.CreateSchedule)
		admin.POST("/:taskId/update", ctrl.UpdateTaskStatus)
	}

	// Public API (no auth)
	userRoutes := r.Group("/api")
	{
		userRoutes.POST("/user/register", ctrl.RegisterUser)
		userRoutes.POST("/admin/register", ctrl.RegAdmin)
		userRoutes.POST("/login", ctrl.LoginUser)
	}

	protected := r.Group("/api")
	// protected.Use(utils.AuthMiddlewareForSwagger())
	{
		protected.GET("/user/schedules", ctrl.GetAllSchedules)
		protected.GET("/user/schedules/today", ctrl.GetTodaySchedules)
		protected.GET("/user/schedules/upcoming", ctrl.GetUpcomingSchedules)
		protected.GET("/user/schedules/missed", ctrl.GetMissedSchedules)
		protected.GET("/user/schedules/completed/today", ctrl.GetTodayCompletedSchedules)
		protected.GET("/user/schedules/:id", ctrl.GetScheduleDetails)
		protected.POST("/user/schedules/:id/start", ctrl.StartVisit)
		protected.POST("/user/schedules/:id/end", ctrl.EndVisit)
		protected.POST("/user/schedules/:id/cancel-start", ctrl.CancelStartVisit)
	}
}

//Adjust
