package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CheckUniqueRequest(c *gin.Context) {

	// err := h.service.ValidateHeaders(c.Request.Header)
	// if err != nil {

	// 	var appErr *apperror.AppError
	// 	if errors.As(err, &appErr) {
	// 		c.AbortWithStatusJSON(400, gin.H{
	// 			"responseCode":    appErr.Code,
	// 			"responseMessage": appErr.Message,
	// 		})
	// 		return
	// 	}

	// 	c.AbortWithStatusJSON(500, gin.H{
	// 		"responseCode":    "5000000",
	// 		"responseMessage": "Internal error",
	// 	})
	// 	return
	// }
	fmt.Println("Middleware: CheckUniqueRequest")

	c.Next()
}
