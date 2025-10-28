package domain

import "context"

type UserRepository interface {
	GetByEmail(ctx context.Context, emial string) (*User, error)
	Create(ctx context.Context, user *User) error
}
