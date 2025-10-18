package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/thanhnamdk2710/auth-service/internal/config"
	"github.com/thanhnamdk2710/auth-service/internal/delivery/grpc"
	"github.com/thanhnamdk2710/auth-service/internal/delivery/http"
	"github.com/thanhnamdk2710/auth-service/internal/repository/postgres"
	"github.com/thanhnamdk2710/auth-service/internal/utils"
)

func main() {
	cfg := config.Load()
	utils.Info(fmt.Sprintf("Starting %s...", cfg.AppName))

	// Connect Database
	db, err := postgres.NewPostgresDB(&cfg.Postgres)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Conn.Close()

	// Run migration
	if err := postgres.RunMigrations(&cfg.Postgres); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	// HTTP server
	go func() {
		defer wg.Done()

		r := http.NewRouter()

		addr := fmt.Sprintf(":%s", cfg.HTTPPort)
		utils.Info(fmt.Sprintf("HTTP server running on %s", addr))

		if err := r.Run(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// gRPC server
	go func() {
		defer wg.Done()
		grpc.StartGRPCServer(cfg)
	}()

	wg.Wait()
}
