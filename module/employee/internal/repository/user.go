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

func (r *UserRepositoryImpl) GetUserSalaryByPeriodID(userID int64, periodID int64) (float64, error) {
	query := `
		SELECT us.amount
		FROM users_salaries us
		JOIN payroll_periods pp ON pp.id = $2
		WHERE us.user_id = $1
		AND us.effective_from <= pp.period_start
		ORDER BY us.effective_from DESC
		LIMIT 1;
	`

	var amount float64
	err := r.DB.Raw(query, userID, periodID).Scan(&amount).Error
	return amount, err
}
