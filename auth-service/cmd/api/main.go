package main

import (
	"fmt"

	"github.com/thanhnamdk2710/auth-service/internal/config"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/cache/redis"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/database/postgres"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/di"
	httpapi "github.com/thanhnamdk2710/auth-service/internal/infrastructure/http"
	"github.com/thanhnamdk2710/auth-service/internal/shared/logger"
)

func main() {
	cfg := config.Load()
	logger.InitLogger(cfg.AppName, "dev")

	// Initialize database connection
	db, err := postgres.NewPostgresDB(&cfg.Postgres)
	if err != nil {
		logger.Fatal("Database connection failed", err)
	}
	defer db.Close()

	// Run database migrations
	if err := postgres.RunMigrations(&cfg.Postgres); err != nil {
		logger.Fatal("Migration failed", err)
	}

	// Initialize cache connection
	redisCache, err := redis.NewClient(&cfg.Redis)
	if err != nil {
		logger.Fatal("Redis connection failed", err)
	}
	defer redisCache.Close()

	// Wire all dependencies (Composition Root)
	c := di.NewContainer(db.Conn, redisCache.Client, cfg)

	// HTTP server
	r := httpapi.NewRouter(c)

	addr := fmt.Sprintf(":%s", cfg.HTTPPort)
	logger.Info(fmt.Sprintf("HTTP server running on %s", addr))

	if err := r.Run(addr); err != nil {
		logger.Fatal("Failed to start server", err)
	}
}
