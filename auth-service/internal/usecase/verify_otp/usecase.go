package verify_otp

import (
	"context"
	"fmt"

	"github.com/thanhnamdk2710/auth-service/internal/domain/repository"
	"github.com/thanhnamdk2710/auth-service/internal/domain/service/otp"
	"github.com/thanhnamdk2710/auth-service/internal/domain/service/token"
	"github.com/thanhnamdk2710/auth-service/internal/domain/valueobject"
)

type Usecase interface {
	Execute(ctx context.Context, input Input) (*Output, error)
}

type service struct {
	userRepo repository.UserRepository
	otpSvc   otp.Service
	tokenSvc token.Service
}

func NewVerifyOTPUsecase(
	userRepo repository.UserRepository,
	otpSvc otp.Service,
	tokenSvc token.Service,
) Usecase {
	return &service{
		userRepo: userRepo,
		otpSvc:   otpSvc,
		tokenSvc: tokenSvc,
	}
}

func (s *service) Execute(ctx context.Context, input Input) (*Output, error) {
	// Verify OTP
	if err := s.otpSvc.Verify(ctx, input.Email, input.OTP); err != nil {
		return nil, err
	}

	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, valueobject.ErrUserNotFound
	}

	// Update user status to active
	if err := s.userRepo.VerifyEmail(ctx, user.ID); err != nil {
		return nil, fmt.Errorf("failed to verify email: %w", err)
	}

	// Delete OTP after successful verification
	if err := s.otpSvc.Delete(ctx, input.Email); err != nil {
		// Log error but don't fail
		fmt.Printf("failed to delete OTP: %v\n", err)
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
