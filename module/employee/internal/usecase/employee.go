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
	SubmitAttendance(userContext entity.UserContext, request entity.SubmitAttendanceRequest) error
	SubmitOvertime(userContext entity.UserContext, request entity.SubmitOvertimeRequest) error
	SubmitReimbursement(userContext entity.UserContext, request entity.SubmitReimbursementRequest) error

	GetPayslipBreakdown(userContext entity.UserContext, periodID int64) (any, error)
}

type EmployeeUseCaseImpl struct {
	employeeRepository EmployeeRepository
	payrollRepository  PayrollRepository
	auditLogRepository AuditLogRepository
}

func NewEmployeeUseCase(
	employeeRepository EmployeeRepository,
	payrollRepository PayrollRepository,
	auditLogRepository AuditLogRepository,
) *EmployeeUseCaseImpl {
	return &EmployeeUseCaseImpl{
		employeeRepository: employeeRepository,
		payrollRepository:  payrollRepository,
		auditLogRepository: auditLogRepository,
	}
}

/*
No rules for late or early check-ins or check-outs; check-in at any time that day counts.
Submissions on the same day should count as one.
Users cannot submit on weekends.
*/
func (e *EmployeeUseCaseImpl) SubmitAttendance(userContext entity.UserContext, request entity.SubmitAttendanceRequest) error {
	if userContext.UserID != request.UserID {
		return errors.New("user context does not match request user ID")
	}

	attandanceDate, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		log.Println(
			"error when parsing Date",
			zap.String("method", "EmployeeUseCaseImpl.SubmitAttendance"),
			zap.Any("user_contex", userContext),
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
			zap.Any("user_contex", userContext),
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
			zap.Any("user_contex", userContext),
			zap.Any("request", request),
			zap.Error(err),
		)
		return errors.New("invalid check-out time format")
	}

	if !e.isSameDay(checkInTime, checkOutTime) || !e.isSameDay(checkInTime, attandanceDate) {
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
		CheckOutTime: checkOutTime,
		CreatedBy:    userContext.Username,
		UpdatedBy:    userContext.Username,
	}

	err = e.employeeRepository.UpsertAttendance(attendance)
	if err != nil {
		log.Println(
			"error when UpsertAttendance",
			zap.String("method", "EmployeeUseCaseImpl.SubmitAttendance"),
			zap.Any("user_contex", userContext),
			zap.Any("request", request),
			zap.Error(err),
		)
		return err
	}

	e.auditLogRepository.Create(entity.AuditLog{
		RequestID: userContext.RequestID,
		IPAddress: userContext.IPAddress,
		Action:    "submit",
		Target:    "attendance",
		TableName: "employee_attendances",
		CreatedBy: userContext.Username,
	}, attendance)

	return nil
}

/*
Overtime must be proposed after they are done working.
They can submit the number of hours taken for that overtime.
Overtime cannot be more than 3 hours per day.
Overtime can be taken any day.
*/
func (e *EmployeeUseCaseImpl) SubmitOvertime(userContext entity.UserContext, request entity.SubmitOvertimeRequest) error {
	if userContext.UserID != request.UserID {
		return errors.New("user context does not match request user ID")
	}

	overtimeDate, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		log.Println(
			"error when parsing Date",
			zap.String("method", "EmployeeUseCaseImpl.SubmitAttendance"),
			zap.Any("user_contex", userContext),
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
		_, err := e.employeeRepository.GetAttendanceByUserAndDate(request.UserID, overtimeDate)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.New("attendance must be submitted before submitting overtime")
			}
			log.Println(
				"error when GetAttendanceByUserAndDate",
				zap.String("method", "EmployeeUseCaseImpl.SubmitOvertime"),
				zap.Any("user_contex", userContext),
				zap.Any("request", request),
				zap.Error(err),
			)
			return err
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
		CreatedBy: userContext.Username,
		UpdatedBy: userContext.Username,
	}

	err = e.employeeRepository.UpsertOvertime(overtime)
	if err != nil {
		log.Println(
			"error when UpsertOvertime",
			zap.String("method", "EmployeeUseCaseImpl.SubmitOvertime"),
			zap.Any("user_contex", userContext),
			zap.Any("request", request),
			zap.Error(err),
		)
		return err
	}

	e.auditLogRepository.Create(entity.AuditLog{
		RequestID: userContext.RequestID,
		IPAddress: userContext.IPAddress,
		Action:    "submit",
		Target:    "overtime",
		TableName: "employee_overtimes",
		CreatedBy: userContext.Username,
	}, overtime)

	return nil
}

