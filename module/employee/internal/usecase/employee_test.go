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

func Test_EmployeeUseCase_SubmitAttendance(t *testing.T) {
	tests := []struct {
		name        string
		userContext entity.UserContext
		request     entity.SubmitAttendanceRequest
		mockFunc    func(
			employeeRepository *mocks.EmployeeRepository,
			payrollRepository *mocks.PayrollRepository,
			auditLogRepository *mocks.AuditLogRepository,
		)
		wantErr error
	}{
		{
			name: "error - user context does not match request user ID",
			userContext: entity.UserContext{
				UserID: 112,
			},
			request: entity.SubmitAttendanceRequest{
				UserID: 332,
			},
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
			},
			wantErr: errors.New("user context does not match request user ID"),
		},
		{
			name: "error - invalid date format",
			request: entity.SubmitAttendanceRequest{
				Date: "2023-13-01", // Invalid month
			},
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
			},
			wantErr: errors.New("invalid date format, must be YYYY-MM-DD"),
		},
		{
			name: "error - period is closed",
			request: entity.SubmitAttendanceRequest{
				Date: "2023-12-01",
			},
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByEntityDate", mock.Anything).
					Return(entity.PayrollPeriod{ID: 1, Status: "closed"}, nil)
			},
			wantErr: errors.New("the attendance cannot be submitted because the payroll period is closed"),
		},
		{
			name: "error - invalid check-in time format",
			request: entity.SubmitAttendanceRequest{
				Date:        "2023-12-01",
				CheckInTime: "invalid-time", // Invalid time format
			},
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByEntityDate", mock.Anything).
					Return(entity.PayrollPeriod{ID: 1, Status: "open"}, nil)
			},
			wantErr: errors.New("invalid check-in time format"),
		},
		{
			name: "error - invalid check-out time format",
			request: entity.SubmitAttendanceRequest{
				Date:         "2023-12-01",
				CheckInTime:  "2023-12-01T08:00:00Z",
				CheckOutTime: "invalid-time", // Invalid time format
			},
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByEntityDate", mock.Anything).
					Return(entity.PayrollPeriod{ID: 1, Status: "open"}, nil)
			},
			wantErr: errors.New("invalid check-out time format"),
		},
		{
			name: "error - check-in and check-out time is on different dates",
			request: entity.SubmitAttendanceRequest{
				Date:         "2023-12-01",
				CheckInTime:  "2023-12-01T08:00:00Z",
				CheckOutTime: "2023-12-02T08:00:00Z",
			},
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByEntityDate", mock.Anything).
					Return(entity.PayrollPeriod{ID: 1, Status: "open"}, nil)
			},
			wantErr: errors.New("check-in and check-out times must be on the same day"),
		},
		{
			name: "error - check-in after check-out time",
			request: entity.SubmitAttendanceRequest{
				Date:         "2023-12-01",
				CheckInTime:  "2023-12-01T18:00:00Z",
				CheckOutTime: "2023-12-01T08:00:00Z",
			},
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByEntityDate", mock.Anything).
					Return(entity.PayrollPeriod{ID: 1, Status: "open"}, nil)
			},
			wantErr: errors.New("check-out time cannot be before check-in time"),
		},
		{
			name: "error - check-in on weekend",
			request: entity.SubmitAttendanceRequest{
				Date:         "2023-12-03",
				CheckInTime:  "2023-12-03T08:00:00Z",
				CheckOutTime: "2023-12-03T18:00:00Z",
			},
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByEntityDate", mock.Anything).
					Return(entity.PayrollPeriod{ID: 1, Status: "open"}, nil)
			},
			wantErr: errors.New("attendance can only be submitted on weekdays (Monday to Friday)"),
		},
		{
			name: "error - upsert attendance",
			request: entity.SubmitAttendanceRequest{
				Date:         "2023-12-01",
				CheckInTime:  "2023-12-01T08:00:00Z",
				CheckOutTime: "2023-12-01T18:00:00Z",
			},
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByEntityDate", mock.Anything).
					Return(entity.PayrollPeriod{ID: 1, Status: "open"}, nil)
				employeeRepository.On("UpsertAttendance", mock.Anything).
					Return(errors.New("database error"))
			},
			wantErr: errors.New("database error"),
		},
		{
			name: "success - submit attendance",
			request: entity.SubmitAttendanceRequest{
				Date:         "2023-12-01",
				CheckInTime:  "2023-12-01T08:00:00Z",
				CheckOutTime: "2023-12-01T18:00:00Z",
			},
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByEntityDate", mock.Anything).
					Return(entity.PayrollPeriod{ID: 1, Status: "open"}, nil)
				employeeRepository.On("UpsertAttendance", mock.Anything).
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

			usecase := usecase.NewEmployeeUseCase(employeeRepository, payrollRepository, auditLogRepository)
			err := usecase.SubmitAttendance(entity.UserContext{}, tt.request)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_EmployeeUseCase_SubmitOvertime(t *testing.T) {
	tests := []struct {
		name        string
		userContext entity.UserContext
		request     entity.SubmitOvertimeRequest
		mockFunc    func(
			employeeRepository *mocks.EmployeeRepository,
			payrollRepository *mocks.PayrollRepository,
			auditLogRepository *mocks.AuditLogRepository,
		)
		wantErr error
	}{
		{
			name: "error - user context does not match request user ID",
			userContext: entity.UserContext{
				UserID: 112,
			},
			request: entity.SubmitOvertimeRequest{
				UserID: 332,
			},
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
			},
			wantErr: errors.New("user context does not match request user ID"),
		},
		{
			name: "error - invalid date format",
			request: entity.SubmitOvertimeRequest{
				Date: "2023-13-01", // Invalid month
			},
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
			},
			wantErr: errors.New("invalid date format, must be YYYY-MM-DD"),
		},
		{
			name: "error - period is closed",
			request: entity.SubmitOvertimeRequest{
				Date: "2023-12-01",
			},
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByEntityDate", mock.Anything).
					Return(entity.PayrollPeriod{ID: 1, Status: "closed"}, errors.New(""))
			},
			wantErr: errors.New("the overtime cannot be submitted because the payroll period is closed"),
		},
		{
			name: "error - have no attendance record",
			request: entity.SubmitOvertimeRequest{
				Date: "2023-12-02",
			},
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByEntityDate", mock.Anything).
					Return(entity.PayrollPeriod{ID: 1, Status: "open"}, nil)
				employeeRepository.On("GetAttendanceByUserAndDate", mock.Anything, mock.Anything).
					Return(entity.EmployeeAttendance{}, gorm.ErrRecordNotFound)
			},
			wantErr: errors.New("attendance must be submitted before submitting overtime"),
		},
		{
			name: "error - get attendance record",
			request: entity.SubmitOvertimeRequest{
				Date: "2023-12-02",
			},
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByEntityDate", mock.Anything).
					Return(entity.PayrollPeriod{ID: 1, Status: "open"}, nil)
				employeeRepository.On("GetAttendanceByUserAndDate", mock.Anything, mock.Anything).
					Return(entity.EmployeeAttendance{}, gorm.ErrInvalidDB)
			},
			wantErr: gorm.ErrInvalidDB,
		},
		{
			name: "error - invalid durations",
			request: entity.SubmitOvertimeRequest{
				Date:      "2023-12-02",
				Durations: 6,
			},
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByEntityDate", mock.Anything).
					Return(entity.PayrollPeriod{ID: 1, Status: "open"}, nil)
				employeeRepository.On("GetAttendanceByUserAndDate", mock.Anything, mock.Anything).
					Return(entity.EmployeeAttendance{}, nil)

			},
			wantErr: errors.New("overtime durations must be between 1 and 3 hours"),
		},
		{
			name: "error - upsert",
			request: entity.SubmitOvertimeRequest{
				Date:      "2023-12-02",
				Durations: 2,
			},
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByEntityDate", mock.Anything).
					Return(entity.PayrollPeriod{ID: 1, Status: "open"}, nil)
				employeeRepository.On("GetAttendanceByUserAndDate", mock.Anything, mock.Anything).
					Return(entity.EmployeeAttendance{}, nil)
				employeeRepository.On("UpsertOvertime", mock.Anything).Return(gorm.ErrInvalidDB)

			},
			wantErr: gorm.ErrInvalidDB,
		},
		{
			name: "success",
			request: entity.SubmitOvertimeRequest{
				Date:      "2023-12-02",
				Durations: 2,
			},
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByEntityDate", mock.Anything).
					Return(entity.PayrollPeriod{ID: 1, Status: "open"}, nil)
				employeeRepository.On("GetAttendanceByUserAndDate", mock.Anything, mock.Anything).
					Return(entity.EmployeeAttendance{}, nil)
				employeeRepository.On("UpsertOvertime", mock.Anything).Return(nil)
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

			usecase := usecase.NewEmployeeUseCase(employeeRepository, payrollRepository, auditLogRepository)
			err := usecase.SubmitOvertime(entity.UserContext{}, tt.request)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_EmployeeUseCase_SubmitReimbursement(t *testing.T) {
	tests := []struct {
		name        string
		userContext entity.UserContext
		request     entity.SubmitReimbursementRequest
		mockFunc    func(
			employeeRepository *mocks.EmployeeRepository,
			payrollRepository *mocks.PayrollRepository,
			auditLogRepository *mocks.AuditLogRepository,
		)
		wantErr error
	}{
		{
			name: "error - user context does not match request user ID",
			userContext: entity.UserContext{
				UserID: 112,
			},
			request: entity.SubmitReimbursementRequest{
				UserID: 332,
			},
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
			},
			wantErr: errors.New("user context does not match request user ID"),
		},
		{
			name: "error - invalid date format",
			request: entity.SubmitReimbursementRequest{
				Date: "2023-13-01", // Invalid month
			},
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
			},
			wantErr: errors.New("invalid date format, must be YYYY-MM-DD"),
		},
		{
			name: "error - period is closed",
			request: entity.SubmitReimbursementRequest{
				Date: "2023-12-01",
			},
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByEntityDate", mock.Anything).
					Return(entity.PayrollPeriod{ID: 1, Status: "closed"}, errors.New(""))
			},
			wantErr: errors.New("the reimbursement cannot be submitted because the payroll period is closed"),
		},
		{
			name: "error - upsert",
			request: entity.SubmitReimbursementRequest{
				Date: "2023-12-02",
			},
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByEntityDate", mock.Anything).
					Return(entity.PayrollPeriod{ID: 1, Status: "open"}, nil)
				employeeRepository.On("UpsertReimbursement", mock.Anything).Return(gorm.ErrInvalidDB)
			},
			wantErr: gorm.ErrInvalidDB,
		},
		{
			name: "success",
			request: entity.SubmitReimbursementRequest{
				Date: "2023-12-02",
			},
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByEntityDate", mock.Anything).
					Return(entity.PayrollPeriod{ID: 1, Status: "open"}, nil)
				employeeRepository.On("UpsertReimbursement", mock.Anything).Return(nil)
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

			usecase := usecase.NewEmployeeUseCase(employeeRepository, payrollRepository, auditLogRepository)
			err := usecase.SubmitReimbursement(entity.UserContext{}, tt.request)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_EmployeeUseCase_GetPayslipBreakdown(t *testing.T) {
	tests := []struct {
		name        string
		userContext entity.UserContext
		mockFunc    func(
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
					Return(entity.PayrollPeriod{}, gorm.ErrInvalidDB)
			},
			wantErr: gorm.ErrInvalidDB,
		},
		{
			name: "error - GetEmployeeBaseSalaryByPeriodStart",
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByID", mock.Anything).
					Return(entity.PayrollPeriod{}, nil)
				employeeRepository.On("GetEmployeeBaseSalaryByPeriodStart", mock.Anything, mock.Anything).
					Return([]entity.EmployeeBaseSalary{}, gorm.ErrInvalidDB)

			},
			wantErr: gorm.ErrInvalidDB,
		},
		{
			name: "error - Base salaries not found",
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByID", mock.Anything).
					Return(entity.PayrollPeriod{}, nil)
				employeeRepository.On("GetEmployeeBaseSalaryByPeriodStart", mock.Anything, mock.Anything).
					Return([]entity.EmployeeBaseSalary{}, nil)

			},
			wantErr: errors.New("base salary not found for the user in this period"),
		},
		{
			name: "error - GetAllAttendanceByTimeRange",
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByID", mock.Anything).
					Return(entity.PayrollPeriod{}, nil)
				employeeRepository.On("GetEmployeeBaseSalaryByPeriodStart", mock.Anything, mock.Anything).
					Return([]entity.EmployeeBaseSalary{{UserID: 1, BaseSalary: 2000}}, nil)
				employeeRepository.On("GetAllAttendanceByTimeRange", mock.Anything, mock.Anything, mock.Anything).
					Return([]entity.EmployeeAttendance{}, gorm.ErrInvalidDB)
			},
			wantErr: gorm.ErrInvalidDB,
		},
		{
			name: "error - GetAllOvertimeByTimeRange",
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByID", mock.Anything).
					Return(entity.PayrollPeriod{}, nil)
				employeeRepository.On("GetEmployeeBaseSalaryByPeriodStart", mock.Anything, mock.Anything).
					Return([]entity.EmployeeBaseSalary{{UserID: 1, BaseSalary: 2000}}, nil)
				employeeRepository.On("GetAllAttendanceByTimeRange", mock.Anything, mock.Anything, mock.Anything).
					Return([]entity.EmployeeAttendance{}, nil)
				employeeRepository.On("GetAllOvertimeByTimeRange", mock.Anything, mock.Anything, mock.Anything).
					Return([]entity.EmployeeOvertime{}, gorm.ErrInvalidDB)
			},
			wantErr: gorm.ErrInvalidDB,
		},
		{
			name: "error - GetAllReimbursementByTimeRange",
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
				auditLogRepository *mocks.AuditLogRepository,
			) {
				payrollRepository.On("GetPeriodByID", mock.Anything).
					Return(entity.PayrollPeriod{}, nil)
				employeeRepository.On("GetEmployeeBaseSalaryByPeriodStart", mock.Anything, mock.Anything).
					Return([]entity.EmployeeBaseSalary{{UserID: 1, BaseSalary: 2000}}, nil)
				employeeRepository.On("GetAllAttendanceByTimeRange", mock.Anything, mock.Anything, mock.Anything).
					Return([]entity.EmployeeAttendance{}, nil)
				employeeRepository.On("GetAllOvertimeByTimeRange", mock.Anything, mock.Anything, mock.Anything).
					Return([]entity.EmployeeOvertime{}, nil)
				employeeRepository.On("GetAllReimbursementByTimeRange", mock.Anything, mock.Anything, mock.Anything).
					Return([]entity.EmployeeReimbursement{}, gorm.ErrInvalidDB)
			},
			wantErr: gorm.ErrInvalidDB,
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
					Return(entity.PayrollPayslip{}, gorm.ErrInvalidDB)

				employeeRepository.On("GetEmployeeBaseSalaryByPeriodStart", mock.Anything, mock.Anything).
					Return([]entity.EmployeeBaseSalary{{UserID: 1, BaseSalary: 2000}}, nil)
				employeeRepository.On("GetAllAttendanceByTimeRange", mock.Anything, mock.Anything, mock.Anything).
					Return([]entity.EmployeeAttendance{}, nil)
				employeeRepository.On("GetAllOvertimeByTimeRange", mock.Anything, mock.Anything, mock.Anything).
					Return([]entity.EmployeeOvertime{}, nil)
				employeeRepository.On("GetAllReimbursementByTimeRange", mock.Anything, mock.Anything, mock.Anything).
					Return([]entity.EmployeeReimbursement{}, nil)
			},
			wantErr: gorm.ErrInvalidDB,
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
				payrollRepository.On("GetPayslip", mock.Anything, mock.Anything).
					Return(entity.PayrollPayslip{}, nil)

				employeeRepository.On("GetEmployeeBaseSalaryByPeriodStart", mock.Anything, mock.Anything).
					Return([]entity.EmployeeBaseSalary{{UserID: 1, BaseSalary: 2000}}, nil)
				employeeRepository.On("GetAllAttendanceByTimeRange", mock.Anything, mock.Anything, mock.Anything).
					Return([]entity.EmployeeAttendance{}, nil)
				employeeRepository.On("GetAllOvertimeByTimeRange", mock.Anything, mock.Anything, mock.Anything).
					Return([]entity.EmployeeOvertime{}, nil)
				employeeRepository.On("GetAllReimbursementByTimeRange", mock.Anything, mock.Anything, mock.Anything).
					Return([]entity.EmployeeReimbursement{}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			employeeRepository := mocks.NewEmployeeRepository(t)
			payrollRepository := mocks.NewPayrollRepository(t)
			auditLogRepository := mocks.NewAuditLogRepository(t)

			tt.mockFunc(employeeRepository, payrollRepository, auditLogRepository)

			usecase := usecase.NewEmployeeUseCase(employeeRepository, payrollRepository, auditLogRepository)
			_, err := usecase.GetPayslipBreakdown(entity.UserContext{}, 123)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
