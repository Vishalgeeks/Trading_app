package availability

import "time"

type Availability struct {
	ID        string
	DayOfWeek int
	StartTime time.Time
	EndTime   time.Time
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AvailabilityResponse struct {
	ID        string    `json:"id"`
	DayOfWeek int       `json:"day_of_week"`
	StartTime string    `json:"start_time"`
	EndTime   string    `json:"end_time"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (a Availability) ToResponse() AvailabilityResponse {
	return AvailabilityResponse{
		ID:        a.ID,
		DayOfWeek: a.DayOfWeek,
		StartTime: a.StartTime.Format("15:04"),
		EndTime:   a.EndTime.Format("15:04"),
		IsActive:  a.IsActive,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

type CreateAvailabilityRequest struct {
	DayOfWeek int    `json:"day_of_week"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type UpdateAvailabilityRequest struct {
	DayOfWeek *int    `json:"day_of_week"`
	StartTime *string `json:"start_time"`
	EndTime   *string `json:"end_time"`
	IsActive  *bool   `json:"is_active"`
}
