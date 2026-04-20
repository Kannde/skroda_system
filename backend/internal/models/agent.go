package models

import (
	"time"

	"github.com/google/uuid"
)

type AgentStatus string

const (
	AgentActive    AgentStatus = "active"
	AgentInactive  AgentStatus = "inactive"
	AgentSuspended AgentStatus = "suspended"
)

type Agent struct {
	ID             uuid.UUID   `json:"id" db:"id"`
	UserID         uuid.UUID   `json:"user_id" db:"user_id"`
	City           string      `json:"city" db:"city"`
	Country        string      `json:"country" db:"country"`
	Bio            string      `json:"bio" db:"bio"`
	Rating         float64     `json:"rating" db:"rating"`
	TotalHandled   int         `json:"total_handled" db:"total_handled"`
	Status         AgentStatus `json:"status" db:"status"`
	CommissionRate float64     `json:"commission_rate" db:"commission_rate"`
	CreatedAt      time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at" db:"updated_at"`
}

type AgentWithUser struct {
	Agent
	User User `json:"user"`
}

type UpdateAgentRequest struct {
	Bio            string  `json:"bio"`
	CommissionRate float64 `json:"commission_rate" validate:"omitempty,min=0,max=10"`
}
