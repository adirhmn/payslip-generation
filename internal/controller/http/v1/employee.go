package v1

import (
	"context"
	"fmt"
	"net/http"
	serverctrl "payslip-generation-system/internal/controller/http"
	"payslip-generation-system/internal/entity/attendance"
	"payslip-generation-system/internal/entity/overtime"
	"payslip-generation-system/internal/entity/reimbursement"
	"time"

	"github.com/gin-gonic/gin"
)

func (v1 *v1Controller) SubmitAttendance(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c.Request.Context(), time.Second*10)
	defer cancelCtx()

	var req struct {
		Date     string `json:"date"`
		PeriodID int    `json:"period_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		serverctrl.ResponseHandler(c, http.StatusBadRequest, nil, fmt.Errorf("invalid input"))
		return
	}

	userID := c.GetInt("user_id")
	isAdmin := c.GetBool("is_admin")
	requestID := c.GetInt("request_log_id")
	if isAdmin {
		serverctrl.ResponseHandler(c, http.StatusForbidden, nil, fmt.Errorf("only employee can perform this action"))
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		serverctrl.ResponseHandler(c, http.StatusBadRequest, nil, fmt.Errorf("invalid input date"))
		return
	}
	attendance := attendance.Attendance{
		UserID:   userID,
		PeriodID: req.PeriodID,
		Date:     date,
	}
	_, err = v1.employeeService.SubmitAttendance(ctx, attendance, requestID)
	if err != nil {
		serverctrl.ResponseHandler(c, http.StatusBadRequest, nil, err)
		return
	}

	serverctrl.ResponseHandler(c, http.StatusOK, "attendance submitted", nil)
}

func (v1 *v1Controller) SubmitOvertime(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c.Request.Context(), time.Second*10)
	defer cancelCtx()

	var req struct {
		Date     string `json:"date"`
		PeriodID int    `json:"period_id"`
		Hours int    `json:"hours"`
		WorkCompleted bool    `json:"work_completed"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		serverctrl.ResponseHandler(c, http.StatusBadRequest, nil, fmt.Errorf("invalid input"))
		return
	}

	userID := c.GetInt("user_id")
	isAdmin := c.GetBool("is_admin")
	requestID := c.GetInt("request_log_id")
	if isAdmin {
		serverctrl.ResponseHandler(c, http.StatusForbidden, nil, fmt.Errorf("only employee can perform this action"))
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		serverctrl.ResponseHandler(c, http.StatusBadRequest, nil, fmt.Errorf("invalid input date"))
		return
	}

	if !req.WorkCompleted{
		serverctrl.ResponseHandler(c, http.StatusBadRequest, nil, fmt.Errorf("you must finish your work before submitting overtime"))
		return
	}

	overtime := overtime.Overtime{
		UserID:   userID,
		PeriodID: req.PeriodID,
		Date:     date,
		Hours:    req.Hours,
	}
	_, err = v1.employeeService.SubmitOvertime(ctx, overtime, requestID)
	if err != nil {
		serverctrl.ResponseHandler(c, http.StatusBadRequest, nil, err)
		return
	}

	serverctrl.ResponseHandler(c, http.StatusOK, "overtime submitted", nil)
}

func (v1 *v1Controller) SubmitReimbursement(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c.Request.Context(), time.Second*10)
	defer cancelCtx()

	var req struct {
		PeriodID int    `json:"period_id"`
		Amount int    `json:"amount"`
		Description string    `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		serverctrl.ResponseHandler(c, http.StatusBadRequest, nil, fmt.Errorf("invalid input"))
		return
	}

	userID := c.GetInt("user_id")
	isAdmin := c.GetBool("is_admin")
	requestID := c.GetInt("request_log_id")
	if isAdmin {
		serverctrl.ResponseHandler(c, http.StatusForbidden, nil, fmt.Errorf("only employee can perform this action"))
		return
	}

	if req.Amount <= 0 {
		serverctrl.ResponseHandler(c, http.StatusBadRequest, nil, fmt.Errorf("amount must be greater than 0"))
		return
	}

	reimbursement := reimbursement.Reimbursement{
		UserID:   userID,
		PeriodID: req.PeriodID,
		Amount:   req.Amount,
		Description: req.Description,
	}
	_, err := v1.employeeService.SubmitReimbursement(ctx, reimbursement, requestID)
	if err != nil {
		serverctrl.ResponseHandler(c, http.StatusBadRequest, nil, err)
		return
	}

	serverctrl.ResponseHandler(c, http.StatusOK, "reimbursement submitted", nil)
}

func (v1 *v1Controller) GeneratePayslips(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c.Request.Context(), time.Second*10)
	defer cancelCtx()

	userID := c.GetInt("user_id")
	isAdmin := c.GetBool("is_admin")
	if isAdmin {
		serverctrl.ResponseHandler(c, http.StatusForbidden, nil, fmt.Errorf("only employee can perform this action"))
		return
	}

	payslips, err := v1.employeeService.GeneratePayslips(ctx, userID)
	if err != nil {
		serverctrl.ResponseHandler(c, http.StatusBadRequest, nil, err)
		return
	}

	serverctrl.ResponseHandler(c, http.StatusOK, payslips, nil)
}