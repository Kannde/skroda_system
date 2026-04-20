package router

import (
	"context"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/skroda/backend/internal/config"
	"github.com/skroda/backend/internal/handlers"
	"github.com/skroda/backend/internal/middleware"
	"github.com/skroda/backend/internal/repository"
	"github.com/skroda/backend/internal/services"
)

func New(cfg *config.Config) *gin.Engine {
	gin.SetMode(cfg.GinMode)

	db, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	txRepo := repository.NewTransactionRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)

	jwtExpiry, _ := time.ParseDuration(cfg.JWTExpiry)
	authSvc := services.NewAuthService(userRepo, cfg.JWTSecret, jwtExpiry)
	escrowSvc := services.NewEscrowService(txRepo, cfg.InspectionPeriodHours)
	paymentSvc := services.NewPaymentService(paymentRepo, txRepo, cfg.MoMoAPIKey)
	agentSvc := services.NewAgentService(userRepo)

	authH := handlers.NewAuthHandler(authSvc)
	txH := handlers.NewTransactionHandler(escrowSvc)
	paymentH := handlers.NewPaymentHandler(paymentSvc)
	agentH := handlers.NewAgentHandler(agentSvc)
	disputeH := handlers.NewDisputeHandler(escrowSvc)
	webhookH := handlers.NewWebhookHandler(paymentSvc)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.RateLimit())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.AppURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	v1 := r.Group("/api/v1")

	// Public routes
	auth := v1.Group("/auth")
	{
		auth.POST("/register", authH.Register)
		auth.POST("/login", authH.Login)
	}

	// Webhooks (no auth — validated by provider signature in prod)
	webhooks := v1.Group("/webhooks")
	{
		webhooks.POST("/momo", webhookH.MoMoCallback)
	}

	// Protected routes
	protected := v1.Group("/")
	protected.Use(middleware.Auth(cfg.JWTSecret))
	{
		protected.GET("/me", authH.Me)

		txs := protected.Group("/transactions")
		{
			txs.POST("", txH.Create)
			txs.GET("", txH.List)
			txs.GET("/:id", txH.Get)
			txs.POST("/:id/confirm", txH.Confirm)
			txs.POST("/:id/cancel", txH.Cancel)
		}

		payments := protected.Group("/payments")
		{
			payments.POST("", paymentH.Initiate)
		}

		agents := protected.Group("/agents")
		{
			agents.GET("", agentH.List)
			agents.GET("/me", agentH.GetProfile)
		}

		disputes := protected.Group("/disputes")
		{
			disputes.POST("", disputeH.Create)
			disputes.POST("/:id/resolve", middleware.RequireRole("agent", "admin"), disputeH.Resolve)
		}
	}

	return r
}
