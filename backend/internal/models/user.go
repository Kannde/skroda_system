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

type AccountStatus string

const (
	StatusActive              AccountStatus = "active"
	StatusSuspended           AccountStatus = "suspended"
	StatusDeactivated         AccountStatus = "deactivated"
	StatusPendingVerification AccountStatus = "pending_verification"
)

type User struct {
	ID                uuid.UUID     `json:"id" db:"id"`
	Email             *string       `json:"email,omitempty" db:"email"`
	Phone             string        `json:"phone" db:"phone"`
	PhoneVerified     bool          `json:"phone_verified" db:"phone_verified"`
	PasswordHash      string        `json:"-" db:"password_hash"`
	FirstName         string        `json:"first_name" db:"first_name"`
	LastName          string        `json:"last_name" db:"last_name"`
	Role              UserRole      `json:"role" db:"role"`
	Status            AccountStatus `json:"status" db:"status"`
	City              *string       `json:"city,omitempty" db:"city"`
	Region            *string       `json:"region,omitempty" db:"region"`
	AvatarURL         *string       `json:"avatar_url,omitempty" db:"avatar_url"`
	TrustScore        float64       `json:"trust_score" db:"trust_score"`
	TotalTransactions int           `json:"total_transactions" db:"total_transactions"`
	CreatedAt         time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at" db:"updated_at"`
}

type PublicProfile struct {
	ID                uuid.UUID `json:"id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	City              *string   `json:"city,omitempty"`
	Role              UserRole  `json:"role"`
	TrustScore        float64   `json:"trust_score"`
	TotalTransactions int       `json:"total_transactions"`
}

func (u *User) ToPublicProfile() PublicProfile {
	return PublicProfile{
		ID:                u.ID,
		FirstName:         u.FirstName,
		LastName:          u.LastName,
		City:              u.City,
		Role:              u.Role,
		TrustScore:        u.TrustScore,
		TotalTransactions: u.TotalTransactions,
	}
}

type RegisterRequest struct {
	Phone     string `json:"phone" binding:"required"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Role      string `json:"role" binding:"required,oneof=buyer seller agent"`
	City      string `json:"city"`
	Region    string `json:"region"`
}

type LoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
