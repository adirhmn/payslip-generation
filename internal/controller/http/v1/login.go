package v1

import (
	"context"
	"net/http"
	"time"

	serverctrl "payslip-generation-system/internal/controller/http"
	transportmodel "payslip-generation-system/internal/entity/transport"

	"github.com/gin-gonic/gin"
)

func (v1 *v1Controller) Login(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c.Request.Context(), time.Second*10)
	defer cancelCtx()

	var req struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
		serverctrl.ResponseHandler(c, http.StatusBadRequest, nil, err)
        return
    }

    token, err := v1.authService.Login(ctx, req.Username, req.Password)
    if err != nil {
		serverctrl.ResponseHandler(c, http.StatusUnauthorized, nil, err)
        return
    }

	serverctrl.ResponseHandler(c, http.StatusOK, transportmodel.LoginResponse{Token: token}, nil)
}