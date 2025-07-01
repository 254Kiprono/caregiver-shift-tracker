package service

import (
	"caregiver-shift-tracker/models"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// CreateSchedule adds a new schedule to the database
func CreateSchedule(db *gorm.DB, schedule *models.Schedule) error {
	return db.Create(schedule).Error
}

func GetAllSchedules(db *gorm.DB, userID int) ([]models.Schedule, error) {
	var schedules []models.Schedule
	err := db.Preload("Tasks").Where("user_id = ?", userID).Find(&schedules).Error
	return schedules, err
}

func GetTodaySchedules(db *gorm.DB, userID int, loc *time.Location) ([]models.Schedule, error) {
	now := time.Now().In(loc)
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).UTC()
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

func GetUpcomingSchedules(db *gorm.DB, userID int) ([]models.Schedule, error) {
	now := time.Now()

	var schedules []models.Schedule
	err := db.Preload("Tasks").
		Where("user_id = ? AND shift_time >= ?", userID, now).
		Where("status = ?", models.SCHEDULE_STATUS_SCHEDULED).
		Find(&schedules).Error

	return schedules, err
}

func GetMissedSchedules(db *gorm.DB, userID int, loc *time.Location) ([]models.Schedule, error) {
	nowInUserTZ := time.Now().In(loc)
	nowUTC := nowInUserTZ.UTC()
	gracePeriod := 5 * time.Minute
	cutoffTime := nowUTC.Add(-gracePeriod)

	startOfDayLocal := time.Date(nowInUserTZ.Year(), nowInUserTZ.Month(), nowInUserTZ.Day(), 0, 0, 0, 0, loc)
	endOfDayLocal := time.Date(nowInUserTZ.Year(), nowInUserTZ.Month(), nowInUserTZ.Day(), 23, 59, 59, 0, loc)
	startOfDayUTC := startOfDayLocal.UTC()
	endOfDayUTC := endOfDayLocal.UTC()

	var schedules []models.Schedule
	err := db.Debug().Preload("Tasks").
		Where("user_id = ? AND end_time BETWEEN ? AND ? AND end_time < ? AND status IN (?)",
			userID,
			startOfDayUTC, endOfDayUTC, cutoffTime,
			[]string{
				models.SCHEDULE_STATUS_SCHEDULED,
				models.SCHEDULE_STATUS_IN_PROGRESS,
			}).
		Find(&schedules).Error

	if err != nil {
		return nil, err
	}

	var missedSchedules []models.Schedule
	for i := range schedules {
		s := &schedules[i]

		if s.Status != models.SCHEDULE_STATUS_MISSED {
			updateErr := db.Model(&s).Update("status", models.SCHEDULE_STATUS_MISSED).Error
			if updateErr != nil {
				fmt.Printf("WARNING: Failed to update status for Schedule ID %d: %v\n", s.ID, updateErr)
			} else {
				s.Status = models.SCHEDULE_STATUS_MISSED
			}
		}
		missedSchedules = append(missedSchedules, *s)
	}

	for i := range missedSchedules {
		missedSchedules[i].ShiftTime = missedSchedules[i].ShiftTime.In(loc)
		if missedSchedules[i].StartTime != nil {
			start := missedSchedules[i].StartTime.In(loc)
			missedSchedules[i].StartTime = &start
		}
		if missedSchedules[i].EndTime != nil {
			end := missedSchedules[i].EndTime.In(loc)
			missedSchedules[i].EndTime = &end
		}
		for j := range missedSchedules[i].Tasks {
			if missedSchedules[i].Tasks[j].CompletedAt != nil {
				completed := missedSchedules[i].Tasks[j].CompletedAt.In(loc)
				missedSchedules[i].Tasks[j].CompletedAt = &completed
			}
		}
	}

	return missedSchedules, nil
}

func GetTodayCompletedSchedules(db *gorm.DB, userID int, loc *time.Location) ([]models.Schedule, error) {
	start := time.Now().UTC().Truncate(24 * time.Hour)
	end := start.Add(24 * time.Hour)
	var schedules []models.Schedule
	err := db.Preload("Tasks").
		Where("user_id = ? AND shift_time BETWEEN ? AND ? AND status = ?", userID, start, end, models.SCHEDULE_STATUS_COMPLETED).
		Find(&schedules).Error
	return schedules, err
}

func CancelStartVisit(db *gorm.DB, scheduleID uint) error {
	return db.Model(&models.Schedule{}).
		Where("id = ?", scheduleID).
		Updates(map[string]interface{}{
			"start_time": nil,
			"start_lat":  nil,
			"start_lon":  nil,
			"status":     models.SCHEDULE_STATUS_SCHEDULED,
		}).Error
}
