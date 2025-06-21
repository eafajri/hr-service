package usecase_test

import (
	"errors"
	"testing"

	"github.com/eafajri/hr-service.git/module/employee/internal/entity"
	"github.com/eafajri/hr-service.git/module/employee/internal/usecase"
	"github.com/eafajri/hr-service.git/module/employee/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func Test_PayrollUseCase_GetPayslip(t *testing.T) {
	tests := []struct {
		name     string
		mockFunc func(
			employeeRepository *mocks.EmployeeRepository,
			payrollRepository *mocks.PayrollRepository,
			auditLogRepository *mocks.AuditLogRepository,
		)
		wantErr error
		wantRes entity.PayrollPayslip
	}{
		{
			name: "error - GetPeriodByID",
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByID", mock.Anything).
					Return(entity.PayrollPeriod{}, gorm.ErrSubQueryRequired)
			},
			wantErr: gorm.ErrSubQueryRequired,
		},
		{
			name: "error - period is open",
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByID", mock.Anything).
					Return(entity.PayrollPeriod{Status: "open"}, nil)
			},
			wantErr: errors.New("the payroll period is still open"),
		},
		{
			name: "error - GetPayslip",
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByID", mock.Anything).
					Return(entity.PayrollPeriod{Status: "closed"}, nil)
				payrollRepository.On("GetPayslip", mock.Anything, mock.Anything).
					Return(entity.PayrollPayslip{}, gorm.ErrSubQueryRequired)
			},
			wantErr: gorm.ErrSubQueryRequired,
		},
		{
			name: "Success",
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByID", mock.Anything).
					Return(entity.PayrollPeriod{Status: "closed"}, nil)
				payrollRepository.On("GetPayslip", mock.Anything, mock.Anything).
					Return(entity.PayrollPayslip{BaseSalary: 21000, TotalTakeHome: 21000}, nil)
			},
			wantRes: entity.PayrollPayslip{
				BaseSalary:    21000,
				TotalTakeHome: 21000,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			employeeRepository := mocks.NewEmployeeRepository(t)
			payrollRepository := mocks.NewPayrollRepository(t)
			auditLogRepository := mocks.NewAuditLogRepository(t)

			tt.mockFunc(employeeRepository, payrollRepository, auditLogRepository)

			usecase := usecase.NewPayrollUseCase(payrollRepository, employeeRepository, auditLogRepository)
			res, err := usecase.GetPayslip(0, 0)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.Equal(t, tt.wantRes, res)
				assert.NoError(t, err)
			}
		})
	}
}

func Test_PayrollUseCase_GetPayslips(t *testing.T) {
	tests := []struct {
		name     string
		mockFunc func(
			employeeRepository *mocks.EmployeeRepository,
			payrollRepository *mocks.PayrollRepository,
			auditLogRepository *mocks.AuditLogRepository,
		)
		wantErr error
		wantRes []entity.PayrollPayslip
	}{
		{
			name: "error - GetPeriodByID",
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByID", mock.Anything).
					Return(entity.PayrollPeriod{}, gorm.ErrSubQueryRequired)
			},
			wantErr: gorm.ErrSubQueryRequired,
		},
		{
			name: "error - period is open",
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByID", mock.Anything).
					Return(entity.PayrollPeriod{Status: "open"}, nil)
			},
			wantErr: errors.New("the payroll period is still open"),
		},
		{
			name: "error - GetPayslips",
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByID", mock.Anything).
					Return(entity.PayrollPeriod{Status: "closed"}, nil)
				payrollRepository.On("GetPayslips", mock.Anything).
					Return([]entity.PayrollPayslip{}, gorm.ErrSubQueryRequired)
			},
			wantErr: gorm.ErrSubQueryRequired,
		},
		{
			name: "Success",
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByID", mock.Anything).
					Return(entity.PayrollPeriod{Status: "closed"}, nil)
				payrollRepository.On("GetPayslips", mock.Anything).
					Return([]entity.PayrollPayslip{{BaseSalary: 21000, TotalTakeHome: 21000}}, nil)
			},
			wantRes: []entity.PayrollPayslip{
				{
					BaseSalary:    21000,
					TotalTakeHome: 21000,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			employeeRepository := mocks.NewEmployeeRepository(t)
			payrollRepository := mocks.NewPayrollRepository(t)
			auditLogRepository := mocks.NewAuditLogRepository(t)

			tt.mockFunc(employeeRepository, payrollRepository, auditLogRepository)

			usecase := usecase.NewPayrollUseCase(payrollRepository, employeeRepository, auditLogRepository)
			res, err := usecase.GetPayslips(0)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.Equal(t, tt.wantRes, res)
				assert.NoError(t, err)
			}
		})
	}
}

