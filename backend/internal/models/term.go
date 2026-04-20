package models

import (
	"time"

	"github.com/google/uuid"
)

type Term struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	TransactionID uuid.UUID  `json:"transaction_id" db:"transaction_id"`
	Description   string     `json:"description" db:"description" validate:"required"`
	BuyerAgreed   bool       `json:"buyer_agreed" db:"buyer_agreed"`
	SellerAgreed  bool       `json:"seller_agreed" db:"seller_agreed"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

type CreateTermRequest struct {
	TransactionID uuid.UUID `json:"transaction_id" validate:"required"`
	Description   string    `json:"description" validate:"required,min=10"`
}
