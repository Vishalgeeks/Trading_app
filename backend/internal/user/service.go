package user

import (
	"context"
	"fmt"
	"regexp"
	"unicode/utf8"

	"github.com/google/uuid"
	"mehndi-booking-backend/internal/auth"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type UserRepository interface {
	CreateUser(ctx context.Context, user User) (User, error)
	GetUserByID(ctx context.Context, id string) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	UpdateUserProfile(ctx context.Context, id string, req UpdateProfileRequest) (User, error)
	DeactivateUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context) ([]User, error)
	ListUsersByRole(ctx context.Context, role Role) ([]User, error)
}

type Service struct {
	repo UserRepository
}

func NewService(repo UserRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUser(ctx context.Context, name, email, phone, passwordHash string, role Role, avatarURL *string) (User, error) {
	if err := validateName(name); err != nil {
		return User{}, err
	}
	if err := validateEmail(email); err != nil {
		return User{}, err
	}
	if err := validateRole(role); err != nil {
		return User{}, err
	}

	hashedPassword, err := auth.HashPassword(passwordHash)
	if err != nil {
		return User{}, err
	}

	user := User{
		Name:         name,
		Email:        email,
		Phone:        strPtr(phone),
		PasswordHash: string(hashedPassword),
		Role:         role,
		AvatarURL:    avatarURL,
		IsActive:     true,
	}

	created, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return User{}, err
	}

	return created, nil
}

func (s *Service) GetUserByID(ctx context.Context, id string) (User, error) {
	if id == "" {
		return User{}, ErrInvalidUUID
	}
	if _, err := uuid.Parse(id); err != nil {
		return User{}, ErrInvalidUUID
	}
	return s.repo.GetUserByID(ctx, id)
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (User, error) {
	if err := validateEmail(email); err != nil {
		return User{}, err
	}
	return s.repo.GetUserByEmail(ctx, email)
}

func (s *Service) UpdateUserProfile(ctx context.Context, id string, req UpdateProfileRequest) (User, error) {
	if id == "" {
		return User{}, ErrInvalidUUID
	}
	if _, err := uuid.Parse(id); err != nil {
		return User{}, ErrInvalidUUID
	}

	if req.Name != nil {
		if err := validateName(*req.Name); err != nil {
			return User{}, err
		}
	}

	return s.repo.UpdateUserProfile(ctx, id, req)
}

func (s *Service) DeactivateUser(ctx context.Context, id string) error {
	if id == "" {
		return ErrInvalidUUID
	}
	if _, err := uuid.Parse(id); err != nil {
		return ErrInvalidUUID
	}
	return s.repo.DeactivateUser(ctx, id)
}

func (s *Service) ListUsers(ctx context.Context) ([]User, error) {
	return s.repo.ListUsers(ctx)
}

func (s *Service) ListUsersByRole(ctx context.Context, role Role) ([]User, error) {
	if err := validateRole(role); err != nil {
		return nil, err
	}
	return s.repo.ListUsersByRole(ctx, role)
}

func validateName(name string) error {
	name = trimSpace(name)
	if name == "" {
		return ErrInvalidName
	}
	if utf8.RuneCountInString(name) > 100 {
		return fmt.Errorf("name exceeds maximum length of 100 characters")
	}
	return nil
}

func validateEmail(email string) error {
	email = trimSpace(email)
	if email == "" {
		return ErrInvalidEmail
	}
	if !emailRegex.MatchString(email) {
		return ErrInvalidEmail
	}
	if len(email) > 255 {
		return fmt.Errorf("email exceeds maximum length of 255 characters")
	}
	return nil
}

func validateRole(role Role) error {
	switch role {
	case RoleClient, RoleAdmin:
		return nil
	default:
		return ErrInvalidRole
	}
}

func trimSpace(s string) string {
	start := 0
	for start < len(s) && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	end := len(s)
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
