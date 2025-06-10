package v1

import (
	modeltransport "payslip-generation-system/internal/entity/transport"

	"github.com/gin-gonic/gin"
)

func ResponseHandler(c *gin.Context, httpStatus int, data interface{}, err error) {
	if err != nil {
		c.JSON(httpStatus, modeltransport.StandardResponse{
			Success: false,
			Error:   err.Error(),
			Data:    data,
		})
		return
	}
	c.JSON(httpStatus, modeltransport.StandardResponse{
		Success: true,
		Error:   "",
		Data:    data,
	})
}