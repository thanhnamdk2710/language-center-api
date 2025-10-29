package main

import (
	"fmt"
	"sync"

	"github.com/thanhnamdk2710/auth-service/internal/config"
	"github.com/thanhnamdk2710/auth-service/internal/container"
	"github.com/thanhnamdk2710/auth-service/internal/delivery/grpc"
	httpDelivery "github.com/thanhnamdk2710/auth-service/internal/delivery/http"
	"github.com/thanhnamdk2710/auth-service/internal/infra/repository/postgres"
	"github.com/thanhnamdk2710/auth-service/internal/infra/repository/redis"
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
	redis, err := redis.NewRedisClient(&cfg.Redis)
	if err != nil {
		logger.Fatal("Redis connection failed", err)
	}
	defer redis.Close()

	// Initialize container with all dependencies
	container := container.NewContainer(db.Conn, cfg)

	var wg sync.WaitGroup
	wg.Add(2)

	// HTTP server
	go func() {
		defer wg.Done()

		r := httpDelivery.NewRouter(container)

		addr := fmt.Sprintf(":%s", cfg.HTTPPort)
		logger.Info(fmt.Sprintf("HTTP server running on %s", addr))

		if err := r.Run(addr); err != nil {
			logger.Fatal("Failed to start server", err)
		}
	}()

	// gRPC server
	go func() {
		defer wg.Done()
		grpc.StartGRPCServer(cfg)
	}()

	wg.Wait()
}
