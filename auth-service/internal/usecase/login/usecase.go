package login

import (
	"context"
	"fmt"

	"github.com/thanhnamdk2710/auth-service/internal/domain"
)

type Usecase interface {
	Execute(ctx context.Context, input Input) (*Output, error)
}

type service struct {
	userRepo    domain.UserRepository
	passwordSvc domain.PasswordService
	tokenSvc    domain.TokenService
}

func NewLoginUsecase(
	userRepo domain.UserRepository,
	passwordSvc domain.PasswordService,
	tokenSvc domain.TokenService,
) Usecase {
	return &service{
		userRepo:    userRepo,
		passwordSvc: passwordSvc,
		tokenSvc:    tokenSvc,
	}
}

func (s *service) Execute(ctx context.Context, input Input) (*Output, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, domain.ErrInvalidCredentials
	}

	// Verify password
	validPassword := s.passwordSvc.Verify(input.Password, user.Password)
	if !validPassword {
		return nil, domain.ErrInvalidCredentials
	}

	// Check if email is verified
	if user.Status == domain.UserStatusPending {
		return nil, domain.ErrEmailNotVerified
	}

	// Check if account is disabled
	if user.Status == domain.UserStatusDisabled {
		return nil, domain.ErrAccountDisabled
	}

	// Generate tokens
	accessToken, err := s.tokenSvc.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.tokenSvc.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &Output{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
