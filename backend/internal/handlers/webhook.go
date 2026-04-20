package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skroda/backend/internal/services"
)

type WebhookHandler struct {
	paymentService *services.PaymentService
}

func NewWebhookHandler(paymentService *services.PaymentService) *WebhookHandler {
	return &WebhookHandler{paymentService: paymentService}
}

func (h *WebhookHandler) MoMoCallback(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.paymentService.HandleMoMoCallback(c.Request.Context(), payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}
