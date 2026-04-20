package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/skroda/backend/internal/models"
	"github.com/skroda/backend/internal/repository"
	"github.com/skroda/backend/internal/services"
	"github.com/skroda/backend/internal/utils"
)

type TransactionHandler struct {
	txnRepo   *repository.TransactionRepository
	auditRepo *repository.AuditRepository
}

func NewTransactionHandler(txnRepo *repository.TransactionRepository, auditRepo *repository.AuditRepository) *TransactionHandler {
	return &TransactionHandler{txnRepo: txnRepo, auditRepo: auditRepo}
}

func (h *TransactionHandler) Create(c *gin.Context) {
	var req models.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sellerID, _ := getUserID(c)

	inviteToken, err := utils.GenerateInviteToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate invite"})
		return
	}

	feeAmount := services.CalculateFee(req.Amount, services.DefaultFeeSchedule)

	currency := req.Currency
	if currency == "" {
		currency = "GHS"
	}
	inspectionHours := req.InspectionHours
	if inspectionHours == 0 {
		inspectionHours = 48
	}
	feePaidBy := req.FeePaidBy
	if feePaidBy == "" {
		feePaidBy = "seller"
	}

	inviteExpiry := time.Now().Add(72 * time.Hour)
	buyerCity := req.BuyerCity

	txn := &models.Transaction{
		ReferenceCode:   utils.GenerateReferenceCode(),
		TransactionType: models.TransactionType(req.TransactionType),
		Title:           req.Title,
		Description:     &req.Description,
		Status:          models.TxStatusDraft,
		SellerID:        sellerID,
		Amount:          req.Amount,
		Currency:        currency,
		FeeAmount:       feeAmount,
		FeePaidBy:       feePaidBy,
		SellerCity:      req.SellerCity,
		BuyerCity:       &buyerCity,
		InspectionHours: inspectionHours,
		InviteToken:     &inviteToken,
		InviteExpiresAt: &inviteExpiry,
	}

	if err := h.txnRepo.Create(c.Request.Context(), txn); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create transaction"})
		return
	}

	ip := c.ClientIP()
	_ = h.auditRepo.Log(c.Request.Context(), "transaction", txn.ID, "created", nil, txn, &sellerID, &ip)

	c.JSON(http.StatusCreated, gin.H{
		"transaction": txn,
		"invite_link": "/join/" + inviteToken,
	})
}

func (h *TransactionHandler) GetByID(c *gin.Context) {
	txnID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction ID"})
		return
	}

	txn, err := h.txnRepo.GetByID(c.Request.Context(), txnID)
	if err != nil {
		if err == repository.ErrTransactionNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch transaction"})
		return
	}

	userID, _ := getUserID(c)
	isParty := txn.SellerID == userID ||
		(txn.BuyerID != nil && *txn.BuyerID == userID) ||
		(txn.AgentID != nil && *txn.AgentID == userID)
	if !isParty {
		c.JSON(http.StatusForbidden, gin.H{"error": "you are not a party to this transaction"})
		return
	}

	c.JSON(http.StatusOK, txn)
}

func (h *TransactionHandler) List(c *gin.Context) {
	userID, _ := getUserID(c)
	status := c.Query("status")

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if limit > 100 {
		limit = 100
	}

	transactions, total, err := h.txnRepo.ListByUser(c.Request.Context(), userID, status, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"total":        total,
		"limit":        limit,
		"offset":       offset,
	})
}

func (h *TransactionHandler) JoinByInvite(c *gin.Context) {
	token := c.Param("token")
	buyerID, _ := getUserID(c)

	txn, err := h.txnRepo.GetByInviteToken(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "invalid or expired invite"})
		return
	}

	if txn.InviteExpiresAt != nil && time.Now().After(*txn.InviteExpiresAt) {
		c.JSON(http.StatusGone, gin.H{"error": "invite has expired"})
		return
	}
	if txn.SellerID == buyerID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot join your own transaction as buyer"})
		return
	}
	if txn.BuyerID != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "transaction already has a buyer"})
		return
	}

	buyerCity := c.Query("city")
	if err := h.txnRepo.AssignBuyer(c.Request.Context(), txn.ID, buyerID, buyerCity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to join transaction"})
		return
	}

	ip := c.ClientIP()
	_ = h.auditRepo.Log(c.Request.Context(), "transaction", txn.ID, "buyer_joined",
		map[string]string{"status": string(txn.Status)},
		map[string]string{"status": "negotiation", "buyer_id": buyerID.String()},
		&buyerID, &ip,
	)

	c.JSON(http.StatusOK, gin.H{"message": "joined transaction successfully", "transaction_id": txn.ID})
}

func (h *TransactionHandler) UpdateStatus(c *gin.Context) {
	txnID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newStatus := models.TransactionStatus(req.Status)

	txn, err := h.txnRepo.GetByID(c.Request.Context(), txnID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}

	if !txn.Status.CanTransitionTo(newStatus) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cannot transition from " + string(txn.Status) + " to " + string(newStatus),
		})
		return
	}

	timestampFields := map[models.TransactionStatus]string{
		models.TxStatusAgreed:        "agreed_at",
		models.TxStatusFunded:        "funded_at",
		models.TxStatusShipped:       "shipped_at",
		models.TxStatusAgentReceived: "agent_received_at",
		models.TxStatusDelivered:     "delivered_at",
		models.TxStatusCompleted:     "completed_at",
		models.TxStatusCancelled:     "cancelled_at",
	}

	tsField, ok := timestampFields[newStatus]
	if !ok {
		tsField = "updated_at"
	}

	if err := h.txnRepo.UpdateStatus(c.Request.Context(), txnID, newStatus, tsField); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update status"})
		return
	}

	userID, _ := getUserID(c)
	ip := c.ClientIP()
	_ = h.auditRepo.Log(c.Request.Context(), "transaction", txnID, "status_changed",
		map[string]string{"status": string(txn.Status)},
		map[string]string{"status": string(newStatus)},
		&userID, &ip,
	)

	c.JSON(http.StatusOK, gin.H{"message": "status updated", "status": newStatus})
}

func getUserID(c *gin.Context) (uuid.UUID, bool) {
	val, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, false
	}
	id, ok := val.(uuid.UUID)
	return id, ok
}
