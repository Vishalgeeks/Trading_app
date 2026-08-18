package favorite

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) CreateFavorite(ctx context.Context, fav Favorite) (Favorite, error) {
	now := time.Now().UTC()
	fav.CreatedAt = now
	fav.UpdatedAt = now

	query := `INSERT INTO favorites (user_id, design_id, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id, user_id, design_id, created_at, updated_at`
	row := r.pool.QueryRow(ctx, query, fav.UserID, fav.DesignID, fav.CreatedAt, fav.UpdatedAt)

	var result Favorite
	err := row.Scan(&result.ID, &result.UserID, &result.DesignID, &result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		return Favorite{}, err
	}

	return result, nil
}

func (r *Repository) GetFavorite(ctx context.Context, userID, designID string) (Favorite, error) {
	query := `SELECT id, user_id, design_id, created_at, updated_at FROM favorites WHERE user_id = $1 AND design_id = $2 LIMIT 1`
	row := r.pool.QueryRow(ctx, query, userID, designID)

	var fav Favorite
	err := row.Scan(&fav.ID, &fav.UserID, &fav.DesignID, &fav.CreatedAt, &fav.UpdatedAt)
	if err != nil {
		return Favorite{}, err
	}

	return fav, nil
}

func (r *Repository) DeleteFavorite(ctx context.Context, userID, designID string) error {
	query := `DELETE FROM favorites WHERE user_id = $1 AND design_id = $2`
	_, err := r.pool.Exec(ctx, query, userID, designID)
	return err
}

func (r *Repository) ListFavoritesByUser(ctx context.Context, userID string) ([]Favorite, error) {
	query := `SELECT id, user_id, design_id, created_at, updated_at FROM favorites WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var favorites []Favorite
	for rows.Next() {
		var fav Favorite
		err := rows.Scan(&fav.ID, &fav.UserID, &fav.DesignID, &fav.CreatedAt, &fav.UpdatedAt)
		if err != nil {
			return nil, err
		}
		favorites = append(favorites, fav)
	}

	return favorites, nil
}

func (r *Repository) CountFavoritesByDesign(ctx context.Context, designID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM favorites WHERE design_id = $1`
	err := r.pool.QueryRow(ctx, query, designID).Scan(&count)
	return count, err
}
