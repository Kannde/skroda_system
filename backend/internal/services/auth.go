package services

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/skroda/backend/internal/middleware"
	"github.com/skroda/backend/internal/models"
	"github.com/skroda/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
	jwtExpiry time.Duration
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string, jwtExpiry time.Duration) *AuthService {
	return &AuthService{userRepo: userRepo, jwtSecret: jwtSecret, jwtExpiry: jwtExpiry}
}

func (s *AuthService) Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
	existing, _ := s.userRepo.GetByPhone(ctx, req.Phone)
	if existing != nil {
		return nil, errors.New("phone number already registered")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	city := req.City
	region := req.Region
	user := &models.User{
		ID:           uuid.New(),
		Phone:        req.Phone,
		PasswordHash: string(hash),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         models.UserRole(req.Role),
		Status:       models.StatusPendingVerification,
		City:         &city,
		Region:       &region,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{Token: token, User: *user}, nil
}

func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	user, err := s.userRepo.GetByPhone(ctx, req.Phone)
	if err != nil || user == nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{Token: token, User: *user}, nil
}

func (s *AuthService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid user id")
	}
	return s.userRepo.GetByID(ctx, uid)
}

func (s *AuthService) generateToken(user *models.User) (string, error) {
	claims := middleware.Claims{
		UserID: user.ID.String(),
		Role:   string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.jwtExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
