package usecase

import (
	"context"

	"github.com/thanhnamdk2710/auth-service/internal/domain"
)

type RegisterUsecase struct{}

func NewRegisterUsecase() RegisterUsecase {
	return RegisterUsecase{}
}

func (u *RegisterUsecase) Register(ctx context.Context, user *domain.User) error {
	return nil
}
