package domain

import (
	"context"
	"time"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*User, error)
	VerifyEmail(ctx context.Context, userID string) error
	Create(ctx context.Context, user *User) error
	IncrementFailedAttempts(ctx context.Context, userID string) error
	ResetFailedAttempts(ctx context.Context, userID string) error
	LockAccount(ctx context.Context, userID string, duration time.Duration) error
}
