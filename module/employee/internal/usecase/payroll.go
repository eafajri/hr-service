package usecase

import (
	"errors"
	"log"

	"github.com/eafajri/hr-service.git/module/employee/internal/entity"
	"go.uber.org/zap"
)

//go:generate mockery --name PayrollUseCase --output ./mocks
type PayrollUseCase interface {
	GetPayslip(userID int64, periodID int64) (entity.PayrollPayslip, error)
	GetPayslips(periodID int64) ([]entity.PayrollPayslip, error)
	ClosePayrollPeriod(userContex entity.UserContext, periodID int64) error
	GeneratePayslipsByPeriodID(userContext entity.UserContext, periodID int64) error
}

type PayrollUseCaseImpl struct {
	payrollRepository  PayrollRepository
	employeeRepository EmployeeRepository
	auditLogRepository AuditLogRepository
}

func NewPayrollUseCase(
	payrollRepository PayrollRepository,
	employeeRepository EmployeeRepository,
	auditLogRepository AuditLogRepository,
) *PayrollUseCaseImpl {
	return &PayrollUseCaseImpl{
		payrollRepository:  payrollRepository,
		employeeRepository: employeeRepository,
		auditLogRepository: auditLogRepository,
	}
}

func (p *PayrollUseCaseImpl) GetPayslip(userID int64, periodID int64) (entity.PayrollPayslip, error) {
	payrollPeriod, err := p.payrollRepository.GetPeriodByID(periodID)
	if err != nil {
		log.Println(
			"error when GetPeriodByID",
			zap.String("method", "PayrollUseCaseImpl.GetPayslip"),
			zap.Int64("user_id", userID),
			zap.Int64("period_id", periodID),
			zap.Error(err),
		)
		return entity.PayrollPayslip{}, err
	}

	if payrollPeriod.Status == "open" {
		return entity.PayrollPayslip{}, errors.New("the payroll period is still open")
	}

	payslip, err := p.payrollRepository.GetPayslip(userID, periodID)
	if err != nil {
		log.Println(
			"error when GetPayslip",
			zap.String("method", "PayrollUseCaseImpl.GetPayslip"),
			zap.Int64("user_id", userID),
			zap.Int64("period_id", periodID),
			zap.Error(err),
		)
		return entity.PayrollPayslip{}, err
	}

	return payslip, nil
}

/*
The summary contains take-home pay of each employee.
The summary contains the total take-home pay of all employees.
*/
func (p *PayrollUseCaseImpl) GetPayslips(periodID int64) ([]entity.PayrollPayslip, error) {
	payrollPeriod, err := p.payrollRepository.GetPeriodByID(periodID)
	if err != nil {
		log.Println(
			"error when GetPeriodByID",
			zap.String("method", "PayrollUseCaseImpl.GetPayslip"),
			zap.Int64("period_id", periodID),
			zap.Error(err),
		)
		return []entity.PayrollPayslip{}, err
	}

	if payrollPeriod.Status == "open" {
		return []entity.PayrollPayslip{}, errors.New("the payroll period is still open")
	}

	payslips, err := p.payrollRepository.GetPayslips(periodID)
	if err != nil {
		log.Println(
			"error when GetPayslips",
			zap.String("method", "PayrollUseCaseImpl.GetPayslips"),
			zap.Int64("period_id", periodID),
			zap.Error(err),
		)
		return nil, err
	}

	var totalTakeHome float64
	for _, payslip := range payslips {
		totalTakeHome += payslip.TotalTakeHome
	}

	return payslips, nil
}

/*
Once payroll is run, attendance, overtime, and reimbursement records from that period cannot affect the payslip.
Payroll for each attendance period can only be run once.
*/
func (p *PayrollUseCaseImpl) ClosePayrollPeriod(userContext entity.UserContext, periodID int64) error {
	payrollPeriod, err := p.payrollRepository.GetPeriodByID(periodID)
	if err != nil {
		log.Println(
			"error when GetPeriodByID",
			zap.String("method", "PayrollUseCaseImpl.GetPayslip"),
			zap.Int64("period_id", periodID),
			zap.Error(err),
		)
		return err
	}

	if payrollPeriod.Status == "closed" {
		return errors.New("the payroll period is already closed")
	}

	err = p.payrollRepository.ClosePayrollPeriod(periodID)
	if err != nil {
		log.Println(
			"error when GetPayslip",
			zap.String("method", "PayrollUseCaseImpl.ClosePayrollPeriod"),
			zap.Int64("period_id", periodID),
			zap.Error(err),
		)
		return err
	}

	p.auditLogRepository.Create(entity.AuditLog{
		RequestID: userContext.RequestID,
		IPAddress: userContext.IPAddress,
		Action:    "update",
		Target:    "reimbursement",
		TableName: "payroll_period",
		CreatedBy: userContext.Username,
	}, payrollPeriod)

	return err
}

