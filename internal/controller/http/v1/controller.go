package v1

import (
	"github.com/gin-gonic/gin"

	// services
	adminsvc "payslip-generation-system/internal/services/admin"
	authsvc "payslip-generation-system/internal/services/auth"
	empsvc "payslip-generation-system/internal/services/employee"
	pingsvc "payslip-generation-system/internal/services/ping"
)

type V1Controller interface {
	Ping(c *gin.Context)
	Login(c *gin.Context)
	AddAttendancePeriod(c *gin.Context)
	SubmitAttendance(c *gin.Context)
	SubmitOvertime(c *gin.Context)
	SubmitReimbursement(c *gin.Context)
	RunPayroll(c *gin.Context)
	GetPayslipSummary(c *gin.Context)
	GeneratePayslips(c *gin.Context)
}

type v1Controller struct {
	pingService                   pingsvc.PingServiceProvider
	authService authsvc.AuthServiceProvider
	adminService adminsvc.AdminServiceProvider
	employeeService empsvc.EmployeeServiceProvider
}

func NewV1Controller(
	pingService pingsvc.PingServiceProvider,
	authService authsvc.AuthServiceProvider,
	adminService adminsvc.AdminServiceProvider,
	employeeService empsvc.EmployeeServiceProvider,
) V1Controller {
	return &v1Controller{
		pingService:                   pingService,
		authService: authService,
		adminService: adminService,
		employeeService: employeeService,
	}
}