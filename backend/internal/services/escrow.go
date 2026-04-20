package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/skroda/backend/internal/models"
	"github.com/skroda/backend/internal/repository"
)

type EscrowService struct {
	txRepo              *repository.TransactionRepository
	inspectionPeriodHrs int
}

func NewEscrowService(txRepo *repository.TransactionRepository, inspectionHrs int) *EscrowService {
	return &EscrowService{txRepo: txRepo, inspectionPeriodHrs: inspectionHrs}
}

func (s *EscrowService) CreateTransaction(ctx context.Context, buyerIDStr string, req *models.CreateTransactionRequest) (*models.Transaction, error) {
	buyerID, err := uuid.Parse(buyerIDStr)
	if err != nil {
		return nil, errors.New("invalid buyer id")
	}

	tx := &models.Transaction{
		ID:         uuid.New(),
		BuyerID:    buyerID,
		SellerID:   req.SellerID,
		Title:      req.Title,
		Description: req.Description,
		Amount:     req.Amount,
		Currency:   req.Currency,
		Status:     models.StatusPending,
		SellerCity: req.SellerCity,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.txRepo.Create(ctx, tx); err != nil {
		return nil, err
	}

	return tx, nil
}

func (s *EscrowService) GetTransaction(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	return s.txRepo.GetByID(ctx, id)
}

func (s *EscrowService) ListUserTransactions(ctx context.Context, userIDStr string) ([]models.Transaction, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user id")
	}
	return s.txRepo.ListByUser(ctx, userID)
}

func (s *EscrowService) ConfirmDelivery(ctx context.Context, txID uuid.UUID, userIDStr string) error {
	tx, err := s.txRepo.GetByID(ctx, txID)
	if err != nil {
		return err
	}

	if tx.Status != models.StatusInspection {
		return errors.New("transaction is not in inspection period")
	}

	userID, _ := uuid.Parse(userIDStr)
	if tx.BuyerID != userID {
		return errors.New("only the buyer can confirm delivery")
	}

	return s.txRepo.UpdateStatus(ctx, txID, models.StatusCompleted)
}

func (s *EscrowService) CancelTransaction(ctx context.Context, txID uuid.UUID, userIDStr string) error {
	tx, err := s.txRepo.GetByID(ctx, txID)
	if err != nil {
		return err
	}

	if tx.Status != models.StatusPending {
		return errors.New("only pending transactions can be cancelled")
	}

	return s.txRepo.UpdateStatus(ctx, txID, models.StatusCancelled)
}

func (s *EscrowService) RaiseDispute(ctx context.Context, userIDStr string, req *models.CreateDisputeRequest) (*models.Dispute, error) {
	tx, err := s.txRepo.GetByID(ctx, req.TransactionID)
	if err != nil {
		return nil, errors.New("transaction not found")
	}

	if tx.Status != models.StatusInspection && tx.Status != models.StatusInProgress {
		return nil, errors.New("disputes can only be raised during inspection or in-progress")
	}

	raisedBy, _ := uuid.Parse(userIDStr)
	dispute := &models.Dispute{
		ID:            uuid.New(),
		TransactionID: req.TransactionID,
		RaisedByID:    raisedBy,
		Reason:        req.Reason,
		Status:        models.DisputeOpen,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.txRepo.UpdateStatus(ctx, req.TransactionID, models.StatusDisputed); err != nil {
		return nil, err
	}

	return dispute, nil
}

func (s *EscrowService) ResolveDispute(ctx context.Context, disputeID uuid.UUID, agentIDStr string, req *models.ResolveDisputeRequest) error {
	_ = disputeID
	_ = agentIDStr
	_ = req
	return nil
}
