package notification

import (
	"context"

	"mehndi-booking-backend/internal/user"
)

type NotificationRepository interface {
	CreateNotification(ctx context.Context, n Notification) (Notification, error)
	GetNotificationByID(ctx context.Context, id string) (Notification, error)
	GetNotificationByIDForUser(ctx context.Context, id string, userID string) (Notification, error)
	ListNotificationsByUser(ctx context.Context, userID string, limit, offset int) ([]Notification, error)
	ListUnreadNotificationsByUser(ctx context.Context, userID string, limit, offset int) ([]Notification, error)
	CountUnreadNotifications(ctx context.Context, userID string) (int, error)
	MarkNotificationAsRead(ctx context.Context, id string, userID string) (Notification, error)
	MarkAllNotificationsAsRead(ctx context.Context, userID string) (int64, error)
}

type UserRepository interface {
	GetUserByID(ctx context.Context, id string) (user.User, error)
	ListUsersByRole(ctx context.Context, role user.Role) ([]user.User, error)
}

type Service struct {
	repo     NotificationRepository
	userRepo UserRepository
}

func NewService(repo NotificationRepository, userRepo UserRepository) *Service {
	return &Service{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *Service) CreateNotification(ctx context.Context, n Notification) (Notification, error) {
	return s.repo.CreateNotification(ctx, n)
}

func (s *Service) GetNotification(ctx context.Context, id string, userID string) (Notification, error) {
	return s.repo.GetNotificationByIDForUser(ctx, id, userID)
}

func (s *Service) ListNotifications(ctx context.Context, userID string, limit, offset int) ([]Notification, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.ListNotificationsByUser(ctx, userID, limit, offset)
}

func (s *Service) ListUnreadNotifications(ctx context.Context, userID string, limit, offset int) ([]Notification, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.ListUnreadNotificationsByUser(ctx, userID, limit, offset)
}

func (s *Service) GetUnreadCount(ctx context.Context, userID string) (int, error) {
	return s.repo.CountUnreadNotifications(ctx, userID)
}

func (s *Service) MarkAsRead(ctx context.Context, id string, userID string) (Notification, error) {
	return s.repo.MarkNotificationAsRead(ctx, id, userID)
}

func (s *Service) MarkAllAsRead(ctx context.Context, userID string) (int64, error) {
	return s.repo.MarkAllNotificationsAsRead(ctx, userID)
}
