package user

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

func (h *Handler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/users/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	u, err := h.service.GetUserByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		if errors.Is(err, ErrInvalidUUID) {
			writeError(w, http.StatusBadRequest, "invalid user id")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, u.ToResponse())
}

func (h *Handler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/users/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	u, err := h.service.UpdateUserProfile(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		if errors.Is(err, ErrInvalidUUID) {
			writeError(w, http.StatusBadRequest, "invalid user id")
			return
		}
		if errors.Is(err, ErrInvalidName) {
			writeError(w, http.StatusBadRequest, "name cannot be empty")
			return
		}
		if errors.Is(err, ErrInvalidEmail) {
			writeError(w, http.StatusBadRequest, "invalid email format")
			return
		}
		if errors.Is(err, ErrDuplicateEmail) {
			writeError(w, http.StatusConflict, "email already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, u.ToResponse())
}

func (h *Handler) UpdateMyProfile(w http.ResponseWriter, r *http.Request) {
	authenticatedUser, ok := middleware.FromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	u, err := h.service.UpdateUserProfile(r.Context(), authenticatedUser.UserID, req)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		if errors.Is(err, ErrInvalidUUID) {
			writeError(w, http.StatusBadRequest, "invalid user id")
			return
		}
		if errors.Is(err, ErrInvalidName) {
			writeError(w, http.StatusBadRequest, "name cannot be empty")
			return
		}
		if errors.Is(err, ErrInvalidEmail) {
			writeError(w, http.StatusBadRequest, "invalid email format")
			return
		}
		if errors.Is(err, ErrDuplicateEmail) {
			writeError(w, http.StatusConflict, "email already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, u.ToResponse())
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
