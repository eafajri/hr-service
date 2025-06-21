package entity

import "time"

type PayrollPeriodStatus string

const (
	PayrollStatusOpen   PayrollPeriodStatus = "open"
	PayrollStatusClosed PayrollPeriodStatus = "closed"
)

type EmployeeBaseSalary struct {
	UserID     int64   `json:"user_id"`
	BaseSalary float64 `json:"base_salary"`
}

type PayrollPeriod struct {
	ID          int64               `gorm:"primaryKey" json:"id"`
	PeriodStart time.Time           `gorm:"type:date;not null" json:"period_start"`
	PeriodEnd   time.Time           `gorm:"type:date;not null" json:"period_end"`
	WorkingDays int                 `gorm:"not null" json:"working_days"`
	Status      PayrollPeriodStatus `gorm:"type:payroll_periods_status;default:'open'" json:"status"`
	UpdatedAt   time.Time           `gorm:"updated_at" json:"updated_at"`
	UpdatedBy   string              `gorm:"updated_by" json:"updated_by"`
	CreatedAt   time.Time           `gorm:"created_at" json:"created_at"`
	CreatedBy   string              `gorm:"created_by" json:"created_by"`
}

func (PayrollPeriod) TableName() string {
	return "payroll_periods"
}

type PayrollPayslip struct {
	ID                 int64     `gorm:"id" json:"id"`
	UserID             int64     `gorm:"user_id" json:"user_id"`
	PayrollPeriodID    int64     `gorm:"payroll_period_id" json:"payroll_period_id"`
	BaseSalary         float64   `gorm:"base_salary" json:"base_salary"`
	AttendanceDays     int       `gorm:"attendance_days" json:"attendance_days"`
	AttendanceHours    int       `gorm:"attendance_hours" json:"attendance_hours"`
	AttendancePay      float64   `gorm:"attendance_pay" json:"attendance_pay"`
	OvertimeHours      int       `gorm:"overtime_hours" json:"overtime_hours"`
	OvertimePay        float64   `gorm:"overtime_pay" json:"overtime_pay"`
	ReimbursementTotal float64   `gorm:"reimbursement_total" json:"reimbursement_total"`
	TotalTakeHome      float64   `gorm:"total_take_home" json:"total_take_home"`
	CreatedAt          time.Time `gorm:"created_at" json:"created_at"`
	CreatedBy          string    `gorm:"created_by" json:"created_by"`
}

func (PayrollPayslip) TableName() string {
	return "payroll_payslips"
}

func (p *PayrollPayslip) GeneratePayslip(periodDetail PayrollPeriod, baseSalaryDetail EmployeeBaseSalary, attendanceRecords []EmployeeAttendance, overtimeRecords []EmployeeOvertime, reimbursementRecords []EmployeeReimbursement) {
	p.UserID = baseSalaryDetail.UserID
	p.PayrollPeriodID = periodDetail.ID
	p.BaseSalary = baseSalaryDetail.BaseSalary

	ratePerDay := baseSalaryDetail.BaseSalary / float64(periodDetail.WorkingDays)
	ratePerHour := ratePerDay / 8

	p.AttendanceDays = len(attendanceRecords)
	for _, record := range attendanceRecords {
		duration := record.CheckOutTime.Sub(record.CheckInTime)
		p.AttendanceHours += int(duration.Hours())
	}
	p.AttendancePay = float64(p.AttendanceHours) * ratePerHour

	for _, record := range overtimeRecords {
		p.OvertimeHours += record.Durations
	}
	p.OvertimePay = float64(p.OvertimeHours) * ratePerHour * 2

	for _, record := range reimbursementRecords {
		p.ReimbursementTotal += record.Amount
	}

	p.TotalTakeHome = p.AttendancePay + p.OvertimePay + p.ReimbursementTotal
}
