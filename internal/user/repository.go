package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrDuplicateEmail  = errors.New("email already exists")
	ErrInvalidName     = errors.New("name cannot be empty")
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrInvalidRole     = errors.New("role must be CLIENT or ADMIN")
	ErrDatabaseFailure = errors.New("database operation failed")
	ErrInvalidUUID     = errors.New("invalid uuid")
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) CreateUser(ctx context.Context, user User) (User, error) {
	query := `
		INSERT INTO users (name, email, phone, password_hash, role, avatar_url, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, name, email, phone, password_hash, role, avatar_url, is_active, created_at, updated_at
	`

	row := r.pool.QueryRow(ctx, query,
		user.Name,
		user.Email,
		user.Phone,
		user.PasswordHash,
		user.Role,
		user.AvatarURL,
		user.IsActive,
	)

	var created User
	if err := row.Scan(
		&created.ID,
		&created.Name,
		&created.Email,
		&created.Phone,
		&created.PasswordHash,
		&created.Role,
		&created.AvatarURL,
		&created.IsActive,
		&created.CreatedAt,
		&created.UpdatedAt,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return User{}, ErrDuplicateEmail
			}
		}
		return User{}, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return created, nil
}

func (r *Repository) GetUserByID(ctx context.Context, id string) (User, error) {
	query := `
		SELECT id, name, email, phone, password_hash, role, avatar_url, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)
	user, err := scanUser(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return user, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (User, error) {
	query := `
		SELECT id, name, email, phone, password_hash, role, avatar_url, is_active, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	row := r.pool.QueryRow(ctx, query, email)
	user, err := scanUser(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return user, nil
}

func (r *Repository) UpdateUserProfile(ctx context.Context, id string, req UpdateProfileRequest) (User, error) {
	query := `
		UPDATE users
		SET name = COALESCE($2, name),
		    phone = COALESCE($3, phone),
		    avatar_url = COALESCE($4, avatar_url),
		    updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, email, phone, password_hash, role, avatar_url, is_active, created_at, updated_at
	`

	row := r.pool.QueryRow(ctx, query, id, req.Name, req.Phone, req.AvatarURL)
	user, err := scanUser(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return user, nil
}

func (r *Repository) DeactivateUser(ctx context.Context, id string) error {
	query := `UPDATE users SET is_active = false, updated_at = NOW() WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *Repository) ListUsers(ctx context.Context) ([]User, error) {
	query := `
		SELECT id, name, email, phone, password_hash, role, avatar_url, is_active, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return users, nil
}

func (r *Repository) ListUsersByRole(ctx context.Context, role Role) ([]User, error) {
	query := `
		SELECT id, name, email, phone, password_hash, role, avatar_url, is_active, created_at, updated_at
		FROM users
		WHERE role = $1
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, role)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return users, nil
}

func scanUser(scanner interface {
	Scan(dest ...interface{}) error
}) (User, error) {
	var user User
	err := scanner.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Phone,
		&user.PasswordHash,
		&user.Role,
		&user.AvatarURL,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return User{}, err
	}
	return user, nil
}
