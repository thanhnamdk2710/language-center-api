package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppName  string
	HTTPPort string
	GRPCPort string
	Postgres PostgresConfig
	Redis    RedisConfig
	SMTP     SMTPConfig
	JWT      JWTConfig
}

func Load() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading .env file: %v", err)
	}

	return &Config{
		AppName:  viper.GetString("APP_NAME"),
		HTTPPort: viper.GetString("HTTP_PORT"),
		GRPCPort: viper.GetString("GRPC_PORT"),
		Postgres: loadPostgresConfig(),
		Redis:    loadRedisConfig(),
		SMTP:     loadSMTPConfig(),
		JWT:      loadJWTConfig(),
	}
}

func getString(key, def string) string {
	if val := viper.GetString(key); val != "" {
		return val
	}
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}

func getBool(key string, def bool) bool {
	if val := viper.GetString(key); val != "" {
		parsed, err := strconv.ParseBool(val)
		if err == nil {
			return parsed
		}
	}
	if val := os.Getenv(key); val != "" {
		parsed, err := strconv.ParseBool(val)
		if err == nil {
			return parsed
		}
	}
	return def
}

func getDuration(key string, def time.Duration) time.Duration {
	if val := viper.GetString(key); val != "" {
		parsed, err := time.ParseDuration(val)
		if err == nil {
			return parsed
		}
	}
	if val := os.Getenv(key); val != "" {
		parsed, err := time.ParseDuration(val)
		if err == nil {
			return parsed
		}
	}
	return def
}
