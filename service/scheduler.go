package service

import (
	"caregiver-shift-tracker/models"
	"time"

	"gorm.io/gorm"
)

func GetAllSchedules(db *gorm.DB, userID int) ([]models.Schedule, error) {
	var schedules []models.Schedule
	err := db.Preload("Tasks").Where("user_id = ?", userID).Find(&schedules).Error
	return schedules, err
}

func GetTodaySchedules(db *gorm.DB, userID int) ([]models.Schedule, error) {
	start := time.Now().Truncate(24 * time.Hour)
	end := start.Add(24 * time.Hour)

	var schedules []models.Schedule
	err := db.Preload("Tasks").
		Where("user_id = ? AND shift_time BETWEEN ? AND ?", userID, start, end).
		Find(&schedules).Error
	return schedules, err
}

func GetScheduleByID(db *gorm.DB, scheduleID uint) (*models.Schedule, error) {
	var schedule models.Schedule
	err := db.Preload("Tasks").First(&schedule, "id = ?", scheduleID).Error
	return &schedule, err
}

func StartVisit(db *gorm.DB, scheduleID uint, lat, lon float64) error {
	now := time.Now()
	return db.Model(&models.Schedule{}).
		Where("id = ?", scheduleID).
		Updates(map[string]interface{}{
			"start_time": now,
			"start_lat":  lat,
			"start_lon":  lon,
			"status":     models.SCHEDULE_STATUS_IN_PROGRESS,
		}).Error
}

func EndVisit(db *gorm.DB, scheduleID uint, lat, lon float64) error {
	now := time.Now()
	return db.Model(&models.Schedule{}).
		Where("id = ?", scheduleID).
		Updates(map[string]interface{}{
			"end_time": now,
			"end_lat":  lat,
			"end_lon":  lon,
			"status":   models.SCHEDULE_STATUS_COMPLETED,
		}).Error
}
