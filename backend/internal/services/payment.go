package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/skroda/backend/internal/models"
	"github.com/skroda/backend/internal/repository"
)

type PaymentService struct {
	paymentRepo *repository.PaymentRepository
	txRepo      *repository.TransactionRepository
	momoAPIKey  string
}

func NewPaymentService(paymentRepo *repository.PaymentRepository, txRepo *repository.TransactionRepository, momoAPIKey string) *PaymentService {
	return &PaymentService{paymentRepo: paymentRepo, txRepo: txRepo, momoAPIKey: momoAPIKey}
}

func (s *PaymentService) Initiate(ctx context.Context, payerIDStr string, req *models.InitiatePaymentRequest) (*models.Payment, error) {
	tx, err := s.txRepo.GetByID(ctx, req.TransactionID)
	if err != nil {
		return nil, err
	}

	payerID, _ := uuid.Parse(payerIDStr)
	payment := &models.Payment{
		ID:            uuid.New(),
		TransactionID: req.TransactionID,
		PayerID:       payerID,
		Amount:        tx.Amount,
		Currency:      tx.Currency,
		Provider:      req.Provider,
		Status:        models.PaymentPending,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *PaymentService) HandleMoMoCallback(ctx context.Context, payload map[string]interface{}) error {
	_ = payload
	return nil
}
