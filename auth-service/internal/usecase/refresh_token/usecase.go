package refresh_token

import (
	"context"
)

type Usecase interface {
	Execute(ctx context.Context) error
}

type service struct{}

func NewRefreshTokenUsecase() Usecase {
	return &service{}
}

func (u *service) Execute(ctx context.Context) error {
	return nil
}
