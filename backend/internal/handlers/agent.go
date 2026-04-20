package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skroda/backend/internal/services"
)

type AgentHandler struct {
	agentService *services.AgentService
}

func NewAgentHandler(agentService *services.AgentService) *AgentHandler {
	return &AgentHandler{agentService: agentService}
}

func (h *AgentHandler) List(c *gin.Context) {
	city := c.Query("city")
	agents, err := h.agentService.ListAgents(c.Request.Context(), city)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, agents)
}

func (h *AgentHandler) GetProfile(c *gin.Context) {
	userID, _ := getUserID(c)
	agent, err := h.agentService.GetAgentByUserID(c.Request.Context(), userID.String())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "agent profile not found"})
		return
	}
	c.JSON(http.StatusOK, agent)
}
