package transport

import (
	"net/http"
	"strconv"

	"github.com/eafajri/hr-service.git/module/employee/internal/entity"
	"github.com/labstack/echo/v4"
)

func (r *Rest) SubmitAttendance(c echo.Context) error {
	userDetail, ok := c.Get("user").(entity.User)
	if !ok {
		return r.standardizeResponse(c, http.StatusUnauthorized, "User ID not found in context", nil)
	}

	var request entity.SubmitAttendanceRequest
	if err := c.Bind(&request); err != nil {
		return r.standardizeResponse(c, http.StatusBadRequest, "Invalid request format", nil)
	}

	err := r.employeeUc.SubmitAttendance(userDetail, request)
	if err != nil {
		return r.standardizeResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return r.standardizeResponse(c, http.StatusOK, "Attendance submitted successfully", nil)
}

func (r *Rest) SubmitOvertime(c echo.Context) error {
	userDetail, ok := c.Get("user").(entity.User)
	if !ok {
		return r.standardizeResponse(c, http.StatusUnauthorized, "User ID not found in context", nil)
	}

	var request entity.SubmitOvertimeRequest
	if err := c.Bind(&request); err != nil {
		return r.standardizeResponse(c, http.StatusBadRequest, "Invalid request format", nil)
	}

	err := r.employeeUc.SubmitOvertime(userDetail, request)
	if err != nil {
		return r.standardizeResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return r.standardizeResponse(c, http.StatusOK, "Overtime submitted successfully", nil)
}

func (r *Rest) SubmitReimbursement(c echo.Context) error {
	userDetail, ok := c.Get("user").(entity.User)
	if !ok {
		return r.standardizeResponse(c, http.StatusUnauthorized, "User ID not found in context", nil)
	}

	var request entity.SubmitReimbursementRequest
	if err := c.Bind(&request); err != nil {
		return r.standardizeResponse(c, http.StatusBadRequest, "Invalid request format", nil)
	}

	err := r.employeeUc.SubmitReimbursement(userDetail, request)
	if err != nil {
		return r.standardizeResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return r.standardizeResponse(c, http.StatusOK, "Reimbursement submitted successfully", nil)
}

func (r *Rest) GetPayslip(c echo.Context) error {
	payrollPeriodID, err := strconv.Atoi(c.Param("period_id"))
	if err != nil {
		return r.standardizeResponse(c, http.StatusBadRequest, "Invalid ID format", nil)
	}

	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		return r.standardizeResponse(c, http.StatusBadRequest, "Invalid ID format", nil)
	}

	response, err := r.payrollUc.GetPayslip(int64(userID), int64(payrollPeriodID))
	if err != nil {
		return r.standardizeResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return r.standardizeResponse(c, http.StatusOK, "Success", response)
}

func (r *Rest) GetPayslips(c echo.Context) error {
	idParam := c.Param("period_id")
	periodID, err := strconv.Atoi(idParam)
	if err != nil {
		return r.standardizeResponse(c, http.StatusBadRequest, "Invalid ID format", nil)
	}

	response, err := r.payrollUc.GetPayslips(int64(periodID))
	if err != nil {
		return r.standardizeResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return r.standardizeResponse(c, http.StatusOK, "Success", response)
}

func (r *Rest) GeneratePayroll(c echo.Context) error {
	userDetail, ok := c.Get("user").(entity.User)
	if !ok {
		return r.standardizeResponse(c, http.StatusUnauthorized, "User ID not found in context", nil)
	}

	idParam := c.Param("period_id")
	periodID, err := strconv.Atoi(idParam)
	if err != nil {
		return r.standardizeResponse(c, http.StatusBadRequest, "Invalid ID format", nil)
	}

	err = r.payrollUc.ClosePayrollPeriod(userDetail, int64(periodID))
	if err != nil {
		return r.standardizeResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return r.standardizeResponse(c, http.StatusOK, "Success", nil)
}

func (r *Rest) GetPayslipBreakdown(c echo.Context) error {
	userDetail, ok := c.Get("user").(entity.User)
	if !ok {
		return r.standardizeResponse(c, http.StatusUnauthorized, "User ID not found in context", nil)
	}

	idParam := c.Param("period_id")
	periodID, err := strconv.Atoi(idParam)
	if err != nil {
		return r.standardizeResponse(c, http.StatusBadRequest, "Invalid ID format", nil)
	}

	response, err := r.employeeUc.GetPayslipBreakdown(userDetail, int64(periodID))
	if err != nil {
		return r.standardizeResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return r.standardizeResponse(c, http.StatusOK, "Success", response)
}
