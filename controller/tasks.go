package controller

import (
	"caregiver-shift-tracker/logger"
	"caregiver-shift-tracker/models"
	"caregiver-shift-tracker/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateTask godoc
// @Summary Create a new task
// @Description Creates a task for a caregiver schedule
// @Tags Tasks
// @Accept json
// @Produce json
// @Param request body models.Task true "Task Info"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks [post]
func (ctrl *Controller) CreateTask(ctx *gin.Context) {
	var req models.Task
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.RespondRaw(ctx, http.StatusBadRequest, gin.H{"error": "Invalid task data", "details": err.Error()})
		return
	}
	if err := service.CreateTask(ctrl.DB, &req); err != nil {
		logger.ErrorLogger.Printf("Failed to create task: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task", "details": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "Task created successfully", "task_id": req.ID})
}

// AssignTasksToSchedule godoc
// @Summary Assign tasks to a schedule
// @Description Assign one or more tasks to a specific schedule ID
// @Tags Tasks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Schedule ID"
// @Param request body object{tasks=[]models.Task} true "List of Tasks"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/assign/{id} [post]
func (ctrl *Controller) AssignTasksToSchedule(ctx *gin.Context) {
	_, err := GetUserIDFromJWT(ctx)
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
// @Description Updates task details by its ID, restricted to the assigned caregiver
// @Tags Tasks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param request body models.Task true "Updated task data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{id} [put]
func (ctrl *Controller) UpdateTask(ctx *gin.Context) {
	userID, err := GetUserIDFromJWT(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
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
	if req.Status != models.TASK_STATUS_COMPLETED && req.Status != models.TASK_STATUS_NOT_COMPLETED {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status: must be 'completed' or 'not_completed'"})
		return
	}
	if req.Status == models.TASK_STATUS_NOT_COMPLETED && (req.Reason == nil || *req.Reason == "") {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Reason required for not_completed status"})
		return
	}
	schedule, err := service.GetScheduleByTaskID(ctrl.DB, uint(taskID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task or schedule not found"})
		return
	}
	if schedule.UserID != uint(userID) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Task not assigned to user"})
		return
	}
	req.ID = uint(taskID)
	if req.Status == models.TASK_STATUS_COMPLETED {
		now := time.Now()
		req.CompletedAt = &now
	}
	err = service.UpdateTask(ctrl.DB, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task", "details": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Task updated"})
}

// UpdateTaskStatus godoc
// @Summary Update task status
// @Description Updates the status of a task (completed or not_completed with reason), restricted to the assigned caregiver
// @Tags Tasks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param taskId path int true "Task ID"
// @Param request body object{status=string,reason=string} true "Task status and optional reason"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{taskId}/update [post]
func (ctrl *Controller) UpdateTaskStatus(ctx *gin.Context) {
	userID, err := GetUserIDFromJWT(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	idParam := ctx.Param("taskId")
	taskID, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}
	var req struct {
		Status string  `json:"status" validate:"required,oneof=completed not_completed"`
		Reason *string `json:"reason"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}
	if req.Status != models.TASK_STATUS_COMPLETED && req.Status != models.TASK_STATUS_NOT_COMPLETED {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status: must be 'completed' or 'not_completed'"})
		return
	}
	if req.Status == models.TASK_STATUS_NOT_COMPLETED && (req.Reason == nil || *req.Reason == "") {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Reason required for not_completed status"})
		return
	}
	schedule, err := service.GetScheduleByTaskID(ctrl.DB, uint(taskID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task or schedule not found"})
		return
	}
	if schedule.UserID != uint(userID) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Task not assigned to user"})
		return
	}
	var completedAt *time.Time
	if req.Status == models.TASK_STATUS_COMPLETED {
		now := time.Now()
		completedAt = &now
	}
	err = service.UpdateTaskStatus(ctrl.DB, uint(taskID), req.Status, req.Reason, completedAt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task status", "details": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Task status updated"})
}
