package models

import (
	"time"
)

const (
	ROLE_ADMIN         = 1
	ROLE_CUSTOMER_CARE = 2
	ROLE_CAREGIVER     = 3
)

// User represents an authenticated user
type User struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	Email        string     `gorm:"type:varchar(100);unique;not null;index" json:"email" validate:"required,email"`
	Mobile       string     `gorm:"type:varchar(100);unique;not null;index" json:"mobile" validate:"required"`
	FullName     string     `gorm:"type:varchar(100);not null" json:"full_name" validate:"required"`
	Password     string     `gorm:"type:varchar(100);not null" json:"password" validate:"required,min=8"`
	RoleID       int        `gorm:"not null;default:3;index" json:"role_id" validate:"required,oneof=1 2 3"`
	RefreshToken *string    `gorm:"type:text" json:"refresh_token,omitempty"`
}

type RegisterUserRequest struct {
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Mobile   string `json:"mobile" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
