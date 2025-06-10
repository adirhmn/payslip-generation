package middleware

import (
	auditsvc "payslip-generation-system/internal/services/audit"

	"github.com/gin-gonic/gin"
)

type httpMiddleware struct {
	auditSvc auditsvc.AuditServiceProvider
}

// HttpMdwProvider is an interface for middleware
type HttpMdwProvider interface {
	JWTMiddleware(secret []byte) gin.HandlerFunc
	LoggingMiddleware() gin.HandlerFunc 
}

// NewMiddleWare is a constructor for middleware
func NewMiddleWare(
	auditService auditsvc.AuditServiceProvider,
) HttpMdwProvider {
	return &httpMiddleware{
		auditSvc: auditService,
	}
}