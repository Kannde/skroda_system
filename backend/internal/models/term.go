package models

import (
	"time"

	"github.com/google/uuid"
)

type TermStatus string

const (
	TermProposed        TermStatus = "proposed"
	TermAccepted        TermStatus = "accepted"
	TermRejected        TermStatus = "rejected"
	TermCounterProposed TermStatus = "counter_proposed"
)

type Term struct {
	ID                 uuid.UUID  `json:"id" db:"id"`
	TransactionID      uuid.UUID  `json:"transaction_id" db:"transaction_id"`
	Version            int        `json:"version" db:"version"`
	ProposedBy         uuid.UUID  `json:"proposed_by" db:"proposed_by"`
	Status             TermStatus `json:"status" db:"status"`
	ItemDescription    string     `json:"item_description" db:"item_description"`
	ItemCondition      string     `json:"item_condition" db:"item_condition"`
	Quantity           int        `json:"quantity" db:"quantity"`
	Amount             float64    `json:"amount" db:"amount"`
	InspectionHours    int        `json:"inspection_hours" db:"inspection_hours"`
	DeliveryDeadline   int        `json:"delivery_deadline" db:"delivery_deadline"`
	AcceptanceCriteria *string    `json:"acceptance_criteria,omitempty" db:"acceptance_criteria"`
	ImageURLs          []string   `json:"image_urls" db:"image_urls"`
	RejectionReason    *string    `json:"rejection_reason,omitempty" db:"rejection_reason"`
	RespondedAt        *time.Time `json:"responded_at,omitempty" db:"responded_at"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
}

type ProposeTermsRequest struct {
	ItemDescription    string   `json:"item_description" binding:"required"`
	ItemCondition      string   `json:"item_condition" binding:"required"`
	Quantity           int      `json:"quantity" binding:"required,min=1"`
	Amount             float64  `json:"amount" binding:"required,gt=0"`
	InspectionHours    int      `json:"inspection_hours" binding:"omitempty,min=24,max=168"`
	DeliveryDeadline   int      `json:"delivery_deadline" binding:"omitempty,min=24,max=168"`
	AcceptanceCriteria string   `json:"acceptance_criteria"`
	ImageURLs          []string `json:"image_urls"`
}

type RespondToTermsRequest struct {
	Accept          bool   `json:"accept"`
	RejectionReason string `json:"rejection_reason"`
}
