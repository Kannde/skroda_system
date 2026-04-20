package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleBuyer  UserRole = "buyer"
	RoleSeller UserRole = "seller"
	RoleAgent  UserRole = "agent"
	RoleAdmin  UserRole = "admin"
)

type User struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	FullName     string     `json:"full_name" db:"full_name" validate:"required,min=2,max=100"`
	Email        string     `json:"email" db:"email" validate:"required,email"`
	Phone        string     `json:"phone" db:"phone" validate:"required"`
	PasswordHash string     `json:"-" db:"password_hash"`
	Role         UserRole   `json:"role" db:"role"`
	City         string     `json:"city" db:"city"`
	Country      string     `json:"country" db:"country"`
	IsVerified   bool       `json:"is_verified" db:"is_verified"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type RegisterRequest struct {
	FullName string   `json:"full_name" validate:"required,min=2,max=100"`
	Email    string   `json:"email" validate:"required,email"`
	Phone    string   `json:"phone" validate:"required"`
	Password string   `json:"password" validate:"required,min=8"`
	Role     UserRole `json:"role" validate:"required,oneof=buyer seller agent"`
	City     string   `json:"city" validate:"required"`
	Country  string   `json:"country" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
