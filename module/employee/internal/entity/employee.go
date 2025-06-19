package entity

import "time"

type EmployeeAttendance struct {
	ID           int64     `gorm:"primaryKey" json:"id"`
	UserID       int64     `gorm:"not null;index" json:"user_id"`
	Date         time.Time `gorm:"type:date;not null;uniqueIndex:idx_user_date" json:"date"`
	CheckInTime  time.Time `json:"check_in_time"`
	CheckOutTime time.Time `json:"check_out_time"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	UpdatedBy    int64     `gorm:"not null;index" json:"updated_by"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy    int64     `gorm:"not null;index" json:"created_by"`
}

func (EmployeeAttendance) TableName() string {
	return "employee_attendances"
}

type EmployeeOvertime struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	UserID    int64     `gorm:"not null;index" json:"user_id"`
	Date      time.Time `gorm:"type:date;not null;uniqueIndex:idx_user_date" json:"date"`
	Durations int       `gorm:"not null;check:durations > 0 AND durations <= 3" json:"durations"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	UpdatedBy    int64     `gorm:"not null;index" json:"updated_by"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy    int64     `gorm:"not null;index" json:"created_by"`
}

func (EmployeeOvertime) TableName() string {
	return "employee_overtimes"
}

type EmployeeReimbursement struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	UserID      int64     `gorm:"not null;index" json:"user_id"`
	Date        time.Time `gorm:"type:date;not null" json:"date"`
	Amount      float64   `gorm:"type:numeric(10,2);default:0" json:"amount"`
	Description string    `json:"description"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	UpdatedBy    int64     `gorm:"not null;index" json:"updated_by"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy    int64     `gorm:"not null;index" json:"created_by"`
}

func (EmployeeReimbursement) TableName() string {
	return "employee_reimbursements"
}
