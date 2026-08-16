package design

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
	ErrDesignNotFound   = errors.New("design not found")
	ErrDuplicateSlug    = errors.New("design slug already exists")
	ErrInvalidName      = errors.New("name is required")
	ErrInvalidSlug      = errors.New("slug is required")
	ErrInvalidImageURL  = errors.New("image_url is required")
	ErrInvalidPrice     = errors.New("price must be greater than or equal to 0")
	ErrInvalidDuration  = errors.New("duration_minutes must be greater than 0")
	ErrCategoryNotFound = errors.New("category not found")
	ErrCategoryInactive = errors.New("category is not active")
	ErrDatabaseFailure  = errors.New("database operation failed")
)

var slugRegex = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) CreateDesign(ctx context.Context, des Design) (Design, error) {
	query := `
		INSERT INTO designs (category_id, name, slug, description, image_url, price, duration_minutes, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, category_id, name, slug, description, image_url, price, duration_minutes, is_active, created_at, updated_at
	`

	row := r.pool.QueryRow(ctx, query,
		des.CategoryID,
		des.Name,
		des.Slug,
		des.Description,
		des.ImageURL,
		des.Price,
		des.DurationMinutes,
		des.IsActive,
	)

	created, err := scanDesign(row)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				if strings.Contains(pgErr.Message, "slug") {
					return Design{}, ErrDuplicateSlug
				}
			}
			if pgErr.Code == "23503" {
				return Design{}, ErrCategoryNotFound
			}
		}
		return Design{}, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return created, nil
}

func (r *Repository) GetDesignByID(ctx context.Context, id string) (Design, error) {
	query := `
		SELECT d.id, d.category_id, d.name, d.slug, d.description, d.image_url, d.price, d.duration_minutes, d.is_active, d.created_at, d.updated_at,
		       c.id, c.name, c.slug
		FROM designs d
		JOIN categories c ON d.category_id = c.id
		WHERE d.id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)
	des, err := scanDesignWithCategory(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Design{}, ErrDesignNotFound
		}
		return Design{}, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return des, nil
}

func (r *Repository) GetDesignBySlug(ctx context.Context, slug string) (Design, error) {
	query := `
		SELECT d.id, d.category_id, d.name, d.slug, d.description, d.image_url, d.price, d.duration_minutes, d.is_active, d.created_at, d.updated_at,
		       c.id, c.name, c.slug
		FROM designs d
		JOIN categories c ON d.category_id = c.id
		WHERE d.slug = $1
	`

	row := r.pool.QueryRow(ctx, query, slug)
	des, err := scanDesignWithCategory(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Design{}, ErrDesignNotFound
		}
		return Design{}, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return des, nil
}

func (r *Repository) ListDesigns(ctx context.Context, categoryID *string, searchQuery *string, page, limit int) ([]Design, error) {
	query := `
		SELECT d.id, d.category_id, d.name, d.slug, d.description, d.image_url, d.price, d.duration_minutes, d.is_active, d.created_at, d.updated_at,
		       c.id, c.name, c.slug
		FROM designs d
		JOIN categories c ON d.category_id = c.id
		WHERE d.is_active = true
	`
	args := []interface{}{}
	argCount := 1

	if categoryID != nil && *categoryID != "" {
		query += fmt.Sprintf(" AND d.category_id = $%d", argCount)
		args = append(args, *categoryID)
		argCount++
	}

	if searchQuery != nil && *searchQuery != "" {
		query += fmt.Sprintf(" AND (d.name ILIKE $%d OR d.description ILIKE $%d)", argCount, argCount+1)
		args = append(args, "%"+*searchQuery+"%", "%"+*searchQuery+"%")
		argCount += 2
	}

	query += " ORDER BY d.created_at DESC"

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, limit)
		argCount++
	}

	if page > 1 && limit > 0 {
		offset := (page - 1) * limit
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, offset)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}
	defer rows.Close()

	var designs []Design
	for rows.Next() {
		des, err := scanDesignWithCategory(rows)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
		}
		designs = append(designs, des)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return designs, nil
}

func (r *Repository) ListActiveDesigns(ctx context.Context) ([]Design, error) {
	return r.ListDesigns(ctx, nil, nil, 0, 0)
}

func (r *Repository) ListDesignsByCategory(ctx context.Context, categoryID string) ([]Design, error) {
	return r.ListDesigns(ctx, &categoryID, nil, 0, 0)
}

func (r *Repository) SearchDesigns(ctx context.Context, query string) ([]Design, error) {
	return r.ListDesigns(ctx, nil, &query, 0, 0)
}

func (r *Repository) UpdateDesign(ctx context.Context, id string, req UpdateDesignRequest) (Design, error) {
	query := `
		UPDATE designs
		SET category_id = COALESCE($2, category_id),
		    name = COALESCE($3, name),
		    slug = COALESCE($4, slug),
		    description = COALESCE($5, description),
		    image_url = COALESCE($6, image_url),
		    price = COALESCE($7, price),
		    duration_minutes = COALESCE($8, duration_minutes),
		    is_active = COALESCE($9, is_active),
		    updated_at = NOW()
		WHERE id = $1
		RETURNING id, category_id, name, slug, description, image_url, price, duration_minutes, is_active, created_at, updated_at
	`

	row := r.pool.QueryRow(ctx, query,
		id,
		req.CategoryID,
		req.Name,
		req.Slug,
		req.Description,
		req.ImageURL,
		req.Price,
		req.DurationMinutes,
		req.IsActive,
	)

	des, err := scanDesign(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Design{}, ErrDesignNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				if strings.Contains(pgErr.Message, "slug") {
					return Design{}, ErrDuplicateSlug
				}
			}
			if pgErr.Code == "23503" {
				return Design{}, ErrCategoryNotFound
			}
		}
		return Design{}, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return des, nil
}

func (r *Repository) DeactivateDesign(ctx context.Context, id string) error {
	query := `UPDATE designs SET is_active = false, updated_at = NOW() WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	if result.RowsAffected() == 0 {
		return ErrDesignNotFound
	}

	return nil
}

