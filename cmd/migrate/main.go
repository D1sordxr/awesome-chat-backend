package main

import (
	"context"
	"database/sql"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"awesome-chat/internal/infrastructure/config/components/postgres"
	"awesome-chat/internal/infrastructure/config/env"
	"awesome-chat/internal/infrastructure/logger"
	"awesome-chat/migrations"
)

const (
	connectTimeout = 2 * time.Minute
	retryInterval  = 2 * time.Second
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log := logger.NewLogger()

	loader := env.NewLoader()
	cfg := postgres.NewConfig(loader)

	if err := loader.Err(); err != nil {
		log.Error("Failed to load config", "error", err.Error())
		os.Exit(1)
	}

	db, err := sql.Open("pgx", cfg.ConnectionString())
	if err != nil {
		log.Error("Failed to open database", "error", err.Error())
		os.Exit(1)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Error("Failed to close database", "error", closeErr.Error())
		}
	}()

	log.Info("Waiting for database", "host", cfg.Host, "database", cfg.Database, "user", cfg.User)

	if err = waitForDatabase(ctx, db); err != nil {
		log.Error("Database is unreachable", "error", err.Error())
		os.Exit(1)
	}

	log.Info("Applying migrations")

	if err = migrations.Up(db); err != nil {
		log.Error("Failed to apply migrations", "error", err.Error())
		os.Exit(1)
	}

	log.Info("Migrations applied")
}

func waitForDatabase(ctx context.Context, db *sql.DB) error {
	ctx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()

	ticker := time.NewTicker(retryInterval)
	defer ticker.Stop()

	for {
		if err := db.PingContext(ctx); err == nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}
