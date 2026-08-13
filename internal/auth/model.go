package auth

import (
	"regexp"
	"time"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	Phone     *string `json:"phone,omitempty"`
	Role      string  `json:"role"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	IsActive  bool    `json:"is_active"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type UserResult struct {
	ID        string
	Name      string
	Email     string
	Phone     *string
	Role      string
	AvatarURL *string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u UserResult) ToResponse() RegisterResponse {
	return RegisterResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Phone:     u.Phone,
		Role:      u.Role,
		AvatarURL: u.AvatarURL,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	}
}
