package login

import (
	"context"
	"fmt"
	"time"

	"github.com/thanhnamdk2710/auth-service/internal/domain/entity"
	"github.com/thanhnamdk2710/auth-service/internal/domain/repository"
	"github.com/thanhnamdk2710/auth-service/internal/domain/service/password"
	"github.com/thanhnamdk2710/auth-service/internal/domain/service/session"
	"github.com/thanhnamdk2710/auth-service/internal/domain/service/token"
	"github.com/thanhnamdk2710/auth-service/internal/domain/valueobject"
)

type Usecase interface {
	Execute(ctx context.Context, input Input) (*Output, error)
}

type service struct {
	userRepo    repository.UserRepository
	passwordSvc password.Service
	tokenSvc    token.Service
	sessionSvc  session.Service
}

func NewLoginUsecase(
	userRepo repository.UserRepository,
	passwordSvc password.Service,
	tokenSvc token.Service,
	sessionSvc session.Service,
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
		return nil, valueobject.ErrInvalidCredentials
	}

	// Check email verification and account status
	if user.Status == valueobject.UserStatusPending {
		return nil, valueobject.ErrEmailNotVerified
	}
	if user.Status == valueobject.UserStatusDisabled {
		return nil, valueobject.ErrAccountDisabled
	}

	// Check if account is locked
	if user.LockedUntil != nil && time.Now().Before(*user.LockedUntil) {
		return nil, valueobject.ErrAccountLocked
	}

	// Verify password
	validPassword := s.passwordSvc.Verify(input.Password, user.Password)
	if !validPassword {
		// Increment failed attemps
		if err := s.userRepo.IncrementFailedAttempts(ctx, user.ID.String()); err != nil {
			// Log error but continue
		}

		// Lock account after 5 failed attempts
		if user.FailedLoginAttempts >= 4 { // Will be 5 after increment
			lockDuration := 30 * time.Minute
			if err := s.userRepo.LockAccount(ctx, user.ID.String(), lockDuration); err != nil {
				// Log error
			}
			return nil, valueobject.ErrAccountLocked
		}

		return nil, valueobject.ErrInvalidCredentials
	}

	// Reset failed attempts on successful login
	if user.FailedLoginAttempts > 0 {
		if err := s.userRepo.ResetFailedAttempts(ctx, user.ID.String()); err != nil {
			// Log error but continue
		}
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

	// Store refresh token session
	session := &entity.Session{
		UserID:       user.ID.String(),
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
		CreatedAt:    time.Now(),
	}
	if err := s.sessionSvc.Store(ctx, session); err != nil {
		// Log error but continue
	}

	return &Output{
		UserID:       user.ID.String(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
