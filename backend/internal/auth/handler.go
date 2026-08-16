package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"mehndi-booking-backend/internal/middleware"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		writeError(w, http.StatusBadRequest, "content-type must be application/json")
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.service.Register(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrMissingName):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, ErrMissingEmail):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, ErrInvalidEmail):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, ErrMissingPassword):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, ErrPasswordTooShort):
			writeError(w, http.StatusBadRequest, err.Error())
		case err.Error() == "email already exists":
			writeError(w, http.StatusConflict, "email already exists")
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		writeError(w, http.StatusBadRequest, "content-type must be application/json")
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.service.Login(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrMissingEmail):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, ErrMissingPassword):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, ErrInvalidCredentials):
			writeError(w, http.StatusUnauthorized, "invalid email or password")
		case errors.Is(err, ErrInactiveUser):
			writeError(w, http.StatusForbidden, "user account is inactive")
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	authenticatedUser, ok := middleware.FromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	resp, err := h.service.GetCurrentUser(r.Context(), authenticatedUser.UserID)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserNotFound):
			writeError(w, http.StatusNotFound, "user not found")
		case errors.Is(err, ErrInactiveUser):
			writeError(w, http.StatusForbidden, "user account is inactive")
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	authenticatedUser, ok := middleware.FromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err := h.service.ChangePassword(r.Context(), authenticatedUser.UserID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		switch {
		case errors.Is(err, ErrMissingPassword):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, ErrPasswordTooShort):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, ErrInvalidCredentials):
			writeError(w, http.StatusUnauthorized, "invalid current password")
		case errors.Is(err, ErrUserNotFound):
			writeError(w, http.StatusNotFound, "user not found")
		case errors.Is(err, ErrInactiveUser):
			writeError(w, http.StatusForbidden, "user account is inactive")
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ChangePasswordResponse{Message: "password changed successfully"})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	_ = h.service.Logout(r.Context())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(LogoutResponse{Message: "logged out successfully"})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
