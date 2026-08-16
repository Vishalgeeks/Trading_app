package category

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrCategoryNotFound   = errors.New("category not found")
	ErrDuplicateName      = errors.New("category name already exists")
	ErrDuplicateSlug      = errors.New("category slug already exists")
	ErrInvalidName        = errors.New("name is required")
	ErrInvalidSlug        = errors.New("slug is required")
	ErrCategoryHasDesigns = errors.New("category has active designs and cannot be deactivated")
	ErrDatabaseFailure    = errors.New("database operation failed")
)

var slugRegex = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) CreateCategory(ctx context.Context, cat Category) (Category, error) {
	query := `
		INSERT INTO categories (name, slug, description, is_active)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, slug, description, is_active, created_at, updated_at
	`

	row := r.pool.QueryRow(ctx, query, cat.Name, cat.Slug, cat.Description, cat.IsActive)
	created, err := scanCategory(row)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				if strings.Contains(pgErr.Message, "name") {
					return Category{}, ErrDuplicateName
				}
				if strings.Contains(pgErr.Message, "slug") {
					return Category{}, ErrDuplicateSlug
				}
				return Category{}, ErrDuplicateName
			}
		}
		return Category{}, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return created, nil
}

func (r *Repository) GetCategoryByID(ctx context.Context, id string) (Category, error) {
	query := `
		SELECT id, name, slug, description, is_active, created_at, updated_at
		FROM categories
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)
	cat, err := scanCategory(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Category{}, ErrCategoryNotFound
		}
		return Category{}, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return cat, nil
}

func (r *Repository) GetCategoryBySlug(ctx context.Context, slug string) (Category, error) {
	query := `
		SELECT id, name, slug, description, is_active, created_at, updated_at
		FROM categories
		WHERE slug = $1
	`

	row := r.pool.QueryRow(ctx, query, slug)
	cat, err := scanCategory(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Category{}, ErrCategoryNotFound
		}
		return Category{}, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return cat, nil
}

func (r *Repository) GetCategory(ctx context.Context, id string) (Category, error) {
	return r.GetCategoryByID(ctx, id)
}

func (r *Repository) ListCategories(ctx context.Context, activeOnly bool) ([]Category, error) {
	query := `
		SELECT id, name, slug, description, is_active, created_at, updated_at
		FROM categories
	`

	var args []interface{}
	if activeOnly {
		query += " WHERE is_active = true"
	}
	query += " ORDER BY name ASC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		cat, err := scanCategory(rows)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
		}
		categories = append(categories, cat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return categories, nil
}

func (r *Repository) ListActiveCategories(ctx context.Context) ([]Category, error) {
	return r.ListCategories(ctx, true)
}

func (r *Repository) UpdateCategory(ctx context.Context, id string, req UpdateCategoryRequest) (Category, error) {
	query := `
		UPDATE categories
		SET name = COALESCE($2, name),
		    slug = COALESCE($3, slug),
		    description = COALESCE($4, description),
		    is_active = COALESCE($5, is_active),
		    updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, slug, description, is_active, created_at, updated_at
	`

	row := r.pool.QueryRow(ctx, query, id, req.Name, req.Slug, req.Description, req.IsActive)
	cat, err := scanCategory(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Category{}, ErrCategoryNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				if strings.Contains(pgErr.Message, "name") {
					return Category{}, ErrDuplicateName
				}
				if strings.Contains(pgErr.Message, "slug") {
					return Category{}, ErrDuplicateSlug
				}
			}
		}
		return Category{}, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return cat, nil
}

func (r *Repository) DeactivateCategory(ctx context.Context, id string) error {
	query := `UPDATE categories SET is_active = false, updated_at = NOW() WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	if result.RowsAffected() == 0 {
		return ErrCategoryNotFound
	}

	return nil
}

func (r *Repository) CountActiveDesignsInCategory(ctx context.Context, categoryID string) (int, error) {
	query := `SELECT COUNT(*) FROM designs WHERE category_id = $1 AND is_active = true`

	var count int
	err := r.pool.QueryRow(ctx, query, categoryID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return count, nil
}

func scanCategory(scanner interface {
	Scan(dest ...interface{}) error
}) (Category, error) {
	var cat Category
	err := scanner.Scan(
		&cat.ID,
		&cat.Name,
		&cat.Slug,
		&cat.Description,
		&cat.IsActive,
		&cat.CreatedAt,
		&cat.UpdatedAt,
	)
	if err != nil {
		return Category{}, err
	}
	return cat, nil
}

func GenerateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = regexp.MustCompile(`[^a-z0-9-]+`).ReplaceAllString(slug, "")
	slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}

func ValidateSlug(slug string) bool {
	if slug == "" || len(slug) > 120 {
		return false
	}
	return slugRegex.MatchString(slug)
}
