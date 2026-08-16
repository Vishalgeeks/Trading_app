package database

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(databaseURL string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		slog.Error("Failed to create database pool", "error", err)
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		slog.Error("Failed to ping database", "error", err)
		pool.Close()
		return nil, err
	}

	slog.Info("Database connection established")
	return pool, nil
}
