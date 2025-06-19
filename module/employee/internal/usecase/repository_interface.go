package usecase

import (
	"time"

	"github.com/eafajri/hr-service.git/module/employee/internal/entity"
)

//go:generate mockery --name UserRepository --output ./mocks
type UserRepository interface {
	GetUserByID(userID int64) (entity.User, error)
	GetUserSalaryByPeriodID(userID int64, periodID int64) (entity.UserSalary, error)
}

//go:generate mockery --name EmployeeRepository --output ./mocks
type EmployeeRepository interface {
	UpsertAttendance(record entity.EmployeeAttendance) error
	UpsertOvertime(record entity.EmployeeOvertime) error
	UpsertReimbursement(record entity.EmployeeReimbursement) error
}

//go:generate mockery --name PayrollPeriodRepository --output ./mocks
type PayrollRepository interface {
	GetPeriodByEntityDate(date time.Time) (entity.PayrollPeriod, error)
	GetPayslip(userID int64, periodID int64) (entity.PayrollPayslip, error)
	GetPayslips(periodID int64) ([]entity.PayrollPayslip, error)

	ClosePayrollPeriod(periodID int64) error
	CreatePayslipsByPeriod(payslips []entity.PayrollPayslip) error
}
