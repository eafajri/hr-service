package usecase

import (
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
