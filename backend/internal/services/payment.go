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
	txn, err := s.txRepo.GetByID(ctx, req.TransactionID)
	if err != nil {
		return nil, err
	}

	payerID, _ := uuid.Parse(payerIDStr)
	payment := &models.Payment{
		TransactionID:  req.TransactionID,
		PaymentType:    models.PayTypeEscrowDeposit,
		PaymentMethod:  models.PaymentMethod(req.PaymentMethod),
		Status:         models.PayPending,
		Amount:         txn.Amount,
		Currency:       txn.Currency,
		PayerID:        payerID,
		IdempotencyKey: req.TransactionID.String() + ":" + payerIDStr + ":" + string(time.Now().UTC().Format("20060102")),
		InitiatedAt:    time.Now(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
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