func Test_PayrollUseCase_ClosePayrollPeriod(t *testing.T) {
	tests := []struct {
		name     string
		mockFunc func(
			employeeRepository *mocks.EmployeeRepository,
			payrollRepository *mocks.PayrollRepository,
			auditLogRepository *mocks.AuditLogRepository,
		)
		wantErr error
	}{
		{
			name: "error - GetPeriodByID",
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByID", mock.Anything).
					Return(entity.PayrollPeriod{}, gorm.ErrSubQueryRequired)
			},
			wantErr: gorm.ErrSubQueryRequired,
		},
		{
			name: "error - period is already closed",
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByID", mock.Anything).
					Return(entity.PayrollPeriod{Status: "closed"}, nil)
			},
			wantErr: errors.New("the payroll period is already closed"),
		},
		{
			name: "error - ClosePayrollPeriod",
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByID", mock.Anything).
					Return(entity.PayrollPeriod{Status: "open"}, nil)
				payrollRepository.On("ClosePayrollPeriod", mock.Anything, mock.Anything).
					Return(gorm.ErrSubQueryRequired)
			},
			wantErr: gorm.ErrSubQueryRequired,
		},
		{
			name: "success",
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByID", mock.Anything).
					Return(entity.PayrollPeriod{Status: "open"}, nil)
				payrollRepository.On("ClosePayrollPeriod", mock.Anything, mock.Anything).
					Return(nil)
				auditLogRepository.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			employeeRepository := mocks.NewEmployeeRepository(t)
			payrollRepository := mocks.NewPayrollRepository(t)
			auditLogRepository := mocks.NewAuditLogRepository(t)

			tt.mockFunc(employeeRepository, payrollRepository, auditLogRepository)

			usecase := usecase.NewPayrollUseCase(payrollRepository, employeeRepository, auditLogRepository)
			err := usecase.ClosePayrollPeriod(entity.UserContext{}, 0)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_PayrollUseCase_GeneratePayslipsByPeriodID(t *testing.T) {
	tests := []struct {
		name     string
		mockFunc func(
			employeeRepository *mocks.EmployeeRepository,
			payrollRepository *mocks.PayrollRepository,
			auditLogRepository *mocks.AuditLogRepository,
		)
		wantErr error
	}{
		{
			name: "error - GetPeriodByID",
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByID", mock.Anything).
					Return(entity.PayrollPeriod{}, gorm.ErrSubQueryRequired)
			},
			wantErr: gorm.ErrSubQueryRequired,
		},
		{
			name: "error - period is still open",
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByID", mock.Anything).
					Return(entity.PayrollPeriod{Status: "open"}, nil)
			},
			wantErr: errors.New("unable to process open period"),
		},
		{
			name: "error - ClosePayrollPeriod",
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByID", mock.Anything).
					Return(entity.PayrollPeriod{Status: "closed"}, nil)
				employeeRepository.On("GetEmployeeBaseSalaryByPeriodStart", mock.Anything, mock.Anything).
					Return([]entity.EmployeeBaseSalary{}, gorm.ErrSubQueryRequired)
			},
			wantErr: gorm.ErrSubQueryRequired,
		},
		{
			name: "error - GetAllAttendanceByTimeRange",
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByID", mock.Anything).
					Return(entity.PayrollPeriod{Status: "closed"}, nil)
				employeeRepository.On("GetEmployeeBaseSalaryByPeriodStart", mock.Anything, mock.Anything).
					Return([]entity.EmployeeBaseSalary{}, nil)
				employeeRepository.On("GetAllAttendanceByTimeRange", mock.Anything, mock.Anything, mock.Anything).
					Return([]entity.EmployeeAttendance{}, gorm.ErrSubQueryRequired)
			},
			wantErr: gorm.ErrSubQueryRequired,
		},
		{
			name: "error - GetAllAttendanceByTimeRange",
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByID", mock.Anything).
					Return(entity.PayrollPeriod{Status: "closed"}, nil)
				employeeRepository.On("GetEmployeeBaseSalaryByPeriodStart", mock.Anything, mock.Anything).
					Return([]entity.EmployeeBaseSalary{}, nil)
				employeeRepository.On("GetAllAttendanceByTimeRange", mock.Anything, mock.Anything, mock.Anything).
					Return([]entity.EmployeeAttendance{}, nil)
				employeeRepository.On("GetAllOvertimeByTimeRange", mock.Anything, mock.Anything, mock.Anything).
					Return([]entity.EmployeeOvertime{}, gorm.ErrSubQueryRequired)
			},
			wantErr: gorm.ErrSubQueryRequired,
		},
		{
			name: "error - GetAllReimbursementByTimeRange",
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByID", mock.Anything).
					Return(entity.PayrollPeriod{Status: "closed"}, nil)
				employeeRepository.On("GetEmployeeBaseSalaryByPeriodStart", mock.Anything, mock.Anything).
					Return([]entity.EmployeeBaseSalary{}, nil)
				employeeRepository.On("GetAllAttendanceByTimeRange", mock.Anything, mock.Anything, mock.Anything).
					Return([]entity.EmployeeAttendance{}, nil)
				employeeRepository.On("GetAllOvertimeByTimeRange", mock.Anything, mock.Anything, mock.Anything).
					Return([]entity.EmployeeOvertime{}, nil)
				employeeRepository.On("GetAllReimbursementByTimeRange", mock.Anything, mock.Anything, mock.Anything).
					Return([]entity.EmployeeReimbursement{}, gorm.ErrSubQueryRequired)
			},
			wantErr: gorm.ErrSubQueryRequired,
		},
		{
			name: "error - CreatePayslipsByPeriod",
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByID", mock.Anything).
					Return(entity.PayrollPeriod{Status: "closed"}, nil)
				employeeRepository.On("GetEmployeeBaseSalaryByPeriodStart", mock.Anything, mock.Anything).
					Return([]entity.EmployeeBaseSalary{}, nil)
				employeeRepository.On("GetAllAttendanceByTimeRange", mock.Anything, mock.Anything, mock.Anything).
					Return([]entity.EmployeeAttendance{}, nil)
				employeeRepository.On("GetAllOvertimeByTimeRange", mock.Anything, mock.Anything, mock.Anything).
					Return([]entity.EmployeeOvertime{}, nil)
				employeeRepository.On("GetAllReimbursementByTimeRange", mock.Anything, mock.Anything, mock.Anything).
					Return([]entity.EmployeeReimbursement{}, nil)
				payrollRepository.On("CreatePayslipsByPeriod", mock.Anything).
					Return(gorm.ErrSubQueryRequired)

			},
			wantErr: gorm.ErrSubQueryRequired,
		},
		{
			name: "success",
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByID", mock.Anything).
					Return(entity.PayrollPeriod{Status: "closed"}, nil)
				employeeRepository.On("GetEmployeeBaseSalaryByPeriodStart", mock.Anything, mock.Anything).
					Return([]entity.EmployeeBaseSalary{{UserID: 12}}, nil)
				employeeRepository.On("GetAllAttendanceByTimeRange", mock.Anything, mock.Anything, mock.Anything).
					Return([]entity.EmployeeAttendance{{ID: 33}}, nil)
				employeeRepository.On("GetAllOvertimeByTimeRange", mock.Anything, mock.Anything, mock.Anything).
					Return([]entity.EmployeeOvertime{{ID: 412}}, nil)
				employeeRepository.On("GetAllReimbursementByTimeRange", mock.Anything, mock.Anything, mock.Anything).
					Return([]entity.EmployeeReimbursement{{ID: 41}}, nil)
				payrollRepository.On("CreatePayslipsByPeriod", mock.Anything).
					Return(nil)
				auditLogRepository.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			employeeRepository := mocks.NewEmployeeRepository(t)
			payrollRepository := mocks.NewPayrollRepository(t)
			auditLogRepository := mocks.NewAuditLogRepository(t)

			tt.mockFunc(employeeRepository, payrollRepository, auditLogRepository)

			usecase := usecase.NewPayrollUseCase(payrollRepository, employeeRepository, auditLogRepository)
			err := usecase.GeneratePayslipsByPeriodID(entity.UserContext{}, 0)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
