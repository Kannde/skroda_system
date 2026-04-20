package models

import (
	"time"

	"github.com/google/uuid"
)

type AgentStatus string

const (
	AgentPendingApproval AgentStatus = "pending_approval"
	AgentActive          AgentStatus = "active"
	AgentSuspended       AgentStatus = "suspended"
	AgentInactive        AgentStatus = "inactive"
)

type AgentTier string

const (
	TierBronze AgentTier = "bronze"
	TierSilver AgentTier = "silver"
	TierGold   AgentTier = "gold"
)

type AgentProfile struct {
	ID                   uuid.UUID   `json:"id" db:"id"`
	UserID               uuid.UUID   `json:"user_id" db:"user_id"`
	BusinessName         string      `json:"business_name" db:"business_name"`
	BusinessAddress      string      `json:"business_address" db:"business_address"`
	BusinessPhone        string      `json:"business_phone" db:"business_phone"`
	City                 string      `json:"city" db:"city"`
	Region               string      `json:"region" db:"region"`
	GPSLat               *float64    `json:"gps_lat,omitempty" db:"gps_lat"`
	GPSLng               *float64    `json:"gps_lng,omitempty" db:"gps_lng"`
	BusinessRegNumber    *string     `json:"business_reg_number,omitempty" db:"business_reg_number"`
	IDDocumentURL        *string     `json:"id_document_url,omitempty" db:"id_document_url"`
	Verified             bool        `json:"verified" db:"verified"`
	VerifiedAt           *time.Time  `json:"verified_at,omitempty" db:"verified_at"`
	Status               AgentStatus `json:"status" db:"status"`
	Tier                 AgentTier   `json:"tier" db:"tier"`
	Rating               float64     `json:"rating" db:"rating"`
	TotalDeliveries      int         `json:"total_deliveries" db:"total_deliveries"`
	SuccessfulDeliveries int         `json:"successful_deliveries" db:"successful_deliveries"`
	BondAmount           float64     `json:"bond_amount" db:"bond_amount"`
	MaxConcurrent        int         `json:"max_concurrent" db:"max_concurrent"`
	ActiveDeliveries     int         `json:"active_deliveries" db:"active_deliveries"`
	CreatedAt            time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time   `json:"updated_at" db:"updated_at"`
}

type RegisterAgentRequest struct {
	BusinessName      string  `json:"business_name" binding:"required"`
	BusinessAddress   string  `json:"business_address" binding:"required"`
	BusinessPhone     string  `json:"business_phone" binding:"required"`
	City              string  `json:"city" binding:"required"`
	Region            string  `json:"region" binding:"required"`
	GPSLat            float64 `json:"gps_lat"`
	GPSLng            float64 `json:"gps_lng"`
	BusinessRegNumber string  `json:"business_reg_number"`
}
