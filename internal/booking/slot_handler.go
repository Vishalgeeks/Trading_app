package booking

import (
	"errors"
	"net/http"
	"time"
)

type SlotHandler struct {
	service *Service
}

func NewSlotHandler(service *Service) *SlotHandler {
	return &SlotHandler{service: service}
}

func (h *SlotHandler) GetAvailableSlots(w http.ResponseWriter, r *http.Request) {
	designID := r.URL.Query().Get("design_id")
	bookingDate := r.URL.Query().Get("date")

	if designID == "" {
		writeError(w, http.StatusBadRequest, "design_id is required")
		return
	}

	if bookingDate == "" {
		writeError(w, http.StatusBadRequest, "date is required")
		return
	}

	_, err := time.Parse("2006-01-02", bookingDate)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid date format, use YYYY-MM-DD")
		return
	}

	slots, err := h.service.CalculateAvailableSlots(r.Context(), designID, bookingDate)
	if err != nil {
		switch {
		case errors.Is(err, ErrBookingConflict):
			writeError(w, http.StatusConflict, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"slots": slots,
		"count": len(slots),
	})
}
