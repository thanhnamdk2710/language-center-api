package register

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
	emailSvc    domain.EmailService
	otpSvc      domain.OTPService
}

func NewRegisterUsecase(
	userRepo domain.UserRepository,
	passwordSvc domain.PasswordService,
	emailSvc domain.EmailService,
	otpSvc domain.OTPService,
) Usecase {
	return &service{
		userRepo:    userRepo,
		passwordSvc: passwordSvc,
		emailSvc:    emailSvc,
		otpSvc:      otpSvc,
	}
}

func (s *service) Execute(ctx context.Context, input Input) (*Output, error) {
	// Check existing user
	existing, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existing != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	// Validate password (business logic)
	if err := s.passwordSvc.Validate(input.Password); err != nil {
		return nil, err
	}

	// Hash password (business logic)
	passwordHash, err := s.passwordSvc.Hash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &domain.User{
		Email:        input.Email,
		PasswordHash: passwordHash,
		Status:       domain.UserStatusPending,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate OTP
	otp, err := s.otpSvc.Generate(ctx, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate OTP: %w", err)
	}

	// Send verification email
	if err := s.emailSvc.SendVerificationEmail(user.Email, otp); err != nil {
		return &Output{UserID: user.ID}, nil
	}

	return &Output{UserID: user.ID}, nil
}
