package user

import (
	"database/sql/driver"
	"time"
)

type Role string

const (
	RoleClient Role = "CLIENT"
	RoleAdmin  Role = "ADMIN"
)

func (r Role) String() string {
	return string(r)
}

func (r *Role) Scan(value interface{}) error {
	if value == nil {
		*r = ""
		return nil
	}
	switch v := value.(type) {
	case string:
		*r = Role(v)
	case []byte:
		*r = Role(string(v))
	default:
		return &scanError{msg: "cannot scan into Role"}
	}
	return nil
}

func (r Role) Value() (driver.Value, error) {
	return string(r), nil
}

type scanError struct {
	msg string
}

func (e *scanError) Error() string {
	return e.msg
}

type User struct {
	ID           string
	Name         string
	Email        string
	Phone        *string
	PasswordHash string
	Role         Role
	AvatarURL    *string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     *string   `json:"phone,omitempty"`
	Role      Role      `json:"role"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Phone:     u.Phone,
		Role:      u.Role,
		AvatarURL: u.AvatarURL,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

type UpdateProfileRequest struct {
	Name      *string `json:"name"`
	Phone     *string `json:"phone"`
	AvatarURL *string `json:"avatar_url"`
}
