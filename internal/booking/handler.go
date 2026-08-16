package booking

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"mehndi-booking-backend/internal/middleware"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.FromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	booking, err := h.service.CreateBooking(r.Context(), user.UserID, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrBookingConflict):
			writeError(w, http.StatusConflict, "slot already booked")
		case errors.Is(err, ErrInvalidStatus):
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusCreated, booking.ToResponse())
}

func (h *Handler) ListBookings(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.FromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	status := r.URL.Query().Get("status")
	date := r.URL.Query().Get("date")
	fromDate := r.URL.Query().Get("from")
	toDate := r.URL.Query().Get("to")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}

	var datePtr *string
	if date != "" {
		datePtr = &date
	}

	var fromDatePtr *string
	if fromDate != "" {
		fromDatePtr = &fromDate
	}

	var toDatePtr *string
	if toDate != "" {
		toDatePtr = &toDate
	}

	bookings, err := h.service.ListUserBookingsPaginated(r.Context(), user.UserID, statusPtr, datePtr, fromDatePtr, toDatePtr, limit, offset)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	total, err := h.service.CountUserBookings(r.Context(), user.UserID, statusPtr, datePtr, fromDatePtr, toDatePtr)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responses := make([]BookingResponse, len(bookings))
	for i, b := range bookings {
		responses[i] = b.ToResponse()
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"bookings": responses,
		"count":    len(responses),
		"total":    total,
		"limit":    limit,
		"offset":   offset,
	})
}

func (h *Handler) GetBooking(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.FromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	booking, err := h.service.GetBooking(r.Context(), id, user.UserID, false)
	if err != nil {
		if errors.Is(err, ErrBookingNotFound) {
			writeError(w, http.StatusNotFound, "booking not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, booking.ToResponse())
}

func (h *Handler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.FromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	booking, err := h.service.CancelBooking(r.Context(), id, user.UserID, false)
	if err != nil {
		switch {
		case errors.Is(err, ErrBookingNotFound):
			writeError(w, http.StatusNotFound, "booking not found")
		case errors.Is(err, ErrInvalidStatus):
			writeError(w, http.StatusConflict, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, booking.ToResponse())
}

func (h *Handler) ListUpcomingBookings(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.FromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	bookings, err := h.service.ListClientUpcomingBookings(r.Context(), user.UserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responses := make([]BookingResponse, len(bookings))
	for i, b := range bookings {
		responses[i] = b.ToResponse()
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"bookings": responses,
		"count":    len(responses),
	})
}

func (h *Handler) ListBookingHistory(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.FromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	bookings, err := h.service.ListClientBookingHistory(r.Context(), user.UserID, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responses := make([]BookingResponse, len(bookings))
	for i, b := range bookings {
		responses[i] = b.ToResponse()
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"bookings": responses,
		"count":    len(responses),
		"limit":    limit,
		"offset":   offset,
	})
}

func (h *Handler) AdminListBookings(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	date := r.URL.Query().Get("date")
	fromDate := r.URL.Query().Get("from")
	toDate := r.URL.Query().Get("to")
	search := r.URL.Query().Get("search")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}

	var datePtr *string
	if date != "" {
		datePtr = &date
	}

	var fromDatePtr *string
	if fromDate != "" {
		fromDatePtr = &fromDate
	}

	var toDatePtr *string
	if toDate != "" {
		toDatePtr = &toDate
	}

	var searchPtr *string
	if search != "" {
		searchPtr = &search
	}

	bookings, err := h.service.ListAdminBookingsPaginated(r.Context(), statusPtr, datePtr, fromDatePtr, toDatePtr, searchPtr, limit, offset)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	total, err := h.service.CountAdminBookings(r.Context(), statusPtr, datePtr, fromDatePtr, toDatePtr, searchPtr)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responses := make([]AdminBookingResponse, len(bookings))
	for i, b := range bookings {
		responses[i] = b.ToAdminResponse()
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"bookings": responses,
		"count":    len(responses),
		"total":    total,
		"limit":    limit,
		"offset":   offset,
	})
}

func (h *Handler) AdminGetBooking(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	booking, err := h.service.GetBooking(r.Context(), id, "", true)
	if err != nil {
		if errors.Is(err, ErrBookingNotFound) {
			writeError(w, http.StatusNotFound, "booking not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, booking.ToAdminResponse())
}

func (h *Handler) AdminUpdateBookingStatus(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	var req UpdateBookingStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	booking, err := h.service.UpdateBookingStatus(r.Context(), id, req.Status)
	if err != nil {
		switch {
		case errors.Is(err, ErrBookingNotFound):
			writeError(w, http.StatusNotFound, "booking not found")
		case errors.Is(err, ErrInvalidStatus):
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusConflict, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, booking.ToAdminResponse())
}

func (h *Handler) AdminUpcomingBookings(w http.ResponseWriter, r *http.Request) {
	bookings, err := h.service.ListAdminUpcomingBookings(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responses := make([]AdminBookingResponse, len(bookings))
	for i, b := range bookings {
		responses[i] = b.ToAdminResponse()
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"bookings": responses,
		"count":    len(responses),
	})
}

func (h *Handler) AdminBookingHistory(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	bookings, err := h.service.ListAdminBookingHistory(r.Context(), limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responses := make([]AdminBookingResponse, len(bookings))
	for i, b := range bookings {
		responses[i] = b.ToAdminResponse()
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"bookings": responses,
		"count":    len(responses),
		"limit":    limit,
		"offset":   offset,
	})
}

func (h *Handler) AdminBookingStats(w http.ResponseWriter, r *http.Request) {
	fromDate := r.URL.Query().Get("from")
	toDate := r.URL.Query().Get("to")

	var fromDatePtr *string
	if fromDate != "" {
		fromDatePtr = &fromDate
	}

	var toDatePtr *string
	if toDate != "" {
		toDatePtr = &toDate
	}

	stats, err := h.service.GetAdminBookingStats(r.Context(), fromDatePtr, toDatePtr)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, stats)
}

func (h *Handler) AdminCancelBooking(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	booking, err := h.service.CancelBooking(r.Context(), id, "", true)
	if err != nil {
		switch {
		case errors.Is(err, ErrBookingNotFound):
			writeError(w, http.StatusNotFound, "booking not found")
		case errors.Is(err, ErrInvalidStatus):
			writeError(w, http.StatusConflict, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, booking.ToAdminResponse())
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
