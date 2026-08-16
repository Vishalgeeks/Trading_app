package notification

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"mehndi-booking-backend/internal/user"
)

var (
	testPool *pgxpool.Pool
	ctx      = context.Background()
)

func getPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	if testPool == nil {
		pool, err := pgxpool.New(ctx, "postgres://postgres:masum7003@localhost:5432/mehndi_booking?sslmode=disable")
		require.NoError(t, err)
		testPool = pool
	}
	return testPool
}

func setupNotificationRepo(t *testing.T) (*Repository, *user.Repository) {
	t.Helper()
	pool := getPool(t)
	return NewRepository(pool), user.NewRepository(pool)
}

func uniqueEmail(prefix string) string {
	return prefix + "-" + time.Now().Format("20060102150405") + "@example.com"
}

func TestRepository_CreateNotification(t *testing.T) {
	repo, userRepo := setupNotificationRepo(t)

	u, err := userRepo.CreateUser(ctx, user.User{
		Name:         "Test User",
		Email:        uniqueEmail("test-notif-1"),
		Phone:        strPtr("1234567890"),
		PasswordHash: "hash",
		Role:         user.RoleClient,
		AvatarURL:    nil,
		IsActive:     true,
	})
	require.NoError(t, err)

	n := Notification{
		UserID:    u.ID,
		Type:      TypeBookingCreated,
		Title:     "Test",
		Message:   "Test message",
		BookingID: nil,
		IsRead:    false,
	}

	created, err := repo.CreateNotification(ctx, n)
	require.NoError(t, err)
	require.NotEmpty(t, created.ID)
	require.Equal(t, TypeBookingCreated, created.Type)
	require.Equal(t, u.ID, created.UserID)
}

func TestService_ListNotifications(t *testing.T) {
	repo, userRepo := setupNotificationRepo(t)
	svc := NewService(repo, userRepo)

	u, err := userRepo.CreateUser(ctx, user.User{
		Name:         "Test User",
		Email:        uniqueEmail("test-notif-2"),
		Phone:        strPtr("1234567890"),
		PasswordHash: "hash",
		Role:         user.RoleClient,
		AvatarURL:    nil,
		IsActive:     true,
	})
	require.NoError(t, err)

	_, err = repo.CreateNotification(ctx, Notification{
		UserID:    u.ID,
		Type:      TypeBookingCreated,
		Title:     "Test",
		Message:   "Test message",
		IsRead:    false,
		BookingID: nil,
	})
	require.NoError(t, err)

	notifications, err := svc.ListNotifications(ctx, u.ID, 20, 0)
	require.NoError(t, err)
	require.Len(t, notifications, 1)
	require.Equal(t, TypeBookingCreated, notifications[0].Type)
}

func TestService_MarkAsRead(t *testing.T) {
	repo, userRepo := setupNotificationRepo(t)
	svc := NewService(repo, userRepo)

	u, err := userRepo.CreateUser(ctx, user.User{
		Name:         "Test User",
		Email:        uniqueEmail("test-notif-3"),
		Phone:        strPtr("1234567890"),
		PasswordHash: "hash",
		Role:         user.RoleClient,
		AvatarURL:    nil,
		IsActive:     true,
	})
	require.NoError(t, err)

	created, err := repo.CreateNotification(ctx, Notification{
		UserID:    u.ID,
		Type:      TypeBookingCreated,
		Title:     "Test",
		Message:   "Test message",
		IsRead:    false,
		BookingID: nil,
	})
	require.NoError(t, err)
	require.False(t, created.IsRead)

	notification, err := svc.MarkAsRead(ctx, created.ID, u.ID)
	require.NoError(t, err)
	require.True(t, notification.IsRead)
	require.NotNil(t, notification.ReadAt)
}

func TestService_MarkAllAsRead(t *testing.T) {
	repo, userRepo := setupNotificationRepo(t)
	svc := NewService(repo, userRepo)

	u, err := userRepo.CreateUser(ctx, user.User{
		Name:         "Test User",
		Email:        uniqueEmail("test-notif-4"),
		Phone:        strPtr("1234567890"),
		PasswordHash: "hash",
		Role:         user.RoleClient,
		AvatarURL:    nil,
		IsActive:     true,
	})
	require.NoError(t, err)

	_, err = repo.CreateNotification(ctx, Notification{
		UserID:    u.ID,
		Type:      TypeBookingCreated,
		Title:     "Test",
		Message:   "Test message",
		IsRead:    false,
		BookingID: nil,
	})
	require.NoError(t, err)

	count, err := svc.MarkAllAsRead(ctx, u.ID)
	require.NoError(t, err)
	require.Equal(t, int64(1), count)

	unread, err := svc.GetUnreadCount(ctx, u.ID)
	require.NoError(t, err)
	require.Equal(t, 0, unread)
}

func TestService_GetUnreadCount(t *testing.T) {
	repo, userRepo := setupNotificationRepo(t)
	svc := NewService(repo, userRepo)

	u, err := userRepo.CreateUser(ctx, user.User{
		Name:         "Test User",
		Email:        uniqueEmail("test-notif-5"),
		Phone:        strPtr("1234567890"),
		PasswordHash: "hash",
		Role:         user.RoleClient,
		AvatarURL:    nil,
		IsActive:     true,
	})
	require.NoError(t, err)

	_, err = repo.CreateNotification(ctx, Notification{
		UserID:    u.ID,
		Type:      TypeBookingCreated,
		Title:     "Test",
		Message:   "Test message",
		IsRead:    false,
		BookingID: nil,
	})
	require.NoError(t, err)

	count, err := svc.GetUnreadCount(ctx, u.ID)
	require.NoError(t, err)
	require.Equal(t, 1, count)
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
