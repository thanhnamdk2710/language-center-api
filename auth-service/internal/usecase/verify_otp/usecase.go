package verify_otp

import (
	"context"
	"fmt"
	"time"

	"github.com/thanhnamdk2710/auth-service/internal/domain/port"
	"github.com/thanhnamdk2710/auth-service/internal/domain/repository"
	"github.com/thanhnamdk2710/auth-service/internal/domain/valueobject"
)

type Usecase interface {
	Execute(ctx context.Context, input Input) (*Output, error)
}

type service struct {
	userRepo repository.UserRepository
	otpSvc   port.OTPService
	tokenSvc port.TokenService
}

func NewVerifyOTPUsecase(
	userRepo repository.UserRepository,
	otpSvc port.OTPService,
	tokenSvc port.TokenService,
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

	if err := user.Activate(time.Now()); err != nil {
		return nil, err
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	// Delete OTP after successful verification
	if err := s.otpSvc.Delete(ctx, input.Email); err != nil {
		// Log error but don't fail
		fmt.Printf("failed to delete OTP: %v\n", err)
	}

	// Generate tokens
	accessToken, err := s.tokenSvc.GenerateAccessToken(user.ID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.tokenSvc.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &Output{
		UserID:       user.ID.String(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
