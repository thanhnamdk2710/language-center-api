package otp

import (
	"context"
)

type Service interface {
	Generate(ctx context.Context, email string) (string, error)
	Verify(ctx context.Context, email, otp string) error
	Delete(ctx context.Context, email string) error
}
