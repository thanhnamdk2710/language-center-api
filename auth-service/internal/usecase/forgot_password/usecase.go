package forgot_password

import (
	"context"
)

type Usecase interface {
	Execute(ctx context.Context) error
}

type service struct{}

func NewForgotPasswordUsecase() Usecase {
	return &service{}
}

func (u *service) Execute(ctx context.Context) error {
	return nil
}
