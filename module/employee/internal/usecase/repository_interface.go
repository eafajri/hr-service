package usecase

import (
	"time"

	"github.com/eafajri/hr-service.git/module/employee/internal/entity"
)

//go:generate mockery --name UserRepository --output ./mocks
type UserRepository interface {
	GetUserByID(userID int64) (entity.User, error)
	GetUserByUsername(username string) (entity.User, error)
}

//go:generate mockery --name EmployeeRepository --output ./mocks
type EmployeeRepository interface {
	UpsertAttendance(record entity.EmployeeAttendance) error
	UpsertOvertime(record entity.EmployeeOvertime) error
	UpsertReimbursement(record entity.EmployeeReimbursement) error

	GetAttendanceByUserAndDate(userID int64, date time.Time) (entity.EmployeeAttendance, error)

	GetAllAttendanceByTimeRange(startTime time.Time, endTime time.Time, userID *int64) ([]entity.EmployeeAttendance, error)
	GetAllOvertimeByTimeRange(startTime time.Time, endTime time.Time, userID *int64) ([]entity.EmployeeOvertime, error)
	GetAllReimbursementByTimeRange(startTime time.Time, endTime time.Time, userID *int64) ([]entity.EmployeeReimbursement, error)

	GetEmployeeBaseSalaryByPeriodStart(periodStartTime time.Time, userID *int64) ([]entity.EmployeeBaseSalary, error)
}

//go:generate mockery --name PayrollRepository --output ./mocks
type PayrollRepository interface {
	GetPeriodByID(periodID int64) (entity.PayrollPeriod, error)
	GetPeriodByEntityDate(date time.Time) (entity.PayrollPeriod, error)
	GetPayslip(userID int64, periodID int64) (entity.PayrollPayslip, error)
	GetPayslips(periodID int64) ([]entity.PayrollPayslip, error)

	ClosePayrollPeriod(periodID int64) error
	CreatePayslipsByPeriod(payslips []entity.PayrollPayslip) error
}

//go:generate mockery --name AuditLogRepository --output ./mocks
type AuditLogRepository interface {
	Create(log entity.AuditLog, payload any) error
}
