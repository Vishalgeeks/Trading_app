package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrMissingName        = fmt.Errorf("name is required")
	ErrMissingEmail       = fmt.Errorf("email is required")
	ErrInvalidEmail       = fmt.Errorf("invalid email format")
	ErrMissingPassword    = fmt.Errorf("password is required")
	ErrPasswordTooShort   = fmt.Errorf("password must be at least 8 characters")
	ErrInvalidCredentials = fmt.Errorf("invalid email or password")
	ErrInactiveUser       = fmt.Errorf("user account is inactive")
	ErrUserNotFound       = fmt.Errorf("user not found")
)

type UserCreator interface {
	CreateUser(ctx context.Context, name, email, phone, passwordHash, role string, avatarURL *string) (UserResult, error)
	GetUserByEmail(ctx context.Context, email string) (UserResult, error)
	GetUserByID(ctx context.Context, id string) (UserResult, error)
	UpdateUserPassword(ctx context.Context, userID, newPasswordHash string) error
}

type Service struct {
	userRepo         UserCreator
	jwtSecret        string
	tokenExpiryHours int
}

func NewService(userRepo UserCreator, jwtSecret string, tokenExpiryHours int) *Service {
	return &Service{
		userRepo:         userRepo,
		jwtSecret:        jwtSecret,
		tokenExpiryHours: tokenExpiryHours,
	}
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

func (s *Service) Login(ctx context.Context, req LoginRequest) (LoginResponse, error) {
	if strings.TrimSpace(req.Email) == "" {
		return LoginResponse{}, ErrMissingEmail
	}
	if req.Password == "" {
		return LoginResponse{}, ErrMissingPassword
	}

	normalizedEmail := strings.ToLower(strings.TrimSpace(req.Email))

	user, err := s.userRepo.GetUserByEmail(ctx, normalizedEmail)
	if err != nil {
		return LoginResponse{}, ErrInvalidCredentials
	}

	if !user.IsActive {
		return LoginResponse{}, ErrInactiveUser
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return LoginResponse{}, ErrInvalidCredentials
	}

	token, err := s.generateJWT(user)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("failed to generate token: %w", err)
	}

	return LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   s.tokenExpiryHours * 3600,
		User:        user.ToResponse(),
	}, nil
}

func (s *Service) GetCurrentUser(ctx context.Context, userID string) (RegisterResponse, error) {
	if strings.TrimSpace(userID) == "" {
		return RegisterResponse{}, fmt.Errorf("invalid user id")
	}

	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return RegisterResponse{}, ErrUserNotFound
	}

	if !user.IsActive {
		return RegisterResponse{}, ErrInactiveUser
	}

	return user.ToResponse(), nil
}

func (s *Service) ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error {
	if strings.TrimSpace(userID) == "" {
		return fmt.Errorf("invalid user id")
	}
	if currentPassword == "" {
		return ErrMissingPassword
	}
	if newPassword == "" {
		return ErrMissingPassword
	}
	if len(newPassword) < 8 {
		return ErrPasswordTooShort
	}

	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if !user.IsActive {
		return ErrInactiveUser
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		return ErrInvalidCredentials
	}

	hashedNewPassword, err := HashPassword(newPassword)
	if err != nil {
		return err
	}

	return s.userRepo.UpdateUserPassword(ctx, userID, hashedNewPassword)
}

func (s *Service) Logout(ctx context.Context) error {
	return nil
}

func (s *Service) generateJWT(user UserResult) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"iat":  now.Unix(),
		"exp":  now.Add(time.Duration(s.tokenExpiryHours) * time.Hour).Unix(),
		"jti":  generateJTI(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func generateJTI() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
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
