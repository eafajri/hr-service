package entity

import "time"

type UserRole string

const (
	RoleEmployee UserRole = "employee"
	RoleAdmin    UserRole = "admin"
)

type User struct {
	ID       int64    `gorm:"primaryKey" json:"id"`
	Username string   `gorm:"username" json:"username"`
	Password string   `gorm:"password" json:"password"`
	Role     UserRole `gorm:"role" json:"role"`
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

type UserContext struct {
	RequestID string   `json:"request_id"`
	IPAddress string   `json:"ip_address"`
	UserID    int64    `json:"user_id"`
	Username  string   `json:"username"`
	Role      UserRole `json:"role"`
}
