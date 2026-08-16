package main

import (
	"context"
	"log/slog"
	"os"

	"mehndi-booking-backend/internal/config"
	"mehndi-booking-backend/internal/database"
	"mehndi-booking-backend/internal/user"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	pool, err := database.NewPool(cfg.DatabaseURL)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	repo := user.NewRepository(pool)
	svc := user.NewService(repo)

	created, err := svc.CreateUser(context.Background(), "Admin User", "admin@example.com", "+919876543210", "admin123", user.RoleAdmin, nil)
	if err != nil {
		slog.Error("Failed to create seed user", "error", err)
		os.Exit(1)
	}

	slog.Info("Seed user created", "id", created.ID, "email", created.Email, "role", created.Role)
}
