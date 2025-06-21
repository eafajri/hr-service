package usecase_test

import (
	"errors"
	"testing"

	"github.com/eafajri/hr-service.git/module/employee/internal/entity"
	"github.com/eafajri/hr-service.git/module/employee/internal/usecase"
	"github.com/eafajri/hr-service.git/module/employee/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_EmployeeUseCase_SubmitAttendance(t *testing.T) {
	tests := []struct {
		name     string
		request  entity.SubmitAttendanceRequest
		mockFunc func(
			employeeRepository *mocks.EmployeeRepository,
			payrollRepository *mocks.PayrollRepository,
		)
		wantErr error
	}{
		{
			name: "error - invalid date format",
			request: entity.SubmitAttendanceRequest{
				Date: "2023-13-01", // Invalid month
			},
			mockFunc: func(
				employeeRepository *mocks.EmployeeRepository,
				payrollRepository *mocks.PayrollRepository,
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
			) {
				payrollRepository.On("GetPeriodByEntityDate", mock.Anything).
					Return(entity.PayrollPeriod{ID: 1, Status: "open"}, nil)
				employeeRepository.On("UpsertAttendance", mock.Anything).
					Return(nil)
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			employeeRepository := mocks.NewEmployeeRepository(t)
			payrollRepository := mocks.NewPayrollRepository(t)

			tt.mockFunc(employeeRepository, payrollRepository)

			usecase := usecase.NewEmployeeUseCase(employeeRepository, payrollRepository)
			err := usecase.SubmitAttendance(entity.User{}, tt.request)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
