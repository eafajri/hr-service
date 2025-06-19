package usecase

import (
	"log"
	"time"

	"github.com/eafajri/hr-service.git/module/employee/internal/entity"
	"go.uber.org/zap"
)

//go:generate mockery --name PayrollUseCase --output ./mocks
type PayrollUseCase interface {
	GetPayslip(userID int64, periodID int64) (entity.PayrollPayslip, error)
	GetPayslips(periodID int64) ([]entity.PayrollPayslip, error)
	ClosePayrollPeriod(periodID int64) error
}

type PayrollUseCaseImpl struct {
	payrollRepository  PayrollRepository
	employeeRepository EmployeeRepository
}

func NewPayrollUseCase(
	payrollRepository PayrollRepository,
	employeeRepository EmployeeRepository,
) *PayrollUseCaseImpl {
	return &PayrollUseCaseImpl{
		payrollRepository:  payrollRepository,
		employeeRepository: employeeRepository,
	}
}

func (p *PayrollUseCaseImpl) GetPayslip(userID int64, periodID int64) (entity.PayrollPayslip, error) {
	if periodID == 0 {
		period, err := p.payrollRepository.GetPeriodByEntityDate(time.Now())
		if err != nil {
			log.Println(
				"error when GetPeriodByEntityDate",
				zap.String("method", "PayrollUseCaseImpl.GetPayslip"),
				zap.Int64("user_id", userID),
				zap.Error(err),
			)
			return entity.PayrollPayslip{}, err
		}

		periodID = period.ID
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

func (p *PayrollUseCaseImpl) GetPayslips(periodID int64) ([]entity.PayrollPayslip, error) {
	if periodID == 0 {
		period, err := p.payrollRepository.GetPeriodByEntityDate(time.Now())
		if err != nil {
			log.Println(
				"error when GetPeriodByEntityDate",
				zap.String("method", "PayrollUseCaseImpl.GetPayslips"),
				zap.Int64("period_id", periodID),
				zap.Error(err),
			)
			return nil, err
		}

		periodID = period.ID
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

	return payslips, nil
}

func (p *PayrollUseCaseImpl) ClosePayrollPeriod(periodID int64) error {
	err := p.payrollRepository.ClosePayrollPeriod(periodID)
	if err != nil {
		log.Println(
			"error when GetPayslip",
			zap.String("method", "PayrollUseCaseImpl.ClosePayrollPeriod"),
			zap.Int64("period_id", periodID),
			zap.Error(err),
		)
		return err
	}

	go p.createPayslipsByPeriodID(periodID)

	return err
}

func (p *PayrollUseCaseImpl) createPayslipsByPeriodID(periodID int64) error {
	// Get employee base salary by period ID
	employeeBaseSalaries, err := p.payrollRepository.GetEmployeeBaseSalaryByPeriodID(periodID)
	if err != nil {
		log.Println(
			"error when GetEmployeeBaseSalaryByPeriodID",
			zap.String("method", "PayrollUseCaseImpl.ClosePayrollPeriod"),
			zap.Int64("period_id", periodID),
			zap.Error(err),
		)
		return err
	}

	attendanceRecordsMap, err := p.getEmployeeBaseSalaryByPeriodID(periodID)
	if err != nil {
		log.Println(
			"error when GetAllAttendanceByPeriodID",
			zap.String("method", "PayrollUseCaseImpl.ClosePayrollPeriod"),
			zap.Int64("period_id", periodID),
			zap.Error(err),
		)
		return err
	}

	overtimeRecordsMap, err := p.getEmployeeOvertimeByPeriodID(periodID)
	if err != nil {
		log.Println(
			"error when GetAllOvertimeByPeriodID",
			zap.String("method", "PayrollUseCaseImpl.ClosePayrollPeriod"),
			zap.Int64("period_id", periodID),
			zap.Error(err),
		)
		return err
	}

	reimbursementRecordsMap, err := p.getEmployeeReimbursementByPeriodID(periodID)
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
		payslip := entity.PayrollPayslip{
			UserID:          employeeBaseSalary.UserID,
			PayrollPeriodID: periodID,
			BaseSalary:      employeeBaseSalary.BaseSalary,
		}

		ratePerDay := employeeBaseSalary.BaseSalary / float64(30)
		ratePerHour := ratePerDay / 8

		attendanceRecords, exists := attendanceRecordsMap[employeeBaseSalary.UserID]
		if exists {
			payslip.AttendanceDays = len(attendanceRecords)
			for _, record := range attendanceRecords {
				duration := record.CheckOutTime.Sub(record.CheckInTime)
				payslip.AttendanceHours += int(duration.Hours())
			}
			payslip.AttendacePay = float64(payslip.AttendanceHours) * ratePerHour
		}

		overtimeRecords, exists := overtimeRecordsMap[employeeBaseSalary.UserID]
		if exists {
			for _, record := range overtimeRecords {
				payslip.OvertimeHours += record.Durations
			}
			// Overtime pay is calculated as double the rate per hour
			payslip.OvertimePay = float64(payslip.OvertimeHours) * ratePerHour * 2
		}

		reimbursementRecords, exists := reimbursementRecordsMap[employeeBaseSalary.UserID]
		if exists {
			for _, record := range reimbursementRecords {
				payslip.ReimbursementTotal += record.Amount
			}
		}

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

	return nil
}

func (p *PayrollUseCaseImpl) getEmployeeBaseSalaryByPeriodID(periodID int64) (map[int64][]entity.EmployeeAttendance, error) {
	attendanceRecords, err := p.employeeRepository.GetAllAttendanceByPeriodID(periodID)
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

func (p *PayrollUseCaseImpl) getEmployeeOvertimeByPeriodID(periodID int64) (map[int64][]entity.EmployeeOvertime, error) {
	// Get employee overtime by period ID
	overtimeRecords, err := p.employeeRepository.GetAllOvertimeByPeriodID(periodID)
	if err != nil {
		log.Println(
			"error when GetAllOvertimeByPeriodID",
			zap.String("method", "PayrollUseCaseImpl.ClosePayrollPeriod"),
			zap.Int64("period_id", periodID),
			zap.Error(err),
		)
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

func (p *PayrollUseCaseImpl) getEmployeeReimbursementByPeriodID(periodID int64) (map[int64][]entity.EmployeeReimbursement, error) {
	reimbursementRecords, err := p.employeeRepository.GetAllReimbursementByPeriodID(periodID)
	if err != nil {
		log.Println(
			"error when GetAllReimbursementByPeriodID",
			zap.String("method", "PayrollUseCaseImpl.ClosePayrollPeriod"),
			zap.Int64("period_id", periodID),
			zap.Error(err),
		)
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
