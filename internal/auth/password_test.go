package auth

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	password := "mySecurePassword123"
	hash, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hash)
	require.NotEqual(t, password, hash)
}

func TestHashPassword_Empty(t *testing.T) {
	hash, err := HashPassword("")
	require.Error(t, err)
	require.Empty(t, hash)
	require.ErrorIs(t, err, ErrEmptyPassword)
}

func TestComparePassword_Correct(t *testing.T) {
	password := "mySecurePassword123"
	hash, err := HashPassword(password)
	require.NoError(t, err)

	err = ComparePassword(hash, password)
	require.NoError(t, err)
}

func TestComparePassword_Incorrect(t *testing.T) {
	password := "mySecurePassword123"
	hash, err := HashPassword(password)
	require.NoError(t, err)

	err = ComparePassword(hash, "wrongPassword")
	require.Error(t, err)
}

func TestComparePassword_EmptyPassword(t *testing.T) {
	hash, err := HashPassword("somePassword")
	require.NoError(t, err)

	err = ComparePassword(hash, "")
	require.Error(t, err)
	require.ErrorIs(t, err, ErrEmptyPassword)
}

func TestHashPassword_Verifiable(t *testing.T) {
	password := "verifiablePassword"
	hash, err := HashPassword(password)
	require.NoError(t, err)

	err = ComparePassword(hash, password)
	require.NoError(t, err)
}
