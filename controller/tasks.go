package controller

import (
	"caregiver-shift-tracker/logger"
	"caregiver-shift-tracker/models"
	"caregiver-shift-tracker/service"
	"caregiver-shift-tracker/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// CreateTask godoc
// @Summary Create a new task
// @Description Allows admin to create a task for a caregiver schedule
// @Tags Tasks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.Task true "Task Info"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks [post]
func (ctrl *Controller) CreateTask(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		logger.RespondRaw(ctx, http.StatusUnauthorized, gin.H{"error": "Authorization header missing or invalid"})
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	userID, roleID, err := utils.ExtractJWT(token, false)
	if err != nil {
		logger.RespondRaw(ctx, http.StatusUnauthorized, gin.H{"error": "Invalid or expired token", "details": err.Error()})
		return
	}
	if roleID != 1 {
		logger.RespondRaw(ctx, http.StatusForbidden, gin.H{"error": "Access denied: Admin role required"})
		return
	}
	var req models.Task
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.RespondRaw(ctx, http.StatusBadRequest, gin.H{"error": "Invalid task data", "details": err.Error()})
		return
	}
	if err := service.CreateTask(ctrl.DB, &req); err != nil {
		logger.ErrorLogger.Printf("Failed to create task for user %d: %v", userID, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task", "details": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "Task created successfully", "task_id": req.ID, "user_id": userID})
}

// AssignTasksToSchedule godoc
// @Summary Assign tasks to a schedule
// @Description Assigns one or more tasks to a specific schedule ID
// @Tags Tasks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Schedule ID"
// @Param request body object{tasks=[]models.Task} true "List of Tasks"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/assign/{id} [post]
func (ctrl *Controller) AssignTasksToSchedule(ctx *gin.Context) {
	idParam := ctx.Param("id")
	scheduleID, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID"})
		return
	}
	var req struct {
		Tasks []models.Task `json:"tasks" validate:"required,dive,required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task data", "details": err.Error()})
		return
	}
	err = service.AssignTasksToSchedule(ctrl.DB, uint(scheduleID), req.Tasks)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign tasks", "details": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Tasks assigned successfully"})
}

// DeleteTask godoc
// @Summary Delete a task
// @Description Deletes a task by its ID
// @Tags Tasks
// @Security BearerAuth
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{id} [delete]
func (ctrl *Controller) DeleteTask(ctx *gin.Context) {
	idParam := ctx.Param("id")
	taskID, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}
	err = service.DeleteTask(ctrl.DB, uint(taskID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task", "details": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}

// UpdateTask godoc
// @Summary Update a task
// @Description Updates task details by its ID
// @Tags Tasks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param request body models.Task true "Updated task data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{id} [put]
func (ctrl *Controller) UpdateTask(ctx *gin.Context) {
	idParam := ctx.Param("id")
	taskID, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}
	var req models.Task
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task data", "details": err.Error()})
		return
	}
	req.ID = uint(taskID)
	err = service.UpdateTask(ctrl.DB, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task", "details": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Task updated"})
}
