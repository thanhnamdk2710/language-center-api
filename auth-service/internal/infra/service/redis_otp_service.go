package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/thanhnamdk2710/auth-service/internal/domain"
)

type redisOTPService struct {
	client     *redis.Client
	expiration time.Duration
}

func NewRedisOTPService(client *redis.Client, expiration time.Duration) domain.OTPService {
	return &redisOTPService{
		client:     client,
		expiration: expiration,
	}
}

func (s *redisOTPService) Generate(ctx context.Context, email string) (string, error) {
	// Generate OTP
	otp, err := GenerateOTP()
	if err != nil {
		return "", err
	}

	// Store in Redis with expiration (default 10 minutes)
	key := fmt.Sprintf("otp:%s", email)
	err = s.client.Set(ctx, key, otp, s.expiration).Err()
	if err != nil {
		return "", fmt.Errorf("failed to store OTP: %w", err)
	}

	return otp, nil
}

func (s *redisOTPService) Verify(ctx context.Context, email, otp string) error {
	key := fmt.Sprintf("otp:%s", email)

	// Get OTP from Redis
	storedOTP, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return domain.ErrOTPExpired
	}
	if err != nil {
		return fmt.Errorf("failed to get OTP: %w", err)
	}

	// Verify OTP
	if storedOTP != otp {
		return domain.ErrOTPInvalid
	}

	return nil
}

func (s *redisOTPService) Delete(ctx context.Context, email string) error {
	key := fmt.Sprintf("otp:%s", email)
	err := s.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete OTP: %w", err)
	}
	return nil
}

func GenerateOTP() (string, error) {
	// Generate a random 6-digit number
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP: %w", err)
	}

	return fmt.Sprintf("%06d", n.Int64()), nil
}
