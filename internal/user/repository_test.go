package user

import (
	"context"
	"os"
	"testing"

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

func TestMain(m *testing.M) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgres://postgres:masum7003@localhost:5432/mehndi_booking?sslmode=disable")
	if err != nil {
		panic(err)
	}
	testPool = pool

	code := m.Run()

	testPool.QueryRow(ctx, "DELETE FROM users WHERE email LIKE $1", "%@example.com")
	// Do not close pool to avoid potential deadlock on test exit
	_ = testPool

	os.Exit(code)
}

func setupRepository(t *testing.T) *Repository {
	t.Helper()
	return NewRepository(getPool(t))
}

func TestRepository_CreateUser(t *testing.T) {
	repo := setupRepository(t)
	ctx := context.Background()

	user, err := repo.CreateUser(ctx, User{
		Name:         "Test User",
		Email:        "test@example.com",
		Phone:        strPtr("+919876543210"),
		PasswordHash: "hashed",
		Role:         RoleClient,
		AvatarURL:    nil,
		IsActive:     true,
	})
	require.NoError(t, err)
	require.NotEmpty(t, user.ID)
	require.Equal(t, "Test User", user.Name)
	require.Equal(t, "test@example.com", user.Email)
	require.Equal(t, RoleClient, user.Role)
	require.True(t, user.IsActive)
	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.UpdatedAt)
}

func TestRepository_CreateUser_DuplicateEmail(t *testing.T) {
	repo := setupRepository(t)
	ctx := context.Background()

	_, err := repo.CreateUser(ctx, User{
		Name:         "User One",
		Email:        "dup@example.com",
		PasswordHash: "hashed",
		Role:         RoleClient,
		IsActive:     true,
	})
	require.NoError(t, err)

	_, err = repo.CreateUser(ctx, User{
		Name:         "User Two",
		Email:        "dup@example.com",
		PasswordHash: "hashed",
		Role:         RoleClient,
		IsActive:     true,
	})
	require.ErrorIs(t, err, ErrDuplicateEmail)
}

func TestRepository_GetUserByID(t *testing.T) {
	repo := setupRepository(t)
	ctx := context.Background()

	created, err := repo.CreateUser(ctx, User{
		Name:         "Get User",
		Email:        "getuser@example.com",
		PasswordHash: "hashed",
		Role:         RoleAdmin,
		IsActive:     true,
	})
	require.NoError(t, err)

	found, err := repo.GetUserByID(ctx, created.ID)
	require.NoError(t, err)
	require.Equal(t, created.ID, found.ID)
	require.Equal(t, "Get User", found.Name)
	require.Equal(t, RoleAdmin, found.Role)
}

func TestRepository_GetUserByID_NotFound(t *testing.T) {
	repo := setupRepository(t)
	ctx := context.Background()

	_, err := repo.GetUserByID(ctx, "00000000-0000-0000-0000-000000000000")
	require.ErrorIs(t, err, ErrUserNotFound)
}

func TestRepository_GetUserByEmail(t *testing.T) {
	repo := setupRepository(t)
	ctx := context.Background()

	_, err := repo.CreateUser(ctx, User{
		Name:         "Email User",
		Email:        "emailuser@example.com",
		PasswordHash: "hashed",
		Role:         RoleClient,
		IsActive:     true,
	})
	require.NoError(t, err)

	found, err := repo.GetUserByEmail(ctx, "emailuser@example.com")
	require.NoError(t, err)
	require.Equal(t, "Email User", found.Name)
}

func TestRepository_GetUserByEmail_NotFound(t *testing.T) {
	repo := setupRepository(t)
	ctx := context.Background()

	_, err := repo.GetUserByEmail(ctx, "missing@example.com")
	require.ErrorIs(t, err, ErrUserNotFound)
}

func TestRepository_UpdateUserProfile(t *testing.T) {
	repo := setupRepository(t)
	ctx := context.Background()

	created, err := repo.CreateUser(ctx, User{
		Name:         "Update Me",
		Email:        "updateme@example.com",
		Phone:        strPtr("+911234567890"),
		PasswordHash: "hashed",
		Role:         RoleClient,
		IsActive:     true,
	})
	require.NoError(t, err)

	updated, err := repo.UpdateUserProfile(ctx, created.ID, UpdateProfileRequest{
		Name:  strPtr("Updated Name"),
		Phone: strPtr("+919876543210"),
	})
	require.NoError(t, err)
	require.Equal(t, "Updated Name", updated.Name)
	require.Equal(t, "+919876543210", *updated.Phone)
	require.Equal(t, RoleClient, updated.Role)
	require.True(t, updated.UpdatedAt.After(created.UpdatedAt) || updated.UpdatedAt.Equal(created.UpdatedAt))
}

func TestRepository_UpdateUserProfile_NotFound(t *testing.T) {
	repo := setupRepository(t)
	ctx := context.Background()

	_, err := repo.UpdateUserProfile(ctx, "00000000-0000-0000-0000-000000000000", UpdateProfileRequest{
		Name: strPtr("New Name"),
	})
	require.ErrorIs(t, err, ErrUserNotFound)
}

func TestRepository_DeactivateUser(t *testing.T) {
	repo := setupRepository(t)
	ctx := context.Background()

	created, err := repo.CreateUser(ctx, User{
		Name:         "Deactivate Me",
		Email:        "deact@example.com",
		PasswordHash: "hashed",
		Role:         RoleClient,
		IsActive:     true,
	})
	require.NoError(t, err)

	err = repo.DeactivateUser(ctx, created.ID)
	require.NoError(t, err)

	found, err := repo.GetUserByID(ctx, created.ID)
	require.NoError(t, err)
	require.False(t, found.IsActive)
}

func TestRepository_DeactivateUser_NotFound(t *testing.T) {
	repo := setupRepository(t)
	ctx := context.Background()

	err := repo.DeactivateUser(ctx, "00000000-0000-0000-0000-000000000000")
	require.ErrorIs(t, err, ErrUserNotFound)
}

func TestRepository_ListUsers(t *testing.T) {
	repo := setupRepository(t)
	ctx := context.Background()

	_, err := repo.CreateUser(ctx, User{Name: "List1", Email: "list1@example.com", PasswordHash: "hashed", Role: RoleClient, IsActive: true})
	require.NoError(t, err)
	_, err = repo.CreateUser(ctx, User{Name: "List2", Email: "list2@example.com", PasswordHash: "hashed", Role: RoleAdmin, IsActive: true})
	require.NoError(t, err)

	users, err := repo.ListUsers(ctx)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(users), 2)
}

func TestRepository_ListUsersByRole(t *testing.T) {
	repo := setupRepository(t)
	ctx := context.Background()

	_, err := repo.CreateUser(ctx, User{Name: "Client1", Email: "client1@example.com", PasswordHash: "hashed", Role: RoleClient, IsActive: true})
	require.NoError(t, err)
	_, err = repo.CreateUser(ctx, User{Name: "Admin1", Email: "admin1@example.com", PasswordHash: "hashed", Role: RoleAdmin, IsActive: true})
	require.NoError(t, err)

	clients, err := repo.ListUsersByRole(ctx, RoleClient)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(clients), 1)

	admins, err := repo.ListUsersByRole(ctx, RoleAdmin)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(admins), 1)
}
