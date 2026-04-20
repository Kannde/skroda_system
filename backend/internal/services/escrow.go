package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/skroda/backend/internal/models"
	"github.com/skroda/backend/internal/repository"
	"github.com/skroda/backend/internal/utils"
)

type EscrowService struct {
	txRepo              *repository.TransactionRepository
	inspectionPeriodHrs int
}

func NewEscrowService(txRepo *repository.TransactionRepository, inspectionHrs int) *EscrowService {
	return &EscrowService{txRepo: txRepo, inspectionPeriodHrs: inspectionHrs}
}

func (s *EscrowService) CreateTransaction(ctx context.Context, sellerIDStr string, req *models.CreateTransactionRequest) (*models.Transaction, error) {
	sellerID, err := uuid.Parse(sellerIDStr)
	if err != nil {
		return nil, errors.New("invalid seller id")
	}

	fee := CalculateFee(req.Amount, DefaultFeeSchedule)
	feePaidBy := req.FeePaidBy
	if feePaidBy == "" {
		feePaidBy = "seller"
	}
	inspectionHours := req.InspectionHours
	if inspectionHours == 0 {
		inspectionHours = s.inspectionPeriodHrs
	}
	currency := req.Currency
	if currency == "" {
		currency = "GHS"
	}

	inviteToken, err := utils.GenerateInviteToken()
	if err != nil {
		return nil, err
	}
	inviteExpires := time.Now().Add(72 * time.Hour)
	desc := req.Description
	buyerCity := req.BuyerCity

	tx := &models.Transaction{
		ID:              uuid.New(),
		ReferenceCode:   utils.GenerateReferenceCode(),
		TransactionType: models.TransactionType(req.TransactionType),
		Title:           req.Title,
		Description:     &desc,
		Status:          models.TxStatusDraft,
		SellerID:        sellerID,
		Amount:          req.Amount,
		Currency:        currency,
		FeeAmount:       fee,
		FeePaidBy:       feePaidBy,
		SellerCity:      req.SellerCity,
		BuyerCity:       &buyerCity,
		InspectionHours: inspectionHours,
		InviteToken:     &inviteToken,
		InviteExpiresAt: &inviteExpires,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
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

func (s *EscrowService) Transition(ctx context.Context, txID uuid.UUID, userIDStr string, target models.TransactionStatus) error {
	tx, err := s.txRepo.GetByID(ctx, txID)
	if err != nil {
		return err
	}

	if !tx.Status.CanTransitionTo(target) {
		return errors.New("invalid status transition from " + string(tx.Status) + " to " + string(target))
	}

	return s.txRepo.UpdateStatus(ctx, txID, target)
}

func (s *EscrowService) ConfirmDelivery(ctx context.Context, txID uuid.UUID, userIDStr string) error {
	tx, err := s.txRepo.GetByID(ctx, txID)
	if err != nil {
		return err
	}

	if tx.Status != models.TxStatusInspection {
		return errors.New("transaction is not in inspection period")
	}

	buyerID, _ := uuid.Parse(userIDStr)
	if tx.BuyerID == nil || *tx.BuyerID != buyerID {
		return errors.New("only the buyer can confirm delivery")
	}

	return s.txRepo.UpdateStatus(ctx, txID, models.TxStatusCompleted)
}

func (s *EscrowService) CancelTransaction(ctx context.Context, txID uuid.UUID, userIDStr string) error {
	tx, err := s.txRepo.GetByID(ctx, txID)
	if err != nil {
		return err
	}

	if !tx.Status.CanTransitionTo(models.TxStatusCancelled) {
		return errors.New("transaction cannot be cancelled in its current state")
	}

	return s.txRepo.UpdateStatus(ctx, txID, models.TxStatusCancelled)
}

func (s *EscrowService) RaiseDispute(ctx context.Context, userIDStr string, req *models.RaiseDisputeRequest) (*models.Dispute, error) {
	tx, err := s.txRepo.GetByID(ctx, req.TransactionID)
	if err != nil {
		return nil, errors.New("transaction not found")
	}

	if !tx.Status.CanTransitionTo(models.TxStatusDisputed) {
		return nil, errors.New("disputes cannot be raised in the current transaction state")
	}

	raisedBy, _ := uuid.Parse(userIDStr)
	now := time.Now()
	dispute := &models.Dispute{
		ID:                 uuid.New(),
		TransactionID:      req.TransactionID,
		RaisedBy:           raisedBy,
		Reason:             models.DisputeReason(req.Reason),
		Description:        req.Description,
		Status:             models.DisputeOpen,
		ResponseDeadline:   now.Add(48 * time.Hour),
		ResolutionDeadline: now.Add(7 * 24 * time.Hour),
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := s.txRepo.UpdateStatus(ctx, req.TransactionID, models.TxStatusDisputed); err != nil {
		return nil, err
	}

	return dispute, nil
}

func (s *EscrowService) ResolveDispute(ctx context.Context, disputeID uuid.UUID, agentIDStr string, resolution models.DisputeStatus, notes string) error {
	_ = disputeID
	_ = agentIDStr
	_ = resolution
	_ = notes
	return nil
}
