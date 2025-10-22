package usecase

import (
	"context"

	"github.com/thanhnamdk2710/auth-service/internal/domain"
)

type LoginUsecase struct{}

func NewLoginUsecase() LoginUsecase {
	return LoginUsecase{}
}

func (u *LoginUsecase) Login(ctx context.Context, user *domain.User) error {
	return nil
}
