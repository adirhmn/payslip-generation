package app

import (
	"payslip-generation-system/config"

	"github.com/gin-gonic/gin"
)

// RegisterHandler registers all the handlers
func (a *appHttp) RegisterHandlers(config *config.Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	handler := gin.New()

	// middlewares
	// handler.Use(gintrace.Middleware(config.App))
	// handler.Use(gin.Logger())
	// handler.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedPaths([]string{"/metrics"})))
	// handler.Use(gin.Recovery())
	// handler.Use(a.middleware.PrometheusMiddleware())


	// v1 group
	a.registerV1(handler.Group("v1"), config)

	return handler
}

func (a *appHttp) registerV1(r *gin.RouterGroup, cfg *config.Config) {
	r.GET("/ping", a.v1Controller.Ping)
	r.POST("/login", a.v1Controller.Login, a.middleware.LoggingMiddleware())
	employeeGroup := r.Group("/employee")
	employeeGroup.Use(a.middleware.LoggingMiddleware())
	employeeGroup.Use(a.middleware.JWTMiddleware([]byte(cfg.JWT.SecretKey)))
	employeeGroup.POST("/submit-attendance", a.v1Controller.SubmitAttendance)
	employeeGroup.POST("/submit-overtime", a.v1Controller.SubmitOvertime)
	employeeGroup.POST("/submit-reimbursement", a.v1Controller.SubmitReimbursement)
	employeeGroup.GET("/generate-payslips", a.v1Controller.GeneratePayslips)

	adminGroup := r.Group("/admin")
	employeeGroup.Use(a.middleware.LoggingMiddleware())
	adminGroup.Use(a.middleware.JWTMiddleware([]byte(cfg.JWT.SecretKey)))
	adminGroup.POST("/add-attendance-period", a.v1Controller.AddAttendancePeriod)
	adminGroup.POST("/run-payroll", a.v1Controller.RunPayroll)
	adminGroup.GET("/get-payslip-summary/:period_id", a.v1Controller.GetPayslipSummary)
}