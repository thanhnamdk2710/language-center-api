package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func InitLogger(serviceName string, env string) {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	level := zerolog.InfoLevel
	if env == "dev" {
		level = zerolog.DebugLevel
	}

	zerolog.SetGlobalLevel(level)
	Logger = zerolog.New(output).
		With().
		Timestamp().
		Str("service", serviceName).
		Str("env", env).
		Logger()
}

// Shortcut helpers
func Debug(msg string, fields ...interface{}) {
	Logger.Debug().Fields(fields).Msg(msg)
}

func Info(msg string, fields ...interface{}) {
	Logger.Info().Fields(fields).Msg(msg)
}

func Warn(msg string, fields ...interface{}) {
	Logger.Warn().Fields(fields).Msg(msg)
}

func Error(msg string, err error) {
	Logger.Error().Err(err).Msg(msg)
}

func Fatal(msg string, err error) {
	Logger.Fatal().Err(err).Msg(msg)
}