func (r *Repository) GetCategory(ctx context.Context, id string) (Category, error) {
	query := `SELECT id, name, slug, description, is_active, created_at, updated_at FROM categories WHERE id = $1`

	row := r.pool.QueryRow(ctx, query, id)
	var cat Category
	err := row.Scan(
		&cat.ID,
		&cat.Name,
		&cat.Slug,
		&cat.Description,
		&cat.IsActive,
		&cat.CreatedAt,
		&cat.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Category{}, ErrCategoryNotFound
		}
		return Category{}, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return cat, nil
}

func scanDesign(scanner interface {
	Scan(dest ...interface{}) error
}) (Design, error) {
	var des Design
	err := scanner.Scan(
		&des.ID,
		&des.CategoryID,
		&des.Name,
		&des.Slug,
		&des.Description,
		&des.ImageURL,
		&des.Price,
		&des.DurationMinutes,
		&des.IsActive,
		&des.CreatedAt,
		&des.UpdatedAt,
	)
	if err != nil {
		return Design{}, err
	}
	return des, nil
}

func scanDesignWithCategory(scanner interface {
	Scan(dest ...interface{}) error
}) (Design, error) {
	var des Design
	var cat Category
	err := scanner.Scan(
		&des.ID,
		&des.CategoryID,
		&des.Name,
		&des.Slug,
		&des.Description,
		&des.ImageURL,
		&des.Price,
		&des.DurationMinutes,
		&des.IsActive,
		&des.CreatedAt,
		&des.UpdatedAt,
		&cat.ID,
		&cat.Name,
		&cat.Slug,
	)
	if err != nil {
		return Design{}, err
	}
	des.Category = &cat
	return des, nil
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
	if slug == "" || len(slug) > 180 {
		return false
	}
	return slugRegex.MatchString(slug)
}
