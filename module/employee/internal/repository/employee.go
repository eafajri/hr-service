package repository

import (
	"time"

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

func (r *EmployeeRepositoryImpl) GetAllAttendanceByTimeRange(startTime time.Time, endTime time.Time, userID *int64) ([]entity.EmployeeAttendance, error) {
	var attendances []entity.EmployeeAttendance
	query := r.DB.Where("date BETWEEN ? AND ?", startTime, endTime)

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	err := query.Find(&attendances).Error
	if err != nil {
		return nil, err
	}
	return attendances, nil
}

func (r *EmployeeRepositoryImpl) GetAllOvertimeByTimeRange(startTime time.Time, endTime time.Time, userID *int64) ([]entity.EmployeeOvertime, error) {
	var overtimes []entity.EmployeeOvertime
	query := r.DB.Where("date BETWEEN ? AND ?", startTime, endTime)

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	err := query.Find(&overtimes).Error
	if err != nil {
		return nil, err
	}

	return overtimes, nil
}

func (r *EmployeeRepositoryImpl) GetAllReimbursementByTimeRange(startTime time.Time, endTime time.Time, userID *int64) ([]entity.EmployeeReimbursement, error) {
	var reimbursements []entity.EmployeeReimbursement
	query := r.DB.Where("date BETWEEN ? AND ?", startTime, endTime)

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	err := query.Find(&reimbursements).Error
	if err != nil {
		return nil, err
	}

	return reimbursements, nil
}

func (r *EmployeeRepositoryImpl) GetEmployeeBaseSalaryByPeriodStart(periodStartTime time.Time, userID *int64) ([]entity.EmployeeBaseSalary, error) {
	var salaries []entity.EmployeeBaseSalary

	baseQuery := `
		SELECT DISTINCT ON (us.user_id) 
			us.user_id AS user_id, us.amount AS base_salary
		FROM user_salaries us
		WHERE us.effective_from <= ?
	`

	args := []interface{}{periodStartTime}

	if userID != nil {
		baseQuery += " AND us.user_id = ?"
		args = append(args, *userID)
	}

	baseQuery += " ORDER BY us.user_id, us.effective_from DESC;"

	err := r.DB.Raw(baseQuery, args...).Scan(&salaries).Error
	return salaries, err
}

// GetAttendanceByUserAndDate implements usecase.EmployeeRepository.
func (r *EmployeeRepositoryImpl) GetAttendanceByUserAndDate(userID int64, date time.Time) (entity.EmployeeAttendance, error) {
	var attendance entity.EmployeeAttendance
	err := r.DB.Where("user_id = ? AND date = ?", userID, date).First(&attendance).Error
	if err != nil {
		return entity.EmployeeAttendance{}, err
	}

	return attendance, nil
}
