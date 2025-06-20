package usecase

import (
	"errors"
	"log"
	"time"

	"github.com/eafajri/hr-service.git/module/employee/internal/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

//go:generate mockery --name ProfileUseCase --output ./mocks
type EmployeeUseCase interface {
	SubmitAttendance(userContex entity.User, request entity.SubmitAttendanceRequest) error
	SubmitOvertime(userContex entity.User, request entity.SubmitOvertimeRequest) error
	SubmitReimbursement(userContex entity.User, request entity.SubmitReimbursementRequest) error
}

type EmployeeUseCaseImpl struct {
	employeeRepository EmployeeRepository
	payrollRepository  PayrollRepository
}

func NewEmployeeUseCase(
	employeeRepository EmployeeRepository,
	payrollRepository PayrollRepository,
) *EmployeeUseCaseImpl {
	return &EmployeeUseCaseImpl{
		employeeRepository: employeeRepository,
		payrollRepository:  payrollRepository,
	}
}

func (e *EmployeeUseCaseImpl) SubmitAttendance(userContex entity.User, request entity.SubmitAttendanceRequest) error {
	attandanceDate, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		log.Println(
			"error when parsing Date",
			zap.String("method", "EmployeeUseCaseImpl.SubmitAttendance"),
			zap.Any("user_contex", userContex),
			zap.Any("request", request),
			zap.Error(err),
		)
		return errors.New("invalid date format, must be YYYY-MM-DD")
	}

	if !e.isPeriodActive(attandanceDate) {
		return errors.New("the attendance cannot be submitted because the payroll period is closed")
	}

	checkInTime, err := time.Parse(time.RFC3339, request.CheckInTime)
	if err != nil {
		log.Println(
			"error when parsing CheckInTime",
			zap.String("method", "EmployeeUseCaseImpl.SubmitAttendance"),
			zap.Any("user_contex", userContex),
			zap.Any("request", request),
			zap.Error(err),
		)
		return errors.New("invalid check-in time format")
	}

	checkOutTime, err := time.Parse(time.RFC3339, request.CheckOutTime)
	if err != nil {
		log.Println(
			"error when parsing CheckOutTime",
			zap.String("method", "EmployeeUseCaseImpl.SubmitAttendance"),
			zap.Any("user_contex", userContex),
			zap.Any("request", request),
			zap.Error(err),
		)
		return errors.New("invalid check-out time format")
	}

	// Ensure check-in and check-out times are in same day
	if checkInTime.Year() != checkOutTime.Year() ||
		checkInTime.Month() != checkOutTime.Month() ||
		checkInTime.Day() != checkOutTime.Day() {
		return errors.New("check-in and check-out times must be on the same day")
	}

	// Ensure check-in time is not in the future
	if checkOutTime.Before(checkInTime) {
		return errors.New("check-out time cannot be before check-in time")
	}

	// Ensure in weekday (Monday to Friday)
	if checkInTime.Weekday() == time.Saturday || checkInTime.Weekday() == time.Sunday {
		return errors.New("attendance can only be submitted on weekdays (Monday to Friday)")
	}

	attendance := entity.EmployeeAttendance{
		UserID:       request.UserID,
		Date:         attandanceDate,
		CheckInTime:  checkInTime,
		CheckOutTime: checkInTime,
		CreatedBy:    userContex.Username,
		UpdatedBy:    userContex.Username,
	}

	err = e.employeeRepository.UpsertAttendance(attendance)
	if err != nil {
		log.Println(
			"error when UpsertAttendance",
			zap.String("method", "EmployeeUseCaseImpl.SubmitAttendance"),
			zap.Any("user_contex", userContex),
			zap.Any("request", request),
			zap.Error(err),
		)
		return err
	}

	return nil
}

func (e *EmployeeUseCaseImpl) SubmitOvertime(userContex entity.User, request entity.SubmitOvertimeRequest) error {
	overtimeDate, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		log.Println(
			"error when parsing Date",
			zap.String("method", "EmployeeUseCaseImpl.SubmitAttendance"),
			zap.Any("user_contex", userContex),
			zap.Any("request", request),
			zap.Error(err),
		)
		return errors.New("invalid date format, must be YYYY-MM-DD")
	}

	if !e.isPeriodActive(overtimeDate) {
		return errors.New("the overtime cannot be submitted because the payroll period is closed")
	}

	// When weekdays, It need to ensure that attendance is submitted
	shouldCheckAttendance := overtimeDate.Weekday() != time.Saturday || overtimeDate.Weekday() != time.Sunday
	if shouldCheckAttendance {
		attendance, err := e.employeeRepository.GetAttendanceByUserAndDate(request.UserID, overtimeDate)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.New("attendance must be submitted before submitting overtime")
			}
			log.Println(
				"error when GetAttendanceByUserAndDate",
				zap.String("method", "EmployeeUseCaseImpl.SubmitOvertime"),
				zap.Any("user_contex", userContex),
				zap.Any("request", request),
				zap.Error(err),
			)
			return err
		}

		if attendance.ID == 0 {
			return errors.New("attendance must be submitted before submitting overtime")
		}
	}

	// Only accept durations between 1 and 3 hours
	if request.Durations < 1 || request.Durations > 3 {
		return errors.New("overtime durations must be between 1 and 3 hours")
	}

	overtime := entity.EmployeeOvertime{
		UserID:    request.UserID,
		Date:      overtimeDate,
		Durations: int(request.Durations),
		CreatedBy: userContex.Username,
		UpdatedBy: userContex.Username,
	}

	err = e.employeeRepository.UpsertOvertime(overtime)
	if err != nil {
		log.Println(
			"error when UpsertOvertime",
			zap.String("method", "EmployeeUseCaseImpl.SubmitOvertime"),
			zap.Any("user_contex", userContex),
			zap.Any("request", request),
			zap.Error(err),
		)
		return err
	}

	return nil
}

func (e *EmployeeUseCaseImpl) SubmitReimbursement(userContex entity.User, request entity.SubmitReimbursementRequest) error {
	reimbursementDate, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		log.Println(
			"error when parsing Date",
			zap.String("method", "EmployeeUseCaseImpl.SubmitAttendance"),
			zap.Any("user_contex", userContex),
			zap.Any("request", request),
			zap.Error(err),
		)
		return errors.New("invalid date format, must be YYYY-MM-DD")
	}

	if !e.isPeriodActive(reimbursementDate) {
		return errors.New("the reimbursement cannot be submitted because the payroll period is closed")
	}

	reimbursement := entity.EmployeeReimbursement{

		UserID:      request.UserID,
		Date:        reimbursementDate,
		Amount:      request.Amount,
		Description: request.Description,
		CreatedBy:   userContex.Username,
		UpdatedBy:   userContex.Username,
	}

	err = e.employeeRepository.UpsertReimbursement(reimbursement)
	if err != nil {
		log.Println(
			"error when UpsertReimbursement",
			zap.String("method", "EmployeeUseCaseImpl.SubmitReimbursement"),
			zap.Any("user_contex", userContex),
			zap.Any("request", request),
			zap.Error(err),
		)
		return err
	}

	return nil
}

func (e *EmployeeUseCaseImpl) isPeriodActive(date time.Time) bool {
	period, err := e.payrollRepository.GetPeriodByEntityDate(date)
	if err != nil {
		log.Println(
			"error when GetPeriodByEntityDate",
			zap.String("method", "EmployeeUseCaseImpl.isPeriodActive"),
			zap.Error(err),
		)
		return false
	}

	if period.ID == 0 {
		return false
	}

	if period.Status == "closed" {
		return false
	}

	return true
}
