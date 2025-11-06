package valueobject

import "errors"

var (
	// User
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailNotVerified   = errors.New("email not verified")
	ErrAccountDisabled    = errors.New("account has been disabled")
	ErrAccountLocked      = errors.New("account is locked due to too many failed login attempts")

	// Token
	ErrInvalidToken = errors.New("invalid token")
	ErrOTPExpired   = errors.New("OTP has expired")
	ErrOTPInvalid   = errors.New("OTP is invalid")

	// Password
	ErrInvalidPassword = errors.New("password must contain uppercase, lowercase, number, and special character")
)
