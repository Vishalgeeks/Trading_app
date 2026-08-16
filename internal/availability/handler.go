package availability

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateAvailability(w http.ResponseWriter, r *http.Request) {
	var req CreateAvailabilityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.DayOfWeek < 0 || req.DayOfWeek > 6 {
		writeError(w, http.StatusBadRequest, "day_of_week must be between 0 and 6")
		return
	}

	startTime, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid start_time format, use HH:MM")
		return
	}

	endTime, err := time.Parse("15:04", req.EndTime)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid end_time format, use HH:MM")
		return
	}

	if !endTime.After(startTime) {
		writeError(w, http.StatusBadRequest, "end_time must be after start_time")
		return
	}

	av := Availability{
		DayOfWeek: req.DayOfWeek,
		StartTime: startTime,
		EndTime:   endTime,
		IsActive:  true,
	}

	created, err := h.service.CreateAvailability(r.Context(), av)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusCreated, created.ToResponse())
}

func (h *Handler) ListAvailability(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active") == "true"

	availabilities, err := h.service.ListAvailability(r.Context(), activeOnly)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responses := make([]AvailabilityResponse, len(availabilities))
	for i, av := range availabilities {
		responses[i] = av.ToResponse()
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"availabilities": responses,
		"count":          len(responses),
	})
}

func (h *Handler) GetAvailability(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	av, err := h.service.GetAvailability(r.Context(), idStr)
	if err != nil {
		if errors.Is(err, ErrAvailabilityNotFound) {
			writeError(w, http.StatusNotFound, "availability not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, av.ToResponse())
}

func (h *Handler) UpdateAvailability(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	var req UpdateAvailabilityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.StartTime != nil {
		startTime, err := time.Parse("15:04", *req.StartTime)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid start_time format, use HH:MM")
			return
		}
		startStr := startTime.Format("15:04")
		req.StartTime = &startStr
	}

	if req.EndTime != nil {
		endTime, err := time.Parse("15:04", *req.EndTime)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid end_time format, use HH:MM")
			return
		}
		endStr := endTime.Format("15:04")
		req.EndTime = &endStr
	}

	av, err := h.service.UpdateAvailability(r.Context(), idStr, req)
	if err != nil {
		if errors.Is(err, ErrAvailabilityNotFound) {
			writeError(w, http.StatusNotFound, "availability not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, av.ToResponse())
}

func (h *Handler) DeleteAvailability(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	err := h.service.DeleteAvailability(r.Context(), idStr)
	if err != nil {
		if errors.Is(err, ErrAvailabilityNotFound) {
			writeError(w, http.StatusNotFound, "availability not found")
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
