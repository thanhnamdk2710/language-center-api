package login

import (
	"context"
)

type Usecase interface {
	Execute(ctx context.Context) error
}

type service struct{}

func NewLoginUsecase() Usecase {
	return &service{}
}

func (u *service) Execute(ctx context.Context) error {
	return nil
}
