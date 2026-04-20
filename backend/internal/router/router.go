package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/skroda/backend/internal/config"
	"github.com/skroda/backend/internal/handlers"
	"github.com/skroda/backend/internal/middleware"
	"github.com/skroda/backend/internal/services"
)

type Handlers struct {
	Auth        *handlers.AuthHandler
	Transaction *handlers.TransactionHandler
}

func Setup(cfg *config.Config, h *Handlers, authService *services.AuthService) *gin.Engine {
	gin.SetMode(cfg.GinMode)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLogger())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.AppURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "skroda-api"})
	})

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", h.Auth.Register)
			auth.POST("/login", h.Auth.Login)
		}

		protected := v1.Group("")
		protected.Use(middleware.AuthRequired(authService))
		protected.Use(middleware.RateLimit(100, 1*time.Minute))
		{
			protected.GET("/me", h.Auth.Me)

			txn := protected.Group("/transactions")
			{
				txn.POST("", h.Transaction.Create)
				txn.GET("", h.Transaction.List)
				txn.GET("/:id", h.Transaction.GetByID)
				txn.PATCH("/:id/status", h.Transaction.UpdateStatus)
				txn.POST("/join/:token", h.Transaction.JoinByInvite)
			}
		}
	}

	return r
}
