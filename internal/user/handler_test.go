package user

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandler_GetUserProfile(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)
	h := NewHandler(svc)

	user, err := svc.CreateUser(context.Background(), "John Doe", "john@example.com", "+919876543210", "password123", RoleClient, nil)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/api/users/"+user.ID, nil)
	w := httptest.NewRecorder()

	h.GetUserProfile(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var resp UserResponse
	err = json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)
	require.Equal(t, user.ID, resp.ID)
	require.Equal(t, "John Doe", resp.Name)
	require.Equal(t, RoleClient, resp.Role)
}

func TestHandler_GetUserProfile_InvalidUUID(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)
	h := NewHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/users/not-a-uuid", nil)
	w := httptest.NewRecorder()

	h.GetUserProfile(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_GetUserProfile_NotFound(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)
	h := NewHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/users/00000000-0000-0000-0000-000000000000", nil)
	w := httptest.NewRecorder()

	h.GetUserProfile(w, req)
	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestHandler_UpdateUserProfile(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)
	h := NewHandler(svc)

	user, err := svc.CreateUser(context.Background(), "John Doe", "john@example.com", "+919876543210", "password123", RoleClient, nil)
	require.NoError(t, err)

	body := map[string]string{
		"name":  "Jane Doe",
		"phone": "+911234567890",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPatch, "/api/users/"+user.ID, bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.UpdateUserProfile(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var resp UserResponse
	err = json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)
	require.Equal(t, "Jane Doe", resp.Name)
	require.Equal(t, "+911234567890", *resp.Phone)
	require.Equal(t, RoleClient, resp.Role)
}

func TestHandler_UpdateUserProfile_InvalidRequestBody(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)
	h := NewHandler(svc)

	req := httptest.NewRequest(http.MethodPatch, "/api/users/00000000-0000-0000-0000-000000000000", bytes.NewReader([]byte("invalid")))
	w := httptest.NewRecorder()

	h.UpdateUserProfile(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_UpdateUserProfile_NotFound(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)
	h := NewHandler(svc)

	body := map[string]string{
		"name": "Jane Doe",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPatch, "/api/users/00000000-0000-0000-0000-000000000000", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.UpdateUserProfile(w, req)
	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestHandler_UpdateUserProfile_InvalidName(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)
	h := NewHandler(svc)

	user, err := svc.CreateUser(context.Background(), "John", "john@example.com", "+919876543210", "password123", RoleClient, nil)
	require.NoError(t, err)

	body := map[string]string{
		"name": "",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPatch, "/api/users/"+user.ID, bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.UpdateUserProfile(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_UpdateUserProfile_DuplicateEmail(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)
	h := NewHandler(svc)

	user1, err := svc.CreateUser(context.Background(), "John", "john@example.com", "+919876543210", "password123", RoleClient, nil)
	require.NoError(t, err)

	_, err = svc.CreateUser(context.Background(), "Jane", "jane@example.com", "+919876543210", "password123", RoleClient, nil)
	require.NoError(t, err)

	body := map[string]string{
		"name":  "Jane",
		"phone": "+911234567890",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPatch, "/api/users/"+user1.ID, bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.UpdateUserProfile(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
