package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/skroda/backend/internal/models"
	"github.com/skroda/backend/internal/repository"
)

type AgentService struct {
	userRepo *repository.UserRepository
}

func NewAgentService(userRepo *repository.UserRepository) *AgentService {
	return &AgentService{userRepo: userRepo}
}

func (s *AgentService) ListAgents(ctx context.Context, city string) ([]models.User, error) {
	return s.userRepo.ListByRoleAndCity(ctx, models.RoleAgent, city)
}

func (s *AgentService) GetAgentByUserID(ctx context.Context, userID string) (*models.User, error) {
	id, _ := uuid.Parse(userID)
	return s.userRepo.GetByID(ctx, id)
}
