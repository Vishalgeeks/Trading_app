package design

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"mehndi-booking-backend/internal/category"
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

func setupDesignRepository(t *testing.T) (*Repository, *category.Repository, func()) {
	t.Helper()
	pool := getPool(t)
	catRepo := category.NewRepository(pool)
	designRepo := NewRepository(pool)

	cleanup := func() {
		ctx := context.Background()
		pool.Exec(ctx, "DELETE FROM designs WHERE slug LIKE $1", "design-test%")
		pool.Exec(ctx, "DELETE FROM categories WHERE name LIKE $1 OR slug LIKE $2 OR slug LIKE $3", "TestCategory%", "test-cat%", "test-category%")
	}

	return designRepo, catRepo, cleanup
}

func uniqueSlug(prefix string) string {
	return prefix + "-" + fmt.Sprintf("%d", time.Now().UnixNano())
}

func TestMain(m *testing.M) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgres://postgres:masum7003@localhost:5432/mehndi_booking?sslmode=disable")
	if err != nil {
		panic(err)
	}
	testPool = pool
	pool.Exec(ctx, "DELETE FROM designs WHERE slug LIKE $1", "design-test%")
	pool.Exec(ctx, "DELETE FROM categories WHERE name LIKE $1 OR slug LIKE $2 OR slug LIKE $3", "TestCategory%", "test-cat%", "test-category%")

	code := m.Run()

	pool.Exec(ctx, "DELETE FROM designs WHERE slug LIKE $1", "design-test%")
	pool.Exec(ctx, "DELETE FROM categories WHERE name LIKE $1 OR slug LIKE $2 OR slug LIKE $3", "TestCategory%", "test-cat%", "test-category%")
	os.Exit(code)
}

func createTestCategory(t *testing.T, ctx context.Context, catRepo *category.Repository) category.Category {
	t.Helper()
	uid := fmt.Sprintf("%d", time.Now().UnixNano())
	cat, err := catRepo.CreateCategory(ctx, category.Category{
		Name:     "TestCategory-" + t.Name() + "-" + uid,
		Slug:     "test-cat-" + uid,
		IsActive: true,
	})
	require.NoError(t, err)
	return cat
}

func TestRepository_CreateDesign(t *testing.T) {
	repo, catRepo, cleanup := setupDesignRepository(t)
	defer cleanup()
	ctx := context.Background()

	cat := createTestCategory(t, ctx, catRepo)

	des, err := repo.CreateDesign(ctx, Design{
		CategoryID:      cat.ID,
		Name:            "Bridal Full Hand-test",
		Slug:            "design-test-" + uniqueSlug(""),
		Description:     strPtr("Full hand bridal design"),
		ImageURL:        "https://example.com/design.jpg",
		Price:           "2500.00",
		DurationMinutes: 180,
		IsActive:        true,
	})
	require.NoError(t, err)
	require.Equal(t, "Bridal Full Hand-test", des.Name)
	require.Equal(t, "2500.00", des.Price)
	require.Equal(t, 180, des.DurationMinutes)
	require.True(t, des.IsActive)
}

func TestRepository_CreateDesign_DuplicateSlug(t *testing.T) {
	repo, catRepo, cleanup := setupDesignRepository(t)
	defer cleanup()
	ctx := context.Background()

	cat := createTestCategory(t, ctx, catRepo)

	slug := "design-test-" + uniqueSlug("")
	_, err := repo.CreateDesign(ctx, Design{
		CategoryID:      cat.ID,
		Name:            "Design 1-test",
		Slug:            slug,
		ImageURL:        "https://example.com/design1.jpg",
		Price:           "1000.00",
		DurationMinutes: 60,
		IsActive:        true,
	})
	require.NoError(t, err)

	_, err = repo.CreateDesign(ctx, Design{
		CategoryID:      cat.ID,
		Name:            "Design 2-test",
		Slug:            slug,
		ImageURL:        "https://example.com/design2.jpg",
		Price:           "2000.00",
		DurationMinutes: 90,
		IsActive:        true,
	})
	require.ErrorIs(t, err, ErrDuplicateSlug)
}

func TestRepository_GetDesignByID(t *testing.T) {
	repo, catRepo, cleanup := setupDesignRepository(t)
	defer cleanup()
	ctx := context.Background()

	cat := createTestCategory(t, ctx, catRepo)
	created, err := repo.CreateDesign(ctx, Design{
		CategoryID:      cat.ID,
		Name:            "Bridal Full Hand-test",
		Slug:            "design-test-" + uniqueSlug(""),
		ImageURL:        "https://example.com/design.jpg",
		Price:           "2500.00",
		DurationMinutes: 180,
		IsActive:        true,
	})
	require.NoError(t, err)

	found, err := repo.GetDesignByID(ctx, created.ID)
	require.NoError(t, err)
	require.Equal(t, "Bridal Full Hand-test", found.Name)
	require.NotNil(t, found.Category)
	require.Equal(t, cat.ID, found.Category.ID)
}

func TestRepository_GetDesignByID_NotFound(t *testing.T) {
	repo, _, cleanup := setupDesignRepository(t)
	defer cleanup()
	ctx := context.Background()

	_, err := repo.GetDesignByID(ctx, "00000000-0000-0000-0000-000000000000")
	require.ErrorIs(t, err, ErrDesignNotFound)
}

