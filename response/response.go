package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Total   interface{} `json:"total,omitempty"`
}

func Success(c *gin.Context, data interface{}, total interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Success: true,
		Message: "success",
		Data:    data,
		Total:   total,
	})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"success": false,
		"message": message,
	})
}
