package response

import (
	"github.com/gin-gonic/gin"
)

// Success response
func Success(c *gin.Context, status int, data interface{}) {
	c.JSON(status, gin.H{
		"data": data,
	})
}

// Error response
func Error(c *gin.Context, status int, code string, message interface{}) {
	c.JSON(status, gin.H{
		"error": gin.H{
			"code":    code,
			"message": message,
		},
	})
}
