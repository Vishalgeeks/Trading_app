package category

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

var testPool *pgxpool.Pool

func getPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	if testPool == nil {
		ctx := context.Background()
		pool, err := pgxpool.New(ctx, "postgres://postgres:masum7003@localhost:5432/mehndi_booking?sslmode=disable")
		require.NoError(t, err)
		testPool = pool
	}
	return testPool
}

func setupRepository(t *testing.T) *Repository {
	t.Helper()
	return NewRepository(getPool(t))
}

func uniqueName(prefix string) string {
	return prefix + "-" + time.Now().Format("20060102150405")
}

func uniqueSlug(prefix string) string {
	return prefix + "-" + time.Now().Format("20060102150405")
}

func cleanupCategories(t *testing.T) {
	t.Helper()
	ctx := context.Background()
	getPool(t).Exec(ctx, "DELETE FROM categories WHERE slug LIKE $1", "%-test")
}

func TestMain(m *testing.M) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgres://postgres:masum7003@localhost:5432/mehndi_booking?sslmode=disable")
	if err != nil {
		panic(err)
	}
	testPool = pool
	pool.Exec(ctx, "DELETE FROM categories WHERE slug LIKE $1", "%-test")

	code := m.Run()

	pool.Exec(ctx, "DELETE FROM categories WHERE slug LIKE $1", "%-test")
	os.Exit(code)
}

func TestRepository_CreateCategory(t *testing.T) {
	repo := setupRepository(t)
	ctx := context.Background()

	name := uniqueName("bridal")
	cat, err := repo.CreateCategory(ctx, Category{
		Name:     name,
		Slug:     uniqueSlug("bridal"),
		IsActive: true,
	})
	require.NoError(t, err)
	require.Equal(t, name, cat.Name)
	require.True(t, cat.IsActive)
}

func TestRepository_CreateCategory_DuplicateName(t *testing.T) {
	repo := setupRepository(t)
	ctx := context.Background()

	name := uniqueName("dupname")
	_, err := repo.CreateCategory(ctx, Category{
		Name:     name,
		Slug:     uniqueSlug("dupname-1"),
		IsActive: true,
	})
	require.NoError(t, err)

	_, err = repo.CreateCategory(ctx, Category{
		Name:     name,
		Slug:     uniqueSlug("dupname-2"),
		IsActive: true,
	})
	require.ErrorIs(t, err, ErrDuplicateName)
}

func TestRepository_CreateCategory_DuplicateSlug(t *testing.T) {
	repo := setupRepository(t)
	ctx := context.Background()

	slug := uniqueSlug("dupslug")
	_, err := repo.CreateCategory(ctx, Category{
		Name:     uniqueName("dupslug-1"),
		Slug:     slug,
		IsActive: true,
	})
	require.NoError(t, err)

	_, err = repo.CreateCategory(ctx, Category{
		Name:     uniqueName("dupslug-2"),
		Slug:     slug,
		IsActive: true,
	})
	require.ErrorIs(t, err, ErrDuplicateSlug)
}

func TestRepository_GetCategoryByID(t *testing.T) {
	repo := setupRepository(t)
	ctx := context.Background()

	name := uniqueName("getbyid")
	created, err := repo.CreateCategory(ctx, Category{
		Name:     name,
		Slug:     uniqueSlug("getbyid"),
		IsActive: true,
	})
	require.NoError(t, err)

	found, err := repo.GetCategoryByID(ctx, created.ID)
	require.NoError(t, err)
	require.Equal(t, name, found.Name)
}

func TestRepository_GetCategoryByID_NotFound(t *testing.T) {
	repo := setupRepository(t)
	ctx := context.Background()

	_, err := repo.GetCategoryByID(ctx, "00000000-0000-0000-0000-000000000000")
	require.ErrorIs(t, err, ErrCategoryNotFound)
}

func TestRepository_ListActiveCategories(t *testing.T) {
	repo := setupRepository(t)
	ctx := context.Background()

	_, _ = repo.CreateCategory(ctx, Category{Name: uniqueName("active"), Slug: uniqueSlug("active"), IsActive: true})
	_, _ = repo.CreateCategory(ctx, Category{Name: uniqueName("inactive"), Slug: uniqueSlug("inactive"), IsActive: false})

	cats, err := repo.ListActiveCategories(ctx)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(cats), 1)
}

func TestRepository_UpdateCategory(t *testing.T) {
	repo := setupRepository(t)
	ctx := context.Background()

	name := uniqueName("update")
	created, err := repo.CreateCategory(ctx, Category{
		Name:     name,
		Slug:     uniqueSlug("update"),
		IsActive: true,
	})
	require.NoError(t, err)

	desc := "Updated description"
	updated, err := repo.UpdateCategory(ctx, created.ID, UpdateCategoryRequest{
		Description: &desc,
	})
	require.NoError(t, err)
	require.Equal(t, "Updated description", *updated.Description)
}

func TestRepository_DeactivateCategory(t *testing.T) {
	repo := setupRepository(t)
	ctx := context.Background()

	name := uniqueName("deactivate")
	created, err := repo.CreateCategory(ctx, Category{
		Name:     name,
		Slug:     uniqueSlug("deactivate"),
		IsActive: true,
	})
	require.NoError(t, err)

	err = repo.DeactivateCategory(ctx, created.ID)
	require.NoError(t, err)

	found, err := repo.GetCategoryByID(ctx, created.ID)
	require.NoError(t, err)
	require.False(t, found.IsActive)
}
