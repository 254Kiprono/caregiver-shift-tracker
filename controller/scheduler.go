package controller

import (
	"caregiver-shift-tracker/models"
	"caregiver-shift-tracker/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (c *Controller) GetAllSchedules(ctx *gin.Context) {
	userID := ctx.GetInt("user_id")

	schedules, err := service.GetAllSchedules(c.DB, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch schedules"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"schedules": schedules})
}

func (ctrl *Controller) GetTodaySchedules(ctx *gin.Context) {
	userID := ctx.GetInt("user_id")

	schedules, err := service.GetTodaySchedules(ctrl.DB, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch today's schedules"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"schedules": schedules})
}

func (ctrl *Controller) GetScheduleDetails(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID"})
		return
	}

	schedule, err := service.GetScheduleByID(ctrl.DB, uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
		return
	}

	ctx.JSON(http.StatusOK, schedule)
}

func (ctrl *Controller) StartVisit(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID"})
		return
	}

	var req models.VisitLocationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid location data"})
		return
	}

	err = service.StartVisit(ctrl.DB, uint(id), req.Latitude, req.Longitude)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start visit"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Visit started"})
}

func (ctrl *Controller) EndVisit(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID"})
		return
	}

	var req models.VisitLocationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid location data"})
		return
	}

	err = service.EndVisit(ctrl.DB, uint(id), req.Latitude, req.Longitude)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to end visit"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Visit ended"})
}
