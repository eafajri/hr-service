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
	userDetail, ok := c.Get("user").(entity.User)
	if !ok {
		return r.standardizeResponse(c, http.StatusUnauthorized, "User ID not found in context", nil)
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return r.standardizeResponse(c, http.StatusBadRequest, "Invalid ID format", nil)
	}

	payrollPeriodID, err := strconv.Atoi(c.QueryParam("payroll_period_id"))
	if err != nil {
		return r.standardizeResponse(c, http.StatusBadRequest, "Invalid payroll period", nil)
	}

	if userDetail.ID != int64(id) {
		return r.standardizeResponse(c, http.StatusForbidden, "Access denied: you can only access your own payslip", nil)
	}

	response, err := r.payrollUc.GetPayslip(userDetail.ID, int64(payrollPeriodID))
	if err != nil {
		return r.standardizeResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return r.standardizeResponse(c, http.StatusOK, "Reimbursement submitted successfully", response)
}
