package repository

import (
	"github.com/eafajri/hr-service.git/module/employee/internal/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EmployeeRepositoryImpl struct {
	DB *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) *EmployeeRepositoryImpl {
	return &EmployeeRepositoryImpl{
		DB: db,
	}
}

func (r *EmployeeRepositoryImpl) UpsertAttendance(attendance entity.EmployeeAttendance) error {
	return r.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}, {Name: "date"}}, // match on user_id + date
		DoUpdates: clause.AssignmentColumns([]string{
			"check_in_time", "check_out_time", "updated_at", "updated_by",
		}),
	}).Create(&attendance).Error
}

func (r *EmployeeRepositoryImpl) UpsertOvertime(overtime entity.EmployeeOvertime) error {
	return r.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}, {Name: "date"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"durations", "updated_at", "updated_by",
		}),
	}).Create(&overtime).Error
}

func (r *EmployeeRepositoryImpl) UpsertReimbursement(reimbursement entity.EmployeeReimbursement) error {
	return r.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}, {Name: "date"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"amount", "description", "updated_at", "updated_by",
		}),
	}).Create(&reimbursement).Error
}
