package notification

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotificationNotFound = errors.New("notification not found")
	ErrDatabaseFailure      = errors.New("database operation failed")
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) CreateNotification(ctx context.Context, n Notification) (Notification, error) {
	query := `
		INSERT INTO notifications (user_id, type, title, message, booking_id, is_read)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, type, title, message, booking_id, is_read, created_at, read_at
	`

	row := r.pool.QueryRow(ctx, query,
		n.UserID,
		n.Type,
		n.Title,
		n.Message,
		n.BookingID,
		n.IsRead,
	)

	created, err := scanNotification(row)
	if err != nil {
		return Notification{}, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return created, nil
}

func (r *Repository) GetNotificationByID(ctx context.Context, id string) (Notification, error) {
	query := `
		SELECT id, user_id, type, title, message, booking_id, is_read, created_at, read_at
		FROM notifications
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)
	notification, err := scanNotification(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Notification{}, ErrNotificationNotFound
		}
		return Notification{}, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return notification, nil
}

func (r *Repository) GetNotificationByIDForUser(ctx context.Context, id string, userID string) (Notification, error) {
	query := `
		SELECT id, user_id, type, title, message, booking_id, is_read, created_at, read_at
		FROM notifications
		WHERE id = $1 AND user_id = $2
	`

	row := r.pool.QueryRow(ctx, query, id, userID)
	notification, err := scanNotification(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Notification{}, ErrNotificationNotFound
		}
		return Notification{}, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return notification, nil
}

func (r *Repository) ListNotificationsByUser(ctx context.Context, userID string, limit, offset int) ([]Notification, error) {
	query := `
		SELECT id, user_id, type, title, message, booking_id, is_read, created_at, read_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		notification, err := scanNotification(rows)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
		}
		notifications = append(notifications, notification)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return notifications, nil
}

func (r *Repository) ListUnreadNotificationsByUser(ctx context.Context, userID string, limit, offset int) ([]Notification, error) {
	query := `
		SELECT id, user_id, type, title, message, booking_id, is_read, created_at, read_at
		FROM notifications
		WHERE user_id = $1 AND is_read = FALSE
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		notification, err := scanNotification(rows)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
		}
		notifications = append(notifications, notification)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return notifications, nil
}

func (r *Repository) CountUnreadNotifications(ctx context.Context, userID string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM notifications
		WHERE user_id = $1 AND is_read = FALSE
	`

	var count int
	row := r.pool.QueryRow(ctx, query, userID)
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return count, nil
}

func (r *Repository) MarkNotificationAsRead(ctx context.Context, id string, userID string) (Notification, error) {
	query := `
		UPDATE notifications
		SET is_read = TRUE, read_at = NOW()
		WHERE id = $1 AND user_id = $2
		RETURNING id, user_id, type, title, message, booking_id, is_read, created_at, read_at
	`

	row := r.pool.QueryRow(ctx, query, id, userID)
	notification, err := scanNotification(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Notification{}, ErrNotificationNotFound
		}
		return Notification{}, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return notification, nil
}

func (r *Repository) MarkAllNotificationsAsRead(ctx context.Context, userID string) (int64, error) {
	query := `
		UPDATE notifications
		SET is_read = TRUE, read_at = NOW()
		WHERE user_id = $1 AND is_read = FALSE
	`

	result, err := r.pool.Exec(ctx, query, userID)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return result.RowsAffected(), nil
}

func scanNotification(scanner interface {
	Scan(dest ...interface{}) error
}) (Notification, error) {
	var n Notification
	var bookingID sql.NullString
	var readAt sql.NullTime

	err := scanner.Scan(
		&n.ID,
		&n.UserID,
		&n.Type,
		&n.Title,
		&n.Message,
		&bookingID,
		&n.IsRead,
		&n.CreatedAt,
		&readAt,
	)
	if err != nil {
		return Notification{}, err
	}

	if bookingID.Valid {
		n.BookingID = &bookingID.String
	}
	if readAt.Valid {
		n.ReadAt = &readAt.Time
	}

	return n, nil
}
