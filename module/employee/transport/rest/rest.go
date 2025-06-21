package transport

import (
	"net/http"
	"time"

	moduleConfig "github.com/eafajri/hr-service.git/module/employee/config"
	"github.com/eafajri/hr-service.git/module/employee/internal/entity"
	"github.com/eafajri/hr-service.git/module/employee/internal/repository"
	"github.com/eafajri/hr-service.git/module/employee/internal/usecase"
	"github.com/labstack/echo/v4"
)

type Rest struct {
	userUc     usecase.UserUseCase
	employeeUc usecase.EmployeeUseCase
	payrollUc  usecase.PayrollUseCase
}

func StartRest(echoInstance *echo.Echo) {

	moduleDependencies := moduleConfig.NewModuleDependencies()

	var (
		userRepository     = repository.NewUserRepository(&moduleDependencies.Database)
		employeeRepository = repository.NewEmployeeRepository(&moduleDependencies.Database)
		payrollRepository  = repository.NewPayrollRepository(&moduleDependencies.Database)
		auditLogRepository = repository.NewAuditLogRepository(&moduleDependencies.Database)
	)

	restHandler := &Rest{
		userUc:     usecase.NewUserUseCase(userRepository),
		employeeUc: usecase.NewEmployeeUseCase(employeeRepository, payrollRepository, auditLogRepository),
		payrollUc:  usecase.NewPayrollUseCase(payrollRepository, employeeRepository, auditLogRepository),
	}

	publicApi := echoInstance.Group("/public")
	publicApi.GET("/check", restHandler.CheckHealth)

	employeeApi := echoInstance.Group("/private/employee")
	employeeApi.Use(BasicAuthMiddleware(restHandler.userUc))
	employeeApi.POST("/attendance/submit", restHandler.SubmitAttendance)
	employeeApi.POST("/overtime/submit", restHandler.SubmitOvertime)
	employeeApi.POST("/reimbursement/submit", restHandler.SubmitReimbursement)
	employeeApi.GET("/payslips/:period_id", restHandler.GetPayslipBreakdown)

	adminApi := echoInstance.Group("/private/admin")
	adminApi.Use(BasicAuthMiddleware(restHandler.userUc))
	adminApi.Use(AdminPrevilageMiddleware(restHandler.userUc))
	adminApi.POST("/generate-payroll/:period_id", restHandler.GeneratePayroll)
	adminApi.GET("/payslips/:period_id", restHandler.GetPayslips)
	adminApi.GET("/payslips/:period_id/:user_id", restHandler.GetPayslip)
}

func (h *Rest) CheckHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok", "timestamp": time.Now().Format(time.RFC3339)})
}

func (r *Rest) standardizeResponse(c echo.Context, statusCode int, message string, data interface{}) error {
	response := entity.Response{
		Meta: entity.Meta{
			StatusCode: statusCode,
			Message:    message,
		},
		Data: data,
	}

	return c.JSON(statusCode, response)
}
