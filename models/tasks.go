package models

import (
	"time"
)

const (
	TASK_STATUS_COMPLETED     = "completed"
	TASK_STATUS_NOT_COMPLETED = "not_completed"
)

type Task struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	ScheduleID  uint       `gorm:"not null;index:idx_schedule_task" json:"schedule_id"`
	Description string     `gorm:"type:varchar(200);not null" json:"description" validate:"required"`
	Status      string     `gorm:"type:enum('completed','not_completed');default:'not_completed'" json:"status" validate:"required,oneof=completed not_completed"`
	Reason      *string    `gorm:"type:text" json:"reason,omitempty"`
	CompletedAt *time.Time `gorm:"type:datetime" json:"completed_at,omitempty"`
}
