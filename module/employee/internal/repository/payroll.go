package repository

import (
	"time"

	"github.com/eafajri/hr-service.git/module/employee/internal/entity"
	"gorm.io/gorm"
)

type PayrollRepositoryImpl struct {
	DB *gorm.DB
}

func NewPayrollRepository(db *gorm.DB) *PayrollRepositoryImpl {
	return &PayrollRepositoryImpl{
		DB: db,
	}
}

func (r *PayrollRepositoryImpl) GetPeriodByID(periodID int64) (entity.PayrollPeriod, error) {
	var period entity.PayrollPeriod
	err := r.DB.First(&period, periodID).Error
	if err != nil {
		return entity.PayrollPeriod{}, err
	}
	return period, nil
}

func (r *PayrollRepositoryImpl) GetPeriodByEntityDate(date time.Time) (entity.PayrollPeriod, error) {
	var period entity.PayrollPeriod
	err := r.DB.Where("period_start <= ? AND period_end >= ?", date, date).First(&period).Error
	return period, err
}

func (r *PayrollRepositoryImpl) GetPayslip(userID int64, periodID int64) (entity.PayrollPayslip, error) {
	var payslip entity.PayrollPayslip
	err := r.DB.Where("user_id = ? AND payroll_period_id = ?", userID, periodID).First(&payslip).Error
	return payslip, err
}

func (r *PayrollRepositoryImpl) GetPayslips(periodID int64) ([]entity.PayrollPayslip, error) {
	var payslips []entity.PayrollPayslip
	err := r.DB.Where("period_id = ?", periodID).Find(&payslips).Error
	return payslips, err
}

func (r *PayrollRepositoryImpl) ClosePayrollPeriod(periodID int64) error {
	return r.DB.Exec("UPDATE payroll_periods SET status = 'closed' WHERE id = ?", periodID).Error
}

func (r *PayrollRepositoryImpl) CreatePayslipsByPeriod(payslips []entity.PayrollPayslip) error {
	return r.DB.CreateInBatches(payslips, 100).Error
}
