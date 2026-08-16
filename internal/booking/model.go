package booking

import "time"

type Booking struct {
	ID          string
	UserID      string
	DesignID    string
	BookingDate time.Time
	StartTime   time.Time
	EndTime     time.Time
	Status      string
	Notes       *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DesignName  string
	UserName    string
	UserEmail   string
	UserPhone   *string
}

type BookingResponse struct {
	ID          string    `json:"id"`
	DesignID    string    `json:"design_id"`
	DesignName  string    `json:"design_name"`
	UserName    string    `json:"user_name"`
	UserEmail   string    `json:"user_email"`
	UserPhone   *string   `json:"user_phone,omitempty"`
	BookingDate string    `json:"booking_date"`
	StartTime   string    `json:"start_time"`
	EndTime     string    `json:"end_time"`
	Status      string    `json:"status"`
	Notes       *string   `json:"notes,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AdminBookingResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	UserName    string    `json:"user_name"`
	UserEmail   string    `json:"user_email"`
	UserPhone   *string   `json:"user_phone,omitempty"`
	DesignID    string    `json:"design_id"`
	DesignName  string    `json:"design_name"`
	BookingDate string    `json:"booking_date"`
	StartTime   string    `json:"start_time"`
	EndTime     string    `json:"end_time"`
	Status      string    `json:"status"`
	Notes       *string   `json:"notes,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (b Booking) ToResponse() BookingResponse {
	return BookingResponse{
		ID:          b.ID,
		DesignID:    b.DesignID,
		DesignName:  b.DesignName,
		UserName:    b.UserName,
		UserEmail:   b.UserEmail,
		UserPhone:   b.UserPhone,
		BookingDate: b.BookingDate.Format("2006-01-02"),
		StartTime:   b.StartTime.Format("15:04"),
		EndTime:     b.EndTime.Format("15:04"),
		Status:      b.Status,
		Notes:       b.Notes,
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
	}
}

func (b Booking) ToAdminResponse() AdminBookingResponse {
	return AdminBookingResponse{
		ID:          b.ID,
		UserID:      b.UserID,
		UserName:    b.UserName,
		UserEmail:   b.UserEmail,
		UserPhone:   b.UserPhone,
		DesignID:    b.DesignID,
		DesignName:  b.DesignName,
		BookingDate: b.BookingDate.Format("2006-01-02"),
		StartTime:   b.StartTime.Format("15:04"),
		EndTime:     b.EndTime.Format("15:04"),
		Status:      b.Status,
		Notes:       b.Notes,
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
	}
}

type CreateBookingRequest struct {
	DesignID    string `json:"design_id"`
	BookingDate string `json:"booking_date"`
	StartTime   string `json:"start_time"`
	Notes       string `json:"notes"`
}

type UpdateBookingStatusRequest struct {
	Status string `json:"status"`
}

type Slot struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}
