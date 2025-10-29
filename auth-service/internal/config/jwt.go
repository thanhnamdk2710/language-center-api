package config

import (
	"time"
)

type JWTConfig struct {
	SecretKey            string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

func loadJWTConfig() JWTConfig {
	return JWTConfig{
		SecretKey:            getString("JWT_HOST", "change-me-in-production"),
		AccessTokenDuration:  getDuration("JWT_ACCESS_TOKEN_DURATION", 15*time.Minute),
		RefreshTokenDuration: getDuration("JWT_REFRESH_TOKEN_DURATION", 7*24*time.Hour),
	}
}
