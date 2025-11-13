package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse represents the standard API response structure
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorData  `json:"errors,omitempty"`
}

// ErrorData represents error details in the response
type ErrorData struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

// JSON sends a JSON response with the standard structure
func JSON(c *gin.Context, status int, success bool, data interface{}, errData *ErrorData) {
	c.JSON(status, APIResponse{
		Success: success,
		Data:    data,
		Error:   errData,
	})
}

// Success sends a successful response
func Success(c *gin.Context, status int, data interface{}) {
	JSON(c, status, true, data, nil)
}

// Error sends an error response
func Error(c *gin.Context, status int, code, message string) {
	JSON(c, status, false, nil, &ErrorData{
		Code:    code,
		Message: message,
	})
}

// ValidationError sends a validation error response
func ValidationError(c *gin.Context, fieldErrors map[string]string) {
	JSON(c, http.StatusBadRequest, false, nil, &ErrorData{
		Code:    "VALIDATION_ERROR",
		Message: "Validation failed",
		Fields:  fieldErrors,
	})
}
