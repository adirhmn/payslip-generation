package app

import (
	"log"
	"time"

	v1 "payslip-generation-system/internal/controller/http/v1"
	grace "payslip-generation-system/internal/grace"

	"payslip-generation-system/internal/app/middleware"
	"payslip-generation-system/internal/httpclient"

	"payslip-generation-system/internal/postgres"

	"payslip-generation-system/config"

	// common

	// services
	adminsvc "payslip-generation-system/internal/services/admin"
	audsvc "payslip-generation-system/internal/services/audit"
	authsvc "payslip-generation-system/internal/services/auth"
	empsvc "payslip-generation-system/internal/services/employee"
	pingsvc "payslip-generation-system/internal/services/ping"

	// repositories
	attrepo "payslip-generation-system/internal/repositories/attendance"
	audrepo "payslip-generation-system/internal/repositories/audit"
	ovtrepo "payslip-generation-system/internal/repositories/overtime"
	payrepo "payslip-generation-system/internal/repositories/payslip"
	pingrepo "payslip-generation-system/internal/repositories/ping"
	reimbursrepo "payslip-generation-system/internal/repositories/reimbursement"
	userrepo "payslip-generation-system/internal/repositories/user"
)

type appHttp struct {
	middleware   middleware.HttpMdwProvider
	v1Controller v1.V1Controller
}

// RegisterHandlers registers the http handlers
func NewAppHTTP(
	config *config.Config,
) *appHttp {
	// postgres
	database, err := postgres.New(&postgres.Config{
		ServiceName:   config.Database.ServiceName,
		Dsn:           config.Database.DSN,
		MaxConn:       config.Database.MaxOpenConn,
		MaxIdle:       config.Database.MaxIdleConn,
	})
	if err != nil {
		log.Fatalf("error init postgres %s", err.Error())
	}

	// init repositories
	pingRepo := pingrepo.NewPingRepository(database)
	userRepo := userrepo.NewUserRepository(database)
	attendanceRepo := attrepo.NewAttendanceRepository(database)
	overtimeRepo := ovtrepo.NewOvertimeRepository(database)
	reimbursementRepo := reimbursrepo.NewReimbursementRepository(database)
	payslipRepo := payrepo.NewPayslipRepository(database)
	auditRepo := audrepo.NewAuditRepository(database)

	// init repositories
	pingService := pingsvc.NewPingService(pingRepo)
	authService:= authsvc.NewAuthService(userRepo, []byte(config.JWT.SecretKey))
	auditService:= audsvc.NewAuditService(auditRepo)
	adminService := adminsvc.NewAdminService(attendanceRepo, payslipRepo, reimbursementRepo, overtimeRepo, userRepo, auditService)
	employeeService := empsvc.NewEmployeeService(attendanceRepo, overtimeRepo, reimbursementRepo, payslipRepo, auditService)

	// init controllers
	v1Controller := v1.NewV1Controller(
		pingService,
		authService,
		adminService,
		employeeService,
	)

	middleware := middleware.NewMiddleWare(
		auditService,
	) 

	return &appHttp{
		middleware:   middleware,
		v1Controller: v1Controller,
	}
}

// Run runs the http app
func (a *appHttp) Run(config *config.Config) {
	// run http server
	grace.Serve(
		config.Port,
		a.RegisterHandlers(config),
	)
}

func SetupHttpClient(cfg *config.Config) httpclient.Client {
	httpClientCfg := &httpclient.Config{
		Timeout: cfg.HTTPClient.TimeoutMS,
		Transport: struct {
			DisableKeepAlives   bool
			MaxIdleConns        int
			MaxConnsPerHost     int
			MaxIdleConnsPerHost int
			IdleConnTimeout     time.Duration
		}{
			DisableKeepAlives:   cfg.HTTPClient.DisableKeepAlives,
			MaxIdleConns:        cfg.HTTPClient.MaxIdleConns,
			MaxConnsPerHost:     cfg.HTTPClient.MaxConnsPerHost,
			MaxIdleConnsPerHost: cfg.HTTPClient.MaxIdleConnsPerHost,
			IdleConnTimeout:     time.Duration(cfg.HTTPClient.IdleConnTimeout) * time.Second,
		},
	}

	return httpclient.New(httpClientCfg)
}
