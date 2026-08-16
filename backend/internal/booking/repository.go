package booking

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrBookingNotFound = errors.New("booking not found")
	ErrBookingConflict = errors.New("booking conflicts with existing booking")
	ErrInvalidStatus   = errors.New("invalid booking status")
	ErrDatabaseFailure = errors.New("database operation failed")
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) CreateBooking(ctx context.Context, booking Booking) (Booking, error) {
	query := `
		INSERT INTO bookings (user_id, design_id, booking_date, start_time, end_time, status, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, design_id, booking_date, start_time, end_time, status, notes, created_at, updated_at
	`

	row := r.pool.QueryRow(ctx, query,
		booking.UserID,
		booking.DesignID,
		booking.BookingDate,
		booking.StartTime,
		booking.EndTime,
		booking.Status,
		booking.Notes,
	)

	created, err := scanBooking(row)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return Booking{}, ErrBookingConflict
			}
		}
		return Booking{}, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return created, nil
}

func (r *Repository) GetBookingByID(ctx context.Context, id string) (Booking, error) {
	query := `
		SELECT b.id, b.user_id, b.design_id, b.booking_date, b.start_time, b.end_time, b.status, b.notes, b.created_at, b.updated_at,
		       d.name, u.name, u.email, u.phone
		FROM bookings b
		JOIN designs d ON b.design_id = d.id
		JOIN users u ON b.user_id = u.id
		WHERE b.id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)
	booking, err := scanBookingWithDetails(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Booking{}, ErrBookingNotFound
		}
		return Booking{}, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return booking, nil
}

func (r *Repository) GetBookingByIDForUser(ctx context.Context, id string, userID string) (Booking, error) {
	query := `
		SELECT b.id, b.user_id, b.design_id, b.booking_date, b.start_time, b.end_time, b.status, b.notes, b.created_at, b.updated_at,
		       d.name, u.name, u.email, u.phone
		FROM bookings b
		JOIN designs d ON b.design_id = d.id
		JOIN users u ON b.user_id = u.id
		WHERE b.id = $1 AND b.user_id = $2
	`

	row := r.pool.QueryRow(ctx, query, id, userID)
	booking, err := scanBookingWithDetails(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Booking{}, ErrBookingNotFound
		}
		return Booking{}, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return booking, nil
}

func (r *Repository) ListBookingsByUser(ctx context.Context, userID string, status *string, date *string) ([]Booking, error) {
	query := `
		SELECT b.id, b.user_id, b.design_id, b.booking_date, b.start_time, b.end_time, b.status, b.notes, b.created_at, b.updated_at,
		       d.name, u.name, u.email, u.phone
		FROM bookings b
		JOIN designs d ON b.design_id = d.id
		JOIN users u ON b.user_id = u.id
		WHERE b.user_id = $1
	`
	args := []interface{}{userID}
	argCount := 1

	if status != nil && *status != "" {
		query += fmt.Sprintf(" AND b.status = $%d", argCount+1)
		args = append(args, *status)
		argCount++
	}

	if date != nil && *date != "" {
		query += fmt.Sprintf(" AND b.booking_date = $%d", argCount+1)
		args = append(args, *date)
		argCount++
	}

	query += " ORDER BY b.booking_date DESC, b.start_time DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}
	defer rows.Close()

	var bookings []Booking
	for rows.Next() {
		booking, err := scanBookingWithDetails(rows)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
		}
		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return bookings, nil
}

func (r *Repository) ListAdminBookings(ctx context.Context, status *string, date *string, fromDate *string, toDate *string) ([]Booking, error) {
	query := `
		SELECT b.id, b.user_id, b.design_id, b.booking_date, b.start_time, b.end_time, b.status, b.notes, b.created_at, b.updated_at,
		       d.name, u.name, u.email, u.phone
		FROM bookings b
		JOIN designs d ON b.design_id = d.id
		JOIN users u ON b.user_id = u.id
		WHERE 1=1
	`
	args := []interface{}{}
	argCount := 0

	if status != nil && *status != "" {
		argCount++
		query += fmt.Sprintf(" AND b.status = $%d", argCount)
		args = append(args, *status)
	}

	if date != nil && *date != "" {
		argCount++
		query += fmt.Sprintf(" AND b.booking_date = $%d", argCount)
		args = append(args, *date)
	}

	if fromDate != nil && *fromDate != "" {
		argCount++
		query += fmt.Sprintf(" AND b.booking_date >= $%d", argCount)
		args = append(args, *fromDate)
	}

	if toDate != nil && *toDate != "" {
		argCount++
		query += fmt.Sprintf(" AND b.booking_date <= $%d", argCount)
		args = append(args, *toDate)
	}

	query += " ORDER BY b.booking_date DESC, b.start_time DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}
	defer rows.Close()

	var bookings []Booking
	for rows.Next() {
		booking, err := scanBookingWithDetails(rows)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
		}
		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return bookings, nil
}

func (r *Repository) CheckBookingOverlap(ctx context.Context, designID string, bookingDate time.Time, startTime, endTime time.Time, excludeBookingID *string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM bookings
			WHERE design_id = $1
			  AND booking_date = $2::date
			  AND status IN ('PENDING', 'CONFIRMED')
			  AND start_time < $4
			  AND end_time > $3
		`
	args := []interface{}{designID, bookingDate.Format("2006-01-02"), startTime, endTime}
	argCount := 4

	if excludeBookingID != nil && *excludeBookingID != "" {
		query += fmt.Sprintf(" AND id != $%d", argCount+1)
		args = append(args, *excludeBookingID)
	}

	query += ")"

	var exists bool
	row := r.pool.QueryRow(ctx, query, args...)
	if err := row.Scan(&exists); err != nil {
		return false, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return exists, nil
}

func (r *Repository) UpdateBookingStatus(ctx context.Context, id string, status string) (Booking, error) {
	query := `
		UPDATE bookings
		SET status = $2, updated_at = NOW()
		WHERE id = $1
		RETURNING id, user_id, design_id, booking_date, start_time, end_time, status, notes, created_at, updated_at
	`

	row := r.pool.QueryRow(ctx, query, id, status)
	booking, err := scanBooking(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Booking{}, ErrBookingNotFound
		}
		return Booking{}, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return booking, nil
}

func (r *Repository) CancelBooking(ctx context.Context, id string) (Booking, error) {
	query := `
		UPDATE bookings
		SET status = 'CANCELLED', updated_at = NOW()
		WHERE id = $1 AND status NOT IN ('CANCELLED', 'COMPLETED')
		RETURNING id, user_id, design_id, booking_date, start_time, end_time, status, notes, created_at, updated_at
	`

	row := r.pool.QueryRow(ctx, query, id)
	booking, err := scanBooking(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Booking{}, ErrBookingNotFound
		}
		return Booking{}, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return booking, nil
}

func (r *Repository) ListClientBookingsPaginated(ctx context.Context, userID string, status *string, date *string, fromDate *string, toDate *string, limit int, offset int) ([]Booking, error) {
	query := `
		SELECT b.id, b.user_id, b.design_id, b.booking_date, b.start_time, b.end_time, b.status, b.notes, b.created_at, b.updated_at,
		       d.name, u.name, u.email, u.phone
		FROM bookings b
		JOIN designs d ON b.design_id = d.id
		JOIN users u ON b.user_id = u.id
		WHERE b.user_id = $1
	`
	args := []interface{}{userID}
	argCount := 1

	if status != nil && *status != "" {
		query += fmt.Sprintf(" AND b.status = $%d", argCount+1)
		args = append(args, *status)
		argCount++
	}

	if date != nil && *date != "" {
		query += fmt.Sprintf(" AND b.booking_date = $%d", argCount+1)
		args = append(args, *date)
		argCount++
	}

	if fromDate != nil && *fromDate != "" {
		query += fmt.Sprintf(" AND b.booking_date >= $%d", argCount+1)
		args = append(args, *fromDate)
		argCount++
	}

	if toDate != nil && *toDate != "" {
		query += fmt.Sprintf(" AND b.booking_date <= $%d", argCount+1)
		args = append(args, *toDate)
		argCount++
	}

	query += fmt.Sprintf(" ORDER BY b.booking_date DESC, b.start_time DESC LIMIT $%d OFFSET $%d", argCount+1, argCount+2)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}
	defer rows.Close()

	var bookings []Booking
	for rows.Next() {
		booking, err := scanBookingWithDetails(rows)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
		}
		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return bookings, nil
}

