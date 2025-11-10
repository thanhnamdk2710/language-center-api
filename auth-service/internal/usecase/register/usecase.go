package register

import (
	"context"
	"errors"

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
	emailSvc    port.EmailService
	otpSvc      port.OTPService
}

func NewRegisterUsecase(
	userRepo repository.UserRepository,
	passwordSvc port.PasswordService,
	emailSvc port.EmailService,
	otpSvc port.OTPService,
) Usecase {
	return &service{
		userRepo:    userRepo,
		passwordSvc: passwordSvc,
		emailSvc:    emailSvc,
		otpSvc:      otpSvc,
	}
}

func (s *service) Execute(ctx context.Context, input Input) (*Output, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	// Check existing user
	existing, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		logger.Error("failed to check existing user", err)
		return nil, errors.New("registration failed, please try again later")
	}
	if existing != nil {
		// User already exists
		switch existing.Status {
		case valueobject.UserStatusPending:
			// Resend OTP for pending users
			otp, err := s.otpSvc.Generate(ctx, existing.Email.String())
			if err != nil {
				logger.Error("failed to generate OTP for resend", err)
				return nil, errors.New("failed to resend verification email, please try again later")
			}

			if err := s.emailSvc.SendVerificationEmail(existing.Email.String(), otp); err != nil {
				logger.Error("failed to send verification email for resend", err)
				return nil, errors.New("failed to send verification email, please try again later")
			}

			return &Output{
				UserID:  existing.ID.String(),
				Message: "Verification email resent. Please check your mail.",
			}, nil
		case valueobject.UserStatusActive:
			return nil, valueobject.ErrUserAlreadyExists
		case valueobject.UserStatusDisabled:
			return nil, errors.New("your account has been disabled, please contact support")
		}
	}

	// Validate password (business logic)
	if err := s.passwordSvc.Validate(input.Password); err != nil {
		return nil, err
	}

	// Hash password (business logic)
	passwordHash, err := s.passwordSvc.Hash(input.Password)
	if err != nil {
		logger.Error("failed to hash password", err)
		return nil, errors.New("registration failed, please try again later")
	}

	// Create user
	user, err := entity.NewUser(input.Email, passwordHash)
	if err != nil {
		logger.Error("failed to create user entity", err)
		return nil, errors.New("registration failed, please try again later")
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		logger.Error("failed to save user to database", err)
		return nil, errors.New("registration failed, please try again later")
	}

	// Generate OTP
	otp, err := s.otpSvc.Generate(ctx, user.Email.String())
	if err != nil {
		// Rollback: delete user
		if deleteErr := s.userRepo.Delete(ctx, user.ID.String()); deleteErr != nil {
			logger.Error("failed to rollback user creation", deleteErr)
		}
		logger.Error("failed to generate OTP", err)
		return nil, errors.New("registration failed, please try again later")
	}

	// Send verification email
	if err := s.emailSvc.SendVerificationEmail(user.Email.String(), otp); err != nil {
		// Rollback: delete user and OTP
		if deleteErr := s.userRepo.Delete(ctx, user.ID.String()); deleteErr != nil {
			logger.Error("failed to rollback user creation", deleteErr)
		}
		if deleteErr := s.otpSvc.Delete(ctx, user.Email.String()); deleteErr != nil {
			logger.Error("failed to delete OTP after email failure", deleteErr)
		}
		logger.Error("failed to send verification email", err)
		return nil, errors.New("failed to send verification email, please check your email address and try again")
	}

	logger.Info("user registered successfully")

	return &Output{
		UserID:  user.ID.String(),
		Message: "Registration successful! Please check your email for the verification code.",
	}, nil
}
