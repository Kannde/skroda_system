package main

import (
	"log"

	"github.com/skroda/backend/internal/config"
	"github.com/skroda/backend/internal/database"
	"github.com/skroda/backend/internal/handlers"
	"github.com/skroda/backend/internal/repository"
	"github.com/skroda/backend/internal/router"
	"github.com/skroda/backend/internal/services"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	pool, err := database.NewPostgresPool(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()
	log.Println("Connected to PostgreSQL")

	userRepo := repository.NewUserRepository(pool)
	txnRepo := repository.NewTransactionRepository(pool)
	auditRepo := repository.NewAuditRepository(pool)

	authService := services.NewAuthService(cfg.JWTSecret, cfg.JWTExpiry)

	authHandler := handlers.NewAuthHandler(userRepo, authService)
	txnHandler := handlers.NewTransactionHandler(txnRepo, auditRepo)

	h := &router.Handlers{
		Auth:        authHandler,
		Transaction: txnHandler,
	}

	r := router.Setup(cfg, h, authService)

	log.Printf("Skroda API starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
