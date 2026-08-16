package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mehndi-booking-backend/internal/auth"
	"mehndi-booking-backend/internal/availability"
	"mehndi-booking-backend/internal/booking"
	"mehndi-booking-backend/internal/category"
	"mehndi-booking-backend/internal/config"
	"mehndi-booking-backend/internal/database"
	"mehndi-booking-backend/internal/design"
	"mehndi-booking-backend/internal/middleware"
	"mehndi-booking-backend/internal/notification"
	"mehndi-booking-backend/internal/router"
	"mehndi-booking-backend/internal/user"
)

type userRepoAdapter struct {
	repo *user.Repository
}

func (a *userRepoAdapter) CreateUser(ctx context.Context, name, email, phone, passwordHash, role string, avatarURL *string) (auth.UserResult, error) {
	u, err := a.repo.CreateUser(ctx, user.User{
		Name:         name,
		Email:        email,
		Phone:        strPtr(phone),
		PasswordHash: passwordHash,
		Role:         user.Role(role),
		AvatarURL:    avatarURL,
		IsActive:     true,
	})
	if err != nil {
		return auth.UserResult{}, err
	}
	return auth.UserResult{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Phone:     u.Phone,
		Role:      string(u.Role),
		AvatarURL: u.AvatarURL,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}

func (a *userRepoAdapter) GetUserByEmail(ctx context.Context, email string) (auth.UserResult, error) {
	u, err := a.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return auth.UserResult{}, err
	}
	return auth.UserResult{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		Phone:        u.Phone,
		Role:         string(u.Role),
		AvatarURL:    u.AvatarURL,
		IsActive:     u.IsActive,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}, nil
}

func (a *userRepoAdapter) GetUserByID(ctx context.Context, id string) (auth.UserResult, error) {
	u, err := a.repo.GetUserByID(ctx, id)
	if err != nil {
		return auth.UserResult{}, err
	}
	return auth.UserResult{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		Phone:        u.Phone,
		Role:         string(u.Role),
		AvatarURL:    u.AvatarURL,
		IsActive:     u.IsActive,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}, nil
}

func (a *userRepoAdapter) UpdateUserPassword(ctx context.Context, userID, newPasswordHash string) error {
	return a.repo.UpdateUserPassword(ctx, userID, newPasswordHash)
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	logger.Info("Starting server", "port", cfg.AppPort)

	pool, err := database.NewPool(cfg.DatabaseURL)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	userRepo := user.NewRepository(pool)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	authRepo := &userRepoAdapter{repo: userRepo}
	authService := auth.NewService(authRepo, cfg.JWTSecret, cfg.TokenExpiryHours)
	authHandler := auth.NewHandler(authService)

	categoryRepo := category.NewRepository(pool)
	categoryService := category.NewService(categoryRepo)
	categoryHandler := category.NewHandler(categoryService)

	designRepo := design.NewRepository(pool)
	designService := design.NewService(designRepo, categoryRepo)
	designHandler := design.NewHandler(designService)

	availabilityRepo := availability.NewRepository(pool)
	availabilityService := availability.NewService(availabilityRepo)
	availabilityHandler := availability.NewHandler(availabilityService)

	bookingRepo := booking.NewRepository(pool)

	notificationRepo := notification.NewRepository(pool)
	notificationService := notification.NewService(notificationRepo, userRepo)

	bookingService := booking.NewService(bookingRepo, availabilityRepo, notificationRepo, userRepo)
	bookingHandler := booking.NewHandler(bookingService)
	slotHandler := booking.NewSlotHandler(bookingService)

	notificationHandler := notification.NewHandler(notificationService)

	middleware.SetJWTSecret(cfg.JWTSecret)

	r := router.NewRouter(userHandler, authHandler, categoryHandler, designHandler, availabilityHandler, bookingHandler, slotHandler, notificationHandler)

	srv := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("Server listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown error", "error", err)
	}

	logger.Info("Server stopped")
}
