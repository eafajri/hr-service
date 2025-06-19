package usecase

import (
	"errors"
	"log"
	"time"

	"github.com/eafajri/hr-service.git/module/employee/internal/entity"
	"go.uber.org/zap"
)

//go:generate mockery --name ProfileUseCase --output ./mocks
type EmployeeUseCase interface {
	SubmitAttendance(userID int64, date time.Time, checkInTime, checkOutTime *time.Time) error
	SubmitOvertime(userID int64, date time.Time, durations int) error
	SubmitReimbursement(userID int64, date time.Time, amount float64, description string) error
}

type EmployeeUseCaseImpl struct {
	employeeRepository EmployeeRepository
	userUseCase        UserUseCase
}

func NewEmployeeUseCase(
	employeeRepository EmployeeRepository,
	userUseCase UserUseCase,
) *EmployeeUseCaseImpl {
	return &EmployeeUseCaseImpl{
		employeeRepository: employeeRepository,
		userUseCase:        userUseCase,
	}
}

func (e *EmployeeUseCaseImpl) SubmitAttendance(userID int64, date time.Time, checkInTime time.Time, checkOutTime time.Time) error {
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
		UserID:       userID,
		Date:         date,
		CheckInTime:  checkInTime,
		CheckOutTime: checkOutTime,
	}

	err := e.employeeRepository.UpsertAttendance(attendance)
	if err != nil {
		log.Println(
			"error when UpsertAttendance",
			zap.String("method", "EmployeeUseCaseImpl.SubmitAttendance"),
			zap.Int64("user_id", userID),
			zap.Time("date", date),
			zap.Time("check_in_time", checkInTime),
			zap.Time("check_out_time", checkOutTime),
			zap.Error(err),
		)
		return err
	}

	return nil
}

func (e *EmployeeUseCaseImpl) SubmitOvertime(userID int64, date time.Time, durations int) error {
	// When weekdays, It need to ensure that attendance is submitted
	shouldCheckAttendance := date.Weekday() != time.Saturday && date.Weekday() != time.Sunday
	if shouldCheckAttendance {
		attendance, err := e.employeeRepository.GetAttendanceByUserAndDate(userID, date)
		if err != nil {
			log.Println(
				"error when GetAttendanceByUserAndDate",
				zap.String("method", "EmployeeUseCaseImpl.SubmitOvertime"),
				zap.Int64("user_id", userID),
				zap.Time("date", date),
				zap.Error(err),
			)
			return err
		}

		if attendance.ID == 0 {
			return errors.New("attendance must be submitted before submitting overtime")
		}
	}

	// Only accept durations between 1 and 3 hours
	if durations < 1 || durations > 3 {
		return errors.New("overtime durations must be between 1 and 3 hours")
	}

	overtime := entity.EmployeeOvertime{
		UserID:    userID,
		Date:      date,
		Durations: durations,
	}

	err := e.employeeRepository.UpsertOvertime(overtime)
	if err != nil {
		log.Println(
			"error when UpsertOvertime",
			zap.String("method", "EmployeeUseCaseImpl.SubmitOvertime"),
			zap.Int64("user_id", userID),
			zap.Time("date", date),
			zap.Int("durations", durations),
			zap.Error(err),
		)
		return err
	}

	return nil
}

func (e *EmployeeUseCaseImpl) SubmitReimbursement(userID int64, date time.Time, amount float64, description string) error {
	reimbursement := entity.EmployeeReimbursement{
		UserID:      userID,
		Date:        date,
		Amount:      amount,
		Description: description,
	}

	err := e.employeeRepository.UpsertReimbursement(reimbursement)
	if err != nil {
		log.Println(
			"error when UpsertReimbursement",
			zap.String("method", "EmployeeUseCaseImpl.SubmitReimbursement"),
			zap.Int64("user_id", userID),
			zap.Time("date", date),
			zap.Error(err),
		)
		return err
	}

	return nil
}