func (r *Repository) ListAdminBookingsPaginated(ctx context.Context, status *string, date *string, fromDate *string, toDate *string, search *string, limit int, offset int) ([]Booking, error) {
	query := `
		SELECT b.id, b.user_id, b.design_id, b.booking_date, b.start_time, b.end_time, b.status, b.notes, b.created_at, b.updated_at,
		       d.name, u.name, u.email, u.phone
		FROM bookings b
		JOIN designs d ON b.design_id = d.id
		JOIN users u ON b.user_id = u.id
		WHERE 1=1
	`
	args := []interface{}{}
	argCount := 0

	if status != nil && *status != "" {
		argCount++
		query += fmt.Sprintf(" AND b.status = $%d", argCount)
		args = append(args, *status)
	}

	if date != nil && *date != "" {
		argCount++
		query += fmt.Sprintf(" AND b.booking_date = $%d", argCount)
		args = append(args, *date)
	}

	if fromDate != nil && *fromDate != "" {
		argCount++
		query += fmt.Sprintf(" AND b.booking_date >= $%d", argCount)
		args = append(args, *fromDate)
	}

	if toDate != nil && *toDate != "" {
		argCount++
		query += fmt.Sprintf(" AND b.booking_date <= $%d", argCount)
		args = append(args, *toDate)
	}

	if search != nil && *search != "" {
		argCount++
		query += fmt.Sprintf(" AND (u.name ILIKE $%d OR u.email ILIKE $%d OR d.name ILIKE $%d)", argCount, argCount, argCount)
		args = append(args, "%"+*search+"%")
	}

	query += fmt.Sprintf(" ORDER BY b.booking_date DESC, b.start_time DESC LIMIT $%d OFFSET $%d", argCount+1, argCount+2)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}
	defer rows.Close()

	var bookings []Booking
	for rows.Next() {
		booking, err := scanBookingWithDetails(rows)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
		}
		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return bookings, nil
}