func (p *PayrollUseCaseImpl) GeneratePayslipsByPeriodID(userContext entity.UserContext, periodID int64) error {
	periodDetails, err := p.payrollRepository.GetPeriodByID(periodID)
	if err != nil {
		log.Println(
			"error when GetPeriodByID",
			zap.String("method", "PayrollUseCaseImpl.ClosePayrollPeriod"),
			zap.Int64("period_id", periodID),
			zap.Error(err),
		)
		return err
	}

	if periodDetails.Status == "open" {
		return errors.New("unable to process open period")
	}

	employeeBaseSalaries, err := p.employeeRepository.GetEmployeeBaseSalaryByPeriodStart(periodDetails.PeriodStart, nil)
	if err != nil {
		log.Println(
			"error when GetEmployeeBaseSalaryByPeriodID",
			zap.String("method", "PayrollUseCaseImpl.ClosePayrollPeriod"),
			zap.Int64("period_id", periodID),
			zap.Error(err),
		)
		return err
	}

	attendanceRecordsMap, err := p.getEmployeeAttendanceByPeriodID(periodDetails)
	if err != nil {
		log.Println(
			"error when GetAllAttendanceByPeriodID",
			zap.String("method", "PayrollUseCaseImpl.ClosePayrollPeriod"),
			zap.Int64("period_id", periodID),
			zap.Error(err),
		)
		return err
	}

	overtimeRecordsMap, err := p.getEmployeeOvertimeByPeriodID(periodDetails)
	if err != nil {
		log.Println(
			"error when GetAllOvertimeByPeriodID",
			zap.String("method", "PayrollUseCaseImpl.ClosePayrollPeriod"),
			zap.Int64("period_id", periodID),
			zap.Error(err),
		)
		return err
	}

	reimbursementRecordsMap, err := p.getEmployeeReimbursementByPeriodID(periodDetails)
	if err != nil {
		log.Println(
			"error when GetAllReimbursementByPeriodID",
			zap.String("method", "PayrollUseCaseImpl.ClosePayrollPeriod"),
			zap.Int64("period_id", periodID),
			zap.Error(err),
		)
		return err
	}

	payslips := make([]entity.PayrollPayslip, 0, len(employeeBaseSalaries))
	for _, employeeBaseSalary := range employeeBaseSalaries {
		payslip := entity.PayrollPayslip{}
		payslip.GeneratePayslip(periodDetails, employeeBaseSalary, attendanceRecordsMap[employeeBaseSalary.UserID], overtimeRecordsMap[employeeBaseSalary.UserID], reimbursementRecordsMap[employeeBaseSalary.UserID], userContext.Username)

		payslips = append(payslips, payslip)
	}

	err = p.payrollRepository.CreatePayslipsByPeriod(payslips)
	if err != nil {
		log.Println(
			"error when CreatePayslipsByPeriod",
			zap.String("method", "PayrollUseCaseImpl.ClosePayrollPeriod"),
			zap.Int64("period_id", periodID),
			zap.Error(err),
		)
		return err
	}

	p.auditLogRepository.Create(entity.AuditLog{
		RequestID: userContext.RequestID,
		IPAddress: userContext.IPAddress,
		Action:    "create",
		Target:    "payslips",
		TableName: "payroll_peyslips",
		CreatedBy: userContext.Username,
	}, payslips)

	return nil
}

func (p *PayrollUseCaseImpl) getEmployeeAttendanceByPeriodID(periodDetails entity.PayrollPeriod) (map[int64][]entity.EmployeeAttendance, error) {
	attendanceRecords, err := p.employeeRepository.GetAllAttendanceByTimeRange(periodDetails.PeriodStart, periodDetails.PeriodEnd, nil)
	if err != nil {
		return map[int64][]entity.EmployeeAttendance{}, err
	}

	attendanceRecordsMap := make(map[int64][]entity.EmployeeAttendance)
	for _, record := range attendanceRecords {
		if _, exists := attendanceRecordsMap[record.UserID]; !exists {
			attendanceRecordsMap[record.UserID] = []entity.EmployeeAttendance{}
		}
		attendanceRecordsMap[record.UserID] = append(attendanceRecordsMap[record.UserID], record)
	}

	return attendanceRecordsMap, nil
}

func (p *PayrollUseCaseImpl) getEmployeeOvertimeByPeriodID(periodDetails entity.PayrollPeriod) (map[int64][]entity.EmployeeOvertime, error) {
	// Get employee overtime by period ID
	overtimeRecords, err := p.employeeRepository.GetAllOvertimeByTimeRange(periodDetails.PeriodStart, periodDetails.PeriodEnd, nil)
	if err != nil {
		return map[int64][]entity.EmployeeOvertime{}, err
	}

	overtimeRecordsMap := make(map[int64][]entity.EmployeeOvertime)
	for _, record := range overtimeRecords {
		if _, exists := overtimeRecordsMap[record.UserID]; !exists {
			overtimeRecordsMap[record.UserID] = []entity.EmployeeOvertime{}
		}
		overtimeRecordsMap[record.UserID] = append(overtimeRecordsMap[record.UserID], record)
	}

	return overtimeRecordsMap, nil
}

func (p *PayrollUseCaseImpl) getEmployeeReimbursementByPeriodID(periodDetails entity.PayrollPeriod) (map[int64][]entity.EmployeeReimbursement, error) {
	reimbursementRecords, err := p.employeeRepository.GetAllReimbursementByTimeRange(periodDetails.PeriodStart, periodDetails.PeriodEnd, nil)
	if err != nil {
		return map[int64][]entity.EmployeeReimbursement{}, err
	}

	reimbursementRecordsMap := make(map[int64][]entity.EmployeeReimbursement)
	for _, record := range reimbursementRecords {
		if _, exists := reimbursementRecordsMap[record.UserID]; !exists {
			reimbursementRecordsMap[record.UserID] = []entity.EmployeeReimbursement{}
		}
		reimbursementRecordsMap[record.UserID] = append(reimbursementRecordsMap[record.UserID], record)
	}

	return reimbursementRecordsMap, nil
}
