package entity

import "time"

type EmployeeAttendance struct {
	ID           int64     `gorm:"id" json:"id"`
	UserID       int64     `gorm:"user_id" json:"user_id"`
	Date         time.Time `gorm:"date" json:"date"`
	CheckInTime  time.Time `gorm:"check_in_time" json:"check_in_time"`
	CheckOutTime time.Time `gorm:"check_out_time" json:"check_out_time"`
	UpdatedAt    time.Time `gorm:"updated_at" json:"updated_at"`
	UpdatedBy    string    `gorm:"updated_by" json:"updated_by"`
	CreatedAt    time.Time `gorm:"created_at" json:"created_at"`
	CreatedBy    string    `gorm:"created_by" json:"created_by"`
}

func (EmployeeAttendance) TableName() string {
	return "employee_attendances"
}

type EmployeeOvertime struct {
	ID        int64     `gorm:"id" json:"id"`
	UserID    int64     `gorm:"user_id" json:"user_id"`
	Date      time.Time `gorm:"date" json:"date"`
	Durations int       `gorm:"durations" json:"durations"`
	UpdatedAt time.Time `gorm:"updated_at" json:"updated_at"`
	UpdatedBy string    `gorm:"updated_by" json:"updated_by"`
	CreatedAt time.Time `gorm:"created_at" json:"created_at"`
	CreatedBy string    `gorm:"created_by" json:"created_by"`
}

func (EmployeeOvertime) TableName() string {
	return "employee_overtimes"
}

type EmployeeReimbursement struct {
	ID          int64     `gorm:"id" json:"id"`
	UserID      int64     `gorm:"user_id" json:"user_id"`
	Date        time.Time `gorm:"date" json:"date"`
	Amount      float64   `gorm:"amount" json:"amount"`
	Description string    `gorm:"description" json:"description"`
	UpdatedAt   time.Time `gorm:"updated_at" json:"updated_at"`
	UpdatedBy   string    `gorm:"updated_by" json:"updated_by"`
	CreatedAt   time.Time `gorm:"created_at" json:"created_at"`
	CreatedBy   string    `gorm:"created_by" json:"created_by"`
}

func (EmployeeReimbursement) TableName() string {
	return "employee_reimbursements"
}
