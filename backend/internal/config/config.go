package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port    string
	GinMode string

	DatabaseURL string
	RedisURL    string

	JWTSecret string
	JWTExpiry time.Duration

	MoMoAPIKey      string
	MoMoAPISecret   string
	MoMoCallbackURL string

	ResendAPIKey string
	SMSAPIKey    string
	SMSAPIURL    string

	AppURL                string
	InspectionPeriodHours int
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	jwtExpiry, err := time.ParseDuration(getEnv("JWT_EXPIRY", "24h"))
	if err != nil {
		jwtExpiry = 24 * time.Hour
	}

	inspectionHours, _ := strconv.Atoi(getEnv("INSPECTION_PERIOD_HOURS", "48"))

	return &Config{
		Port:    getEnv("PORT", "8080"),
		GinMode: getEnv("GIN_MODE", "debug"),

		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:vitaminB@localhost:5432/skroda_main_db?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),

		JWTSecret: getEnv("JWT_SECRET", "change-this-in-production"),
		JWTExpiry: jwtExpiry,

		MoMoAPIKey:      getEnv("MOMO_API_KEY", ""),
		MoMoAPISecret:   getEnv("MOMO_API_SECRET", ""),
		MoMoCallbackURL: getEnv("MOMO_CALLBACK_URL", ""),

		ResendAPIKey: getEnv("RESEND_API_KEY", ""),
		SMSAPIKey:    getEnv("SMS_API_KEY", ""),
		SMSAPIURL:    getEnv("SMS_API_URL", ""),

		AppURL:                getEnv("APP_URL", "http://localhost:5173"),
		InspectionPeriodHours: inspectionHours,
	}, nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
