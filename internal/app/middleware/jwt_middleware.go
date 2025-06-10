package middleware

import (
	"fmt"
	"net/http"
	"strings"

	serverctrl "payslip-generation-system/internal/controller/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (m *httpMiddleware) JWTMiddleware(secret []byte) gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenStr := c.GetHeader("Authorization")
        if tokenStr == "" {
			serverctrl.ResponseHandler(c, http.StatusUnauthorized, nil, fmt.Errorf("missing token"))
            c.Abort()
            return
        }

        parts := strings.Split(tokenStr, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			serverctrl.ResponseHandler(c, http.StatusUnauthorized, nil, fmt.Errorf("invalid authorization format"))
			c.Abort()
			return
		}

        token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
            return secret, nil
        })

        if err != nil || !token.Valid {
			serverctrl.ResponseHandler(c, http.StatusUnauthorized, nil, fmt.Errorf("invalid token"))
            c.Abort()
            return
        }

        claims := token.Claims.(jwt.MapClaims)
        c.Set("user_id", int(claims["user_id"].(float64)))
        c.Set("is_admin", claims["is_admin"].(bool))

        c.Next()
    }
}
