package notification

import "time"

type NotificationType string

const (
	TypeBookingCreated         NotificationType = "BOOKING_CREATED"
	TypeBookingConfirmed       NotificationType = "BOOKING_CONFIRMED"
	TypeBookingCancelled       NotificationType = "BOOKING_CANCELLED"
	TypeBookingCompleted       NotificationType = "BOOKING_COMPLETED"
	TypeNewBooking             NotificationType = "NEW_BOOKING"
	TypeClientCancelledBooking NotificationType = "CLIENT_CANCELLED_BOOKING"
)

type Notification struct {
	ID        string
	UserID    string
	Type      NotificationType
	Title     string
	Message   string
	BookingID *string
	IsRead    bool
	CreatedAt time.Time
	ReadAt    *time.Time
}

type NotificationResponse struct {
	ID        string     `json:"id"`
	Type      string     `json:"type"`
	Title     string     `json:"title"`
	Message   string     `json:"message"`
	BookingID *string    `json:"booking_id,omitempty"`
	IsRead    bool       `json:"is_read"`
	CreatedAt time.Time  `json:"created_at"`
	ReadAt    *time.Time `json:"read_at,omitempty"`
}

func (n Notification) ToResponse() NotificationResponse {
	return NotificationResponse{
		ID:        n.ID,
		Type:      string(n.Type),
		Title:     n.Title,
		Message:   n.Message,
		BookingID: n.BookingID,
		IsRead:    n.IsRead,
		CreatedAt: n.CreatedAt,
		ReadAt:    n.ReadAt,
	}
}
