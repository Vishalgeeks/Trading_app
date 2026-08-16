package availability

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrAvailabilityNotFound = errors.New("availability not found")
	ErrDatabaseFailure      = errors.New("database operation failed")
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) CreateAvailability(ctx context.Context, av Availability) (Availability, error) {
	query := `
		INSERT INTO admin_availability (day_of_week, start_time, end_time, is_active)
		VALUES ($1, $2, $3, $4)
		RETURNING id, day_of_week, start_time, end_time, is_active, created_at, updated_at
	`

	startTime, _ := time.Parse("15:04", av.StartTime.Format("15:04"))
	endTime, _ := time.Parse("15:04", av.EndTime.Format("15:04"))

	row := r.pool.QueryRow(ctx, query, av.DayOfWeek, startTime, endTime, av.IsActive)
	created, err := scanAvailability(row)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return Availability{}, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
		}
		return Availability{}, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return created, nil
}

func (r *Repository) GetAvailabilityByID(ctx context.Context, id string) (Availability, error) {
	query := `
		SELECT id, day_of_week, start_time, end_time, is_active, created_at, updated_at
		FROM admin_availability
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)
	av, err := scanAvailability(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Availability{}, ErrAvailabilityNotFound
		}
		return Availability{}, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return av, nil
}

func (r *Repository) ListAvailability(ctx context.Context, activeOnly bool) ([]Availability, error) {
	query := `
		SELECT id, day_of_week, start_time, end_time, is_active, created_at, updated_at
		FROM admin_availability
	`

	if activeOnly {
		query += " WHERE is_active = true"
	}
	query += " ORDER BY day_of_week, start_time"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}
	defer rows.Close()

	var availabilities []Availability
	for rows.Next() {
		av, err := scanAvailability(rows)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
		}
		availabilities = append(availabilities, av)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return availabilities, nil
}

func (r *Repository) GetAvailabilityForDay(ctx context.Context, dayOfWeek int) ([]Availability, error) {
	query := `
		SELECT id, day_of_week, start_time, end_time, is_active, created_at, updated_at
		FROM admin_availability
		WHERE day_of_week = $1 AND is_active = true
		ORDER BY start_time
	`

	rows, err := r.pool.Query(ctx, query, dayOfWeek)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}
	defer rows.Close()

	var availabilities []Availability
	for rows.Next() {
		av, err := scanAvailability(rows)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
		}
		availabilities = append(availabilities, av)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return availabilities, nil
}

func (r *Repository) UpdateAvailability(ctx context.Context, id string, req UpdateAvailabilityRequest) (Availability, error) {
	query := `
		UPDATE admin_availability
		SET day_of_week = COALESCE($2, day_of_week),
		    start_time = COALESCE($3, start_time),
		    end_time = COALESCE($4, end_time),
		    is_active = COALESCE($5, is_active),
		    updated_at = NOW()
		WHERE id = $1
		RETURNING id, day_of_week, start_time, end_time, is_active, created_at, updated_at
	`

	row := r.pool.QueryRow(ctx, query, id, req.DayOfWeek, req.StartTime, req.EndTime, req.IsActive)
	av, err := scanAvailability(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Availability{}, ErrAvailabilityNotFound
		}
		return Availability{}, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return av, nil
}

func (r *Repository) DeactivateAvailability(ctx context.Context, id string) error {
	query := `UPDATE admin_availability SET is_active = false, updated_at = NOW() WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	if result.RowsAffected() == 0 {
		return ErrAvailabilityNotFound
	}

	return nil
}

func (r *Repository) DeleteAvailability(ctx context.Context, id string) error {
	query := `DELETE FROM admin_availability WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	if result.RowsAffected() == 0 {
		return ErrAvailabilityNotFound
	}

	return nil
}

func scanAvailability(scanner interface {
	Scan(dest ...interface{}) error
}) (Availability, error) {
	var av Availability
	err := scanner.Scan(
		&av.ID,
		&av.DayOfWeek,
		&av.StartTime,
		&av.EndTime,
		&av.IsActive,
		&av.CreatedAt,
		&av.UpdatedAt,
	)
	if err != nil {
		return Availability{}, err
	}
	return av, nil
}
