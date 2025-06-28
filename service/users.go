package service

import (
	"caregiver-shift-tracker/models"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RegisterUser(db *gorm.DB, newUser *models.User) (int, error) {
	var existingUser models.User
	if err := db.Where("email = ?", newUser.Email).First(&existingUser).Error; err == nil {
		return 0, gorm.ErrDuplicatedKey
	}

	if err := db.Create(newUser).Error; err != nil {
		return 0, err
	}

	return int(newUser.ID), nil
}

func LoginUser(db *gorm.DB, email, password string) (*models.User, error) {
	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	return &user, nil
}
