package user

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type mockRepo struct {
	users map[string]User
}

func newMockRepo() *mockRepo {
	return &mockRepo{users: make(map[string]User)}
}

func (m *mockRepo) CreateUser(ctx context.Context, user User) (User, error) {
	for _, u := range m.users {
		if u.Email == user.Email {
			return User{}, ErrDuplicateEmail
		}
	}
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	m.users[user.ID] = user
	return user, nil
}

func (m *mockRepo) GetUserByID(ctx context.Context, id string) (User, error) {
	user, ok := m.users[id]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return user, nil
}

func (m *mockRepo) GetUserByEmail(ctx context.Context, email string) (User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return User{}, ErrUserNotFound
}

func (m *mockRepo) UpdateUserProfile(ctx context.Context, id string, req UpdateProfileRequest) (User, error) {
	user, ok := m.users[id]
	if !ok {
		return User{}, ErrUserNotFound
	}
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}
	if req.AvatarURL != nil {
		user.AvatarURL = req.AvatarURL
	}
	user.UpdatedAt = time.Now()
	m.users[id] = user
	return user, nil
}

func (m *mockRepo) DeactivateUser(ctx context.Context, id string) error {
	user, ok := m.users[id]
	if !ok {
		return ErrUserNotFound
	}
	user.IsActive = false
	user.UpdatedAt = time.Now()
	m.users[id] = user
	return nil
}

func (m *mockRepo) ListUsers(ctx context.Context) ([]User, error) {
	var users []User
	for _, u := range m.users {
		users = append(users, u)
	}
	return users, nil
}

func (m *mockRepo) ListUsersByRole(ctx context.Context, role Role) ([]User, error) {
	var users []User
	for _, u := range m.users {
		if u.Role == role {
			users = append(users, u)
		}
	}
	return users, nil
}

func TestService_CreateUser(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	user, err := svc.CreateUser(context.Background(), "John Doe", "john@example.com", "+919876543210", "password123", RoleClient, nil)
	require.NoError(t, err)
	require.NotEmpty(t, user.ID)
	require.Equal(t, "John Doe", user.Name)
	require.Equal(t, "john@example.com", user.Email)
	require.Equal(t, RoleClient, user.Role)
	require.True(t, user.IsActive)

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("password123"))
	require.NoError(t, err)
}

func TestService_CreateUser_InvalidName(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	_, err := svc.CreateUser(context.Background(), "", "john@example.com", "+919876543210", "password123", RoleClient, nil)
	require.ErrorIs(t, err, ErrInvalidName)
}

func TestService_CreateUser_InvalidEmail(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	_, err := svc.CreateUser(context.Background(), "John", "invalid-email", "+919876543210", "password123", RoleClient, nil)
	require.ErrorIs(t, err, ErrInvalidEmail)
}

func TestService_CreateUser_InvalidRole(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	_, err := svc.CreateUser(context.Background(), "John", "john@example.com", "+919876543210", "password123", Role("STAFF"), nil)
	require.ErrorIs(t, err, ErrInvalidRole)
}

func TestService_CreateUser_DuplicateEmail(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	_, err := svc.CreateUser(context.Background(), "John", "john@example.com", "+919876543210", "password123", RoleClient, nil)
	require.NoError(t, err)

	_, err = svc.CreateUser(context.Background(), "Jane", "john@example.com", "+919876543210", "password123", RoleClient, nil)
	require.ErrorIs(t, err, ErrDuplicateEmail)
}

func TestService_GetUserByID(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	user, err := svc.CreateUser(context.Background(), "John", "john@example.com", "+919876543210", "password123", RoleClient, nil)
	require.NoError(t, err)

	found, err := svc.GetUserByID(context.Background(), user.ID)
	require.NoError(t, err)
	require.Equal(t, user.ID, found.ID)
}

func TestService_GetUserByID_InvalidUUID(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	_, err := svc.GetUserByID(context.Background(), "not-a-uuid")
	require.ErrorIs(t, err, ErrInvalidUUID)
}

func TestService_GetUserByID_NotFound(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	_, err := svc.GetUserByID(context.Background(), "00000000-0000-0000-0000-000000000000")
	require.ErrorIs(t, err, ErrUserNotFound)
}

func TestService_UpdateUserProfile(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	user, err := svc.CreateUser(context.Background(), "John", "john@example.com", "+919876543210", "password123", RoleClient, nil)
	require.NoError(t, err)

	updated, err := svc.UpdateUserProfile(context.Background(), user.ID, UpdateProfileRequest{
		Name:  strPtr("Jane Doe"),
		Phone: strPtr("+911234567890"),
	})
	require.NoError(t, err)
	require.Equal(t, "Jane Doe", updated.Name)
	require.Equal(t, "+911234567890", *updated.Phone)
	require.Equal(t, RoleClient, updated.Role)
	require.NotEmpty(t, updated.PasswordHash)
}

func TestService_UpdateUserProfile_InvalidUUID(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	_, err := svc.UpdateUserProfile(context.Background(), "not-a-uuid", UpdateProfileRequest{Name: strPtr("Jane")})
	require.ErrorIs(t, err, ErrInvalidUUID)
}

func TestService_UpdateUserProfile_InvalidName(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	user, err := svc.CreateUser(context.Background(), "John", "john@example.com", "+919876543210", "password123", RoleClient, nil)
	require.NoError(t, err)

	emptyName := ""
	_, err = svc.UpdateUserProfile(context.Background(), user.ID, UpdateProfileRequest{Name: &emptyName})
	require.ErrorIs(t, err, ErrInvalidName)
}

func TestService_DeactivateUser(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	user, err := svc.CreateUser(context.Background(), "John", "john@example.com", "+919876543210", "password123", RoleClient, nil)
	require.NoError(t, err)

	err = svc.DeactivateUser(context.Background(), user.ID)
	require.NoError(t, err)

	found, err := svc.GetUserByID(context.Background(), user.ID)
	require.NoError(t, err)
	require.False(t, found.IsActive)
}

func TestService_ListUsers(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	_, err := svc.CreateUser(context.Background(), "User1", "user1@example.com", "+919876543210", "password123", RoleClient, nil)
	require.NoError(t, err)
	_, err = svc.CreateUser(context.Background(), "User2", "user2@example.com", "+919876543210", "password123", RoleAdmin, nil)
	require.NoError(t, err)

	users, err := svc.ListUsers(context.Background())
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(users), 2)
}

func TestService_ListUsersByRole(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	_, err := svc.CreateUser(context.Background(), "Client1", "client1@example.com", "+919876543210", "password123", RoleClient, nil)
	require.NoError(t, err)
	_, err = svc.CreateUser(context.Background(), "Admin1", "admin1@example.com", "+919876543210", "password123", RoleAdmin, nil)
	require.NoError(t, err)

	clients, err := svc.ListUsersByRole(context.Background(), RoleClient)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(clients), 1)

	admins, err := svc.ListUsersByRole(context.Background(), RoleAdmin)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(admins), 1)
}
