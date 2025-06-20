package repository

import (
	"github.com/eafajri/hr-service.git/module/employee/internal/entity"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		DB: db,
	}
}

func (r *UserRepositoryImpl) GetUserByID(userID int64) (entity.User, error) {
	var user entity.User
	err := r.DB.First(&user, userID).Error

	return user, err
}

func (r *UserRepositoryImpl) GetUserByUsername(username string) (entity.User, error) {
	var user entity.User
	err := r.DB.Where("username = ?", username).First(&user).Error

	return user, err
}
