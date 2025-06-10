package v1

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	serverctrl "payslip-generation-system/internal/controller/http"
	"payslip-generation-system/internal/entity/attendance"

	"github.com/gin-gonic/gin"
)

func (v1 *v1Controller) AddAttendancePeriod(c *gin.Context){
	ctx, cancelCtx := context.WithTimeout(c.Request.Context(), time.Second*10)
	defer cancelCtx()

	var req struct {
        StartDate string `json:"start_date"`
        EndDate   string `json:"end_date"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
		serverctrl.ResponseHandler(c, http.StatusBadRequest, nil, fmt.Errorf("invalid input"))
        return
    }

    userID := c.GetInt("user_id")
    isAdmin := c.GetBool("is_admin")
    requestID := c.GetInt("request_log_id")
    if !isAdmin {
		serverctrl.ResponseHandler(c, http.StatusForbidden, nil, fmt.Errorf("only admin can perform this action"))
        return
    }

    startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		serverctrl.ResponseHandler(c, http.StatusBadRequest, nil, fmt.Errorf("invalid input start_date"))
		return
	}

    endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		serverctrl.ResponseHandler(c, http.StatusBadRequest, nil, fmt.Errorf("invalid input end_date"))
		return
	}

    attendancePeriod := attendance.AttendancePeriod{
        StartDate: startDate,
        EndDate: endDate,
	}
    _, err = v1.adminService.AddPeriod(ctx, attendancePeriod, userID,requestID)
    if err != nil {
        serverctrl.ResponseHandler(c, http.StatusBadRequest, nil, err)
        return
    }

	serverctrl.ResponseHandler(c, http.StatusOK,"attendance period created", nil)
}

func (v1 *v1Controller) RunPayroll(c *gin.Context){
	ctx, cancelCtx := context.WithTimeout(c.Request.Context(), time.Second*10)
	defer cancelCtx()

	var req struct {
        PeriodID int `json:"period_id"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
		serverctrl.ResponseHandler(c, http.StatusBadRequest, nil, fmt.Errorf("invalid input"))
        return
    }

    userID := c.GetInt("user_id")
    isAdmin := c.GetBool("is_admin")
    requestID := c.GetInt("request_log_id")
    if !isAdmin {
		serverctrl.ResponseHandler(c, http.StatusForbidden, nil, fmt.Errorf("only admin can perform this action"))
        return
    }

    err := v1.adminService.RunPayroll(ctx, req.PeriodID, userID, requestID)
    if err != nil {
		serverctrl.ResponseHandler(c, http.StatusBadRequest, nil, err)
        return
    }

	serverctrl.ResponseHandler(c, http.StatusOK,"payroll already run", nil)
}

func (v1 *v1Controller) GetPayslipSummary(c *gin.Context){
	ctx, cancelCtx := context.WithTimeout(c.Request.Context(), time.Second*10)
	defer cancelCtx()

	periodIDStr := c.Param("period_id")

	periodID, err := strconv.Atoi(periodIDStr)
	if err != nil {
		serverctrl.ResponseHandler(c, http.StatusBadRequest, nil, fmt.Errorf("invalid period_id"))
		return
	}

    isAdmin := c.GetBool("is_admin")
    if !isAdmin {
		serverctrl.ResponseHandler(c, http.StatusForbidden, nil, fmt.Errorf("only admin can perform this action"))
        return
    }

    summary, err := v1.adminService.GetPayslipSummary(ctx, periodID)
    if err != nil {
		serverctrl.ResponseHandler(c, http.StatusBadRequest, nil, err)
        return
    }

	serverctrl.ResponseHandler(c, http.StatusOK, summary, nil)
}

