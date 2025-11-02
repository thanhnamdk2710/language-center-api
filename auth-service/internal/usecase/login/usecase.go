package login

import (
	"context"
	"fmt"
	"time"

	"github.com/thanhnamdk2710/auth-service/internal/domain"
)

type Usecase interface {
	Execute(ctx context.Context, input Input) (*Output, error)
}

type service struct {
	userRepo    domain.UserRepository
	passwordSvc domain.PasswordService
	tokenSvc    domain.TokenService
	sessionSvc  domain.SessionService
}

func NewLoginUsecase(
	userRepo domain.UserRepository,
	passwordSvc domain.PasswordService,
	tokenSvc domain.TokenService,
	sessionSvc domain.SessionService,
) Usecase {
	return &service{
		userRepo:    userRepo,
		passwordSvc: passwordSvc,
		tokenSvc:    tokenSvc,
		sessionSvc:  sessionSvc,
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

	// Check email verification and account status
	if user.Status == domain.UserStatusPending {
		return nil, domain.ErrEmailNotVerified
	}
	if user.Status == domain.UserStatusDisabled {
		return nil, domain.ErrAccountDisabled
	}

	// Check if account is locked
	if user.LockedUntil != nil && time.Now().Before(*user.LockedUntil) {
		return nil, domain.ErrAccountLocked
	}

	// Verify password
	validPassword := s.passwordSvc.Verify(input.Password, user.Password)
	if !validPassword {
		// Increment failed attemps
		if err := s.userRepo.IncrementFailedAttempts(ctx, user.ID); err != nil {
			// Log error but continue
		}

		// Lock account after 5 failed attempts
		if user.FailedLoginAttempts >= 4 { // Will be 5 after increment
			lockDuration := 30 * time.Minute
			if err := s.userRepo.LockAccount(ctx, user.ID, lockDuration); err != nil {
				// Log error
			}
			return nil, domain.ErrAccountLocked
		}

		return nil, domain.ErrInvalidCredentials
	}

	// Reset failed attempts on successful login
	if user.FailedLoginAttempts > 0 {
		if err := s.userRepo.ResetFailedAttempts(ctx, user.ID); err != nil {
			// Log error but continue
		}
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

	// Store refresh token session
	session := &domain.Session{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
		CreatedAt:    time.Now(),
	}
	if err := s.sessionSvc.Store(ctx, session); err != nil {
		// Log error but continue
	}

	return &Output{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
