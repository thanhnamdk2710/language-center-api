package usecase

import (
	"context"

	"github.com/thanhnamdk2710/auth-service/internal/domain"
)

type RefreshTokenUsecase struct{}

func NewRefreshTokenUsecase() RefreshTokenUsecase {
	return RefreshTokenUsecase{}
}

func (u *RefreshTokenUsecase) RefreshToken(ctx context.Context, user *domain.User) error {
	return nil
}
