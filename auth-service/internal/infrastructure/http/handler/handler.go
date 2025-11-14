package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/thanhnamdk2710/auth-service/internal/domain/valueobject"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/http/response"
	"github.com/thanhnamdk2710/auth-service/internal/shared/validation"
)

// BaseHandler provides common handler functionality
type BaseHandler struct{}

// NewBaseHandler creates a new base handler
func NewBaseHandler() BaseHandler {
	return BaseHandler{}
}

// errorMapping maps domain errors to HTTP error responses
var errorMapping = map[error]struct {
	status  int
	code    string
	message string
}{
	// Authentication errors
	valueobject.ErrEmailNotVerified:   {http.StatusUnauthorized, "EMAIL_NOT_VERIFIED", "Please verify your email before logging in"},
	valueobject.ErrAccountDisabled:    {http.StatusUnauthorized, "ACCOUNT_DISABLED", "Your account has been disabled"},
	valueobject.ErrAccountLocked:      {http.StatusUnauthorized, "ACCOUNT_LOCKED", "Your account locked due to too many failed attempts"},
	valueobject.ErrInvalidCredentials: {http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password"},

	// Registration errors
	valueobject.ErrEmailAlreadyExists: {http.StatusConflict, "EMAIL_EXISTS", "Email already exists"},
	valueobject.ErrInvalidPassword:    {http.StatusBadRequest, "INVALID_PASSWORD", "Password must contain uppercase, lowercase, number, and special character"},

	// OTP errors
	valueobject.ErrOTPExpired: {http.StatusUnauthorized, "OTP_EXPIRED", "OTP has expired"},
	valueobject.ErrOTPInvalid: {http.StatusUnauthorized, "OTP_INVALID", "Invalid OTP"},

	// User errors
	valueobject.ErrUserNotFound:     {http.StatusNotFound, "USER_NOT_FOUND", "User not found"},
	valueobject.ErrUserNotInPending: {http.StatusBadRequest, "USER_NOT_IN_PENDING", "User is not in pending status"},

	// Session errors
	valueobject.ErrSessionExpired: {http.StatusUnauthorized, "SESSION_EXPIRED", "Session has expired"},
	valueobject.ErrSessionInvalid: {http.StatusUnauthorized, "SESSION_INVALID", "Invalid session"},
}

// HandleError converts domain errors to HTTP responses
func (h *BaseHandler) HandleError(c *gin.Context, err error) {
	if err == nil {
		response.Error(c, http.StatusInternalServerError, "UNKNOWN_ERROR", "Unknown error")
		return
	}

	// Handle validation errors
	if verr, ok := err.(validator.ValidationErrors); ok {
		response.ValidationError(c, validation.ParseValidationErrors(verr))
		return
	}

	// Handle domain errors
	for domainErr, httpErr := range errorMapping {
		if errors.Is(err, domainErr) {
			response.Error(c, httpErr.status, httpErr.code, httpErr.message)
			return
		}
	}

	// Default to internal server error
	response.Error(c, http.StatusInternalServerError, "SERVER_ERROR", "Internal server error")
}
