package booking

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

	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}

	var datePtr *string
	if date != "" {
		datePtr = &date
	}

	bookings, err := h.service.ListUserBookings(r.Context(), user.UserID, statusPtr, datePtr)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responses := make([]BookingResponse, len(bookings))
	for i, booking := range bookings {
		responses[i] = booking.ToResponse()
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"bookings": responses,
		"count":    len(responses),
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

func (h *Handler) AdminListBookings(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	date := r.URL.Query().Get("date")
	fromDate := r.URL.Query().Get("from")
	toDate := r.URL.Query().Get("to")

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

	bookings, err := h.service.ListAdminBookings(r.Context(), statusPtr, datePtr, fromDatePtr, toDatePtr)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responses := make([]AdminBookingResponse, len(bookings))
	for i, booking := range bookings {
		responses[i] = booking.ToAdminResponse()
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"bookings": responses,
		"count":    len(responses),
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

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
