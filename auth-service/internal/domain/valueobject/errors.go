package valueobject

import "errors"

var (
	// User
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserNotInPending   = errors.New("user is not in pending status")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailNotVerified   = errors.New("email not verified")
	ErrAccountDisabled    = errors.New("account has been disabled")
	ErrAccountLocked      = errors.New("account is locked due to too many failed login attempts")

	// Token
	ErrInvalidToken = errors.New("invalid token")
	ErrOTPExpired   = errors.New("otp has expired")
	ErrOTPInvalid   = errors.New("otp is invalid")

	// Session
	ErrSessionExpired = errors.New("session has expired")
	ErrSessionInvalid = errors.New("invalid session: empty refresh token")

	// Password
	ErrInvalidPassword = errors.New("password must contain uppercase, lowercase, number, and special character")
)