func TestRepository_ListActiveDesigns(t *testing.T) {
	repo, catRepo, cleanup := setupDesignRepository(t)
	defer cleanup()
	ctx := context.Background()

	cat := createTestCategory(t, ctx, catRepo)

	_, _ = repo.CreateDesign(ctx, Design{
		CategoryID:      cat.ID,
		Name:            "Active Design-test",
		Slug:            "design-test-" + uniqueSlug("active"),
		ImageURL:        "https://example.com/active.jpg",
		Price:           "1000.00",
		DurationMinutes: 60,
		IsActive:        true,
	})
	_, _ = repo.CreateDesign(ctx, Design{
		CategoryID:      cat.ID,
		Name:            "Inactive Design-test",
		Slug:            "design-test-" + uniqueSlug("inactive"),
		ImageURL:        "https://example.com/inactive.jpg",
		Price:           "2000.00",
		DurationMinutes: 90,
		IsActive:        false,
	})

	designs, err := repo.ListActiveDesigns(ctx)
	require.NoError(t, err)
	require.Len(t, designs, 1)
}

func TestRepository_ListDesignsByCategory(t *testing.T) {
	repo, catRepo, cleanup := setupDesignRepository(t)
	defer cleanup()
	ctx := context.Background()

	cat1 := createTestCategory(t, ctx, catRepo)
	cat2 := createTestCategory(t, ctx, catRepo)

	_, _ = repo.CreateDesign(ctx, Design{
		CategoryID:      cat1.ID,
		Name:            "Design in Cat1-test",
		Slug:            "design-test-" + uniqueSlug("cat1"),
		ImageURL:        "https://example.com/cat1.jpg",
		Price:           "1000.00",
		DurationMinutes: 60,
		IsActive:        true,
	})
	_, _ = repo.CreateDesign(ctx, Design{
		CategoryID:      cat2.ID,
		Name:            "Design in Cat2-test",
		Slug:            "design-test-" + uniqueSlug("cat2"),
		ImageURL:        "https://example.com/cat2.jpg",
		Price:           "2000.00",
		DurationMinutes: 90,
		IsActive:        true,
	})

	designs, err := repo.ListDesignsByCategory(ctx, cat1.ID)
	require.NoError(t, err)
	require.Len(t, designs, 1)
}

func TestRepository_SearchDesigns(t *testing.T) {
	repo, catRepo, cleanup := setupDesignRepository(t)
	defer cleanup()
	ctx := context.Background()

	cat := createTestCategory(t, ctx, catRepo)

	_, _ = repo.CreateDesign(ctx, Design{
		CategoryID:      cat.ID,
		Name:            "Bridal Full Hand-test",
		Slug:            "design-test-" + uniqueSlug("bridal"),
		ImageURL:        "https://example.com/bridal.jpg",
		Price:           "2500.00",
		DurationMinutes: 180,
		IsActive:        true,
	})
	_, _ = repo.CreateDesign(ctx, Design{
		CategoryID:      cat.ID,
		Name:            "Arabic Light-test",
		Slug:            "design-test-" + uniqueSlug("arabic"),
		ImageURL:        "https://example.com/arabic.jpg",
		Price:           "1500.00",
		DurationMinutes: 120,
		IsActive:        true,
	})

	designs, err := repo.SearchDesigns(ctx, "bridal")
	require.NoError(t, err)
	require.Len(t, designs, 1)
}

func TestRepository_UpdateDesign(t *testing.T) {
	repo, catRepo, cleanup := setupDesignRepository(t)
	defer cleanup()
	ctx := context.Background()

	cat := createTestCategory(t, ctx, catRepo)
	created, err := repo.CreateDesign(ctx, Design{
		CategoryID:      cat.ID,
		Name:            "Original Design-test",
		Slug:            "design-test-" + uniqueSlug("update"),
		ImageURL:        "https://example.com/original.jpg",
		Price:           "1000.00",
		DurationMinutes: 60,
		IsActive:        true,
	})
	require.NoError(t, err)

	newPrice := "1500.00"
	updated, err := repo.UpdateDesign(ctx, created.ID, UpdateDesignRequest{
		Price: &newPrice,
	})
	require.NoError(t, err)
	require.Equal(t, "1500.00", updated.Price)
}

func TestRepository_DeactivateDesign(t *testing.T) {
	repo, catRepo, cleanup := setupDesignRepository(t)
	defer cleanup()
	ctx := context.Background()

	cat := createTestCategory(t, ctx, catRepo)
	created, err := repo.CreateDesign(ctx, Design{
		CategoryID:      cat.ID,
		Name:            "Deactivate Test-test",
		Slug:            "design-test-" + uniqueSlug("deactivate"),
		ImageURL:        "https://example.com/deactivate.jpg",
		Price:           "1000.00",
		DurationMinutes: 60,
		IsActive:        true,
	})
	require.NoError(t, err)

	err = repo.DeactivateDesign(ctx, created.ID)
	require.NoError(t, err)

	found, err := repo.GetDesignByID(ctx, created.ID)
	require.NoError(t, err)
	require.False(t, found.IsActive)
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
