package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/thanhnamdk2710/auth-service/internal/config"
	"github.com/thanhnamdk2710/auth-service/internal/shared/logger"
)

type RedisCache struct {
	Client *redis.Client
	TTL    time.Duration
}

func NewRedisCache(cfg *config.RedisConfig) (*RedisCache, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect redis: %w", err)
	}

	logger.Info("Connected to Redis successfully")
	return &RedisCache{
		Client: rdb,
		TTL:    time.Duration(cfg.TTL) * time.Second,
	}, nil
}

func (r *RedisCache) Close() error {
	return r.Client.Close()
}
