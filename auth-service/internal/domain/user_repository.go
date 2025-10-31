package domain

import "context"

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*User, error)
	VerifyEmail(ctx context.Context, userID string) error
	Create(ctx context.Context, user *User) error
}