/*
Employees can attach the amount of money that needs to be reimbursed.
Employees can attach a description to that reimbursement.
*/
func (e *EmployeeUseCaseImpl) SubmitReimbursement(userContext entity.UserContext, request entity.SubmitReimbursementRequest) error {
	if userContext.UserID != request.UserID {
		return errors.New("user context does not match request user ID")
	}

	reimbursementDate, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		log.Println(
			"error when parsing Date",
			zap.String("method", "EmployeeUseCaseImpl.SubmitAttendance"),
			zap.Any("user_contex", userContext),
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
		CreatedBy:   userContext.Username,
		UpdatedBy:   userContext.Username,
	}

	err = e.employeeRepository.UpsertReimbursement(reimbursement)
	if err != nil {
		log.Println(
			"error when UpsertReimbursement",
			zap.String("method", "EmployeeUseCaseImpl.SubmitReimbursement"),
			zap.Any("user_contex", userContext),
			zap.Any("request", request),
			zap.Error(err),
		)
		return err
	}

	e.auditLogRepository.Create(entity.AuditLog{
		RequestID: userContext.RequestID,
		IPAddress: userContext.IPAddress,
		Action:    "submit",
		Target:    "reimbursement",
		TableName: "employee_reimbursements",
		CreatedBy: userContext.Username,
	}, reimbursement)

	return nil
}

func (e *EmployeeUseCaseImpl) GetPayslipBreakdown(userContext entity.UserContext, periodID int64) (any, error) {
	periodDetails, err := e.payrollRepository.GetPeriodByID(periodID)
	if err != nil {
		log.Println(
			"error when GetPeriodByEntityDate",
			zap.String("method", "EmployeeUseCaseImpl.isPeriodActive"),
			zap.Error(err),
		)
		return nil, err
	}

	baseSalaries, err := e.employeeRepository.GetEmployeeBaseSalaryByPeriodStart(periodDetails.PeriodStart, &userContext.UserID)
	if err != nil {
		log.Println(
			"error when GetBaseSalaryByUserID",
			zap.String("method", "EmployeeUseCaseImpl.GetPayslipSummary"),
			zap.Int64("user_id", userContext.UserID),
			zap.Error(err),
		)
		return nil, err
	}

	if len(baseSalaries) != 1 {
		return nil, errors.New("base salary not found for the user in this period")
	}
	baseSalaryDetail := baseSalaries[0]

	attendanceRecords, err := e.employeeRepository.GetAllAttendanceByTimeRange(periodDetails.PeriodStart, periodDetails.PeriodEnd, &userContext.UserID)
	if err != nil {
		return nil, err
	}

	overtimeRecords, err := e.employeeRepository.GetAllOvertimeByTimeRange(periodDetails.PeriodStart, periodDetails.PeriodEnd, &userContext.UserID)
	if err != nil {
		return nil, err
	}

	reimbursementRecords, err := e.employeeRepository.GetAllReimbursementByTimeRange(periodDetails.PeriodStart, periodDetails.PeriodEnd, &userContext.UserID)
	if err != nil {
		return map[int64][]entity.EmployeeReimbursement{}, err
	}

	payslip := entity.PayrollPayslip{}
	payslip.GeneratePayslip(periodDetails, baseSalaryDetail, attendanceRecords, overtimeRecords, reimbursementRecords, userContext.Username)

	return map[string]interface{}{
		"summary":        payslip,
		"period_detail":  periodDetails,
		"attendances":    attendanceRecords,
		"overtimes":      overtimeRecords,
		"reimbursements": reimbursementRecords,
	}, nil
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

	if period.Status == "closed" {
		return false
	}

	return true
}

func (e *EmployeeUseCaseImpl) isSameDay(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() &&
		t1.Month() == t2.Month() &&
		t1.Day() == t2.Day()
}
