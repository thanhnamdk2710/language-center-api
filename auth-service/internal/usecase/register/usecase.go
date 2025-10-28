package register

import (
	"context"
	"fmt"

	"github.com/thanhnamdk2710/auth-service/internal/domain"
)

type Usecase interface {
	Execute(ctx context.Context, input Input) error
}

type serivce struct {
	userRepo domain.UserRepository
}

type Input struct {
	Email        string
	PasswordHash string
}

func NewRegisterUsecase(repo domain.UserRepository) Usecase {
	return &serivce{userRepo: repo}
}

func (s *serivce) Execute(ctx context.Context, input Input) error {
	existing, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return fmt.Errorf("failed to check existing user: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("email already registered")
	}

	user := &domain.User{
		Email:        input.Email,
		PasswordHash: input.PasswordHash,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Send OTP mail by Kafka

	return nil
}
