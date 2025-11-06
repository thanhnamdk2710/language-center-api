package repository

import (
	"context"
	"time"

	"github.com/thanhnamdk2710/auth-service/internal/domain/entity"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	VerifyEmail(ctx context.Context, userID string) error
	Create(ctx context.Context, user *entity.User) error
	IncrementFailedAttempts(ctx context.Context, userID string) error
	ResetFailedAttempts(ctx context.Context, userID string) error
	LockAccount(ctx context.Context, userID string, duration time.Duration) error
}
