package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	IsActive     bool      `json:"is_active"`
	Roles        []Role    `json:"roles,omitempty"`
}

type UserAuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserResponse struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type AddRoleToUserRequest struct {
	RoleID int `json:"role_id"`
}