func (r *Repository) ListClientUpcomingBookings(ctx context.Context, userID string, limit int) ([]Booking, error) {
	query := `
		SELECT b.id, b.user_id, b.design_id, b.booking_date, b.start_time, b.end_time, b.status, b.notes, b.created_at, b.updated_at,
		       d.name, u.name, u.email, u.phone
		FROM bookings b
		JOIN designs d ON b.design_id = d.id
		JOIN users u ON b.user_id = u.id
		WHERE b.user_id = $1
		  AND b.status IN ('PENDING', 'CONFIRMED')
		  AND (b.booking_date > CURRENT_DATE OR (b.booking_date = CURRENT_DATE AND b.start_time > CURRENT_TIME))
		ORDER BY b.booking_date ASC, b.start_time ASC
		LIMIT $2
	`

	rows, err := r.pool.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}
	defer rows.Close()

	var bookings []Booking
	for rows.Next() {
		booking, err := scanBookingWithDetails(rows)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
		}
		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return bookings, nil
}

func (r *Repository) ListClientBookingHistory(ctx context.Context, userID string, limit int, offset int) ([]Booking, error) {
	query := `
		SELECT b.id, b.user_id, b.design_id, b.booking_date, b.start_time, b.end_time, b.status, b.notes, b.created_at, b.updated_at,
		       d.name, u.name, u.email, u.phone
		FROM bookings b
		JOIN designs d ON b.design_id = d.id
		JOIN users u ON b.user_id = u.id
		WHERE b.user_id = $1
		  AND b.status IN ('COMPLETED', 'CANCELLED')
		ORDER BY b.booking_date DESC, b.start_time DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}
	defer rows.Close()

	var bookings []Booking
	for rows.Next() {
		booking, err := scanBookingWithDetails(rows)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
		}
		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return bookings, nil
}

func (r *Repository) ListAdminUpcomingBookings(ctx context.Context, limit int) ([]Booking, error) {
	query := `
		SELECT b.id, b.user_id, b.design_id, b.booking_date, b.start_time, b.end_time, b.status, b.notes, b.created_at, b.updated_at,
		       d.name, u.name, u.email, u.phone
		FROM bookings b
		JOIN designs d ON b.design_id = d.id
		JOIN users u ON b.user_id = u.id
		WHERE b.status IN ('PENDING', 'CONFIRMED')
		  AND (b.booking_date > CURRENT_DATE OR (b.booking_date = CURRENT_DATE AND b.start_time > CURRENT_TIME))
		ORDER BY b.booking_date ASC, b.start_time ASC
		LIMIT $1
	`

	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}
	defer rows.Close()

	var bookings []Booking
	for rows.Next() {
		booking, err := scanBookingWithDetails(rows)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
		}
		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return bookings, nil
}

func (r *Repository) ListAdminBookingHistory(ctx context.Context, limit int, offset int) ([]Booking, error) {
	query := `
		SELECT b.id, b.user_id, b.design_id, b.booking_date, b.start_time, b.end_time, b.status, b.notes, b.created_at, b.updated_at,
		       d.name, u.name, u.email, u.phone
		FROM bookings b
		JOIN designs d ON b.design_id = d.id
		JOIN users u ON b.user_id = u.id
		WHERE b.status IN ('COMPLETED', 'CANCELLED')
		ORDER BY b.booking_date DESC, b.start_time DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}
	defer rows.Close()

	var bookings []Booking
	for rows.Next() {
		booking, err := scanBookingWithDetails(rows)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
		}
		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return bookings, nil
}

