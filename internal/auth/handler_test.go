package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandler_Register_Success(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewService(repo)
	h := NewHandler(svc)

	body := RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Phone:    "+919876543210",
		Password: "securepass123",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Register(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var resp RegisterResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)
	require.Equal(t, "John Doe", resp.Name)
	require.Equal(t, "john@example.com", resp.Email)
	require.Equal(t, "CLIENT", resp.Role)
}

func TestHandler_Register_DuplicateEmail(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewService(repo)
	h := NewHandler(svc)

	body := RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "securepass123",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Register(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	req2 := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(bodyBytes))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()

	h.Register(w2, req2)
	require.Equal(t, http.StatusConflict, w2.Code)
}

func TestHandler_Register_InvalidEmail(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewService(repo)
	h := NewHandler(svc)

	body := RegisterRequest{
		Name:     "John Doe",
		Email:    "invalid-email",
		Password: "securepass123",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Register(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_Register_MissingName(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewService(repo)
	h := NewHandler(svc)

	body := RegisterRequest{
		Email:    "john@example.com",
		Password: "securepass123",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Register(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_Register_MissingPassword(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewService(repo)
	h := NewHandler(svc)

	body := RegisterRequest{
		Name:  "John Doe",
		Email: "john@example.com",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Register(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_Register_ShortPassword(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewService(repo)
	h := NewHandler(svc)

	body := RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "short",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Register(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_Register_MalformedJSON(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewService(repo)
	h := NewHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Register(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}
