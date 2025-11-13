package presenter

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/thanhnamdk2710/auth-service/internal/domain/valueobject"
	"github.com/thanhnamdk2710/auth-service/internal/shared/validation"
)

// HTTPError represents an HTTP error response
type HTTPError struct {
	Status  int               `json:"-"`
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

// ErrorPresenter converts domain errors to HTTP errors
type ErrorPresenter struct {
	errorMap map[error]HTTPError
}

func NewErrorPresenter() *ErrorPresenter {
	return &ErrorPresenter{
		errorMap: map[error]HTTPError{
			// Authentication errors
			valueobject.ErrEmailNotVerified:   {http.StatusUnauthorized, "EMAIL_NOT_VERIFIED", "Please verify your email before logging in", nil},
			valueobject.ErrAccountDisabled:    {http.StatusUnauthorized, "ACCOUNT_DISABLED", "Your account has been disabled", nil},
			valueobject.ErrAccountLocked:      {http.StatusUnauthorized, "ACCOUNT_LOCKED", "Your account locked due to too many failed attempts", nil},
			valueobject.ErrInvalidCredentials: {http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password", nil},
			
			// Registration errors
			valueobject.ErrEmailAlreadyExists: {http.StatusConflict, "EMAIL_EXISTS", "Email already exists", nil},
			valueobject.ErrInvalidPassword:    {http.StatusBadRequest, "INVALID_PASSWORD", "Password must contain uppercase, lowercase, number, and special character", nil},
			
			// OTP errors
			valueobject.ErrOTPExpired: {http.StatusUnauthorized, "OTP_EXPIRED", "OTP has expired", nil},
			valueobject.ErrOTPInvalid: {http.StatusUnauthorized, "OTP_INVALID", "Invalid OTP", nil},
			
			// User errors
			valueobject.ErrUserNotFound:     {http.StatusNotFound, "USER_NOT_FOUND", "User not found", nil},
			valueobject.ErrUserNotInPending: {http.StatusBadRequest, "USER_NOT_IN_PENDING", "User is not in pending status", nil},
			
			// Session errors
			valueobject.ErrSessionExpired: {http.StatusUnauthorized, "SESSION_EXPIRED", "Session has expired", nil},
			valueobject.ErrSessionInvalid: {http.StatusUnauthorized, "SESSION_INVALID", "Invalid session", nil},
		},
	}
}

// PresentError converts any error to HTTPError
func (p *ErrorPresenter) PresentError(err error) HTTPError {
	if err == nil {
		return HTTPError{http.StatusInternalServerError, "UNKNOWN_ERROR", "Unknown error", nil}
	}

	// Handle validation errors
	if verr, ok := err.(validator.ValidationErrors); ok {
		return HTTPError{
			Status:  http.StatusBadRequest,
			Code:    "VALIDATION_ERROR",
			Message: "Validation failed",
			Fields:  validation.ParseValidationErrors(verr),
		}
	}

	// Handle domain errors
	for domainErr, httpErr := range p.errorMap {
		if errors.Is(err, domainErr) {
			return httpErr
		}
	}

	// Default to internal server error
	return HTTPError{http.StatusInternalServerError, "SERVER_ERROR", "Internal server error", nil}
}
