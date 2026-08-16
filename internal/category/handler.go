package category

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	cat, err := h.service.CreateCategory(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidName):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, ErrInvalidSlug):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, ErrDuplicateName):
			writeError(w, http.StatusConflict, "category name already exists")
		case errors.Is(err, ErrDuplicateSlug):
			writeError(w, http.StatusConflict, "category slug already exists")
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusCreated, cat.ToResponse())
}

func (h *Handler) GetCategory(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "invalid category id")
		return
	}

	cat, err := h.service.GetCategory(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			writeError(w, http.StatusNotFound, "category not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if !cat.IsActive {
		writeError(w, http.StatusNotFound, "category not found")
		return
	}

	writeJSON(w, http.StatusOK, cat.ToResponse())
}

func (h *Handler) ListCategories(w http.ResponseWriter, r *http.Request) {
	cats, err := h.service.ListCategories(r.Context(), true)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responses := make([]CategoryResponse, len(cats))
	for i, cat := range cats {
		responses[i] = cat.ToResponse()
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"categories": responses,
		"count":      len(responses),
	})
}

func (h *Handler) AdminListCategories(w http.ResponseWriter, r *http.Request) {
	cats, err := h.service.ListCategories(r.Context(), false)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responses := make([]CategoryResponse, len(cats))
	for i, cat := range cats {
		responses[i] = cat.ToResponse()
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"categories": responses,
		"count":      len(responses),
	})
}

func (h *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/admin/categories/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "invalid category id")
		return
	}

	var req UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	cat, err := h.service.UpdateCategory(r.Context(), id, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrCategoryNotFound):
			writeError(w, http.StatusNotFound, "category not found")
		case errors.Is(err, ErrInvalidName):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, ErrInvalidSlug):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, ErrDuplicateName):
			writeError(w, http.StatusConflict, "category name already exists")
		case errors.Is(err, ErrDuplicateSlug):
			writeError(w, http.StatusConflict, "category slug already exists")
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, cat.ToResponse())
}

func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/admin/categories/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "invalid category id")
		return
	}

	err := h.service.DeactivateCategory(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, ErrCategoryNotFound):
			writeError(w, http.StatusNotFound, "category not found")
		case errors.Is(err, ErrCategoryHasDesigns):
			writeError(w, http.StatusConflict, "category has active designs and cannot be deactivated")
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
