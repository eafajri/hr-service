package entity

import "time"

type UserRole string

const (
	RoleEmployee UserRole = "employee"
	RoleAdmin    UserRole = "admin"
)

type User struct {
	ID       int64    `gorm:"primaryKey" json:"id"`
	Username string   `gorm:"size:255;not null" json:"name"`
	Password string   `gorm:"not null" json:"password"`
	Role     UserRole `gorm:"type:user_role;not null" json:"role"`
}

func (EmployeeAttendance) User() string {
	return "users"
}

type UserSalary struct {
	ID            int64     `gorm:"primaryKey" json:"id"`
	UserID        int64     `gorm:"not null;index" json:"user_id"`
	Amount        float64   `gorm:"type:numeric(10,2);default:0" json:"amount"`
	EffectiveFrom time.Time `gorm:"autoCreateTime" json:"effective_from"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (EmployeeAttendance) UserSalary() string {
	return "user_salaries"
}
