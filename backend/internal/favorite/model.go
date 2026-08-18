package favorite

import (
	"errors"
	"time"
)

type Favorite struct {
	ID        string
	UserID    string
	DesignID  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type FavoriteResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	DesignID  string    `json:"design_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (f Favorite) ToResponse() FavoriteResponse {
	return FavoriteResponse{
		ID:        f.ID,
		UserID:    f.UserID,
		DesignID:  f.DesignID,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}

var ErrFavoriteAlreadyExists = errors.New("favorite already exists")
