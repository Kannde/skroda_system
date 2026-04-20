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
	var req models.CreateDisputeRequest
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

	var req models.ResolveDisputeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	agentID, _ := c.Get("user_id")
	if err := h.escrowService.ResolveDispute(c.Request.Context(), id, agentID.(string), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "dispute resolved"})
}
