package middleware

import (
	"fmt"
	"payslip-generation-system/internal/entity/audit"

	"github.com/gin-gonic/gin"
)

func (m *httpMiddleware) LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
        log := audit.RequestLog{
            URL: c.Request.URL.String(),
            Method: c.Request.Method,
            IPAddress: c.ClientIP(),
        }
        requestLogID, err:= m.auditSvc.RecordRequestLog(c, log)
        if err != nil {
            fmt.Println(err)
        }
        c.Set("request_log_id", requestLogID)

        c.Next()
    }
}