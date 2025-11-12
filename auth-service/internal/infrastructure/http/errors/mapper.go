package errors

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/thanhnamdk2710/auth-service/internal/domain/valueobject"
	"github.com/thanhnamdk2710/auth-service/internal/shared/validation"
)

type HTTPError struct {
	Status  int
	Code    string
	Message string
	Fields  map[string]string
}

var errorMap = map[error]HTTPError{
	// Login
	valueobject.ErrEmailNotVerified:   {http.StatusUnauthorized, "EMAIL_NOT_VERIFIED", "Please verify your email before logging in", nil},
	valueobject.ErrAccountDisabled:    {http.StatusUnauthorized, "ACCOUNT_DISABLED", "Your account has been disabled", nil},
	valueobject.ErrAccountLocked:      {http.StatusUnauthorized, "ACCOUNT_LOCKED", "Your account locked due to too many failed attempts", nil},
	valueobject.ErrEmailAlreadyExists: {http.StatusUnauthorized, "EMAIL_EXISTS", "Email already exists", nil},

	// Register
	valueobject.ErrInvalidCredentials: {http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password", nil},
	valueobject.ErrInvalidPassword:    {http.StatusUnauthorized, "INVALID_PASSWORD", "Password must contain uppercase, lowercase, number, and special character", nil},

	// Verify OTP
	valueobject.ErrOTPExpired:   {http.StatusUnauthorized, "OTP_EXPIRED", "Password must contain uppercase, lowercase, number, and special character", nil},
	valueobject.ErrOTPInvalid:   {http.StatusUnauthorized, "OTP_INVALID", "Password must contain uppercase, lowercase, number, and special character", nil},
	valueobject.ErrUserNotFound: {http.StatusNotFound, "USER_NOT_FOUND", "User not found", nil},
}

func MapError(err error) HTTPError {
	if err == nil {
		return HTTPError{http.StatusInternalServerError, "UNKNOWN_ERROR", "Unknown error", nil}
	}

	if verr, ok := err.(validator.ValidationErrors); ok {
		return HTTPError{
			Status:  http.StatusBadRequest,
			Code:    "VALIDATION_ERROR",
			Message: "Validation failed",
			Fields:  validation.ParseValidationErrors(verr),
		}
	}

	for domainErr, httpErr := range errorMap {
		if errors.Is(err, domainErr) {
			return httpErr
		}
	}

	return HTTPError{http.StatusInternalServerError, "SERVER_ERROR", "Internal server error", nil}
}
