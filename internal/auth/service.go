package auth

import (
	"context"
	"fmt"
	"strings"
)

var (
	ErrMissingName      = fmt.Errorf("name is required")
	ErrMissingEmail     = fmt.Errorf("email is required")
	ErrInvalidEmail     = fmt.Errorf("invalid email format")
	ErrMissingPassword  = fmt.Errorf("password is required")
	ErrPasswordTooShort = fmt.Errorf("password must be at least 8 characters")
)

type UserCreator interface {
	CreateUser(ctx context.Context, name, email, phone, passwordHash, role string, avatarURL *string) (UserResult, error)
	GetUserByEmail(ctx context.Context, email string) (UserResult, error)
}

type Service struct {
	userRepo UserCreator
}

func NewService(userRepo UserCreator) *Service {
	return &Service{userRepo: userRepo}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (RegisterResponse, error) {
	if strings.TrimSpace(req.Name) == "" {
		return RegisterResponse{}, ErrMissingName
	}
	if strings.TrimSpace(req.Email) == "" {
		return RegisterResponse{}, ErrMissingEmail
	}
	if !isValidEmail(req.Email) {
		return RegisterResponse{}, ErrInvalidEmail
	}
	if req.Password == "" {
		return RegisterResponse{}, ErrMissingPassword
	}
	if len(req.Password) < 8 {
		return RegisterResponse{}, ErrPasswordTooShort
	}

	normalizedEmail := strings.ToLower(strings.TrimSpace(req.Email))

	_, err := s.userRepo.GetUserByEmail(ctx, normalizedEmail)
	if err == nil {
		return RegisterResponse{}, fmt.Errorf("email already exists")
	}
	if err != fmt.Errorf("user not found") && err.Error() != "user not found" {
		return RegisterResponse{}, err
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return RegisterResponse{}, err
	}

	createdUser, err := s.userRepo.CreateUser(ctx, strings.TrimSpace(req.Name), normalizedEmail, req.Phone, hashedPassword, "CLIENT", nil)
	if err != nil {
		return RegisterResponse{}, err
	}

	return createdUser.ToResponse(), nil
}

func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