func (r *Repository) SearchAdminBookings(ctx context.Context, search string, limit int, offset int) ([]Booking, error) {
	query := `
		SELECT b.id, b.user_id, b.design_id, b.booking_date, b.start_time, b.end_time, b.status, b.notes, b.created_at, b.updated_at,
		       d.name, u.name, u.email, u.phone
		FROM bookings b
		JOIN designs d ON b.design_id = d.id
		JOIN users u ON b.user_id = u.id
		WHERE u.name ILIKE $1 OR u.email ILIKE $1 OR d.name ILIKE $1
		ORDER BY b.booking_date DESC, b.start_time DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, "%"+search+"%", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}
	defer rows.Close()

	var bookings []Booking
	for rows.Next() {
		booking, err := scanBookingWithDetails(rows)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
		}
		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return bookings, nil
}

func (r *Repository) GetAdminBookingStats(ctx context.Context, fromDate *string, toDate *string) (BookingStats, error) {
	query := `
		SELECT
			COUNT(*) AS total,
			COUNT(*) FILTER (WHERE status = 'PENDING') AS pending,
			COUNT(*) FILTER (WHERE status = 'CONFIRMED') AS confirmed,
			COUNT(*) FILTER (WHERE status = 'COMPLETED') AS completed,
			COUNT(*) FILTER (WHERE status = 'CANCELLED') AS cancelled,
			COUNT(*) FILTER (WHERE status IN ('PENDING', 'CONFIRMED') AND (booking_date > CURRENT_DATE OR (booking_date = CURRENT_DATE AND start_time > CURRENT_TIME))) AS upcoming
		FROM bookings
		WHERE 1=1
	`
	args := []interface{}{}
	argCount := 0

	if fromDate != nil && *fromDate != "" {
		argCount++
		query += fmt.Sprintf(" AND booking_date >= $%d", argCount)
		args = append(args, *fromDate)
	}

	if toDate != nil && *toDate != "" {
		argCount++
		query += fmt.Sprintf(" AND booking_date <= $%d", argCount)
		args = append(args, *toDate)
	}

	var stats BookingStats
	row := r.pool.QueryRow(ctx, query, args...)
	if err := row.Scan(&stats.Total, &stats.Pending, &stats.Confirmed, &stats.Completed, &stats.Cancelled, &stats.Upcoming); err != nil {
		return BookingStats{}, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return stats, nil
}

func (r *Repository) CountClientBookings(ctx context.Context, userID string, status *string, date *string, fromDate *string, toDate *string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM bookings
		WHERE user_id = $1
	`
	args := []interface{}{userID}
	argCount := 1

	if status != nil && *status != "" {
		query += fmt.Sprintf(" AND status = $%d", argCount+1)
		args = append(args, *status)
		argCount++
	}

	if date != nil && *date != "" {
		query += fmt.Sprintf(" AND booking_date = $%d", argCount+1)
		args = append(args, *date)
		argCount++
	}

	if fromDate != nil && *fromDate != "" {
		query += fmt.Sprintf(" AND booking_date >= $%d", argCount+1)
		args = append(args, *fromDate)
		argCount++
	}

	if toDate != nil && *toDate != "" {
		query += fmt.Sprintf(" AND booking_date <= $%d", argCount+1)
		args = append(args, *toDate)
		argCount++
	}

	var count int
	row := r.pool.QueryRow(ctx, query, args...)
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return count, nil
}

