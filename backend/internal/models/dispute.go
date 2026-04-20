package models

import (
	"time"

	"github.com/google/uuid"
)

type DisputeStatus string

const (
	DisputeOpen           DisputeStatus = "open"
	DisputeUnderReview    DisputeStatus = "under_review"
	DisputeResolvedBuyer  DisputeStatus = "resolved_buyer"
	DisputeResolvedSeller DisputeStatus = "resolved_seller"
	DisputeEscalated      DisputeStatus = "escalated"
)

type DisputeReason string

const (
	ReasonNotAsDescribed DisputeReason = "item_not_as_described"
	ReasonDamaged        DisputeReason = "item_damaged"
	ReasonNotDelivered   DisputeReason = "item_not_delivered"
	ReasonWrongItem      DisputeReason = "wrong_item"
	ReasonPartial        DisputeReason = "partial_delivery"
	ReasonSellerNoShip   DisputeReason = "seller_no_ship"
	ReasonOther          DisputeReason = "other"
)

type Dispute struct {
	ID                 uuid.UUID     `json:"id" db:"id"`
	TransactionID      uuid.UUID     `json:"transaction_id" db:"transaction_id"`
	RaisedBy           uuid.UUID     `json:"raised_by" db:"raised_by"`
	Reason             DisputeReason `json:"reason" db:"reason"`
	Description        string        `json:"description" db:"description"`
	Status             DisputeStatus `json:"status" db:"status"`
	Evidence           interface{}   `json:"evidence" db:"evidence"`
	ResolvedBy         *uuid.UUID    `json:"resolved_by,omitempty" db:"resolved_by"`
	ResolutionNotes    *string       `json:"resolution_notes,omitempty" db:"resolution_notes"`
	ResolvedAt         *time.Time    `json:"resolved_at,omitempty" db:"resolved_at"`
	ResponseDeadline   time.Time     `json:"response_deadline" db:"response_deadline"`
	ResolutionDeadline time.Time     `json:"resolution_deadline" db:"resolution_deadline"`
	CreatedAt          time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time     `json:"updated_at" db:"updated_at"`

	Messages []DisputeMessage `json:"messages,omitempty" db:"-"`
}

type DisputeMessage struct {
	ID          uuid.UUID   `json:"id" db:"id"`
	DisputeID   uuid.UUID   `json:"dispute_id" db:"dispute_id"`
	SenderID    uuid.UUID   `json:"sender_id" db:"sender_id"`
	Message     string      `json:"message" db:"message"`
	Attachments interface{} `json:"attachments" db:"attachments"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
}

type RaiseDisputeRequest struct {
	TransactionID uuid.UUID `json:"transaction_id" binding:"required"`
	Reason        string    `json:"reason" binding:"required,oneof=item_not_as_described item_damaged item_not_delivered wrong_item partial_delivery seller_no_ship other"`
	Description   string    `json:"description" binding:"required,min=20"`
}
