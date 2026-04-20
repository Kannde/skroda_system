package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/skroda/backend/internal/models"
	"github.com/skroda/backend/internal/services"
)

type DisputeHandler struct {
	escrowService *services.EscrowService
}

func NewDisputeHandler(escrowService *services.EscrowService) *DisputeHandler {
	return &DisputeHandler{escrowService: escrowService}
}

func (h *DisputeHandler) Create(c *gin.Context) {
	var req models.RaiseDisputeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	dispute, err := h.escrowService.RaiseDispute(c.Request.Context(), userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dispute)
}

func (h *DisputeHandler) Resolve(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid dispute id"})
		return
	}

	var body struct {
		Resolution string `json:"resolution" binding:"required,oneof=resolved_buyer resolved_seller"`
		Notes      string `json:"notes" binding:"required,min=10"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	agentID, _ := c.Get("user_id")
	if err := h.escrowService.ResolveDispute(
		c.Request.Context(), id, agentID.(string),
		models.DisputeStatus(body.Resolution), body.Notes,
	); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "dispute resolved"})
}
