package repository

import (
	"context"

	"github.com/thanhnamdk2710/auth-service/internal/domain/entity"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
}
