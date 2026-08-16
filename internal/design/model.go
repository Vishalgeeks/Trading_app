package design

import "time"

type Design struct {
	ID              string
	CategoryID      string
	Name            string
	Slug            string
	Description     *string
	ImageURL        string
	Price           string
	DurationMinutes int
	IsActive        bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Category        *Category
}

type DesignResponse struct {
	ID              string    `json:"id"`
	Category        Category  `json:"category"`
	Name            string    `json:"name"`
	Slug            string    `json:"slug"`
	Description     *string   `json:"description,omitempty"`
	ImageURL        string    `json:"image_url"`
	Price           string    `json:"price"`
	DurationMinutes int       `json:"duration_minutes"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Category struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description *string   `json:"description,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (d Design) ToResponse() DesignResponse {
	category := Category{}
	if d.Category != nil {
		category = *d.Category
	}

	return DesignResponse{
		ID:              d.ID,
		Category:        category,
		Name:            d.Name,
		Slug:            d.Slug,
		Description:     d.Description,
		ImageURL:        d.ImageURL,
		Price:           d.Price,
		DurationMinutes: d.DurationMinutes,
		IsActive:        d.IsActive,
		CreatedAt:       d.CreatedAt,
		UpdatedAt:       d.UpdatedAt,
	}
}

type CreateDesignRequest struct {
	CategoryID      string  `json:"category_id"`
	Name            string  `json:"name"`
	Slug            string  `json:"slug"`
	Description     *string `json:"description"`
	ImageURL        string  `json:"image_url"`
	Price           string  `json:"price"`
	DurationMinutes int     `json:"duration_minutes"`
}

type UpdateDesignRequest struct {
	CategoryID      *string `json:"category_id"`
	Name            *string `json:"name"`
	Slug            *string `json:"slug"`
	Description     *string `json:"description"`
	ImageURL        *string `json:"image_url"`
	Price           *string `json:"price"`
	DurationMinutes *int    `json:"duration_minutes"`
	IsActive        *bool   `json:"is_active"`
}

type CategoryResult struct {
	ID   string
	Name string
	Slug string
}
