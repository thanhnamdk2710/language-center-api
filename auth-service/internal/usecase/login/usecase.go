package login

import (
	"context"
	"fmt"
	"time"

	"github.com/thanhnamdk2710/auth-service/internal/domain/entity"
	"github.com/thanhnamdk2710/auth-service/internal/domain/port"
	"github.com/thanhnamdk2710/auth-service/internal/domain/repository"
	"github.com/thanhnamdk2710/auth-service/internal/domain/valueobject"
	"github.com/thanhnamdk2710/auth-service/internal/shared/logger"
)

type Usecase interface {
	Execute(ctx context.Context, input Input) (*Output, error)
}

type service struct {
	userRepo    repository.UserRepository
	passwordSvc port.PasswordService
	tokenSvc    port.TokenService
	sessionSvc  port.SessionService
}

func NewLoginUsecase(
	userRepo repository.UserRepository,
	passwordSvc port.PasswordService,
	tokenSvc port.TokenService,
	sessionSvc port.SessionService,
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

	now := time.Now()

	if err := user.EnsureCanLogin(now); err != nil {
		return nil, err
	}

	// Verify password
	if !s.passwordSvc.Verify(input.Password, user.Password) {
		if isLocked := user.RecordFailedLogin(now); isLocked {
			if updateErr := s.userRepo.Update(ctx, user); updateErr != nil {
				logger.Error("failed to update user", updateErr)
			}
			return nil, valueobject.ErrAccountLocked
		}

		if err := s.userRepo.Update(ctx, user); err != nil {
			logger.Error("failed to update user", err)
		}
		return nil, valueobject.ErrInvalidCredentials
	}

	// Reset failed attempts on successful login
	if user.FailedLoginAttempts > 0 {
		if err := s.userRepo.Update(ctx, user); err != nil {
			logger.Error("failed to update user", err)
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
	session, err := entity.NewSession(user.ID, refreshToken, now.Add(7*24*time.Hour))
	if err != nil {
		return nil, fmt.Errorf("failed to initial refresh token: %w", err)
	}

	if err := s.sessionSvc.Store(ctx, session); err != nil {
		logger.Error("failed to store session", err)
		// Log error but continue
	}

	return &Output{
		UserID:       user.ID.String(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
