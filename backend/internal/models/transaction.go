package models

import (
	"time"

	"github.com/google/uuid"
)

type TransactionStatus string

const (
	StatusPending    TransactionStatus = "pending"
	StatusFunded     TransactionStatus = "funded"
	StatusInProgress TransactionStatus = "in_progress"
	StatusInspection TransactionStatus = "inspection"
	StatusCompleted  TransactionStatus = "completed"
	StatusDisputed   TransactionStatus = "disputed"
	StatusCancelled  TransactionStatus = "cancelled"
	StatusRefunded   TransactionStatus = "refunded"
)

type Transaction struct {
	ID               uuid.UUID         `json:"id" db:"id"`
	BuyerID          uuid.UUID         `json:"buyer_id" db:"buyer_id"`
	SellerID         uuid.UUID         `json:"seller_id" db:"seller_id"`
	AgentID          *uuid.UUID        `json:"agent_id,omitempty" db:"agent_id"`
	Title            string            `json:"title" db:"title" validate:"required,min=5,max=200"`
	Description      string            `json:"description" db:"description"`
	Amount           float64           `json:"amount" db:"amount" validate:"required,gt=0"`
	Currency         string            `json:"currency" db:"currency"`
	Status           TransactionStatus `json:"status" db:"status"`
	BuyerCity        string            `json:"buyer_city" db:"buyer_city"`
	SellerCity       string            `json:"seller_city" db:"seller_city"`
	InspectionEndsAt *time.Time        `json:"inspection_ends_at,omitempty" db:"inspection_ends_at"`
	CreatedAt        time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at" db:"updated_at"`
}

type CreateTransactionRequest struct {
	SellerID    uuid.UUID `json:"seller_id" validate:"required"`
	Title       string    `json:"title" validate:"required,min=5,max=200"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount" validate:"required,gt=0"`
	Currency    string    `json:"currency" validate:"required,len=3"`
	SellerCity  string    `json:"seller_city" validate:"required"`
}

type TransactionWithDetails struct {
	Transaction
	Buyer  User   `json:"buyer"`
	Seller User   `json:"seller"`
	Agent  *User  `json:"agent,omitempty"`
	Terms  []Term `json:"terms"`
}
