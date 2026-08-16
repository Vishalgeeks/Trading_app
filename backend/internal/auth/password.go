package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	DefaultCost = 12
)

var ErrEmptyPassword = fmt.Errorf("password cannot be empty")

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", ErrEmptyPassword
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashed), nil
}

func ComparePassword(hash string, password string) error {
	if password == "" {
		return ErrEmptyPassword
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return fmt.Errorf("invalid password")
	}

	return nil
}
