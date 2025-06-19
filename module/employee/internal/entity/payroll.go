package entity

import "time"

type PayrollPeriodStatus string

const (
	PayrollStatusOpen   PayrollPeriodStatus = "open"
	PayrollStatusClosed PayrollPeriodStatus = "closed"
)

type PayrollPeriod struct {
	ID          int64               `gorm:"primaryKey" json:"id"`
	PeriodStart time.Time           `gorm:"type:date;not null" json:"period_start"`
	PeriodEnd   time.Time           `gorm:"type:date;not null" json:"period_end"`
	WorkingDays int                 `gorm:"not null" json:"working_days"`
	Status      PayrollPeriodStatus `gorm:"type:payroll_periods_status;default:'open'" json:"status"`
	UpdatedAt   time.Time           `gorm:"autoUpdateTime" json:"updated_at"`
	UpdatedBy   int64               `gorm:"not null;index" json:"updated_by"`
	CreatedAt   time.Time           `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy   int64               `gorm:"not null;index" json:"created_by"`
}

func (PayrollPeriod) TableName() string {
	return "payroll_periods"
}

type PayrollPayslip struct {
	ID                 int64     `gorm:"primaryKey" json:"id"`
	UserID             int64     `gorm:"not null;index:idx_user_period,unique" json:"user_id"`
	PayrollPeriodID    int64     `gorm:"not null;index:idx_user_period,unique" json:"payroll_period_id"`
	BaseSalary         float64   `gorm:"type:numeric(10,2);not null" json:"base_salary"`
	AttendanceDays     int       `gorm:"not null" json:"attendance_days"`
	AttendanceHours    int       `gorm:"not null" json:"attendance_hours"`
	AttendacePay       float64   `gorm:"type:numeric(10,2);not null" json:"attendance_pay"`
	OvertimeHours      int       `gorm:"not null" json:"overtime_hours"`
	OvertimePay        float64   `gorm:"type:numeric(10,2);not null" json:"overtime_pay"`
	ReimbursementTotal float64   `gorm:"type:numeric(10,2);not null" json:"reimbursement_total"`
	TotalTakeHome      float64   `gorm:"type:numeric(10,2);not null" json:"total_take_home"`
	CreatedAt          time.Time `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy          int64     `gorm:"not null;index" json:"created_by"`
}

func (PayrollPayslip) TableName() string {
	return "payroll_payslips"
}

type EmployeeBaseSalary struct {
	UserID     int64   `json:"user_id"`
	BaseSalary float64 `json:"base_salary"`
}
