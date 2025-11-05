package main

import (
	"fmt"

	httpapi "github.com/thanhnamdk2710/auth-service/internal/adapter/http"
	"github.com/thanhnamdk2710/auth-service/internal/config"
	di "github.com/thanhnamdk2710/auth-service/internal/container"
	"github.com/thanhnamdk2710/auth-service/internal/infra/repository/postgres"
	redisrepo "github.com/thanhnamdk2710/auth-service/internal/infra/repository/redis"
	"github.com/thanhnamdk2710/auth-service/internal/shared/logger"
)

func main() {
	cfg := config.Load()
	logger.InitLogger(cfg.AppName, "dev")

	// Connect Database
	db, err := postgres.NewPostgresDB(&cfg.Postgres)
	if err != nil {
		logger.Fatal("Database connection failed", err)
	}
	defer db.Conn.Close()

	// Run migration
	if err := postgres.RunMigrations(&cfg.Postgres); err != nil {
		logger.Fatal("Migration failed", err)
	}

	// Connect Redis
	rc, err := redisrepo.NewRedisClient(&cfg.Redis)
	if err != nil {
		logger.Fatal("Redis connection failed", err)
	}
	defer rc.Close()

	// Initialize container with all dependencies
	c := di.NewContainer(db.Conn, rc.Client, cfg)

	// HTTP server
	r := httpapi.NewRouter(c)

	addr := fmt.Sprintf(":%s", cfg.HTTPPort)
	logger.Info(fmt.Sprintf("HTTP server running on %s", addr))

	if err := r.Run(addr); err != nil {
		logger.Fatal("Failed to start server", err)
	}
}