func (r *Repository) CountAdminBookings(ctx context.Context, status *string, date *string, fromDate *string, toDate *string, search *string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM bookings b
		JOIN designs d ON b.design_id = d.id
		JOIN users u ON b.user_id = u.id
		WHERE 1=1
	`
	args := []interface{}{}
	argCount := 0

	if status != nil && *status != "" {
		argCount++
		query += fmt.Sprintf(" AND b.status = $%d", argCount)
		args = append(args, *status)
	}

	if date != nil && *date != "" {
		argCount++
		query += fmt.Sprintf(" AND b.booking_date = $%d", argCount)
		args = append(args, *date)
	}

	if fromDate != nil && *fromDate != "" {
		argCount++
		query += fmt.Sprintf(" AND b.booking_date >= $%d", argCount)
		args = append(args, *fromDate)
	}

	if toDate != nil && *toDate != "" {
		argCount++
		query += fmt.Sprintf(" AND b.booking_date <= $%d", argCount)
		args = append(args, *toDate)
	}

	if search != nil && *search != "" {
		argCount++
		query += fmt.Sprintf(" AND (u.name ILIKE $%d OR u.email ILIKE $%d OR d.name ILIKE $%d)", argCount, argCount, argCount)
		args = append(args, "%"+*search+"%")
	}

	var count int
	row := r.pool.QueryRow(ctx, query, args...)
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return count, nil
}

func scanBooking(scanner interface {
	Scan(dest ...interface{}) error
}) (Booking, error) {
	var b Booking
	err := scanner.Scan(
		&b.ID,
		&b.UserID,
		&b.DesignID,
		&b.BookingDate,
		&b.StartTime,
		&b.EndTime,
		&b.Status,
		&b.Notes,
		&b.CreatedAt,
		&b.UpdatedAt,
	)
	if err != nil {
		return Booking{}, err
	}
	return b, nil
}

func scanBookingWithDetails(scanner interface {
	Scan(dest ...interface{}) error
}) (Booking, error) {
	var b Booking
	err := scanner.Scan(
		&b.ID,
		&b.UserID,
		&b.DesignID,
		&b.BookingDate,
		&b.StartTime,
		&b.EndTime,
		&b.Status,
		&b.Notes,
		&b.CreatedAt,
		&b.UpdatedAt,
		&b.DesignName,
		&b.UserName,
		&b.UserEmail,
		&b.UserPhone,
	)
	if err != nil {
		return Booking{}, err
	}
	return b, nil
}

func (r *Repository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.pool.BeginTx(ctx, pgx.TxOptions{})
}
