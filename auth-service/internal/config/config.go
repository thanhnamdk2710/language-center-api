package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppName  string
	HTTPPort string
	GRPCPort string
}

func Load() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading .env file: %v", err)
	}

	cfg := &Config{
		AppName:  viper.GetString("APP_NAME"),
		HTTPPort: viper.GetString("HTTP_PORT"),
		GRPCPort: viper.GetString("GRPC_PORT"),
	}
	return cfg
}
