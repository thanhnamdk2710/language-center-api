package redis

import (
	"context"
	"fmt"
	"time"
)

type SessionRepository struct {
	rdb *RedisClient
}

func NewSessionRepository(r *RedisClient) *SessionRepository {
	return &SessionRepository{rdb: r}
}

// SaveRefreshToken save refresh token by userID
func (s *SessionRepository) SaveRefreshToken(ctx context.Context, userID, token string, ttl time.Duration) error {
	key := fmt.Sprintf("session:%s", userID)
	return s.rdb.Client.Set(ctx, key, token, ttl).Err()
}

// GetRefreshToken get refresh token
func (s *SessionRepository) GetRefreshToken(ctx context.Context, userID string) (string, error) {
	key := fmt.Sprintf("session:%s", userID)
	return s.rdb.Client.Get(ctx, key).Result()
}

// DeleteRefreshToken delete session
func (s *SessionRepository) DeleteRefreshToken(ctx context.Context, userID string) error {
	key := fmt.Sprintf("session:%s", userID)
	return s.rdb.Client.Del(ctx, key).Err()
}

// BlacklistToken add access token into blacklist to until expire
func (s *SessionRepository) BlacklistToken(ctx context.Context, token string, expiry time.Duration) error {
	key := fmt.Sprintf("token:blacklist:%s", token)
	return s.rdb.Client.Set(ctx, key, true, expiry).Err()
}

// IsTokenBlacklisted check token exists in blacklist
func (s *SessionRepository) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	key := fmt.Sprintf("token:blacklist:%s", token)
	exists, err := s.rdb.Client.Exists(ctx, key).Result()
	return exists == 1, err
}
