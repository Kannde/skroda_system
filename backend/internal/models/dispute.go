package models

import (
	"time"

	"github.com/google/uuid"
)

type DisputeStatus string

const (
	DisputeOpen     DisputeStatus = "open"
	DisputeReview   DisputeStatus = "under_review"
	DisputeResolved DisputeStatus = "resolved"
	DisputeClosed   DisputeStatus = "closed"
)

type DisputeResolution string

const (
	ResolutionBuyer  DisputeResolution = "buyer"
	ResolutionSeller DisputeResolution = "seller"
	ResolutionSplit  DisputeResolution = "split"
)

type Dispute struct {
	ID            uuid.UUID          `json:"id" db:"id"`
	TransactionID uuid.UUID          `json:"transaction_id" db:"transaction_id"`
	RaisedByID    uuid.UUID          `json:"raised_by_id" db:"raised_by_id"`
	AgentID       *uuid.UUID         `json:"agent_id,omitempty" db:"agent_id"`
	Reason        string             `json:"reason" db:"reason" validate:"required,min=20"`
	Status        DisputeStatus      `json:"status" db:"status"`
	Resolution    *DisputeResolution `json:"resolution,omitempty" db:"resolution"`
	ResolutionNote string            `json:"resolution_note,omitempty" db:"resolution_note"`
	CreatedAt     time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at" db:"updated_at"`
	ResolvedAt    *time.Time         `json:"resolved_at,omitempty" db:"resolved_at"`
}

type CreateDisputeRequest struct {
	TransactionID uuid.UUID `json:"transaction_id" validate:"required"`
	Reason        string    `json:"reason" validate:"required,min=20"`
}

type ResolveDisputeRequest struct {
	Resolution     DisputeResolution `json:"resolution" validate:"required,oneof=buyer seller split"`
	ResolutionNote string            `json:"resolution_note" validate:"required,min=10"`
}
