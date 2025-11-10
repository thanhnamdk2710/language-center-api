package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/thanhnamdk2710/auth-service/internal/domain/entity"
	"github.com/thanhnamdk2710/auth-service/internal/domain/port"
)

type sessionService struct {
	client *redis.Client
}

func NewSessionService(client *redis.Client) port.SessionService {
	return &sessionService{
		client: client,
	}
}

func (s *sessionService) Store(ctx context.Context, sess *entity.Session) error {
	key := fmt.Sprintf("session:%s", sess.RefreshToken)
	userKey := fmt.Sprintf("user_sessions:%s", sess.UserID)

	// Serialize session
	data, err := json.Marshal(sess)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	// Calculate TTL
	ttl := time.Until(sess.ExpiresAt)
	if ttl <= 0 {
		return fmt.Errorf("session already expired")
	}

	// Store session by refresh token
	if err := s.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to store session: %w", err)
	}

	// Add to user's session set (for logout all devices)
	if err := s.client.SAdd(ctx, userKey, sess.RefreshToken).Err(); err != nil {
		return fmt.Errorf("failed to add to user session: %w", err)
	}

	// Set expiration on user session set
	if err := s.client.Expire(ctx, userKey, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set expiration on user sessions: %w", err)
	}

	return nil
}

func (s *sessionService) Get(ctx context.Context, refreshToken string) (*entity.Session, error) {
	key := fmt.Sprintf("session:%s", refreshToken)

	data, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Session not found or expired
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	var sess entity.Session
	if err := json.Unmarshal([]byte(data), &sess); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	return &sess, nil
}

func (s *sessionService) Delete(ctx context.Context, refreshToken string) error {
	// Get session first to get user ID
	sess, err := s.Get(ctx, refreshToken)
	if err != nil {
		return err
	}
	if sess == nil {
		return nil // Already deleted
	}

	key := fmt.Sprintf("session:%s", refreshToken)
	userKey := fmt.Sprintf("user_sessions:%s", sess.UserID)

	// Delete session
	if err := s.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	// Remove from user's session set
	if err := s.client.SRem(ctx, userKey, refreshToken).Err(); err != nil {
		return fmt.Errorf("failed to remove from user session: %w", err)
	}

	return nil
}

func (s *sessionService) DeleteAllByUserID(ctx context.Context, userID string) error {
	userKey := fmt.Sprintf("user_session:%s", userID)

	// Get all refresh tokens from this user
	refreshTokens, err := s.client.SMembers(ctx, userKey).Result()
	if err != nil {
		return fmt.Errorf("failed to get user sessions: %w", err)
	}

	// Delete each session
	for _, token := range refreshTokens {
		key := fmt.Sprintf("session:%s", token)
		if err := s.client.Del(ctx, key).Err(); err != nil {
			// Log error but continue
			fmt.Printf("failed to delete session %s: %v", token, err)
		}
	}

	// Delete user sessions set
	if err := s.client.Del(ctx, userKey).Err(); err != nil {
		return fmt.Errorf("failed to delete user sessions set: %w", err)
	}

	return nil
}
