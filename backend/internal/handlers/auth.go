package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skroda/backend/internal/models"
	"github.com/skroda/backend/internal/repository"
	"github.com/skroda/backend/internal/services"
)

type AuthHandler struct {
	userRepo    *repository.UserRepository
	authService *services.AuthService
}

func NewAuthHandler(userRepo *repository.UserRepository, authService *services.AuthService) *AuthHandler {
	return &AuthHandler{userRepo: userRepo, authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := h.authService.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process registration"})
		return
	}

	user := &models.User{
		Phone:        req.Phone,
		PasswordHash: hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         models.UserRole(req.Role),
		Status:       models.StatusActive,
	}
	if req.City != "" {
		user.City = &req.City
	}
	if req.Region != "" {
		user.Region = &req.Region
	}

	if err := h.userRepo.Create(c.Request.Context(), user); err != nil {
		switch err {
		case repository.ErrDuplicatePhone:
			c.JSON(http.StatusConflict, gin.H{"error": "phone number already registered"})
		case repository.ErrDuplicateEmail:
			c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create account"})
		}
		return
	}

	token, err := h.authService.GenerateToken(user.ID, user.Phone, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "account created but failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, models.AuthResponse{Token: token, User: *user})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userRepo.GetByPhone(c.Request.Context(), req.Phone)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid phone or password"})
		return
	}

	if !h.authService.CheckPassword(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid phone or password"})
		return
	}

	if user.Status == models.StatusSuspended {
		c.JSON(http.StatusForbidden, gin.H{"error": "account suspended, contact support"})
		return
	}

	token, err := h.authService.GenerateToken(user.ID, user.Phone, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{Token: token, User: *user})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := getUserID(c)
	user, err := h.userRepo.GetByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}
