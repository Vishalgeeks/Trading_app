package favorite

import (
	"encoding/json"
	"errors"
	"net/http"

	"mehndi-booking-backend/internal/middleware"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateFavorite(w http.ResponseWriter, r *http.Request) {
	authenticatedUser, ok := middleware.FromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req CreateFavoriteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	fav, err := h.service.CreateFavorite(r.Context(), authenticatedUser.UserID, req.DesignID)
	if err != nil {
		if errors.Is(err, ErrFavoriteAlreadyExists) {
			writeError(w, http.StatusConflict, "favorite already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusCreated, fav.ToResponse())
}

func (h *Handler) DeleteFavorite(w http.ResponseWriter, r *http.Request) {
	authenticatedUser, ok := middleware.FromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req DeleteFavoriteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err := h.service.RemoveFavorite(r.Context(), authenticatedUser.UserID, req.DesignID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListFavorites(w http.ResponseWriter, r *http.Request) {
	authenticatedUser, ok := middleware.FromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	favorites, err := h.service.GetUserFavorites(r.Context(), authenticatedUser.UserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responses := make([]FavoriteResponse, len(favorites))
	for i, fav := range favorites {
		responses[i] = fav.ToResponse()
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"favorites": responses,
		"count":     len(responses),
	})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

type CreateFavoriteRequest struct {
	DesignID string `json:"design_id"`
}

type DeleteFavoriteRequest struct {
	DesignID string `json:"design_id"`
}
