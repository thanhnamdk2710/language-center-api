package domain

import (
	"context"
	"errors"
)

var (
	ErrOTPExpired = errors.New("OTP has expired")
	ErrOTPInvalid = errors.New("OTP is invalid")
)

type OTPService interface {
	Generate(ctx context.Context, email string) (string, error)
	Verify(ctx context.Context, email, otp string) error
	Delete(ctx context.Context, email string) error
}
