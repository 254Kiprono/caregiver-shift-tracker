package controller

import (
	"caregiver-shift-tracker/logger"
	"caregiver-shift-tracker/models"
	"caregiver-shift-tracker/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateSchedule godoc
// @Summary Create a schedule
// @Description Create a new schedule for a caregiver
// @Tags Schedules
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.Schedule true "Schedule Info"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/create/schedule [post]
func (ctrl *Controller) CreateSchedule(ctx *gin.Context) {
	_, err := GetUserIDFromJWT(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req models.Schedule
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	if req.StartTime != nil && req.EndTime != nil && req.StartTime.After(*req.EndTime) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Shift start time must be before end time"})
		return
	}

	if err := service.CreateSchedule(ctrl.DB, &req); err != nil {
		logger.ErrorLogger.Printf("Failed to create schedule: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create schedule", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":     "Schedule created successfully",
		"schedule_id": req.ID,
		"user_id":     req.UserID,
	})
}

// GetAllSchedules godoc
// @Summary Get all schedules
// @Description Fetch all schedules for the authenticated caregiver
// @Tags Schedules
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{} "List of schedules"
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/user/schedules [get]
func (ctrl *Controller) GetAllSchedules(ctx *gin.Context) {
	userID, err := GetUserIDFromJWT(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	schedules, err := service.GetAllSchedules(ctrl.DB, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch schedules"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"schedules": schedules})
}

// GetTodaySchedules godoc
// @Summary Get today's schedules
// @Description Fetch today's schedules for the authenticated caregiver
// @Tags Schedules
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{} "List of today's schedules"
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/user/schedules/today [get]
func (ctrl *Controller) GetTodaySchedules(ctx *gin.Context) {
	userID, err := GetUserIDFromJWT(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	schedules, err := service.GetTodaySchedules(ctrl.DB, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch today's schedules"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"schedules": schedules})
}

// GetScheduleDetails godoc
// @Summary Get schedule details
// @Description Fetch a specific schedule by ID for the authenticated caregiver
// @Tags Schedules
// @Security BearerAuth
// @Produce json
// @Param id path int true "Schedule ID"
// @Success 200 {object} models.Schedule "Schedule details"
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string "Not found"
// @Router /api/user/schedules/{id} [get]
func (ctrl *Controller) GetScheduleDetails(ctx *gin.Context) {
	userID, err := GetUserIDFromJWT(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
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
	if schedule.UserID != uint(userID) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Schedule not assigned to user"})
		return
	}
	ctx.JSON(http.StatusOK, schedule)
}

// StartVisit godoc
// @Summary Start visit
// @Description Start a visit for a specific schedule by ID
// @Tags Schedules
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Schedule ID"
// @Param request body models.VisitLocationRequest true "Start location coordinates"
// @Success 200 {object} map[string]string "Visit started"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string "Server error"
// @Router /api/user/schedules/{id}/start [post]
func (ctrl *Controller) StartVisit(ctx *gin.Context) {
	userID, err := GetUserIDFromJWT(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
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
	if schedule.UserID != uint(userID) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Schedule not assigned to user"})
		return
	}
	var req models.VisitLocationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid location data"})
		return
	}
	if req.Latitude < -90 || req.Latitude > 90 || req.Longitude < -180 || req.Longitude > 180 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude or longitude"})
		return
	}
	err = service.StartVisit(ctrl.DB, uint(id), req.Latitude, req.Longitude)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start visit"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Visit started"})
}

// EndVisit godoc
// @Summary End visit
// @Description End a visit for a specific schedule by ID
// @Tags Schedules
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Schedule ID"
// @Param request body models.VisitLocationRequest true "End location coordinates"
// @Success 200 {object} map[string]string "Visit ended"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string "Server error"
// @Router /api/user/schedules/{id}/end [post]
func (ctrl *Controller) EndVisit(ctx *gin.Context) {
	userID, err := GetUserIDFromJWT(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
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
	if schedule.UserID != uint(userID) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Schedule not assigned to user"})
		return
	}
	var req models.VisitLocationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid location data"})
		return
	}
	if req.Latitude < -90 || req.Latitude > 90 || req.Longitude < -180 || req.Longitude > 180 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude or longitude"})
		return
	}
	err = service.EndVisit(ctrl.DB, uint(id), req.Latitude, req.Longitude)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to end visit"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Visit ended"})
}

// GetUpcomingSchedules godoc
// @Summary Get upcoming schedules
// @Description Fetch all upcoming schedules for the authenticated caregiver (from today onward)
// @Tags Schedules
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{} "List of upcoming schedules"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/user/schedules/upcoming [get]
func (ctrl *Controller) GetUpcomingSchedules(ctx *gin.Context) {
	userID, err := GetUserIDFromJWT(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	schedules, err := service.GetUpcomingSchedules(ctrl.DB, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch upcoming schedules"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"schedules": schedules})
}

// GetMissedSchedules godoc
// @Summary Get missed schedules
// @Description Fetch all missed schedules for the authenticated caregiver (end time passed and not completed)
// @Tags Schedules
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{} "List of missed schedules"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/user/schedules/missed [get]
func (ctrl *Controller) GetMissedSchedules(ctx *gin.Context) {
	userID, err := GetUserIDFromJWT(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	schedules, err := service.GetMissedSchedules(ctrl.DB, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch missed schedules"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"schedules": schedules})
}

// GetTodayCompletedSchedules godoc
// @Summary Get today's completed schedules
// @Description Fetch all completed schedules for today for the authenticated caregiver
// @Tags Schedules
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{} "List of today's completed schedules"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/user/schedules/completed/today [get]
func (ctrl *Controller) GetTodayCompletedSchedules(ctx *gin.Context) {
	userID, err := GetUserIDFromJWT(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	schedules, err := service.GetTodayCompletedSchedules(ctrl.DB, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch completed schedules"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"schedules": schedules})
}

// CancelStartVisit godoc
// @Summary Cancel start visit (undo clock-in)
// @Description Allows caregiver to cancel their clock-in (reset start time and location)
// @Tags Schedules
// @Security BearerAuth
// @Produce json
// @Param id path int true "Schedule ID"
// @Success 200 {object} map[string]string "Clock-in canceled"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/user/schedules/{id}/cancel-start [post]
func (ctrl *Controller) CancelStartVisit(ctx *gin.Context) {
	userID, err := GetUserIDFromJWT(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	idParam := ctx.Param("id")
	scheduleID, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID"})
		return
	}

	schedule, err := service.GetScheduleByID(ctrl.DB, uint(scheduleID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
		return
	}

	if schedule.UserID != uint(userID) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Schedule not assigned to user"})
		return
	}

	err = service.CancelStartVisit(ctrl.DB, uint(scheduleID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel clock-in"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Clock-in canceled successfully"})
}
