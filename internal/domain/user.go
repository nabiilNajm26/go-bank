package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	Email           string     `json:"email" db:"email"`
	PasswordHash    string     `json:"-" db:"password_hash"`
	FullName        string     `json:"full_name" db:"full_name"`
	Phone           *string    `json:"phone,omitempty" db:"phone"`
	ProfileImageURL *string    `json:"profile_image_url,omitempty" db:"profile_image_url"`
	IsVerified      bool       `json:"is_verified" db:"is_verified"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name" validate:"required,min=3,max=255"`
	Phone    string `json:"phone,omitempty" validate:"omitempty,e164"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserRequest struct {
	FullName string `json:"full_name,omitempty" validate:"omitempty,min=3,max=255"`
	Phone    string `json:"phone,omitempty" validate:"omitempty,e164"`
}

type AuthResponse struct {
	User         *User  `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}