package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppName  string
	HTTPPort string
	GRPCPort string
	Postgres PostgresConfig
	Redis    RedisConfig
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	DB       int
	Password string
	TTL      int
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
		Postgres: PostgresConfig{
			Host:     viper.GetString("POSTGRES_HOST"),
			Port:     viper.GetString("POSTGRES_PORT"),
			User:     viper.GetString("POSTGRES_USER"),
			Password: viper.GetString("POSTGRES_PASSWORD"),
			DBName:   viper.GetString("POSTGRES_DB"),
			SSLMode:  viper.GetString("POSTGRES_SSLMODE"),
		},
		Redis: RedisConfig{
			Host:     viper.GetString("REDIS_HOST"),
			Port:     viper.GetString("REDIS_PORT"),
			DB:       viper.GetInt("REDIS_DB"),
			Password: viper.GetString("REDIS_PASSWORD"),
			TTL:      viper.GetInt("REDIS_TTL"),
		},
	}
}
