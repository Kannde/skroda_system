package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/skroda/backend/internal/config"
	"github.com/skroda/backend/internal/router"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := config.Load()

	r := router.New(cfg)

	log.Printf("Starting Skroda server on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
