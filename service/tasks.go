package service

import (
	"caregiver-shift-tracker/models"
	"time"

	"gorm.io/gorm"
)

// CreateTask creates a new task in the database
func CreateTask(db *gorm.DB, task *models.Task) error {
	return db.Create(task).Error
}

// AssignTasksToSchedule assigns a list of tasks to a specific schedule
func AssignTasksToSchedule(db *gorm.DB, scheduleID uint, tasks []models.Task) error {
	var schedule models.Schedule
	if err := db.First(&schedule, "id = ?", scheduleID).Error; err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for i := range tasks {
			tasks[i].ScheduleID = scheduleID
			if err := tx.Create(&tasks[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// DeleteTask deletes a specific task from the database
func DeleteTask(db *gorm.DB, taskID uint) error {
	return db.Delete(&models.Task{}, "id = ?", taskID).Error
}

// UpdateTask updates a specific task in the database
func UpdateTask(db *gorm.DB, task *models.Task) error {
	return db.Save(task).Error
}

// GetScheduleByTaskID retrieves the schedule associated with a task
func GetScheduleByTaskID(db *gorm.DB, taskID uint) (*models.Schedule, error) {
	var task models.Task
	if err := db.First(&task, "id = ?", taskID).Error; err != nil {
		return nil, err
	}
	var schedule models.Schedule
	if err := db.First(&schedule, "id = ?", task.ScheduleID).Error; err != nil {
		return nil, err
	}
	return &schedule, nil
}

// UpdateTaskStatus updates the status and reason of a task
func UpdateTaskStatus(db *gorm.DB, taskID uint, status string, reason *string, completedAt *time.Time) error {
	return db.Model(&models.Task{}).
		Where("id = ?", taskID).
		Updates(map[string]interface{}{
			"status":       status,
			"reason":       reason,
			"completed_at": completedAt,
		}).Error
}
