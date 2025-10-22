package usecase

import (
	"context"

	"github.com/thanhnamdk2710/auth-service/internal/domain"
)

type ForgotPasswordUsecase struct{}

func NewForgotPasswordUsecase() ForgotPasswordUsecase {
	return ForgotPasswordUsecase{}
}

func (u *ForgotPasswordUsecase) ForgotPassword(ctx context.Context, user *domain.User) error {
	return nil
}
