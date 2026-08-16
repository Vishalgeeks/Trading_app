package design

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateDesign(w http.ResponseWriter, r *http.Request) {
	var req CreateDesignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	des, err := h.service.CreateDesign(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidName):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, ErrInvalidSlug):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, ErrInvalidImageURL):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, ErrInvalidPrice):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, ErrInvalidDuration):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, ErrCategoryNotFound):
			writeError(w, http.StatusNotFound, "category not found")
		case errors.Is(err, ErrCategoryInactive):
			writeError(w, http.StatusBadRequest, "category is not active")
		case errors.Is(err, ErrDuplicateSlug):
			writeError(w, http.StatusConflict, "design slug already exists")
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusCreated, des.ToResponse())
}

func (h *Handler) GetDesign(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/designs/")
	if id == "" || strings.Contains(id, "/") {
		if strings.Contains(id, "/category/") {
			return
		}
		writeError(w, http.StatusBadRequest, "invalid design id")
		return
	}

	des, err := h.service.GetDesign(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrDesignNotFound) {
			writeError(w, http.StatusNotFound, "design not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, des.ToResponse())
}

func (h *Handler) ListDesigns(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	categoryID := r.URL.Query().Get("category_id")
	searchQuery := r.URL.Query().Get("q")

	var catIDPtr *string
	if categoryID != "" {
		catIDPtr = &categoryID
	}

	var searchPtr *string
	if searchQuery != "" {
		searchPtr = &searchQuery
	}

	designs, err := h.service.ListDesigns(r.Context(), catIDPtr, searchPtr, page, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responses := make([]DesignResponse, len(designs))
	for i, des := range designs {
		responses[i] = des.ToResponse()
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"designs": responses,
		"count":   len(responses),
	})
}

func (h *Handler) SearchDesigns(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if strings.TrimSpace(query) == "" {
		writeError(w, http.StatusBadRequest, "search query is required")
		return
	}

	designs, err := h.service.SearchDesigns(r.Context(), query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responses := make([]DesignResponse, len(designs))
	for i, des := range designs {
		responses[i] = des.ToResponse()
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"designs": responses,
		"count":   len(responses),
	})
}

func (h *Handler) ListDesignsByCategory(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/designs/category/")
	categoryID := strings.TrimSpace(path)
	if categoryID == "" {
		writeError(w, http.StatusBadRequest, "invalid category id")
		return
	}

	designs, err := h.service.ListDesignsByCategory(r.Context(), categoryID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responses := make([]DesignResponse, len(designs))
	for i, des := range designs {
		responses[i] = des.ToResponse()
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"designs": responses,
		"count":   len(responses),
	})
}

func (h *Handler) AdminCreateDesign(w http.ResponseWriter, r *http.Request) {
	h.CreateDesign(w, r)
}

func (h *Handler) AdminListDesigns(w http.ResponseWriter, r *http.Request) {
	h.ListDesigns(w, r)
}

func (h *Handler) AdminGetDesign(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/admin/designs/")
	if id == "" || strings.Contains(id, "/") {
		writeError(w, http.StatusBadRequest, "invalid design id")
		return
	}

	des, err := h.service.GetDesign(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrDesignNotFound) {
			writeError(w, http.StatusNotFound, "design not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, des.ToResponse())
}

func (h *Handler) AdminUpdateDesign(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/admin/designs/")
	if id == "" || strings.Contains(id, "/") {
		writeError(w, http.StatusBadRequest, "invalid design id")
		return
	}

	var req UpdateDesignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	des, err := h.service.UpdateDesign(r.Context(), id, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrDesignNotFound):
			writeError(w, http.StatusNotFound, "design not found")
		case errors.Is(err, ErrInvalidName):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, ErrInvalidSlug):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, ErrInvalidImageURL):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, ErrInvalidPrice):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, ErrInvalidDuration):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, ErrCategoryNotFound):
			writeError(w, http.StatusNotFound, "category not found")
		case errors.Is(err, ErrCategoryInactive):
			writeError(w, http.StatusBadRequest, "category is not active")
		case errors.Is(err, ErrDuplicateSlug):
			writeError(w, http.StatusConflict, "design slug already exists")
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, des.ToResponse())
}

func (h *Handler) AdminDeleteDesign(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/admin/designs/")
	if id == "" || strings.Contains(id, "/") {
		writeError(w, http.StatusBadRequest, "invalid design id")
		return
	}

	err := h.service.DeactivateDesign(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrDesignNotFound) {
			writeError(w, http.StatusNotFound, "design not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
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
