package domain

import "errors"

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrEmailNotVerified   = errors.New("email not verified")
	ErrAccountDisabled    = errors.New("account has been disables")
	ErrAccountLocked      = errors.New("account is locked due to many failed login attempts")
)
