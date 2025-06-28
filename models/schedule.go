package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	SCHEDULE_STATUS_SCHEDULED   = "scheduled"
	SCHEDULE_STATUS_IN_PROGRESS = "in_progress"
	SCHEDULE_STATUS_COMPLETED   = "completed"
	SCHEDULE_STATUS_CANCELLED   = "cancelled"
)

type Schedule struct {
	gorm.Model
	UserID     uint       `gorm:"not null;index:idx_user_schedule" json:"user_id"`
	ClientName string     `gorm:"type:varchar(100);not null" json:"client_name" validate:"required"`
	Location   string     `gorm:"type:varchar(200);not null" json:"location" validate:"required"`
	ShiftTime  time.Time  `gorm:"type:datetime;not null;index" json:"shift_time" validate:"required"`
	Status     string     `gorm:"type:enum('scheduled','in_progress','completed','cancelled');default:'scheduled'" json:"status" validate:"required,oneof=scheduled in_progress completed cancelled"`
	StartTime  *time.Time `gorm:"type:datetime" json:"start_time"`
	EndTime    *time.Time `gorm:"type:datetime" json:"end_time"`
	StartLat   *float64   `gorm:"type:decimal(10,8)" json:"start_lat"`
	StartLon   *float64   `gorm:"type:decimal(11,8)" json:"start_lon"`
	EndLat     *float64   `gorm:"type:decimal(10,8)" json:"end_lat"`
	EndLon     *float64   `gorm:"type:decimal(11,8)" json:"end_lon"`

	Tasks []Task `gorm:"foreignKey:ScheduleID;constraint:OnDelete:CASCADE" json:"tasks"`
}

type VisitLocationRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}
