package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	AppName  string
	HTTPPort string
	GRPCPort string
	Postgres PostgresConfig
	Redis    RedisConfig
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
