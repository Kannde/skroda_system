package models

import (
	"time"

	"github.com/google/uuid"
)

type TransactionStatus string

const (
	TxStatusDraft         TransactionStatus = "draft"
	TxStatusNegotiation   TransactionStatus = "negotiation"
	TxStatusAgreed        TransactionStatus = "agreed"
	TxStatusFunded        TransactionStatus = "funded"
	TxStatusShipped       TransactionStatus = "shipped"
	TxStatusAgentReceived TransactionStatus = "agent_received"
	TxStatusDelivered     TransactionStatus = "delivered"
	TxStatusInspection    TransactionStatus = "inspection"
	TxStatusCompleted     TransactionStatus = "completed"
	TxStatusDisputed      TransactionStatus = "disputed"
	TxStatusCancelled     TransactionStatus = "cancelled"
	TxStatusRefunded      TransactionStatus = "refunded"
)

type TransactionType string

const (
	TxTypeGoods    TransactionType = "goods"
	TxTypeServices TransactionType = "services"
)

var ValidTransitions = map[TransactionStatus][]TransactionStatus{
	TxStatusDraft:         {TxStatusNegotiation, TxStatusCancelled},
	TxStatusNegotiation:   {TxStatusAgreed, TxStatusCancelled},
	TxStatusAgreed:        {TxStatusFunded, TxStatusCancelled},
	TxStatusFunded:        {TxStatusShipped, TxStatusCancelled, TxStatusRefunded},
	TxStatusShipped:       {TxStatusAgentReceived, TxStatusDisputed},
	TxStatusAgentReceived: {TxStatusDelivered, TxStatusDisputed},
	TxStatusDelivered:     {TxStatusInspection},
	TxStatusInspection:    {TxStatusCompleted, TxStatusDisputed},
	TxStatusDisputed:      {TxStatusCompleted, TxStatusRefunded},
}

func (s TransactionStatus) CanTransitionTo(target TransactionStatus) bool {
	allowed, ok := ValidTransitions[s]
	if !ok {
		return false
	}
	for _, a := range allowed {
		if a == target {
			return true
		}
	}
	return false
}

type Transaction struct {
	ID              uuid.UUID         `json:"id" db:"id"`
	ReferenceCode   string            `json:"reference_code" db:"reference_code"`
	TransactionType TransactionType   `json:"transaction_type" db:"transaction_type"`
	Title           string            `json:"title" db:"title"`
	Description     *string           `json:"description,omitempty" db:"description"`
	Status          TransactionStatus `json:"status" db:"status"`

	SellerID uuid.UUID  `json:"seller_id" db:"seller_id"`
	BuyerID  *uuid.UUID `json:"buyer_id,omitempty" db:"buyer_id"`
	AgentID  *uuid.UUID `json:"agent_id,omitempty" db:"agent_id"`

	Amount    float64 `json:"amount" db:"amount"`
	Currency  string  `json:"currency" db:"currency"`
	FeeAmount float64 `json:"fee_amount" db:"fee_amount"`
	FeePaidBy string  `json:"fee_paid_by" db:"fee_paid_by"`

	SellerCity string  `json:"seller_city" db:"seller_city"`
	BuyerCity  *string `json:"buyer_city,omitempty" db:"buyer_city"`

	InspectionHours  int        `json:"inspection_hours" db:"inspection_hours"`
	InspectionEndsAt *time.Time `json:"inspection_ends_at,omitempty" db:"inspection_ends_at"`

	AgreedAt        *time.Time `json:"agreed_at,omitempty" db:"agreed_at"`
	FundedAt        *time.Time `json:"funded_at,omitempty" db:"funded_at"`
	ShippedAt       *time.Time `json:"shipped_at,omitempty" db:"shipped_at"`
	AgentReceivedAt *time.Time `json:"agent_received_at,omitempty" db:"agent_received_at"`
	DeliveredAt     *time.Time `json:"delivered_at,omitempty" db:"delivered_at"`
	CompletedAt     *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	CancelledAt     *time.Time `json:"cancelled_at,omitempty" db:"cancelled_at"`

	InviteToken     *string    `json:"invite_token,omitempty" db:"invite_token"`
	InviteExpiresAt *time.Time `json:"invite_expires_at,omitempty" db:"invite_expires_at"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	Seller *PublicProfile `json:"seller,omitempty" db:"-"`
	Buyer  *PublicProfile `json:"buyer,omitempty" db:"-"`
	Agent  *PublicProfile `json:"agent,omitempty" db:"-"`
	Terms  []Term         `json:"terms,omitempty" db:"-"`
}

type CreateTransactionRequest struct {
	Title           string  `json:"title" binding:"required"`
	Description     string  `json:"description"`
	TransactionType string  `json:"transaction_type" binding:"required,oneof=goods services"`
	Amount          float64 `json:"amount" binding:"required,gt=0"`
	Currency        string  `json:"currency" binding:"omitempty"`
	SellerCity      string  `json:"seller_city" binding:"required"`
	BuyerCity       string  `json:"buyer_city"`
	InspectionHours int     `json:"inspection_hours" binding:"omitempty,min=24,max=168"`
	FeePaidBy       string  `json:"fee_paid_by" binding:"omitempty,oneof=seller buyer split"`

	ItemDescription    string   `json:"item_description" binding:"required"`
	ItemCondition      string   `json:"item_condition" binding:"required,oneof=new used_like_new used_good used_fair"`
	Quantity           int      `json:"quantity" binding:"required,min=1"`
	AcceptanceCriteria string   `json:"acceptance_criteria"`
	ImageURLs          []string `json:"image_urls"`
}
