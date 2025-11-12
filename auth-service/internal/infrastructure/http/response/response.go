package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorData  `json:"errors,omitempty"`
}

type ErrorData struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func JSON(c *gin.Context, status int, success bool, data interface{}, errData *ErrorData) {
	c.JSON(status, Response{
		Success: success,
		Data:    data,
		Error:   errData,
	})
}

// Success response
func Success(c *gin.Context, status int, data interface{}) {
	JSON(c, status, true, data, nil)
}

// Error response
func Error(c *gin.Context, status int, code, message string) {
	JSON(c, status, false, nil, &ErrorData{Code: code, Message: message})
}

func ValidationError(c *gin.Context, fieldErrors map[string]string) {
	JSON(c, http.StatusBadRequest, false, nil, &ErrorData{
		Code:    "VALIDATION_ERROR",
		Message: "Validation failed",
		Fields:  fieldErrors,
	})
}
