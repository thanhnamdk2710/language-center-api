package config

import (
	"time"

	"github.com/spf13/viper"
)

type RedisConfig struct {
	Host     string
	Port     string
	DB       int
	Password string
	TTL      time.Duration
}

func loadRedisConfig() RedisConfig {
	return RedisConfig{
		Host:     getString("REDIS_HOST", "localhost"),
		Port:     getString("REDIS_PORT", "6379"),
		Password: getString("REDIS_PASSWORD", ""),
		DB:       viper.GetInt("REDIS_DB"),
		TTL:      time.Duration(viper.GetInt("REDIS_TTL")) * time.Second,
	}
}
