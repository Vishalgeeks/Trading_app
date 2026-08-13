package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type mockUserRepo struct {
	users map[string]UserResult
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]UserResult)}
}

func (m *mockUserRepo) CreateUser(ctx context.Context, name, email, phone, passwordHash, role string, avatarURL *string) (UserResult, error) {
	for _, u := range m.users {
		if u.Email == email {
			return UserResult{}, errors.New("email already exists")
		}
	}
	user := UserResult{
		ID:        "mock-id",
		Name:      name,
		Email:     email,
		Phone:     strPtr(phone),
		Role:      role,
		AvatarURL: avatarURL,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m.users[email] = user
	return user, nil
}

func (m *mockUserRepo) GetUserByEmail(ctx context.Context, email string) (UserResult, error) {
	u, ok := m.users[email]
	if !ok {
		return UserResult{}, errors.New("user not found")
	}
	return u, nil
}

func TestService_Register_Success(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewService(repo)

	resp, err := svc.Register(context.Background(), RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Phone:    "+919876543210",
		Password: "securepass123",
	})
	require.NoError(t, err)
	require.Equal(t, "John Doe", resp.Name)
	require.Equal(t, "john@example.com", resp.Email)
	require.Equal(t, "CLIENT", resp.Role)
	require.NotEmpty(t, resp.ID)
}

func TestService_Register_DuplicateEmail(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewService(repo)

	_, err := svc.Register(context.Background(), RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Phone:    "+919876543210",
		Password: "securepass123",
	})
	require.NoError(t, err)

	_, err = svc.Register(context.Background(), RegisterRequest{
		Name:     "Jane Doe",
		Email:    "john@example.com",
		Phone:    "+919876543210",
		Password: "securepass123",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "email already exists")
}

func TestService_Register_InvalidEmail(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewService(repo)

	_, err := svc.Register(context.Background(), RegisterRequest{
		Name:     "John Doe",
		Email:    "invalid-email",
		Phone:    "+919876543210",
		Password: "securepass123",
	})
	require.ErrorIs(t, err, ErrInvalidEmail)
}

func TestService_Register_MissingName(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewService(repo)

	_, err := svc.Register(context.Background(), RegisterRequest{
		Email:    "john@example.com",
		Password: "securepass123",
	})
	require.ErrorIs(t, err, ErrMissingName)
}

func TestService_Register_MissingPassword(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewService(repo)

	_, err := svc.Register(context.Background(), RegisterRequest{
		Name:  "John Doe",
		Email: "john@example.com",
	})
	require.ErrorIs(t, err, ErrMissingPassword)
}

func TestService_Register_ShortPassword(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewService(repo)

	_, err := svc.Register(context.Background(), RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "short",
	})
	require.ErrorIs(t, err, ErrPasswordTooShort)
}

func TestService_Register_IgnoresRole(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewService(repo)

	resp, err := svc.Register(context.Background(), RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "securepass123",
	})
	require.NoError(t, err)
	require.Equal(t, "CLIENT", resp.Role)
}
